package server

import (
	"context"
	"github.com/CassioRoos/grpc_currency/protos/currency"
	"github.com/CassioRoos/grpc_currency/services"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type Currency struct {
	rates *services.ExchangeRates
	log   hclog.Logger
	// A list of our subscriptions
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

func NewCurrency(r *services.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{
		rates:         r,
		log:           l,
		subscriptions: make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)}
	go c.HandleUpdate()
	return c
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handler GetRate", "base", rr.GetBase(), "destination", rr.Destination)
	// Some validations were made in order to test how it works in grpc
	// BASE AND DESTINATION Should not be the same
	if rr.Base == rr.Destination {
		// return a status to further more add details to enhance the error
		stat := status.Newf(
			//It works like the error on HTTP
			codes.InvalidArgument,
			"Base currency %s can not be the same as the destination %s currency",
			rr.Base.String(),
			rr.Destination.String())
		// enhancing the message
		stat, err := stat.WithDetails(rr)
		if err != nil {
			return nil, err
		}
		return nil, stat.Err()
	}
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.Destination.String())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{
		Base:        rr.Base,
		Destination: rr.Destination,
		Rate:        rate,
	}, nil
}

func (c *Currency) HandleUpdate() {
	// set a time period to monitor the "CHANGES"
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		//loop over subscribed clients
		for k, v := range c.subscriptions {
			//loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.Base.String(), rr.Destination.String())
				if err != nil {
					c.log.Error("Unable to get update rate",
						"base", rr.Base.String(),
						"destinaton", rr.Destination.String())
				}
				err = k.Send(
					// Now we need to wrap the message, because it only support one of kind
					&currency.StreamingRateResponse{
						// Here we are setting the message kind to rate response
						Message: &currency.StreamingRateResponse_RateResponse{
							// here is the actual message
							RateResponse: &currency.RateResponse{
								Base:        rr.Base,
								Destination: rr.Destination,
								Rate:        r,
							},
						},
					},
				)

				if err != nil {
					c.log.Error("Unable to send update rate",
						"base", rr.Base.String(),
						"destinaton", rr.Destination.String(),
						"error", err)
				}
			}
		}
	}
}

func (c *Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {
	// subscription is blocking
	for {
		rr, err := src.Recv()
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		// any other error means the transport between the server and client is unavailable
		if err != nil {
			c.log.Info("Unable to read from client", "error", err)
			break
		}
		c.log.Info("Handle client request", "request", rr)

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*currency.RateRequest{}
		}
		var validadtionError *status.Status
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {
				validadtionError := status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe for a currency as subscription already exists")
				validadtionError, err := validadtionError.WithDetails(rr)
				if err != nil {
					c.log.Error("Unable to add metadata to error", "error", err)
					break
				}
				break
			}
		}
		if validadtionError != nil {
			src.Send(
				// Unlike the previous case here the message is an ERROR
				&currency.StreamingRateResponse{
					Message: &currency.StreamingRateResponse_Error{
						Error: validadtionError.Proto(),
					},
				})
			continue
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}
	return nil
}
