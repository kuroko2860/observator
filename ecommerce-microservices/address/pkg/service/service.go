package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

// NewAddressService returns a new implementation of AddressService
func NewAddressService() AddressService {
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
	tracer := otel.Tracer("address-service")
	ctx, span := tracer.Start(ctx, "GetAddress-service")
	defer span.End()

	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("user_id", userID).
		Logger()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", userID),
	)

	logger.Debug().Msg("Getting address for user")

	if userID == "" {
		err := errors.New("user ID is required")
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("User ID is required")
		return Address{}, err
	}

	address, exists := s.addresses[userID]
	if !exists {
		err := errors.New("address not found for user: " + userID)
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Address not found")
		return Address{}, err
	}

	logger.Info().
		Str("street", address.Street).
		Str("city", address.City).
		Msg("Address found successfully")

	return address, nil
}
