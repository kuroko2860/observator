package main

import (
	"context"
	"flag"
	"net"
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
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"kltn/ecommerce-microservices/inventory/pkg/handler"
	"kltn/ecommerce-microservices/inventory/pkg/service"
	pb "kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	// Parse command line flags
	var (
		httpAddr = flag.String("http.addr", ":8083", "HTTP listen address")
		grpcAddr = flag.String("grpc.addr", ":50051", "gRPC listen address")
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
		natsWriter := &logging.NATSLogWriter{ServiceName: "inventory-service"}
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
	shutdown, err := tracing.InitTracer("inventory-service", *natsURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize tracer")
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// Create the service
	svc := service.NewInventoryService()

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	pb.RegisterInventoryServiceServer(grpcServer, service.NewGRPCServer(svc))

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			log.Fatal().Err(err).Str("addr", *grpcAddr).Msg("Failed to listen for gRPC")
		}
		log.Info().Str("transport", "gRPC").Str("addr", *grpcAddr).Msg("Starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Create and register HTTP handlers
	inventoryHandler := handler.NewInventoryHandler(svc)

	// Create Echo instance
	e := echo.New()

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware("inventory-service"))
	e.Use(logging.CreateLoggingMiddleware("inventory-service"))

	// Register routes
	inventoryHandler.RegisterRoutes(e)

	// Start HTTP server in a goroutine
	go func() {
		log.Info().Str("transport", "HTTP").Str("addr", *httpAddr).Msg("Starting HTTP server")
		if err := e.Start(*httpAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server startup failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down servers...")

	// Gracefully shutdown the HTTP server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("HTTP server forced to shutdown")
	}

	// Stop gRPC server
	grpcServer.GracefulStop()

	log.Info().Msg("Servers exited")
}
