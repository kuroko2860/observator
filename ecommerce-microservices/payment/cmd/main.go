package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"kltn/ecommerce-microservices/payment/pkg/handler"
	"kltn/ecommerce-microservices/payment/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	// Parse command line flags
	var (
		httpAddr = flag.String("http.addr", ":8082", "HTTP listen address")
		natsURL  = flag.String("nats.url", "nats://nats:4222", "NATS server URL")
	)
	flag.Parse()

	// Initialize NATS connection first (before configuring zerolog)
	err := logging.InitNATS(*natsURL)
	if err != nil {
		// Use standard log here since zerolog isn't configured yet
		log.Error().Err(err).Msg("Failed to connect to NATS")
	} else {
		defer logging.CloseNATS()

		// Set up custom zerolog writer that sends logs to NATS
		natsWriter := &logging.NATSLogWriter{ServiceName: "address-service"}
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		multi := zerolog.MultiLevelWriter(consoleWriter, natsWriter)

		// Configure zerolog with the multi-writer
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()
	}

	// Set log level
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Initialize the tracer
	shutdown, err := tracing.InitTracer("payment-service", *natsURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize tracer")
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// Create the service
	svc := service.NewPaymentService()

	// Create and register handlers
	paymentHandler := handler.NewPaymentHandler(svc)

	// Create Echo instance
	e := echo.New()

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware("payment-service"))
	e.Use(logging.CreateLoggingMiddleware("payment-service"))

	// Register routes
	paymentHandler.RegisterRoutes(e)

	// Start server in a goroutine
	go func() {
		log.Info().Str("transport", "HTTP").Str("addr", *httpAddr).Msg("Starting server")
		if err := e.Start(*httpAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
