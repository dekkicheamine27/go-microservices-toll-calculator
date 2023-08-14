package main

import (
	"time"

	"github.com/go/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"error": err,
		}).Info("Aggregate distance")
	}(time.Now())
	return l.next.AggregateDistance(distance)
}

func (l *LogMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var(
			totalDistance float64
			totalAmount float64

		)

		if inv != nil {
			totalDistance = inv.TotalDistance
            totalAmount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"error":       err,
			"obuId": id,
			"totalDistance": totalDistance,
            "totalAmount": totalAmount,
		}).Info("CalculateInvoice")
	}(time.Now())
	return l.next.CalculateInvoice(id)
}
