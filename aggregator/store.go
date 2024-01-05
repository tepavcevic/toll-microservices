package main

import "github.com/tepavcevic/toll-microservices/types"

type Storer interface {
	Insert(types.Distance) error
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
