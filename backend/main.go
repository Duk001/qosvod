package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"
	"encoding/json"

	"github.com/google/uuid"

	"database/sql"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	_ "github.com/denisenkom/go-mssqldb"
)


var DATABASE_LOGIN string 
var DATABASE_PASSWORD string
var TRANSKODER_ADDRESS string

var Bandwidth int = 10000000

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

type userCredentials struct {
	Username            string
	Token               string
	TokenExpirationTime time.Time
}

type loggedUsers struct {
	users []userCredentials
}

func (lu *loggedUsers) Remove(index int) {
	listLength := len(lu.users)
	if listLength < index || index < 0 {
		return
	}
	lu.users[index] = lu.users[listLength-1]
	lu.users = lu.users[:listLength-1]
}

func (lu *loggedUsers) _generateToken() string {

	Token := uuid.New().String()
	if len(lu.users) > 0 {
		return Token
	}
	flag := false
	for _, user := range lu.users {
		if Token == user.Token {
			flag = true
			break
		}
	}
	if flag {
		Token = lu._generateToken()
	}
	return Token
}
func (lu *loggedUsers) Add(username string, expirationTime int) string {
	//expiration time in hours
	addedTime := time.Duration(expirationTime)
	tokenExpirationTime := time.Now().Add(time.Hour * addedTime)
	for _, user := range lu.users {
		if username == user.Username {
			user.TokenExpirationTime = tokenExpirationTime
			return user.Token
		}
	}

	var newUserCredentials userCredentials
	newUserCredentials.Username = username
	newUserCredentials.Token = lu._generateToken()
	newUserCredentials.TokenExpirationTime = tokenExpirationTime
	lu.users = append(lu.users, newUserCredentials)

	return newUserCredentials.Token

}
func (lu *loggedUsers) FindByToken(token string) *userCredentials {
	var nullUser userCredentials
	for i, user := range lu.users {
		if token == user.Token {
			if user.TokenExpirationTime.After(time.Now()) {
				return &user
			} else {
				lu.Remove(i)
				return &nullUser
			}
		}
	}
	return &nullUser
}
func (lu *loggedUsers) DeleteByToken(token string) bool {
	for i, user := range lu.users {
		if token == user.Token {
			lu.Remove(i)
			return true
		}
	}
	return false
}

type credentials struct {
	Name string
	Key  string
}

type videoQuality struct {
	Bitrate      string `json:"Bitrate"`
	Resolution   string `json:"Resolution"`
	bitrateValue int
}

type filmVideoQuality struct {
	vq map[string]*[]videoQuality
}

func (fvq *filmVideoQuality) Initiate(filmID string) {
	if fvq.vq == nil {
		fvq.vq = make(map[string]*[]videoQuality)
	} else if _, ok := fvq.vq[filmID]; ok {
		return
	}

	sort := func(data []videoQuality) []videoQuality { //? Insertion Sort	descending
		cmpBitrate := func(a, b string) bool {
			i, _ := strconv.Atoi(strings.Replace(a, "k", "", -1))
			j, _ := strconv.Atoi(strings.Replace(b, "k", "", -1))
			return i > j
		}
		i := 1
		for i < len(data) {
			j := i
			for j >= 1 && cmpBitrate(data[j].Bitrate, data[j-1].Bitrate) {
				data[j], data[j-1] = data[j-1], data[j]
				j -= 1
			}
			i += 1
		}
		return data
	}
	vq, _ := getFilmQuality(filmID)
	vq = sort(vq)
	fvq.vq[filmID] = &vq
}
func (fvq *filmVideoQuality) Get(filmID string) []videoQuality {
	data := fvq.vq[filmID]
	return *data
}

