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
	svc := NewLogMiddleware(NewInvoiceAggregator(store))

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
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(dist); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, dist)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
