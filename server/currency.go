package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"grpc_currency/protos/currency"
)

type Currency struct {
	log hclog.Logger
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handler GetRate", "base", rr.GetBase(), "destination", rr.Destination)
	return &currency.RateResponse{Rate: 5.37}, nil
}
