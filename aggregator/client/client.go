package client

import (
	"context"

	"github.com/tepavcevic/toll-microservices/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
