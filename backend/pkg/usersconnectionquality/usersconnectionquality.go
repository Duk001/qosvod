package usersconnectionquality

import (
	"fmt"
	"math"

	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Duk001/qosvod/backend/pkg/videoquality"
	_ "github.com/denisenkom/go-mssqldb"
)

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
type UsersConnectionQuality struct {
	uscMap map[string]*usc
}

func (uscq *UsersConnectionQuality) Initiate(token, filmID string, fvq videoquality.FilmVideoQuality) {
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
func (uscq *UsersConnectionQuality) Delete(token string) {
	delete(uscq.uscMap, token)
}
func (uscq *UsersConnectionQuality) GetBandwidth(token string) (int, bool) {
	u, ok := uscq.uscMap[token]
	return u.CurrentBandwidth, ok
}
func (uscq *UsersConnectionQuality) GetBuffer(token string) (int, bool) {
	u, ok := uscq.uscMap[token]
	return int(u.BufferLength), ok
}

func (uscq *UsersConnectionQuality) ConnectionQualityString(token string) string {
	out := ""
	u, ok := uscq.uscMap[token]
	if !ok {
		return out
	}
	out = fmt.Sprintf("Current Bandwidth: %d\n Segments Hit: %d\n Average Bandwitdh: %d\n Standard Deviation: %d", u.CurrentBandwidth, u.SegmentsHit, u.AverageBandwidth, u.StandardDeviation)
	return out
}
func (uscq *UsersConnectionQuality) ChangeBandwidth(token string, bandwidth string, buffer string) bool {
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
func (uscq *UsersConnectionQuality) SetVideoQuality(token string, qualityList []videoquality.VideoQuality, lowBufferState, highBufferState int) (int, string) {
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
		nextBitrate := qualityList[nextBitrateIndex].BitrateValue
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

	u.lastSetBitrate = vq.BitrateValue //bitrateValue * 1000
	return u.lastSetBitrate, containerName
}
