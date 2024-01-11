package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tepavcevic/toll-microservices/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(e string) *HTTPClient {
	return &HTTPClient{
		Endpoint: e,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with non 200 status code: %v", res.StatusCode)
	}

	return nil
}