type userVideoQualityLogger struct {
	BandwidthArray              []int `json:"BandwidthArray"`
	AdjustedBandwidthArray      []int `json:"AdjustedBandwidthArray"`
	LocalStandardDeviationArray []int `json:"LocalStandardDeviationArray"`
	SetBitrateArray             []int `json:"SetBitrateArray"`
}
type videoQualityLogger struct {
	usqMap map[string]*userVideoQualityLogger
}

func (vql *videoQualityLogger) Initiate(token string) {
	if vql.usqMap == nil {
		vql.usqMap = make(map[string]*userVideoQualityLogger)
	}
	//var tmpArray []int
	u := userVideoQualityLogger{}
	vql.usqMap[token] = &u
}
func (vql *videoQualityLogger) Update(token string, bandwidth, adjustedBandwidth, standardDeviation int) {
	u, ok := vql.usqMap[token]
	if !ok {
		return
	}
	u.BandwidthArray = append(u.BandwidthArray, bandwidth)
	u.AdjustedBandwidthArray = append(u.AdjustedBandwidthArray, adjustedBandwidth)
	u.LocalStandardDeviationArray = append(u.LocalStandardDeviationArray, standardDeviation)
	vql.usqMap[token] = u
}
func (vql *videoQualityLogger) UpdateBitRate(token string, bitRate int) {
	u, ok := vql.usqMap[token]
	if !ok {
		return
	}
	u.SetBitrateArray = append(u.SetBitrateArray, bitRate)
	vql.usqMap[token] = u
}
func (vql *videoQualityLogger) Jsonify() []byte {
	jsonData, _ := json.Marshal(vql.usqMap)
	return jsonData
}

type usc struct {
	CurrentBandwidth  int
	FilmId            string
	Date              time.Time
	SegmentsHit       int
	LastBandwidth     int
	AverageBandwidth  int
	BandwidthSum      int
	BandwidthArray    []int
	StandardDeviation int
	FilmMaxBitrate    int
	BufferLength      float32
	qualityLevel      int
	lastSetBitrate    int
	lastBufferLength  float32
}
type usersConectionQuality struct {
	uscMap map[string]*usc
}

func (uscq *usersConectionQuality) Initiate(token, filmID string, fvq filmVideoQuality) {
	if uscq.uscMap == nil {
		uscq.uscMap = make(map[string]*usc)
	}
	var tmpArray []int
	filmQuality := fvq.Get(filmID)
	bitrateValueString := strings.ReplaceAll(filmQuality[0].Bitrate, "k", "")
	filmMaxBitrate, err := strconv.Atoi(bitrateValueString)
	filmMaxBitrate *= 1000
	if err != nil {
		filmMaxBitrate = -1
		log.Println("Error converting max video bitrate: \n", err)
	}
	u := usc{-1, filmID, time.Now(), 0, -1, -1, 0, tmpArray, 0, filmMaxBitrate, 0, 0, 0, 0}

	uscq.uscMap[token] = &u
	fmt.Println("TEST")
}
func (uscq *usersConectionQuality) Delete(token string) {
	delete(uscq.uscMap, token)
}
func (uscq *usersConectionQuality) GetBandwidth(token string) (int, bool) {
	u, ok := uscq.uscMap[token]
	return u.CurrentBandwidth, ok
}
func (uscq *usersConectionQuality) GetBuffer(token string) (int, bool) {
	u, ok := uscq.uscMap[token]
	return int(u.BufferLength), ok
}

