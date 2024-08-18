package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/McFlanky/microservices-fullstack-example/api/data"
	"github.com/McFlanky/microservices-fullstack-example/api/handlers"
	protos "github.com/McFlanky/microservices-fullstack-example/currency/protos/currency"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l := hclog.Default()
	v := data.NewValidation()

	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create client
	cc := protos.NewCurrencyClient(conn)

	// create database instance
	db := data.NewProductsDB(cc, l)

	// create the handlers
	ph := handlers.NewProducts(l, v, db)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// <---------------- API HANDLERS ----------------->
	getRtr := sm.Methods(http.MethodGet).Subrouter()
	getRtr.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRtr.HandleFunc("/products", ph.ListAll)

	getRtr.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")

	putRtr := sm.Methods(http.MethodPut).Subrouter()
	putRtr.HandleFunc("/products", ph.UpdateProducts)
	putRtr.Use(ph.MiddlewareValidateProduct)

	postRtr := sm.Methods(http.MethodPost).Subrouter()
	postRtr.HandleFunc("/products", ph.AddProduct)
	postRtr.Use(ph.MiddlewareValidateProduct)

	deleteRtr := sm.Methods(http.MethodDelete).Subrouter()
	deleteRtr.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	// <--------------- API DOCS HANDLERS --------------->
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRtr.Handle("/docs", sh)
	getRtr.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS handler
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := &http.Server{
		Addr:         ":8080",                                          // bind address
		Handler:      ch(sm),                                           // default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // logger for server
		ReadTimeout:  1 * time.Second,                                  // max time to read request from client
		WriteTimeout: 1 * time.Second,                                  // max time to write response to client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Received terminate, graceful shutdown", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
