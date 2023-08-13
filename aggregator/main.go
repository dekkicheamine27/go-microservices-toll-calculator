package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/go/truck-toll-calculator/types"
)

func main() {
	listenAdrr := flag.String("listenaddr", ":5000", "the listen of the HTTP server")

	var (
		store = NewMemoryStore()
		svc = NewInvoiceAggregator(store)
		
	)
	svc = NewLogMiddleware(svc)

	makeHTTPTransport(*listenAdrr, svc)

}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("server is running on ", listenAddr)
	http.HandleFunc("/aggregate", hadleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)

}

func hadleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest , map[string]string{"error": err.Error()})
			return
		}

		if err:=  svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		    return
		}

		 
	}
}

func writeJSON( rw http.ResponseWriter, status int, v any) error{
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
