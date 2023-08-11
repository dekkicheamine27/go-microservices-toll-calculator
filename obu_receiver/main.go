package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go/truck-toll-calculator/types"
	"github.com/gorilla/websocket"
)

func main() {

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	fmt.Println("receiver working")
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		pr  DataProducer
		err error
		kafkaTopic = "obudata"
	)
	pr, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	pr = NewLogMiddleware(pr)
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  pr,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()

}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read err:", err)
		}

		fmt.Println("receive data:", data)

		if err := dr.produceData(data); err != nil {
			fmt.Println("produce err:", err)
		}

		//fmt.Printf("obu id : [%d] <lat: %2f, lng: %2f>\n", data.OBUID, data.Lat, data.Lng)
		//dr.msgch <- data

	}

}
