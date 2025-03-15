package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"kltn/ecommerce-microservices/pkg/tracing"
)

// CheckoutService describes the service
type CheckoutService interface {
	UserCheckout(ctx context.Context, userID string, items []string) (string, error)
}

type basicCheckoutService struct {
	orderServiceURL string
	httpClient      *http.Client
}

// NewBasicCheckoutService returns a naive, stateless implementation of CheckoutService
func NewBasicCheckoutService(orderServiceURL string, client *http.Client) CheckoutService {
	return &basicCheckoutService{
		orderServiceURL: orderServiceURL,
		httpClient:      client,
	}
}

// UserCheckout implements CheckoutService
func (s *basicCheckoutService) UserCheckout(ctx context.Context, userID string, items []string) (string, error) {
	// Create a span for the checkout operation
	tracer := tracing.Tracer("checkout-service")
	ctx, span := tracer.Start(ctx, "UserCheckout")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.Int("items.count", len(items)),
	)

	if len(items) == 0 {
		span.RecordError(errors.New("no items in cart"))
		return "", errors.New("no items in cart")
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
		span.RecordError(err)
		return "", err
	}

	// Call the Order Service to create an order
	req, err := http.NewRequestWithContext(ctx, "POST", s.orderServiceURL+"/orders", bytes.NewBuffer(jsonBody))
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-ID", userID)

	// Make the request to order service
	resp, err := s.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		err = errors.New("failed to create order")
		span.RecordError(err)
		return "", err
	}

	span.AddEvent("Order Service call completed")

	// Simulate calling the order service
	// In a real implementation, you would use the httpClient to make a request
	orderID := "order-123" // This would come from the Order Service

	span.SetAttributes(attribute.String("order.id", orderID))
	return orderID, nil
}
