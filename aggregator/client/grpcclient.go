package client

import (
	"context"

	"github.com/tepavcevic/toll-microservices/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(e string) (*GRPCClient, error) {
	conn, err := grpc.Dial(e, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Endpoint: e,
		client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, distance)

	return err
}

func (c *GRPCClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: 1234.33213,
		TotalAmount:   534658734.432423,
	}, nil
}
