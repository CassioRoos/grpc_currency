package server

import (
	"context"
	"github.com/CassioRoos/grpc_currency/protos/healthcheck"
	"github.com/hashicorp/go-hclog"
)

type healthCheck struct {
	log hclog.Logger
}

func NewHealthCheck(l hclog.Logger) *healthCheck{
	return &healthCheck{log: l}
}

func (f *healthCheck)Check(context.Context, *healthcheck.HealthCheckParam) (*healthcheck.HealthCheckReturn, error){
	return &healthcheck.HealthCheckReturn{Message: "WORKING"}, nil
}