package main

import (
	"net"
	"os"

	protos "github.com/McFlanky/microservices-fullstack-example/currency/protos/currency"
	"github.com/McFlanky/microservices-fullstack-example/currency/protos/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	cs := server.NewCurrency(log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	// register the reflection service which allows clients to determine the methods for this gRPC service
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
