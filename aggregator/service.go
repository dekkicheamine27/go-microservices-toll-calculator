package main

import "github.com/go/truck-toll-calculator/types"

const basPrice = 3.25

type Aggregator interface{
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface{
	Insert(types.Distance) error
	GetDistance(int) (float64, error)

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

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error){
    distance, err := i.store.GetDistance(id)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID: id,
		TotalDistance: distance,
		TotalAmount: basPrice * distance,

	}, nil
}

