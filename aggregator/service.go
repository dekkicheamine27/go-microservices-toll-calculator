package main

import "github.com/go/truck-toll-calculator/types"

type Aggregator interface{
	AggregateDistance(types.Distance) error
}

type Storer interface{
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
    return &InvoiceAggregator{store: store}
}


 

func (i *InvoiceAggregator) AggregateDistance(dis types.Distance) error{
    return i.store.Insert(dis)
}