func (uscq *usersConectionQuality) ConnectionQualityString(token string) string {
	out := ""
	u, ok := uscq.uscMap[token]
	if !ok {
		return out
	}
	out = fmt.Sprintf("Current Bandwidth: %d\n Segments Hit: %d\n Average Bandwitdh: %d\n Standard Deviation: %d", u.CurrentBandwidth, u.SegmentsHit, u.AverageBandwidth, u.StandardDeviation)
	return out
}
func (uscq *usersConectionQuality) ChangeBandwidth(token string, bandwidth string, buffer string) bool {
	//var err error
	numberOfLocalSegments := 10
	u, ok := uscq.uscMap[token]
	if !ok {
		return ok
	}
	if bandwidth != "-1" {

		if u.SegmentsHit > 0 {
			u.LastBandwidth = u.CurrentBandwidth
		}

		tmpBandwidth, err := strconv.Atoi(bandwidth)
		if err != nil {
			tmpBandwidth = 0
		}
		bufferFloat, err := strconv.ParseFloat(buffer, 32)
		if err != nil {
			bufferFloat = 0
		}
		u.BufferLength = float32(bufferFloat)

		if tmpBandwidth > u.FilmMaxBitrate*2 {
			u.CurrentBandwidth = u.FilmMaxBitrate * 2
			ok = false
		} else {
			u.CurrentBandwidth = tmpBandwidth
			ok = false
		}
		// Cap bandwidth

		u.BandwidthSum += u.CurrentBandwidth
		u.SegmentsHit += 1

		if u.SegmentsHit > numberOfLocalSegments {
			u.BandwidthArray[u.SegmentsHit%numberOfLocalSegments] = u.CurrentBandwidth
		} else {
			u.BandwidthArray = append(u.BandwidthArray, u.CurrentBandwidth)
		}

		u.AverageBandwidth = func() int {
			average := 0
			for _, elem := range u.BandwidthArray {
				average += elem
			}
			return average / len(u.BandwidthArray)
		}() //u.BandwidthSum/u.SegmentsHit

		u.StandardDeviation = func() int {
			variance := 0
			for _, elem := range u.BandwidthArray {
				variance += (u.AverageBandwidth - elem) * (u.AverageBandwidth - elem)
			}
			variance /= len(u.BandwidthArray)
			return int(math.Sqrt(float64(variance)))
		}()
		//u.StandardDeviation = int(math.Sqrt(float64(u.VarianceSum/u.SegmentsHit)))

		uscq.uscMap[token] = u
	}
	return ok
}
func (uscq *usersConectionQuality) SetVideoQuality(token string, qualityList []videoQuality, lowBufferState, highBufferState int) (int, string) {
	//segmentLength := 2
	segmentsHitThreshold := 5
	//lowBufferState := 12
	//highBufferState := 21
	containerName := ""
	u, ok := uscq.uscMap[token]
	if !ok {
		return 0, containerName
	}
	bandwidth := u.CurrentBandwidth
	//fmt.Println("Buffer: ", u.BufferLength)
	if u.SegmentsHit < segmentsHitThreshold {

		setBitRate := 0
		for i, vq := range qualityList {
			value, err := strconv.Atoi(strings.Replace(vq.Bitrate, "k", "", -1))

			if err != nil {
				value = 0
			}
			if value < bandwidth/1000 {
				setBitRate = value
				containerName = strings.ToLower(u.FilmId + "-" + strings.Replace(vq.Resolution, ":", "-", -1) + "-" + vq.Bitrate)
				fmt.Println("Change bitrate to: ", value, "k")
				u.qualityLevel = len(qualityList) - i
				break
			}
		}
		u.lastSetBitrate = setBitRate * 1000
		return u.lastSetBitrate, containerName
	}

	qualityLevel := u.qualityLevel
	maxQualityLevel := len(qualityList) - 1

	if u.BufferLength < float32(lowBufferState) { //! low state
		if u.CurrentBandwidth > u.lastSetBitrate {

			if u.CurrentBandwidth > 2*u.lastSetBitrate {
				qualityLevel = qualityLevel + 1
			}

		} else {
			if u.CurrentBandwidth < u.lastSetBitrate/2 {
				qualityLevel = qualityLevel - 3
			} else {
				qualityLevel = qualityLevel - 1
			}
		}
	} else if u.BufferLength >= float32(lowBufferState) && u.BufferLength < float32(highBufferState) { //! optimum state
		if u.CurrentBandwidth > u.lastSetBitrate*800 {
			qualityLevel = qualityLevel + 0
		} else if u.CurrentBandwidth > u.lastSetBitrate/800 {
			qualityLevel = qualityLevel - 0
		}

	} else { //! high state
		nextBitrateIndex := (maxQualityLevel - (qualityLevel + 1)) % u.FilmMaxBitrate
		if nextBitrateIndex < 0 {
			nextBitrateIndex = 0
		}
		nextBitrate := qualityList[nextBitrateIndex].bitrateValue
		if u.CurrentBandwidth > u.lastSetBitrate && u.CurrentBandwidth > nextBitrate && nextBitrate != u.lastSetBitrate {
			if u.CurrentBandwidth > 2*u.lastSetBitrate {
				qualityLevel = qualityLevel + 2
			} else {
				qualityLevel = qualityLevel + 1
			}

		}

	}

	if qualityLevel > maxQualityLevel {
		qualityLevel = maxQualityLevel
	} else if qualityLevel < 0 {
		qualityLevel = 0
	}

	vq := qualityList[maxQualityLevel-qualityLevel]


	containerName = strings.ToLower(u.FilmId + "-" + strings.Replace(vq.Resolution, ":", "-", -1) + "-" + vq.Bitrate)
	u.qualityLevel = qualityLevel

	u.lastSetBitrate = vq.bitrateValue //bitrateValue * 1000
	return u.lastSetBitrate, containerName
}

