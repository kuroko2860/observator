package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	pb "kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/inventory/pkg/service"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	var (
		grpcAddr  = flag.String("grpc.addr", ":8083", "gRPC listen address")
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
		shutdown, err := tracing.InitTracer("inventory-service", *zipkinURL)
		if err != nil {
			level.Error(logger).Log("msg", "Failed to initialize tracer", "err", err)
			os.Exit(1)
		}
		defer shutdown(context.Background())
	}

	// Create the service
	var svc service.InventoryService
	{
		svc = service.NewBasicInventoryService()
	}

	// Create the gRPC server
	var server *grpc.Server
	{
		server = grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
		pb.RegisterInventoryServiceServer(server, service.NewGRPCServer(svc))
		reflection.Register(server)
	}

	// Start the server
	errs := make(chan error)
	go func() {
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		errs <- server.Serve(lis)
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
	server.GracefulStop()
}