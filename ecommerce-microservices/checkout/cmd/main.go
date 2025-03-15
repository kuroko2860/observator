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
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"kltn/ecommerce-microservices/checkout/pkg/endpoint"
	"kltn/ecommerce-microservices/checkout/pkg/service"
	"kltn/ecommerce-microservices/checkout/pkg/transport"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	var (
		httpAddr        = flag.String("http.addr", ":8080", "HTTP listen address")
		orderServiceURL = flag.String("order.url", "http://localhost:8081", "Order service URL")
		zipkinURL       = flag.String("zipkin.url", "http://localhost:9411/api/v2/spans", "Zipkin server URL")
		natsURL = flag.String("nats.url", "nats://localhost:4222", "NATS server URL")
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
		shutdown, err := tracing.InitTracer("checkout-service", *zipkinURL)
		if err != nil {
			level.Error(logger).Log("msg", "Failed to initialize tracer", "err", err)
			os.Exit(1)
		}
		defer shutdown(context.Background())
	}

	// Create an instrumented HTTP client for service-to-service communication
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   5 * time.Second,
	}

	// Create the service
	var svc service.CheckoutService
	{
		svc = service.NewBasicCheckoutService(*orderServiceURL, httpClient)
	}

	// Create the endpoints
	var endpoints endpoint.Endpoints
	{
		endpoints = endpoint.MakeEndpoints(svc, logger)
	}

	// Initialize NATS connection
	err := logging.InitNATS(*natsURL)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to connect to NATS", "err", err)
	} else {
		defer logging.CloseNATS()
	}

	// Create the HTTP handler
	var h http.Handler
	{
		h = transport.NewHTTPHandler(endpoints, logger)
		// Add the logging middleware with service name
		h = logging.HTTPMiddleware(logger, "checkout-service")(h)
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