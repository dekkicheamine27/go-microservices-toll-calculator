package main

import (

	"log"
	"math"
	"math/rand"
	"time"

	"github.com/go/truck-toll-calculator/types"
	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:3000/ws"




func gencoord() float64 {
	i := float64(rand.Intn(100))
	f := rand.Float64()
	return i + f
}

func lanAndLong() (float64, float64){
	return gencoord(), gencoord()
}

func main() {
    obuIds := genObuData(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for  {
		lat, lng := lanAndLong()
		for i := 0; i < len(obuIds); i++ {
			data:= types.OBUData{
				OBUID: obuIds[i],
				Lat: lat,
				Lng: lng,
  
			}	
			if err:=conn.WriteJSON(data); err !=nil{
				log.Fatal(err)
			}
			
		}

		

		time.Sleep(time.Second)
		
	}



}

func genObuData(n int) []int {
	obuIds := make([]int, n)
	for i := 0; i < n; i++ {
		obuIds[i] = rand.Intn(math.MaxInt)
	}
	return obuIds
}

func init(){
	rand.Seed(time.Now().UnixNano() )
}