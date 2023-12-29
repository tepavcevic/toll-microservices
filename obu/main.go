package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tepavcevic/toll-microservices/types"
)

const (
	sendInterval = 60
	wsEndpoint   = "ws://127.0.0.1:30000/ws"
)

func sendOBUData(conn *websocket.Conn, data types.OBUData) error {
	return conn.WriteJSON(data)
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

	conn, _, _ := websocket.DefaultDialer.Dial(wsEndpoint, nil)

	for _, id := range obuIDS {
		time.Sleep(time.Second * sendInterval)
		location := genLocation()
		data := types.OBUData{
			OBUID: id,
			Lat:   location.Lat,
			Long:  location.Long,
		}
		if err := sendOBUData(conn, data); err != nil {
			log.Fatal(err)
		}
	}
}
