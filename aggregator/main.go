package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tepavcevic/toll-microservices/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3500", "the listen addres for aggregator service")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("Aggregator service running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dist types.Distance

		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := svc.AggregateDistance(dist); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
