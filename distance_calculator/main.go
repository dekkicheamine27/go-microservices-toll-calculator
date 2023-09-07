package main

import (
	"fmt"
	"log"

	"github.com/go/truck-toll-calculator/aggregator/client"
)

const (
	topic       = "obudata"
	aggEndPoint = "http://localhost:5000/aggregate"
)

func main() {

	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafKaConsumer(topic, svc, client.NewHTTPClient(aggEndPoint))
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.start()

	fmt.Println("distance call")
}
