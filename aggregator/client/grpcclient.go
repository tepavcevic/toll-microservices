package client

import (
	"github.com/tepavcevic/toll-microservices/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(e string) (*GRPCClient, error) {
	conn, err := grpc.Dial(e, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Endpoint:         e,
		AggregatorClient: c,
	}, nil
}
