package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	pb "kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// InventoryService describes the service
type InventoryService interface {
	UpdateInventory(ctx context.Context, orderID string, items []string) error
	VerifyInventory(ctx context.Context, items []string) (bool, error)
}

// GRPCServer is the gRPC server for the inventory service
type GRPCServer struct {
	pb.UnimplementedInventoryServiceServer
	svc InventoryService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(svc InventoryService) *GRPCServer {
	return &GRPCServer{svc: svc}
}

// UpdateInventory implements the gRPC method
func (s *GRPCServer) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.UpdateInventoryResponse, error) {
	// Create a span for this gRPC method
	tracer := tracing.Tracer("inventory-grpc")
	ctx, span := tracer.Start(ctx, "UpdateInventory-gRPC")
	defer span.End()
	
	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("order_id", req.OrderId).
		Int("items_count", len(req.Items)).
		Logger()

	logger.Debug().Msg("gRPC UpdateInventory request received")

	err := s.svc.UpdateInventory(ctx, req.OrderId, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
		logger.Error().Err(err).Msg("Failed to update inventory")
		return &pb.UpdateInventoryResponse{Error: err.Error()}, nil
	}

	logger.Info().Msg("Inventory updated successfully")
	return &pb.UpdateInventoryResponse{}, nil
}

// VerifyInventory implements the gRPC method
func (s *GRPCServer) VerifyInventory(ctx context.Context, req *pb.VerifyInventoryRequest) (*pb.VerifyInventoryResponse, error) {
	// Create a span for this gRPC method
	tracer := tracing.Tracer("inventory-grpc")
	ctx, span := tracer.Start(ctx, "VerifyInventory-gRPC")
	defer span.End()
	
	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Int("items_count", len(req.Items)).
		Logger()

	logger.Debug().Msg("gRPC VerifyInventory request received")

	available, err := s.svc.VerifyInventory(ctx, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
		logger.Error().Err(err).Msg("Failed to verify inventory")
		return &pb.VerifyInventoryResponse{Error: err.Error()}, nil
	}

	logger.Info().Bool("available", available).Msg("Inventory verification completed")
	return &pb.VerifyInventoryResponse{Available: available}, nil
}

type basicInventoryService struct {
	// In a real implementation, this would have a database connection
	inventory map[string]int
}

// NewInventoryService returns a new implementation of InventoryService
func NewInventoryService() InventoryService {
	// Initialize with some dummy inventory
	inventory := make(map[string]int)
	inventory["item1"] = 10
	inventory["item2"] = 5
	inventory["item3"] = 15

	return &basicInventoryService{
		inventory: inventory,
	}
}

// UpdateInventory implements InventoryService
func (s *basicInventoryService) UpdateInventory(ctx context.Context, orderID string, items []string) error {
	// Create a span for the update inventory operation
	tracer := tracing.Tracer("inventory-service")
	ctx, span := tracer.Start(ctx, "UpdateInventory")
	defer span.End()

	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("order_id", orderID).
		Int("items_count", len(items)).
		Logger()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.Int("items.count", len(items)),
	)

	logger.Debug().Msg("Starting inventory update")

	if len(items) == 0 {
		err := errors.New("no items to update")
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
		logger.Error().Err(err).Msg("Update failed")
		return err
	}

	// Check if all items are available
	for _, item := range items {
		if s.inventory[item] <= 0 {
			err := errors.New("item out of stock: " + item)
			// Add error tag to span
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.String("error.message", err.Error()))
			span.RecordError(err)
			
			logger.Error().Err(err).Str("item", item).Msg("Item out of stock")
			return err
		}
	}

	// Update inventory
	for _, item := range items {
		s.inventory[item]--
		logger.Debug().Str("item", item).Int("remaining", s.inventory[item]).Msg("Item quantity updated")
	}

	logger.Info().Str("order_id", orderID).Msg("Inventory updated successfully")
	return nil
}

// VerifyInventory implements InventoryService
func (s *basicInventoryService) VerifyInventory(ctx context.Context, items []string) (bool, error) {
	// Create a span for the verify inventory operation
	tracer := tracing.Tracer("inventory-service")
	ctx, span := tracer.Start(ctx, "VerifyInventory")
	defer span.End()

	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Int("items_count", len(items)).
		Logger()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("items.count", len(items)),
	)

	logger.Debug().Msg("Starting inventory verification")

	if len(items) == 0 {
		err := errors.New("no items to verify")
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
		logger.Error().Err(err).Msg("Verification failed")
		return false, err
	}

	// Check if all items are available
	for _, item := range items {
		if s.inventory[item] <= 0 {
			logger.Info().Str("item", item).Msg("Item not available")
			span.SetAttributes(attribute.Bool("inventory.available", false))
			return false, nil
		}
	}

	logger.Info().Msg("All items available")
	span.SetAttributes(attribute.Bool("inventory.available", true))
	return true, nil
}
