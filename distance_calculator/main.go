package main

import (
	"log"

	"github.com/tepavcevic/toll-microservices/aggregator/client"
)

const (
	aggEndpoint = "http://localhost:3500/aggregate"
	kafkaTopic  = "obudata"
)

func main() {
	aggClient := client.NewClient(aggEndpoint)
	distCalc := NewLogMiddleware(NewCalculatorService())

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, distCalc, aggClient)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
