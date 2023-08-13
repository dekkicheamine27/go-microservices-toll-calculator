package main

import (
	"fmt"

	"github.com/go/truck-toll-calculator/types"
)


type MemoryStore struct{
	data map[int]float64
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(distance types.Distance) error {
	fmt.Println("the distance is insert in storage", distance)
    s.data[distance.OBUID] += distance.Value
	return nil
}