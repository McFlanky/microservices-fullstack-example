package server

import (
	"context"

	protos "github.com/McFlanky/microservices-fullstack-example/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server that implements the methods defined by the CurrencyServer interface
type Currency struct {
	log hclog.Logger
	protos.UnimplementedCurrencyServer
}

// NewCurrency creates a new Currency server
func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

// GetRate implements the CurrencyServer method GetRate and returns the currency exchange rate
// for the 2 given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())
	return &protos.RateResponse{Rate: 0.5}, nil
}
