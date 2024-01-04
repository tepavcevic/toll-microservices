package main

import "github.com/tepavcevic/toll-microservices/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	return nil
}
