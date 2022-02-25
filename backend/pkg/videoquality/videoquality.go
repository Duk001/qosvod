package videoquality

import (
	"database/sql"
	"errors"

	"log"

	"strconv"
	"strings"

	"encoding/json"
	// "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	// _ "github.com/denisenkom/go-mssqldb"
	// "github.com/Duk001/qosvod/backend/pkg/server"
)


func GetFilmQuality(id string,db *sql.DB) ([]VideoQuality, error) { //TODO Add film map with video quality settings
	var data []VideoQuality
	var tmpData string

	row, err := db.Query("select quality from films where id = ?", id)
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
			data = append(data, VideoQuality{splitElem[1], splitElem[0], bitrateValue * 1000})
		}
		return data, nil
	}
	return data, errors.New("Film with id: " + id + " not found in database" + "		->		" + err.Error())
}

type VideoQuality struct {
	Bitrate      string `json:"Bitrate"`
	Resolution   string `json:"Resolution"`
	BitrateValue int
}

type FilmVideoQuality struct {
	vq map[string]*[]VideoQuality
}




func (fvq *FilmVideoQuality) Initiate(filmID string,db *sql.DB) {
	if fvq.vq == nil {
		fvq.vq = make(map[string]*[]VideoQuality)
	} else if _, ok := fvq.vq[filmID]; ok {
		return
	}

	sort := func(data []VideoQuality) []VideoQuality { //? Insertion Sort	descending
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
	vq, _ := GetFilmQuality(filmID,db)
	vq = sort(vq)
	fvq.vq[filmID] = &vq
}
func (fvq *FilmVideoQuality) Get(filmID string) []VideoQuality {
	data := fvq.vq[filmID]
	return *data
}

type userVideoQualityLogger struct {
	BandwidthArray              []int `json:"BandwidthArray"`
	AdjustedBandwidthArray      []int `json:"AdjustedBandwidthArray"`
	LocalStandardDeviationArray []int `json:"LocalStandardDeviationArray"`
	SetBitrateArray             []int `json:"SetBitrateArray"`
}
type VideoQualityLogger struct {
	usqMap map[string]*userVideoQualityLogger
}

func (vql *VideoQualityLogger) Initiate(token string) {
	if vql.usqMap == nil {
		vql.usqMap = make(map[string]*userVideoQualityLogger)
	}
	//var tmpArray []int
	u := userVideoQualityLogger{}
	vql.usqMap[token] = &u
}
func (vql *VideoQualityLogger) Update(token string, bandwidth, adjustedBandwidth, standardDeviation int) {
	u, ok := vql.usqMap[token]
	if !ok {
		return
	}
	u.BandwidthArray = append(u.BandwidthArray, bandwidth)
	u.AdjustedBandwidthArray = append(u.AdjustedBandwidthArray, adjustedBandwidth)
	u.LocalStandardDeviationArray = append(u.LocalStandardDeviationArray, standardDeviation)
	vql.usqMap[token] = u
}
func (vql *VideoQualityLogger) UpdateBitRate(token string, bitRate int) {
	u, ok := vql.usqMap[token]
	if !ok {
		return
	}
	u.SetBitrateArray = append(u.SetBitrateArray, bitRate)
	vql.usqMap[token] = u
}
func (vql *VideoQualityLogger) Jsonify() []byte {
	jsonData, _ := json.Marshal(vql.usqMap)
	return jsonData
}