//var Credentials = credentials{Name: "piqosvod", Key: "7Yd0mfoGmBjfRE3DZCKlr6YlTfH5sNKIvMeL+zkGLiyAjjCwml/kO5k4ZF85PzEePcoKWyGh64HD+5zl8/j/vw=="}
var Credentials credentials
var serviceClient azblob.ServiceClient
var ctx context.Context
var cred *azblob.SharedKeyCredential
var db *sql.DB

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func initDatabaseAbstraction() {

	var err error
	server := os.Getenv("DATABASE_URL")

	var port = 1433
	var user = DATABASE_LOGIN
	var password = DATABASE_PASSWORD
	var database = os.Getenv("DATABASE_NAME")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	db, err = sql.Open("sqlserver",
		connString)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected!")
}

func getListOfFilmIDs() []string {
	var id string
	var listOfIDs []string
	rows, err := db.Query("select id from films")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err == nil {
			listOfIDs = append(listOfIDs, id)
		}
	}
	return listOfIDs
}
func getListOfFilmIDsByCategory(categories string) []string {
	var id string
	var listOfIDs []string
	listOfCategories := strings.Split(categories, ",")
	qms, listOfCategoriesInterface := func() (string, []interface{}) {
		out := ""
		outInterface := []interface{}{}

		for i := 1; i <= len(listOfCategories); i++ {
			tmp := listOfCategories[i-1]
			outInterface = append(outInterface, tmp)
			out += fmt.Sprintf("@p%d,", i)
		}
		return out[:len(out)-1], outInterface
	}()

	q := "SELECT id FROM films WHERE category IN ( " + qms + " );" //TODO mssql: The data types text and nvarchar are incompatible in the equal to operator.
	rows, err := db.Query(q, listOfCategoriesInterface...)
	if err != nil {
		log.Println("getListOfFilmIDsByCategory -> \n", err)
		return listOfIDs
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err == nil {
			listOfIDs = append(listOfIDs, id)
		}
	}
	return listOfIDs
}

