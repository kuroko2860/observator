package service

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	pb "kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// InventoryService describes the service
type InventoryService interface {
	UpdateInventory(ctx context.Context, orderID string, items []string) error
	VerifyInventory(ctx context.Context, items []string) (bool, error)
}

// GRPCServer is the gRPC server implementation
type GRPCServer struct {
	pb.UnimplementedInventoryServiceServer
	svc InventoryService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(svc InventoryService) *GRPCServer {
	return &GRPCServer{svc: svc}
}

// UpdateInventory implements the gRPC UpdateInventory method
func (s *GRPCServer) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.UpdateInventoryResponse, error) {
	err := s.svc.UpdateInventory(ctx, req.OrderId, req.Items)
	if err != nil {
		return &pb.UpdateInventoryResponse{Error: err.Error()}, nil
	}
	return &pb.UpdateInventoryResponse{}, nil
}

// VerifyInventory implements the gRPC VerifyInventory method
func (s *GRPCServer) VerifyInventory(ctx context.Context, req *pb.VerifyInventoryRequest) (*pb.VerifyInventoryResponse, error) {
	available, err := s.svc.VerifyInventory(ctx, req.Items)
	if err != nil {
		return &pb.VerifyInventoryResponse{Error: err.Error()}, nil
	}
	return &pb.VerifyInventoryResponse{Available: available}, nil
}

type basicInventoryService struct {
	// In a real implementation, this would have a database connection
	inventory map[string]int
}

// NewBasicInventoryService returns a naive, stateless implementation of InventoryService
func NewBasicInventoryService() InventoryService {
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

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.Int("items.count", len(items)),
	)

	if len(items) == 0 {
		span.RecordError(errors.New("no items to update"))
		return errors.New("no items to update")
	}

	// Check if all items are available
	for _, item := range items {
		if s.inventory[item] <= 0 {
			err := errors.New("item out of stock: " + item)
			span.RecordError(err)
			return err
		}
	}

	// Update inventory
	for _, item := range items {
		s.inventory[item]--
	}

	return nil
}

// VerifyInventory implements InventoryService
func (s *basicInventoryService) VerifyInventory(ctx context.Context, items []string) (bool, error) {
	// Create a span for the verify inventory operation
	tracer := tracing.Tracer("inventory-service")
	ctx, span := tracer.Start(ctx, "VerifyInventory")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("items.count", len(items)),
	)

	if len(items) == 0 {
		span.RecordError(errors.New("no items to verify"))
		return false, errors.New("no items to verify")
	}

	// Check if all items are available
	for _, item := range items {
		if s.inventory[item] <= 0 {
			span.SetAttributes(attribute.Bool("inventory.available", false))
			return false, nil
		}
	}

	span.SetAttributes(attribute.Bool("inventory.available", true))
	return true, nil
}