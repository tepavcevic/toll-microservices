package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "gateway HTTP server port")

	http.HandleFunc("/invoice", makeAPIFunc(handleGetInvoice))

	logrus.Infof("Gateway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	// writeJSON(w, http.StatusOK, map[string]string{"invoice": "some invoice"})
	return errors.New("random error do boga")
}

func writeJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(data)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
