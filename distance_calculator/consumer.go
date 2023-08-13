package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go/truck-toll-calculator/aggregator/client"
	"github.com/go/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	client *client.Client
}

func NewKafKaConsumer(topic string, svc CalculatorServicer, client *client.Client) (*kafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.

	return &kafkaConsumer{
		consumer:    c,
		calcService: svc,
		client: client,
	}, nil
}

func (c *kafkaConsumer) start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *kafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume err %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serializer err %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculate distance err %s", err)
			continue
		}
		_ = distance

		req := types.Distance{
			Value: distance,
			OBUID: data.OBUID,
			Unix: float64(time.Now().UnixNano()),
			
		}

		if err:= c.client.AggregateDistance(req); err != nil {
			logrus.Errorf("aggregate error", err)
			continue
		}
	}

}
