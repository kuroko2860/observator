package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// OrderService describes the service
type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []string) (string, error)
}

type basicOrderService struct {
	paymentServiceURL   string
	inventoryServiceURL string
	addressServiceURL   string
	httpClient          *http.Client
	inventoryClient     proto.InventoryServiceClient
}

// NewBasicOrderService returns a naive, stateless implementation of OrderService
func NewBasicOrderService(paymentURL, inventoryURL, addressURL string, client *http.Client) OrderService {
	// Set up gRPC connection to inventory service with OpenTelemetry interceptors
	conn, err := grpc.Dial(
		inventoryURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to inventory service: %v", err))
	}

	inventoryClient := proto.NewInventoryServiceClient(conn)

	return &basicOrderService{
		paymentServiceURL:   paymentURL,
		inventoryServiceURL: inventoryURL,
		addressServiceURL:   addressURL,
		httpClient:          client,
		inventoryClient:     inventoryClient,
	}
}

// CreateOrder implements OrderService
func (s *basicOrderService) CreateOrder(ctx context.Context, userID string, items []string) (string, error) {
	// Create a span for the create order operation
	tracer := tracing.Tracer("order-service")
	ctx, span := tracer.Start(ctx, "CreateOrder")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.Int("items.count", len(items)),
	)

	if len(items) == 0 {
		span.RecordError(errors.New("no items in order"))
		return "", errors.New("no items in order")
	}

	// 1. Get user address
	span.AddEvent("Getting user address")
	addressReq, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/address/%s", s.addressServiceURL, userID), nil)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to create address request: %w", err)
	}

	addressResp, err := s.httpClient.Do(addressReq)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to get user address: %w", err)
	}
	defer addressResp.Body.Close()

	if addressResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("address service returned status: %d", addressResp.StatusCode)
		span.RecordError(err)
		return "", err
	}

	// 2. Verify inventory using gRPC
	span.AddEvent("Verifying inventory")
	verifyResp, err := s.inventoryClient.VerifyInventory(ctx, &proto.VerifyInventoryRequest{
		Items: items,
	})
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to verify inventory: %w", err)
	}

	if !verifyResp.Available {
		err = errors.New("items not available in inventory")
		if verifyResp.Error != "" {
			err = errors.New(verifyResp.Error)
		}
		span.RecordError(err)
		return "", err
	}

	// 3. Update inventory using gRPC
	span.AddEvent("Updating inventory")
	orderID := "order-" + userID + "-123"
	updateResp, err := s.inventoryClient.UpdateInventory(ctx, &proto.UpdateInventoryRequest{
		OrderId: orderID,
		Items:   items,
	})
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to update inventory: %w", err)
	}

	if updateResp.Error != "" {
		err = errors.New(updateResp.Error)
		span.RecordError(err)
		return "", err
	}

	// 4. Calculate money
	span.AddEvent("Calculating payment")
	paymentReq, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/calculate-money", s.paymentServiceURL),
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"order_id":"%s","items":%s}`, orderID, convertItemsToJSON(items)))))
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to create payment calculation request: %w", err)
	}

	paymentReq.Header.Set("Content-Type", "application/json")
	paymentResp, err := s.httpClient.Do(paymentReq)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to calculate payment: %w", err)
	}
	defer paymentResp.Body.Close()

	if paymentResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("payment service returned status: %d", paymentResp.StatusCode)
		span.RecordError(err)
		return "", err
	}

	// 5. Apply coupon if available
	span.AddEvent("Applying coupon")
	couponReq, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/apply-coupon", s.paymentServiceURL),
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"order_id":"%s","user_id":"%s"}`, orderID, userID))))
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to create coupon request: %w", err)
	}

	couponReq.Header.Set("Content-Type", "application/json")
	couponResp, err := s.httpClient.Do(couponReq)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to apply coupon: %w", err)
	}
	defer couponResp.Body.Close()

	// We don't fail the order if coupon application fails, just log it
	if couponResp.StatusCode != http.StatusOK {
		span.AddEvent(fmt.Sprintf("Coupon application failed with status: %d", couponResp.StatusCode))
	}

	span.SetAttributes(attribute.String("order.id", orderID))
	return orderID, nil
}

// Helper function to convert items slice to JSON array string
func convertItemsToJSON(items []string) string {
	jsonItems, _ := json.Marshal(items)
	return string(jsonItems)
}
