package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc_currency/protos/currency"
	"grpc_currency/server"
	"net"
	"os"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)
	//this should be disabled in production
	reflection.Register(gs)

	//register our server in order to "Respond" our request
	currency.RegisterCurrencyServer(gs, cs)

	l, err := net.Listen("tcp", ":9098")
	if err != nil {
		log.Error(fmt.Sprintf("Unable to listen port %s", "9098"), "Error", err)
		os.Exit(1)
	}
	gs.Serve(l)
}