func getListOfFilmCategories() []string {
	var name string
	var listOfCategories []string
	rows, err := db.Query("select name from categories")
	if err != nil {

		log.Println("Error reading film categories, returning empty,\nerror msg -> ", err)
		return listOfCategories
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err == nil {
			listOfCategories = append(listOfCategories, name)
		}
	}
	return listOfCategories
}
func postNewFilmCategory(category string) error {
	value := strings.ToLower(category)
	_, err := db.Query("insert into categories values ( @p1 )", value)
	return err
}
func generateUUID(table string) string {
	id := uuid.New().String()
	//id = "id123"
	row, err := db.Query("select id from "+table+" where id = @p1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&id)
	if err == nil {
		id = generateUUID(table)
	}
	return id
}

type filmData struct {
	Name, Id, Owner, Description, Category, Quality string
	visible                                         bool
}

func (fd *filmData) jsonify() []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(fd)
	if err != nil {
		log.Println(err)
	}
	return jsonData
}

func newFilmData(Name, Owner, Description, Category, Quality string) (filmData, error) {

	//var data filmData
	var err error
	id := generateUUID("films")

	data := filmData{Name: Name, Owner: Owner, Description: Description, Category: strings.ToLower(Category), Id: id, Quality: Quality, visible: true}

	return data, err
}

