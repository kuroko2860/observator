package service

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	"kltn/ecommerce-microservices/pkg/tracing"
)

// Address represents a shipping address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// AddressService describes the service
type AddressService interface {
	GetAddress(ctx context.Context, userID string) (Address, error)
}

type basicAddressService struct {
	// In a real implementation, this would have a database connection
	addresses map[string]Address
}

// NewBasicAddressService returns a naive, stateless implementation of AddressService
func NewBasicAddressService() AddressService {
	// Initialize with some dummy addresses
	addresses := make(map[string]Address)
	addresses["user1"] = Address{
		Street:     "123 Main St",
		City:       "Anytown",
		State:      "CA",
		PostalCode: "12345",
		Country:    "USA",
	}
	addresses["user2"] = Address{
		Street:     "456 Oak Ave",
		City:       "Somewhere",
		State:      "NY",
		PostalCode: "67890",
		Country:    "USA",
	}
	
	return &basicAddressService{
		addresses: addresses,
	}
}

// GetAddress implements AddressService
func (s *basicAddressService) GetAddress(ctx context.Context, userID string) (Address, error) {
	// Create a span for the get address operation
	tracer := tracing.Tracer("address-service")
	ctx, span := tracer.Start(ctx, "GetAddress")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", userID),
	)

	if userID == "" {
		span.RecordError(errors.New("user ID is required"))
		return Address{}, errors.New("user ID is required")
	}

	address, exists := s.addresses[userID]
	if !exists {
		err := errors.New("address not found for user: " + userID)
		span.RecordError(err)
		return Address{}, err
	}

	return address, nil
}