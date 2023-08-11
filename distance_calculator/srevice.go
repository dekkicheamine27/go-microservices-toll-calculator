package main

import (
	"fmt"
	"math"

	"github.com/go/truck-toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoint) > 0 {
		fmt.Println(s.prevPoint)
		fmt.Println(data)
		distance = calculateDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Lng)
	}
	s.prevPoint = []float64{data.Lat, data.Lng}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}