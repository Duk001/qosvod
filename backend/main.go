package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"context"
	"database/sql"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Duk001/qosvod/backend/pkg/server"
	// _ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var DATABASE_LOGIN string
var DATABASE_PASSWORD string
var TRANSKODER_ADDRESS string

type credentials struct {
	Name string
	Key  string
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
	if server == "" {
		server = "127.0.0.1"
	}

	var port = 3306
	var user = DATABASE_LOGIN
	var password = DATABASE_PASSWORD
	var database = os.Getenv("DATABASE_NAME")
	if database == "" {
		database = "/qosvod"
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	connString = "root@tcp(127.0.0.1:3306)/qosvod"

	// connString := "root@localhost"
	db, err = sql.Open("mysql", //"sqlserver" -> azure sql
		connString)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected!")
}

// func (s *server)

func handleRequests() {
	var s server.Server
	s.Init(db, ctx, serviceClient, TRANSKODER_ADDRESS)
	s.Name = "qosvod"
	http.HandleFunc("/", s.HomePage)
	http.HandleFunc("/videoManifest", s.VideoManifest)
	http.HandleFunc("/videoSegment", s.VideoSegment)
	http.HandleFunc("/bandwidth", s.BandwidthMonitor)
	http.HandleFunc("/categories", s.FilmsCategoriesListEndpoint)
	http.HandleFunc("/films", s.FilmsIdListEndpoint)
	http.HandleFunc("/filmsByCategory", s.FilmsIdListByCategoryEndpoint)
	http.HandleFunc("/film", s.FilmDataEndpoint)
	http.HandleFunc("/filmFile", s.FilmFileEndpoint)
	http.HandleFunc("/deleteFilm", s.DeleteFilmEndpoint)
	http.HandleFunc("/initFilmSession", s.InitFilmSession)
	http.HandleFunc("/login", s.UserLoginEndpoint)
	http.HandleFunc("/logout", s.LogoutEndpoint)
	http.HandleFunc("/tokenCheck", s.TokenCheckEndpoint)
	http.HandleFunc("/qualityLog", s.UsersVideoQualityLogEndpoint)
	http.HandleFunc("/filmQuality", s.FilmQualityEndpoint)
	http.HandleFunc("/filmPoster", s.FilmPosterEndpoint)
	http.HandleFunc("/ABR_DEBUG", s.ABRTest)

	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		httpPort = "11000"
	}
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))

}

func main() {
	var err error

	DATABASE_LOGIN = "root"         //os.Getenv("DATABASE_LOGIN")
	DATABASE_PASSWORD = "localhost" //os.Getenv("DATABASE_PASSWORD")
	TRANSKODER_ADDRESS = os.Getenv("TRANSKODER_ADDRESS")
	if TRANSKODER_ADDRESS == "" {
		TRANSKODER_ADDRESS = "http://127.0.0.1:11001"
	}
	blob_key := os.Getenv("BLOB_KEY")
	if blob_key == "" {
		blob_key = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
	}
	blob_acc_name := os.Getenv("BLOB_NAME")
	if blob_acc_name == "" {
		blob_acc_name = "devstoreaccount1"
	}
	blob_url := os.Getenv("BLOB_URL")
	if blob_url == "" {
		blob_url = "http://127.0.0.1:10000/devstoreaccount1"
	}

	fmt.Println("Transkoder_Address: ", TRANSKODER_ADDRESS)
	fmt.Println("blob address: ", blob_url)
	Credentials = credentials{Name: blob_acc_name, Key: blob_key}
	initDatabaseAbstraction()

	cred, err = azblob.NewSharedKeyCredential(Credentials.Name, Credentials.Key)
	check(err)
	serviceClient, err = azblob.NewServiceClientWithSharedKey(blob_url, cred, nil)
	check(err)
	ctx = context.Background() // context
	handleRequests()

}
