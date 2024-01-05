package main

import (
	"fmt"

	"github.com/tepavcevic/toll-microservices/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[int]float64{},
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	m.data[distance.OBUID] += distance.Value
	return nil
}

func (m *MemoryStore) Get(obuID int) (float64, error) {
	dist, ok := m.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("no data associated with provided obuID: %v", obuID)
	}

	return dist, nil
}
