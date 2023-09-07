package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go/truck-toll-calculator/aggregator/client"
	"github.com/go/truck-toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAdrr := flag.String("httpAddr", ":5000", "the listen of the HTTP server")
	grpcListenAdrr := flag.String("grpcAdrr", ":5001", "the listen of the GRPC server")
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go func()  {
		log.Fatal(makeGRPCTransport(*grpcListenAdrr, svc))
	}()
	time.Sleep(time.Second * 2)

	c, err := client.NewGRPCClient(*grpcListenAdrr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := c.Client.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 1,
		Value: 58.5,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}

	log.Fatal(makeHTTPTransport(*httpListenAdrr, svc))

}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC server is running on ", listenAddr)
	//make a TCP listeners
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer lis.Close()

	// make a new GRPC natibe server with(options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	//regiter our GRPC server to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(lis)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP server is running on ", listenAddr)
	http.HandleFunc("/aggregate", hadleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddr, nil)

}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu ID"})
			return
		}
		obuID, err := strconv.Atoi(values[0])

		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(obuID)

		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)

		w.Write([]byte("need to return the invoice"))
	}
}

func hadleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
