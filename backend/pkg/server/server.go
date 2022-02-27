package server

import (
	"bytes"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"context"
	"encoding/json"
	"database/sql"
	"reflect"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Duk001/qosvod/backend/pkg/film"
	"github.com/Duk001/qosvod/backend/pkg/loggedusers"
	"github.com/Duk001/qosvod/backend/pkg/usersconnectionquality"
	"github.com/Duk001/qosvod/backend/pkg/videoquality"
	_ "github.com/denisenkom/go-mssqldb"
)
func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Azr version of nopCloser
type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

type Server struct {
	Name                  string
	UsersConectionQuality *usersconnectionquality.UsersConnectionQuality
	VideoQualityLogger    *videoquality.VideoQualityLogger
	FilmVideoQuality      *videoquality.FilmVideoQuality
	LoggedUsers           *loggedusers.LoggedUsers
	db                    *sql.DB
	ctx                   context.Context
	serviceClient         azblob.ServiceClient
	transkoderAddress     string
}

func (s *Server) Init(db *sql.DB, ctx context.Context, serviceClient azblob.ServiceClient, transkoderAddress string) {
	s.db = db
	s.ctx = ctx
	s.serviceClient = serviceClient
	s.transkoderAddress = transkoderAddress
	s.LoggedUsers = &loggedusers.LoggedUsers{}
	s.FilmVideoQuality = &videoquality.FilmVideoQuality{}
	s.VideoQualityLogger = &videoquality.VideoQualityLogger{}
	s.UsersConectionQuality = &usersconnectionquality.UsersConnectionQuality{}
}

func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage! "+s.Name)
	fmt.Println("Endpoint Hit: homePage")
}
func (s *Server) loginUser(login, password string) (string, bool) {
	var dbPassword string
	row, err := s.db.QueryContext(s.ctx, "select password from users where login = @p1", login)
	if err != nil {
		return "", false
		//log.Fatal(err)
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&dbPassword)

	if err != nil {
		return "", false
	}
	if dbPassword == password {
		token := s.LoggedUsers.Add(login, 1)
		return token, true
	}

	return "", false
}

// Endpoint for downloading video manifest file
func (s *Server) VideoManifestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	fmt.Println("Endpoint Hit: videoManifest")
	filmName, ok := r.URL.Query()["film"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	containerName := "manifests"
	container := s.serviceClient.NewContainerClient(containerName)
	blockBlob := container.NewBlockBlobClient(filmName[0] + ".m3u8")
	get, err := blockBlob.Download(s.ctx, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	check(err)

	data := &bytes.Buffer{}

	reader := get.Body(azblob.RetryReaderOptions{}) //
	_, err = data.ReadFrom(reader)
	check(err)
	err = reader.Close()
	check(err)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data.Bytes())
}