type rawFilmData struct {
	Name        string `json:"title"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Quality     string `json:"quality"`
}

func (r *rawFilmData) createFilmData() (filmData, error) {
	data, err := newFilmData(r.Name, r.Owner, r.Description, r.Category, r.Quality)
	return data, err
}
func getDataOfFilmByID(id string) (filmData, error) {
	var data filmData
	row, err := db.Query("select id, name, owner, description, category,quality, visible from films where id = @p1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&data.Id, &data.Name, &data.Owner, &data.Description, &data.Category, &data.Quality, &data.visible)

	if err == nil {
		return data, nil
	}
	return data, errors.New("Film with id: " + id + " not found in database" + "		->		" + err.Error())
}
func getFilmQuality(id string) ([]videoQuality, error) { //TODO Add film map with video quality settings
	var data []videoQuality
	var tmpData string

	row, err := db.Query("select quality from films where id = @p1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&tmpData)

	if err == nil && tmpData != "" {
		tmpData = strings.ReplaceAll(tmpData, " ", "")
		stringArray := strings.Split(tmpData, ",")
		for _, elem := range stringArray {
			//1920:1080-4000k
			splitElem := strings.Split(elem, "-")
			bitrateValue, err := strconv.Atoi(strings.Replace(splitElem[1], "k", "", -1))
			if err != nil {
				bitrateValue = 0
			}
			data = append(data, videoQuality{splitElem[1], splitElem[0], bitrateValue * 1000})
		}
		return data, nil
	}
	return data, errors.New("Film with id: " + id + " not found in database" + "		->		" + err.Error())
}
func postFilmData(data filmData) error {
	var name string
	row, err := db.Query("select name from categories where name = @p1", data.Category)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&name)
	if err != nil {
		err = postNewFilmCategory(data.Category)
		if err != nil {
			return err
		}
		//na = generateUUID(table)
	}
	_, err = db.Query("insert into films values (@p1, @p2, @p3, @p4, @p5, @p6,@p7)", data.Name, data.Id, data.Owner, data.Description, data.Category, data.visible, data.Quality)

	return err
}
func deleteFilm(data filmData) []error {
	var errors []error
	filmQualityList := strings.Split(data.Quality, ",")
	for _, q := range filmQualityList {
		q = strings.Trim(q, " ")
		out := strings.Split(q, "-")
		resolution, bitrate := out[0], out[1]
		containerName := strings.ToLower(data.Id + "-" + strings.Replace(resolution, ":", "-", -1) + "-" + bitrate)
		container := serviceClient.NewContainerClient(containerName)
		_, err := container.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr container : ", containerName)
			errors = append(errors, err)
		}

		containerManifests := serviceClient.NewContainerClient("manifests")
		manifestName := data.Id + ".m3u8"
		blobManifest := containerManifests.NewBlobClient(manifestName)
		_, err = blobManifest.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr blob(manifest) : ", containerName)
			errors = append(errors, err)
		}

		containerPosters := serviceClient.NewContainerClient("posters")
		posterName := data.Id
		blobPoster := containerPosters.NewBlobClient(posterName)
		_, err = blobPoster.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr blob(poster) : ", containerName)
			errors = append(errors, err)
		}

		_, err = db.Query("DELETE FROM films WHERE id = @p1", data.Id)
		if err != nil {
			log.Println(err)
			errors = append(errors, err)
		}

	}
	return errors
}

type server struct {
	name                  string
	UsersConectionQuality usersConectionQuality
	VideoQualityLogger    videoQualityLogger
	FilmVideoQuality      filmVideoQuality
	LoggedUsers           loggedUsers
}

func (s *server) homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage! "+s.name)
	fmt.Println("Endpoint Hit: homePage")
}
func (s *server) loginUser(login, password string) (string, bool) {
	var dbPassword string
	row, err := db.QueryContext(ctx, "select password from users where login = @p1", login)
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
func (s *server) videoManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	fmt.Println("Endpoint Hit: videoManifest")
	filmName, ok := r.URL.Query()["film"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	containerName := "manifests"
	container := serviceClient.NewContainerClient(containerName)
	blockBlob := container.NewBlockBlobClient(filmName[0] + ".m3u8")
	get, err := blockBlob.Download(ctx, nil)
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
func (s *server) videoSegment(w http.ResponseWriter, r *http.Request) {
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
		container := serviceClient.NewContainerClient(containerName)
		blockBlob := container.NewBlockBlobClient(segNumber[0])
		get, err := blockBlob.Download(ctx, nil)
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
func (s *server) bandwidthMonitor(w http.ResponseWriter, r *http.Request) {
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
func (s *server) initFilmSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "POST":
		token := r.Header.Get("token")
		//log.Println(LoggedUsers.FindByToken(token))
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
		s.FilmVideoQuality.Initiate(data.FilmID)
		s.UsersConectionQuality.Initiate(token, data.FilmID, s.FilmVideoQuality)
		s.VideoQualityLogger.Initiate(token)

		w.WriteHeader(http.StatusOK)
	}
}
func (s *server) filmsCategoriesListEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	switch r.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := getListOfFilmCategories()

		dataJson, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dataJson)
	}
}
func (s *server) filmsIdListEndpoint(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: list of films")
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := getListOfFilmIDs()

	dataJson, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dataJson)
}
func (s *server) filmsIdListByCategoryEndpoint(w http.ResponseWriter, r *http.Request) {
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
		data := getListOfFilmIDsByCategory(categories[0])

		dataJson, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dataJson)
	}
}
func (s *server) filmDataEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	//enableCors(&w)
	switch r.Method {
	case "GET":
		idList, ok := r.URL.Query()["id"]
		//enableCors(&w)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := idList[0]
		data, err := getDataOfFilmByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data.jsonify())
	case "POST":
		var data rawFilmData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Decoding error: ", err)
			return
		}
		filmData, err := data.createFilmData()
		if err != nil {
			fmt.Println("Error when creating film data from raw film data: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Print(filmData)
		//!!!!!!
		err = postFilmData(filmData)
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
func (s *server) filmFileEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Println("Not here...")
	case "POST":
		var file_name string
		fmt.Println(r.ContentLength)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

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

		res, err := http.Post(TRANSKODER_ADDRESS+"/upload?name="+file_name, "application/octet-stream", outReader)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			//TODO  Clean database and azr containers
			_, err = db.Query("DELETE FROM films WHERE id = @p1", file_name)
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
func (s *server) deleteFilmEndpoint(w http.ResponseWriter, r *http.Request) {
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
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := idList[0]

		data, err := getDataOfFilmByID(id)
		if user.Username != data.Owner && user.Username != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err != nil {
			log.Println("deleteFilmEndpoint() -> Error geting film data\nerr -> ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		errors := deleteFilm(data)
		if errors != nil {
			for _, err := range errors {
				log.Println("deleteFilmEndpoint() -> Error deleting film \nerr -> ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

	}
}
func (s *server) userLoginEndpoint(w http.ResponseWriter, r *http.Request) {
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
func (s *server) tokenCheckEndpoint(w http.ResponseWriter, r *http.Request) {
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
func (s *server) logoutEndpoint(w http.ResponseWriter, r *http.Request) {
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
func (s *server) usersVideoQualityLogEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, token")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write(s.VideoQualityLogger.Jsonify())
	}
}
func (s *server) filmQualityEndpoint(w http.ResponseWriter, r *http.Request) {
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
			filmQuality, err := getFilmQuality(film_name)
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
func (s *server) filmPosterEndpoint(w http.ResponseWriter, r *http.Request) {
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

		container := serviceClient.NewContainerClient("posters")
		blockBlob := container.NewBlockBlobClient(film_name)
		get, err := blockBlob.Download(ctx, nil)
		if err != nil {
			log.Println("Error downloading filme poster, loading default")
			blockBlob := container.NewBlockBlobClient("default-poster.jpg")
			get, err = blockBlob.Download(ctx, nil)
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

		container := serviceClient.NewContainerClient("posters")
		blockBlob := container.NewBlockBlobClient(film_name)
		//data := bytes.NewReader(manifest) //buf//bytes.NewReader(&buf)
		_, err = blockBlob.Upload(ctx, NopCloser(data), nil)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		//check(err)

	}
}
func (s *server) ABRTest(w http.ResponseWriter, r *http.Request) {
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

// func (s *server)

func handleRequests() {
	var s server
	s.name = "Piqosvod"
	http.HandleFunc("/", s.homePage)
	http.HandleFunc("/videoManifest", s.videoManifest)
	http.HandleFunc("/videoSegment", s.videoSegment)
	http.HandleFunc("/bandwidth", s.bandwidthMonitor)
	http.HandleFunc("/categories", s.filmsCategoriesListEndpoint)
	http.HandleFunc("/films", s.filmsIdListEndpoint)
	http.HandleFunc("/filmsByCategory", s.filmsIdListByCategoryEndpoint)
	http.HandleFunc("/film", s.filmDataEndpoint)
	http.HandleFunc("/filmFile", s.filmFileEndpoint)
	http.HandleFunc("/deleteFilm", s.deleteFilmEndpoint)
	http.HandleFunc("/initFilmSession", s.initFilmSession)
	http.HandleFunc("/login", s.userLoginEndpoint)
	http.HandleFunc("/logout", s.logoutEndpoint)
	http.HandleFunc("/tokenCheck", s.tokenCheckEndpoint)
	http.HandleFunc("/qualityLog", s.usersVideoQualityLogEndpoint)
	http.HandleFunc("/filmQuality", s.filmQualityEndpoint)
	http.HandleFunc("/filmPoster", s.filmPosterEndpoint)
	http.HandleFunc("/ABR_DEBUG", s.ABRTest)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))

}

func main() {
	var err error
	DATABASE_LOGIN = os.Getenv("DATABASE_LOGIN")
	DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")
	TRANSKODER_ADDRESS = os.Getenv("TRANSKODER_ADDRESS")
	blob_key := os.Getenv("BLOB_KEY")
	blob_acc_name := os.Getenv("BLOB_NAME")
	blob_url := os.Getenv("BLOB_URL")

	fmt.Println("Transkoder_Address: ", TRANSKODER_ADDRESS)
	fmt.Println("blob address: ",blob_url)
	Credentials = credentials{Name: blob_acc_name, Key: blob_key}
	initDatabaseAbstraction()

	cred, err = azblob.NewSharedKeyCredential(Credentials.Name, Credentials.Key)
	check(err)
	serviceClient, err = azblob.NewServiceClientWithSharedKey(blob_url, cred, nil)
	check(err)
	ctx = context.Background() // context
	handleRequests()

}
