package client

import (
	"github.com/go/truck-toll-calculator/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	EndPoint string
	Client   types.AggregatorClient
}

func NewGRPCClient(endPoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		EndPoint: endPoint,
		Client:   c,
	}, nil
}
