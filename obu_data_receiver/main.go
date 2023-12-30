package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/tepavcevic/toll-microservices/types"
)

const kafkaTopic = "obudata"

func main() {
	dataReceiver, err := newDataReceiver()
	if err != nil {
		panic(err)
	}
	defer dataReceiver.prod.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range dataReceiver.prod.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	http.HandleFunc("/ws", dataReceiver.wsHandler)
	log.Fatal(http.ListenAndServe(":30000", nil))
	fmt.Println("data receiver working")
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func newDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	topic := kafkaTopic

	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
	dr.prod.Flush(500)

	return err
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
		// dr.msgch <- payload
		dr.produceData(payload)
	}

}
