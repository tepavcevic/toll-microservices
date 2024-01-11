package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tepavcevic/toll-microservices/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "gateway HTTP server port")
	aggregatorServiceAddr := flag.String(
		"aggregatorServiceAddr",
		"http://localhost:3500",
		"aggregator service HTTP address",
	)

	aggClient := client.NewHTTPClient(*aggregatorServiceAddr)
	invHandler := newInvoiceHandler(aggClient)

	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))

	logrus.Infof("Gateway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obuID := r.URL.Query().Get("obu")
	if obuID == "" {
		return errors.New("missing OBUID")
	}
	obuIDInt, err := strconv.Atoi(obuID)
	if err != nil {
		return errors.New("OBUID has to be an integer")
	}

	inv, err := h.client.GetInvoice(context.Background(), obuIDInt)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(data)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(started time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(started),
				"uri":  r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