// Endpoint for downloading video segment
func (s *Server) VideoSegmentEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		fmt.Println("\nEndpoint Hit: videoSegment")
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)

		segNumber, ok := r.URL.Query()["seg"]
		if !ok { // ok == false
			return
		}
		filmName, ok := r.URL.Query()["filmName"]
		if !ok {
			return
		}
		fmt.Println("Segment: ", segNumber[0], "	Film: ", filmName[0])

		videoQualityList := s.FilmVideoQuality.Get(filmName[0])

		setBitrate, containerName := s.UsersConectionQuality.SetVideoQuality(token, videoQualityList, 12, 21)
		if containerName == "" {
			containerName = strings.ToLower(filmName[0] + "-" + strings.Replace(videoQualityList[len(videoQualityList)-1].Resolution, ":", "-", -1) + "-" + videoQualityList[len(videoQualityList)-1].Bitrate)
		}

		s.VideoQualityLogger.UpdateBitRate(token, setBitrate)
		s.VideoQualityLogger.Jsonify()
		container := s.serviceClient.NewContainerClient(containerName)
		blockBlob := container.NewBlockBlobClient(segNumber[0])
		get, err := blockBlob.Download(s.ctx, nil)
		check(err)

		data := &bytes.Buffer{}

		reader := get.Body(azblob.RetryReaderOptions{}) //
		_, err = data.ReadFrom(reader)
		check(err)
		err = reader.Close()
		check(err)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data.Bytes())
	}
}
func (s *Server) BandwidthMonitorEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
	case "POST":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)
		//fmt.Println(r.ContentLength)
		type respData struct {
			Bandwidth string
			Buffer    string
		}
		var data respData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Decoding error: ", err)
			return
		}
		s.UsersConectionQuality.ChangeBandwidth(token, data.Bandwidth, data.Buffer)
		//fmt.Println("Bandwidth: ", data.Bandwidth)
		w.WriteHeader(http.StatusOK)
	}
}
// Endpoint initiates film session, needed for every film session.
func (s *Server) InitFilmSessionEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "POST":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)
		type respData struct {
			FilmID string
		}
		var data respData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Decoding error: ", err)
			return
		}
		s.FilmVideoQuality.Initiate(data.FilmID, s.db)
		s.UsersConectionQuality.Initiate(token, data.FilmID, *s.FilmVideoQuality)
		s.VideoQualityLogger.Initiate(token)

		w.WriteHeader(http.StatusOK)
	}
}
func (s *Server) FilmsCategoriesListEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	switch r.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := film.GetListOfFilmCategories(s.db)

		dataJson, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dataJson)
	}
}
func (s *Server) FilmsIdListEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: list of films")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := film.GetListOfFilmIDs(s.db)

	dataJson, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dataJson)
}
func (s *Server) FilmsIdListByCategoryEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	switch r.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		categories, ok := r.URL.Query()["category"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := film.GetListOfFilmIDsByCategory(categories[0], s.db)

		dataJson, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dataJson)
	}
}
func (s *Server) FilmDataEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	//enableCors(&w)
	switch r.Method {
	case "GET":
		idList, ok := r.URL.Query()["id"]
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		//enableCors(&w)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := idList[0]
		data, err := film.GetDataOfFilmByID(id, s.db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data.Jsonify())
	case "POST":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)


		var data film.RawFilmData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Decoding error: ", err)
			return
		}
		filmData, err := data.CreateFilmData(s.db)
		if err != nil {
			fmt.Println("Error when creating film data from raw film data: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Print(filmData)
		//!!!!!!
		err = film.PostFilmData(filmData, s.db)
		if err != nil {
			fmt.Println("Error when sending film data to azr: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//!!!!!!
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(filmData.Id))
		//return filmData.Id
	}
}
func (s *Server) FilmFileEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Println("Not here...")
	case "POST":
		var file_name string
		fmt.Println(r.ContentLength)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)

		
		if  user.Username != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		file_name_array, ok := r.URL.Query()["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		file_name = file_name_array[0]
		file, _, err := r.FormFile("file")

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()
		out, err := ioutil.ReadAll(file)
		check(err)

		outReader := bytes.NewReader(out)

		res, err := http.Post(s.transkoderAddress+"/upload?name="+file_name, "application/octet-stream", outReader)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = s.db.Query("DELETE FROM films WHERE id = @p1", file_name)
			if err != nil {
				log.Println(err)
			}
			return
		}
		//check(err)
		fmt.Println(res)

		w.WriteHeader(http.StatusOK)
	}
}
func (s *Server) DeleteFilmEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "DELETE":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)

		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)
		idList, ok := r.URL.Query()["id"]
		//enableCors(&w)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := idList[0]

		data, err := film.GetDataOfFilmByID(id, s.db)
		if user.Username != data.Owner && user.Username != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err != nil {
			log.Println("deleteFilmEndpoint() -> Error geting film data\nerr -> ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		errors := film.DeleteFilm(data, s.db, s.serviceClient, s.ctx)
		if errors != nil {
			for _, err := range errors {
				log.Println("deleteFilmEndpoint() -> Error deleting film \nerr -> ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

	}
}
func (s *Server) UserLoginEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch r.Method {
	case "POST":
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

		type loginData struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		var data loginData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Decoding error: ", err)
			return
		}
		token, ok := s.loginUser(data.Login, data.Password)
		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(token))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))

	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (s *Server) TokenCheckEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		if user.Token == "" {
			w.Write([]byte(""))
			return
		}
		w.Write([]byte(token))
	}
}
func (s *Server) LogoutEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		token := r.Header.Get("token")
		ok := s.LoggedUsers.DeleteByToken(token)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	}
}
// Get video quality log of specific user, depreciated 
func (s *Server) UsersVideoQualityLogEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write(s.VideoQualityLogger.Jsonify())
	}
}
// Get list of available video qualities of a film. 
func (s *Server) FilmQualityEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
		switch r.Method {
		case "GET":
			film_name_array, ok := r.URL.Query()["film"]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			film_name := film_name_array[0]
			filmQuality, err := videoquality.GetFilmQuality(film_name, s.db)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
			filmQualityJSON, err := json.Marshal(filmQuality)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(filmQualityJSON)
		}
	}
}
func (s *Server) FilmPosterEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		film_name_array, ok := r.URL.Query()["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		film_name := film_name_array[0]

		container := s.serviceClient.NewContainerClient("posters")
		blockBlob := container.NewBlockBlobClient(film_name)
		get, err := blockBlob.Download(s.ctx, nil)
		if err != nil {
			log.Println("Error downloading filme poster, loading default")
			blockBlob := container.NewBlockBlobClient("default-poster.jpg")
			get, err = blockBlob.Download(s.ctx, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		//check(err)

		data := &bytes.Buffer{}

		reader := get.Body(azblob.RetryReaderOptions{}) //
		_, err = data.ReadFrom(reader)
		check(err)
		err = reader.Close()
		check(err)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data.Bytes())

	case "POST":
		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)

		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)

		film_name_array, ok := r.URL.Query()["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		film_name := film_name_array[0]
		file, _, err := r.FormFile("file")

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()
		out, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		data := bytes.NewReader(out)

		container := s.serviceClient.NewContainerClient("posters")
		blockBlob := container.NewBlockBlobClient(film_name)
		//data := bytes.NewReader(manifest) //buf//bytes.NewReader(&buf)
		_, err = blockBlob.Upload(s.ctx, NopCloser(data), nil)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		//check(err)

	}
}
// Algorithm testing endpoint, allows to test ABR algorithm without downloading real segments.
func (s *Server) ABRTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":

		token := r.Header.Get("token")
		user := s.LoggedUsers.FindByToken(token)
		if reflect.ValueOf(user.Token).IsZero() || time.Now().After(user.TokenExpirationTime) { // login check
			w.WriteHeader(http.StatusForbidden)
			return
		}
		user.TokenExpirationTime = time.Now().Add(time.Hour)
		filmName, ok := r.URL.Query()["filmName"]
		if !ok {
			return
		}
		lowString, ok := r.URL.Query()["low"]
		if !ok {
			return
		}
		highString, ok := r.URL.Query()["high"]
		if !ok {
			return
		}

		low, _ := strconv.Atoi(lowString[0])
		high, _ := strconv.Atoi(highString[0])

		videoQualityList := s.FilmVideoQuality.Get(filmName[0])

		setBitrate, _ := s.UsersConectionQuality.SetVideoQuality(token, videoQualityList, low, high)

		s.VideoQualityLogger.UpdateBitRate(token, setBitrate)
		s.VideoQualityLogger.Jsonify()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		data := fmt.Sprint(setBitrate)
		w.Write([]byte(data))
	}
}
