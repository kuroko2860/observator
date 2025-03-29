package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CheckoutService describes the service
type CheckoutService interface {
	UserCheckout(ctx context.Context, userID string, items []string) (string, error)
}

type checkoutService struct {
	orderServiceURL string
	httpClient      *http.Client
}

// NewCheckoutService returns a new implementation of CheckoutService
func NewCheckoutService(orderServiceURL string, client *http.Client) CheckoutService {
	return &checkoutService{
		orderServiceURL: orderServiceURL,
		httpClient:      client,
	}
}

// UserCheckout implements CheckoutService
func (s *checkoutService) UserCheckout(ctx context.Context, userID string, items []string) (string, error) {
	// Create a span for the checkout operation
	tracer := otel.Tracer("checkout-service")
	ctx, span := tracer.Start(ctx, "UserCheckout-service")
	defer span.End()

	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("user_id", userID).
		Int("items_count", len(items)).
		Logger()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.Int("items.count", len(items)),
	)

	logger.Debug().Msg("Starting checkout process")

	if len(items) == 0 {
		err := errors.New("no items in cart")
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Checkout failed")
		return "", err
	}

	// Create request body with items
	type createOrderRequest struct {
		UserID string   `json:"user_id"`
		Items  []string `json:"items"`
	}

	requestBody := createOrderRequest{
		UserID: userID,
		Items:  items,
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Failed to marshal request")
		return "", err
	}

	logger.Debug().Msg("Calling order service")

	// Call the Order Service to create an order
	req, err := http.NewRequestWithContext(ctx, "POST", s.orderServiceURL+"/orders", bytes.NewBuffer(jsonBody))
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Failed to create request")
		return "", err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-ID", userID)

	// Make the request to order service
	resp, err := s.httpClient.Do(req)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Order service request failed")
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		err = errors.New("failed to create order")
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().
			Err(err).
			Int("status_code", resp.StatusCode).
			Msg("Order service returned error")
		return "", err
	}

	// Parse the response
	var orderResponse struct {
		OrderID string `json:"order_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Failed to decode response")
		return "", err
	}

	orderID := orderResponse.OrderID

	span.SetAttributes(attribute.String("order.id", orderID))
	span.AddEvent("Order Service call completed")

	logger.Info().
		Str("order_id", orderID).
		Msg("Order created successfully")

	return orderID, nil
}
