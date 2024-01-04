package main

import (
	"log"
)

type DistanceCalculator struct {
	consumer DataConsumer
}

var kafkaTopic = "obudata"

func main() {
	distCalc := NewLogMiddleware(NewCalculatorService())

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, distCalc)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
