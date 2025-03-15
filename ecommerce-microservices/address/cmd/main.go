package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"kltn/ecommerce-microservices/address/pkg/endpoint"
	"kltn/ecommerce-microservices/address/pkg/service"
	"kltn/ecommerce-microservices/address/pkg/transport"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	var (
		httpAddr  = flag.String("http.addr", ":8084", "HTTP listen address")
		zipkinURL = flag.String("zipkin.url", "http://localhost:9411/api/v2/spans", "Zipkin server URL")
	)
	flag.Parse()

	// Create a logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Initialize the tracer
	{
		shutdown, err := tracing.InitTracer("address-service", *zipkinURL)
		if err != nil {
			level.Error(logger).Log("msg", "Failed to initialize tracer", "err", err)
			os.Exit(1)
		}
		defer shutdown(context.Background())
	}

	// Create the service
	var svc service.AddressService
	{
		svc = service.NewBasicAddressService()
	}

	// Create the endpoints
	var endpoints endpoint.Endpoints
	{
		endpoints = endpoint.MakeEndpoints(svc, logger)
	}

	// Create the HTTP handler
	var h http.Handler
	{
		h = transport.NewHTTPHandler(endpoints, logger)
		// Add the logging middleware
		h = logging.HTTPMiddleware(logger)(h)
	}

	// Create the HTTP server
	var server *http.Server
	{
		server = &http.Server{
			Addr:    *httpAddr,
			Handler: h,
		}
	}

	// Start the server
	errs := make(chan error)
	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- server.ListenAndServe()
	}()

	// Handle shutdown signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Wait for an error
	level.Info(logger).Log("exit", <-errs)

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}