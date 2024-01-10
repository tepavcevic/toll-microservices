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

	"github.com/tepavcevic/toll-microservices/aggregator/client"
	"github.com/tepavcevic/toll-microservices/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3500", "the listen address for HTTP aggregator service")
	grpcListenAddr := flag.String("grpcAddr", ":4500", "the listen address for gRPC aggregator service")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewLogMiddleware(NewInvoiceAggregator(store))

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc))
	}()

	time.Sleep(time.Second * 2)
	c, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	req := types.AggregateRequest{
		ObuID: 23,
		Value: 34.655,
		Unix:  time.Now().Unix(),
	}
	_, err = c.Aggregate(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(makeHTTPTransport(*httpListenAddr, svc))
}

// Makes a TCP listener,
// makes a gRPC server with options and
// registers our gRPC server implementation to the grpc package.
func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("gRPC aggregator service running on port", listenAddr)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer lis.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))

	return server.Serve(lis)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP aggregator service running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))

	return http.ListenAndServe(listenAddr, nil)
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

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID := r.URL.Query().Get("obu")
		if obuID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obuID"})
			return
		}
		obuIDInt, err := strconv.Atoi(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obuID has to be integer"})
			return
		}

		invoice, err := svc.GetInvoice(obuIDInt)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
