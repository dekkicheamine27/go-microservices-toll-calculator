package main

import (
	"fmt"
	"log"
)

const topic = "obudata"

func main() {

	svc := NewCalculatorService()

	kafkaConsumer, err := NewKafKaConsumer(topic, svc)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.start()

	fmt.Println("distance call")
}
