package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"grpc_currency/protos/currency"
	"grpc_currency/services"
)

type Currency struct {
	rates *services.ExchangeRates
	log   hclog.Logger
}

func NewCurrency(r *services.ExchangeRates, l hclog.Logger) *Currency {

	return &Currency{rates: r, log: l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handler GetRate", "base", rr.GetBase(), "destination", rr.Destination)
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.Destination.String())
	if err !=nil{
		return nil, err
	}
	return &currency.RateResponse{Rate: rate}, nil
}
