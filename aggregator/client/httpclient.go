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
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with non 200 status code: %v", res.StatusCode)
	}

	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := types.GetInvoiceRequest{
		ObuID: int32(id),
	}

	b, err := json.Marshal(&invReq)
	if err != nil {
		return &types.Invoice{}, err
	}

	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", id)
	req, err := http.NewRequest("GET", endpoint, bytes.NewReader(b))
	if err != nil {
		return &types.Invoice{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &types.Invoice{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &types.Invoice{}, fmt.Errorf("server responded with non 200 status code: %v", res.StatusCode)
	}

	var inv types.Invoice

	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		return &types.Invoice{}, err
	}

	return &inv, nil
}
