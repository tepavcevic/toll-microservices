package main

import (
	"math"
	"math/rand"
	"time"
)

const sendInterval = 60

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Location struct {
	Lat  float64
	Long float64
}

func genLocation() Location {
	return Location{
		Lat:  genCoord("lat"),
		Long: genCoord("long"),
	}
}

func genCoord(variant string) float64 {
	switch variant {
	case "lat":
		return rand.Float64()*180 - 90
	case "long":
		return rand.Float64() * 180
	default:
		return rand.Float64()*180 - 90
	}
}

func genOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func main() {
	obuIDS := genOBUIDS(20)
	obuData := make([]OBUData, 20)
	for i, id := range obuIDS {
		time.Sleep(time.Second * sendInterval)
		location := genLocation()
		data := OBUData{
			OBUID: id,
			Lat:   location.Lat,
			Long:  location.Long,
		}
		obuData[i] = data
	}
}
