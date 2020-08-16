package main

import (
	"fmt"
	"github.com/CassioRoos/grpc_currency/protos/currency"
	"github.com/CassioRoos/grpc_currency/protos/healthcheck"
	"github.com/CassioRoos/grpc_currency/server"
	"github.com/CassioRoos/grpc_currency/services"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()
	rates, err := services.NewRates(log)
	if err!=nil{
		log.Error("Could not create rates", "Error", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(rates,log)
	hs := server.NewHealthCheck(log)
	//this should be disabled in production
	reflection.Register(gs)

	//register our server in order to "Respond" our request
	currency.RegisterCurrencyServer(gs, cs)
	healthcheck.RegisterHealthCheckServer(gs, hs)

	l, err := net.Listen("tcp", ":9098")

	if err != nil {
		log.Error(fmt.Sprintf("Unable to listen port %s", "9098"), "Error", err)
		os.Exit(1)
	}
	log.Info("Listen on Port: :9098")
	gs.Serve(l)

}
