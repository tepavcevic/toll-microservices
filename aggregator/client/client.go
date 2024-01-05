package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tepavcevic/toll-microservices/types"
)

type Client struct {
	Endpoint string
}

func NewClient(e string) *Client {
	return &Client{
		Endpoint: e,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
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
