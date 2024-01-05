package main

import (
	"fmt"

	"github.com/tepavcevic/toll-microservices/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	GetInvoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in storage", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) GetInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)
	if err != nil {
		return &types.Invoice{}, err
	}
	invoice := types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   dist * basePrice,
	}

	return &invoice, nil
}
