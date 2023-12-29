package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tepavcevic/toll-microservices/types"
)

func main() {
	dataReceiver := newDataReceiver()

	http.HandleFunc("/ws", dataReceiver.wsHandler)
	log.Fatal(http.ListenAndServe(":30000", nil))
	fmt.Println("data receiver working")
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func newDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}

func (dr *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		WriteBufferSize: 1028,
		ReadBufferSize:  1028,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU client connected!")
	for {
		var payload types.OBUData

		if err := dr.conn.ReadJSON(&payload); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Println("received data:", payload)
		dr.msgch <- payload
	}

}
