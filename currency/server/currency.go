package server

import (
	"context"

	"github.com/McFlanky/microservices-fullstack-example/currency/data"
	protos "github.com/McFlanky/microservices-fullstack-example/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server that implements the methods defined by the CurrencyServer interface
type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
	protos.UnimplementedCurrencyServer
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{log: l}
}

// GetRate implements the CurrencyServer method GetRate and returns the currency exchange rate
// for the 2 given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}
