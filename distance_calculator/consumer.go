package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/tepavcevic/toll-microservices/types"
)

type DataConsumer interface {
	ConsumeData(types.OBUData) error
}

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, calcService CalculatorServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: calcService,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	kc.isRunning = true
	kc.readMessageLoop()
}

func (kc *KafkaConsumer) ConsumeData(data types.OBUData) error {

	return nil
}

func (kc *KafkaConsumer) readMessageLoop() {
	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("consumer error: %v (%v)\n", err, msg)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("json serialization error: %v (%v)\n", err, msg)
			continue
		}

		dist, err := kc.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("error calculating distance: %v (%v)\n", err, msg)
			continue
		}
		fmt.Println("distance:", dist)
	}
}
