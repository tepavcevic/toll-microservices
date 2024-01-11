package main

import (
	"log"

	"github.com/tepavcevic/toll-microservices/aggregator/client"
)

const (
	aggregatorEndpoint = "http://localhost:3500"
	kafkaTopic         = "obudata"
)

func main() {
	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	distCalc := NewLogMiddleware(NewCalculatorService())

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, distCalc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
