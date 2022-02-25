package main

//! https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/storage/azblob#readme
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var BACKEND_SERVER_ADDRESS string

type credentials struct {
	Name string
	Key  string
}

type videoQuality struct {
	Bitrate    string `json:"Bitrate"`
	Resolution string `json:"Resolution"`
}

var Credentials credentials

var serviceClient azblob.ServiceClient
var ctx context.Context
var cred *azblob.SharedKeyCredential

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

func Encoder(file []byte, filmName string) {

	fmt.Println("Encoder")

	var ffmpegPath = "util\\ffmpeg.exe"
	//var ffmpegPath = "/usr/local/bin/ffmpeg"
	//var ffmpegPath = "mwader/static-ffmpeg:4.4.1"	// docker
	var buf bytes.Buffer

	qualityList, err := getFilmQuality(filmName)
	if err != nil {
		log.Println(err)
		return

	}
	for _, vq := range qualityList {

		containerName := strings.ToLower(filmName + "-" + strings.Replace(vq.Resolution, ":", "-", -1) + "-" + vq.Bitrate)
		// containerName = "ttt"
		fmt.Println("Container Name: ", containerName)
		container := serviceClient.NewContainerClient(containerName)
		_, err := container.Create(ctx, nil)
		check(err)

	}
	command := []string{ffmpegPath,
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-muxdelay", "0", //! Possible problem -> choppy audio
		"-f", "hls",
		"-c:v", "libx264",
		"-force_key_frames", "expr:gte(t,n_forced*1)",
		"-hls_time", "4", // old = 1
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", "http://127.0.0.1:11001" + "/transcode" + "?filmName=" + filmName + "&name=" + "%02d.ts",

		"-method", "POST",
		"http://127.0.0.1:11001" + "/manifest" + "?name=" + filmName + ".m3u8",
		//"stream_%v/stream.m3u8",

	}

	cmd := exec.Command(command[0], command[1:]...)

	//resultBuffer := bytes.NewBuffer(make([]byte, 5*1024*1024)) // pre allocate 5MiB buffer

	cmd.Stderr = os.Stderr // bind log stream to stderr
	//cmd.Stdout = resultBuffer // stdout result will be written here
	stdin, err := cmd.StdinPipe() // Open stdin pipe
	check(err)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start() // Start a process on another goroutine
	check(err)

	_, err = stdin.Write(file) // pump audio data to stdin pipe
	check(err)

	err = stdin.Close() // close the stdin, or ffmpeg will wait forever
	check(err)

	n, err := io.Copy(&buf, stdout)
	check(err)

	fmt.Printf("Copied %d bytes\n", n)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)

}
func Transcode(file []byte, name string, segment string, bitrate string, resolution string, container azblob.ContainerClient) {

	//var ffmpegPath = "/usr/local/bin/ffmpeg"
	var ffmpegPath = "util\\ffmpeg.exe"
	//var fontPath = "util\\Roboto-Regular.ttf"
	var buf bytes.Buffer

	//file_name := name+"/"+bitrate+"/"+segment
	//file_name := filepath.Join(".",name,bitrate,segment)
	//fmt.Println("Saving segment to: ",file_name)
	//watermarkText := resolution + " - " + bitrate + " - " + segment
	//watermarkParam := fmt.Sprintf("drawtext=text='%s':x=40:y=H-th-40:fontfile=%s:fontsize=36:fontcolor=white:shadowcolor=black:shadowx=5:shadowy=5", watermarkText, fontPath)
	log.Println("Transcoding seg: ", segment, "   bitrate: ", bitrate)
	command := []string{ffmpegPath,
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-hls_time", "2",
		"-copyts",
		"-muxdelay", "0",
		"-b:v", bitrate,
		"-f", "hls",
		"pipe:1"}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = os.Stderr // bind log stream to stderr

	stdin, err := cmd.StdinPipe() // Open stdin pipe
	check(err)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start() // Start a process on another goroutine
	check(err)

	_, err = stdin.Write(file) // pump segment data to stdin pipe
	check(err)

	err = stdin.Close() // close the stdin, or ffmpeg will wait forever
	check(err)

	n, err := io.Copy(&buf, stdout)
	check(err)
	fmt.Printf("Copied %d bytes\n", n)

	//resultBuffer := bytes.NewBuffer(make([]byte, 5*1024*1024))
	fmt.Println(stdout)

	blockBlob := container.NewBlockBlobClient(segment)

	data := bytes.NewReader(buf.Bytes()) //buf//bytes.NewReader(&buf)
	//data = bytes.NewReader(file) //!DEBUG
	_, err = blockBlob.Upload(ctx, NopCloser(data), nil)
	log.Println(err)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)
}

func ffmpegTest() {

	//var ffmpegPath = "util\\ffmpeg.exe"
	// var ffmpegPath = "/usr/local/bin/ffmpeg" // docker
	var ffmpegPath = "util\\ffmpeg.exe"

	command := []string{ffmpegPath,
		"-version"}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = os.Stderr // bind log stream to stderr

	err := cmd.Start() // Start a process on another goroutine
	check(err)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)
}

func manifestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
	case "POST":
		//file, _, err := r.FormFile("file")
		var out []byte
		var filmName string

		out, err := ioutil.ReadAll(r.Body)
		filmName_array, ok := r.URL.Query()["name"]
		if !ok {
			return
		}
		filmName = filmName_array[0]

		if err != nil {
			fmt.Errorf("Error during reading body: %v", err)
		}
		manifestContent := string(out)
		lines := strings.Split(manifestContent, "\n")
		//var out_lines  []string
		for i, line := range lines {
			newLine := strings.Split(line, "=")
			lines[i] = newLine[len(newLine)-1]
		}
		newManifestString := strings.Join(lines, "\n")
		manifest := []byte(newManifestString)
		container := serviceClient.NewContainerClient("manifests")

		blockBlob := container.NewBlockBlobClient(filmName)
		data := bytes.NewReader(manifest) //buf//bytes.NewReader(&buf)
		_, err = blockBlob.Upload(ctx, NopCloser(data), nil)

		check(err)

		//fmt.Println(out)
		//return

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// r.Body = http.MaxBytesReader(w, r.Body, 5000000000)
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Println("Not here...")
	case "POST":
		fmt.Print(r.ContentLength)
		var file_name string

		file_name_array, ok := r.URL.Query()["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		file_name = file_name_array[0]
		//file, header, err := r.FormFile("file")
		file, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		//defer file.Close()
		//out, err := ioutil.ReadAll(file)
		//check(err)

		Encoder(file, file_name)

		//fmt.Fprintf(w, "File uploaded successfully: ")
		//fmt.Fprintf(w, header.Filename)
	}
}

func getFilmQuality(filmName string) ([]videoQuality, error) {
	tryNumber := 0
	for tryNumber < 12 {
		resp, err := http.Get(BACKEND_SERVER_ADDRESS + "/filmQuality?film=" + filmName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		switch resp.StatusCode {

		case 425:
			if tryNumber > 12 {
				err = errors.New("getFilmQuality() -> Retry limit reached")
				return nil, err
			}
			tryNumber += 1
			time.Sleep(5 * time.Second)
		case 200:

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
				//log.Fatalf("Error fetching url %s: %v", url, err)
			}
			var data []videoQuality
			err = json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		default:
			err = errors.New(fmt.Sprintf("getFilmQuality()\n -> Wrong response code: %d expected 425 or 200", resp.StatusCode))
			return nil, err
		}
	}
	err := errors.New("getFilmQuality()\n -> Internal Error")
	return nil, err
}
func transcodeSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Println(r)
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
	case "POST":
		file, _, err := r.FormFile("file")
		var out []byte
		var segment_name string
		var filmName string
		if err != nil {

			fmt.Fprintln(w, err)
			out, err = ioutil.ReadAll(r.Body)
			segment_name_array, ok := r.URL.Query()["name"]
			if !ok != false {
				return
			}
			filmName_array, ok := r.URL.Query()["filmName"]
			if !ok != false {
				return
			}
			segment_name = segment_name_array[0]
			filmName = filmName_array[0]

			if err != nil {
				fmt.Errorf("Error during reading body: %v", err)
			}
			//fmt.Println(out)
			//return
		} else {

			defer file.Close()
			out, err = ioutil.ReadAll(file)
			segment_name = "output.ts"
			filmName = "testFilm"
			check(err)
		}

		//qualityList := baseVideoQualityList
		qualityList, err := getFilmQuality(filmName)
		if err != nil {
			log.Println(err)
			return

		}

		for _, vq := range qualityList {
			//containerName := filmName+"-"+strings.Replace(vq.resolution,":","-",-1)+"-"+vq.bitrate
			containerName := strings.ToLower(filmName + "-" + strings.Replace(vq.Resolution, ":", "-", -1) + "-" + vq.Bitrate)
			container := serviceClient.NewContainerClient(containerName)

			go Transcode(out, filmName, segment_name, vq.Bitrate, vq.Resolution, container)
		}

		fmt.Fprintf(w, "File Transcoded successfully: ")
		//fmt.Fprintf(w, header.Filename)
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadVideo)
	http.HandleFunc("/transcode", transcodeSegment)
	http.HandleFunc("/manifest", manifestHandler)
	//log.Fatal(http.ListenAndServe(":10001", nil))

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "11001"
	}
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

//?https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/storage/azblob#readme

func main() {
	var err error
	BACKEND_SERVER_ADDRESS = os.Getenv("BACKEND_SERVER_ADDRESS")
	if BACKEND_SERVER_ADDRESS == "" {
		BACKEND_SERVER_ADDRESS = "http://127.0.0.1:11000"
	}
	// blob_key := os.Getenv("BLOB_KEY")
	// blob_acc_name := os.Getenv("BLOB_NAME")
	// blob_url := os.Getenv("BLOB_URL")
	// AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;DefaultEndpointsProtocol=http;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;QueueEndpoint=http://127.0.0.1:10001/devstoreaccount1;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;
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

	fmt.Println("acc: ", blob_acc_name, "\nkey: ", blob_key, "\nurl: ", blob_url, "\nbackend server: ", BACKEND_SERVER_ADDRESS)

	Credentials = credentials{Name: blob_acc_name, Key: blob_key}
	ffmpegTest()

	cred, err = azblob.NewSharedKeyCredential(Credentials.Name, Credentials.Key)
	check(err)
	serviceClient, err = azblob.NewServiceClientWithSharedKey(blob_url, cred, nil)
	check(err)
	ctx = context.Background() // context

	fmt.Println(cred, "	\n	", serviceClient)
	handleRequests()

}
