package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"kltn/ecommerce-microservices/inventory/proto"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// OrderService describes the service
type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []string) (string, error)
}

type orderService struct {
	paymentServiceURL   string
	inventoryServiceURL string
	addressServiceURL   string
	httpClient          *http.Client
	inventoryClient     proto.InventoryServiceClient
}

// NewOrderService returns a new implementation of OrderService
func NewOrderService(paymentURL, inventoryURL, addressURL string, client *http.Client) OrderService {
	// Set up gRPC connection to inventory service with OpenTelemetry interceptors
	conn, err := grpc.Dial(
		inventoryURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatal().Err(err).Str("url", inventoryURL).Msg("Failed to connect to inventory service")
	}

	inventoryClient := proto.NewInventoryServiceClient(conn)

	return &orderService{
		paymentServiceURL:   paymentURL,
		inventoryServiceURL: inventoryURL,
		addressServiceURL:   addressURL,
		httpClient:          client,
		inventoryClient:     inventoryClient,
	}
}

// CreateOrder implements OrderService
func (s *orderService) CreateOrder(ctx context.Context, userID string, items []string) (string, error) {
	// Create a span for the create order operation
	tracer := tracing.Tracer("order-service")
	ctx, span := tracer.Start(ctx, "CreateOrder-service")
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

	logger.Debug().Msg("Starting order creation process")

	if len(items) == 0 {
		err := errors.New("no items in order")
		span.RecordError(err)
		logger.Error().Err(err).Msg("Order creation failed")
		return "", err
	}

	// 1. Get user address
	logger.Debug().Msg("Getting user address")
	span.AddEvent("Getting user address")
	addressReq, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/address/%s", s.addressServiceURL, userID), nil)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to create address request")
		return "", fmt.Errorf("failed to create address request: %w", err)
	}

	addressResp, err := s.httpClient.Do(addressReq)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to get user address")
		return "", fmt.Errorf("failed to get user address: %w", err)
	}
	defer addressResp.Body.Close()

	if addressResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("address service returned status: %d", addressResp.StatusCode)
		span.RecordError(err)
		logger.Error().Err(err).Int("status_code", addressResp.StatusCode).Msg("Address service error")
		return "", err
	}

	// 2. Verify inventory using gRPC
	logger.Debug().Msg("Verifying inventory")
	span.AddEvent("Verifying inventory")
	verifyResp, err := s.inventoryClient.VerifyInventory(ctx, &proto.VerifyInventoryRequest{
		Items: items,
	})
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to verify inventory")
		return "", fmt.Errorf("failed to verify inventory: %w", err)
	}

	if !verifyResp.Available {
		err = errors.New("items not available in inventory")
		if verifyResp.Error != "" {
			err = errors.New(verifyResp.Error)
		}
		span.RecordError(err)
		logger.Error().Err(err).Msg("Inventory verification failed")
		return "", err
	}

	// 3. Update inventory using gRPC
	logger.Debug().Msg("Updating inventory")
	span.AddEvent("Updating inventory")
	orderID := "order-" + userID + "-123"
	updateResp, err := s.inventoryClient.UpdateInventory(ctx, &proto.UpdateInventoryRequest{
		OrderId: orderID,
		Items:   items,
	})
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to update inventory")
		return "", fmt.Errorf("failed to update inventory: %w", err)
	}

	if updateResp.Error != "" {
		err = errors.New(updateResp.Error)
		span.RecordError(err)
		logger.Error().Err(err).Msg("Inventory update failed")
		return "", err
	}

	// 4. Process payment
	logger.Debug().Msg("Processing payment")
	span.AddEvent("Processing payment")
	paymentReq := struct {
		UserID  string   `json:"user_id"`
		OrderID string   `json:"order_id"`
		Items   []string `json:"items"`
	}{
		UserID:  userID,
		OrderID: orderID,
		Items:   items,
	}

	paymentJSON, err := json.Marshal(paymentReq)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to marshal payment request")
		return "", fmt.Errorf("failed to marshal payment request: %w", err)
	}

	paymentReqObj, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/calculate-money", s.paymentServiceURL), bytes.NewBuffer(paymentJSON))
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to create payment request")
		return "", fmt.Errorf("failed to create payment request: %w", err)
	}
	paymentReqObj.Header.Set("Content-Type", "application/json")

	paymentResp, err := s.httpClient.Do(paymentReqObj)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Payment service request failed")
		return "", fmt.Errorf("failed to process payment: %w", err)
	}
	defer paymentResp.Body.Close()

	if paymentResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("payment service returned status: %d", paymentResp.StatusCode)
		span.RecordError(err)
		logger.Error().Err(err).Int("status_code", paymentResp.StatusCode).Msg("Payment service error")
		return "", err
	}

	// Parse the payment response to get the amount
	var paymentRespBody struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(paymentResp.Body).Decode(&paymentRespBody); err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to decode payment response")
		return "", fmt.Errorf("failed to decode payment response: %w", err)
	}

	// 4.1 Apply coupon if available (for this example, we'll use a fixed coupon code)
	logger.Debug().Msg("Applying coupon")
	span.AddEvent("Applying coupon")

	// In a real implementation, the coupon code would come from the request
	// For this example, we'll use a fixed coupon code "DISCOUNT10"
	couponCode := "DISCOUNT10"

	couponReq := struct {
		OrderID    string  `json:"order_id"`
		CouponCode string  `json:"coupon_code"`
		Amount     float64 `json:"amount"`
	}{
		OrderID:    orderID,
		CouponCode: couponCode,
		Amount:     paymentRespBody.Amount,
	}

	couponJSON, err := json.Marshal(couponReq)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to marshal coupon request")
		return "", fmt.Errorf("failed to marshal coupon request: %w", err)
	}

	couponReqObj, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/apply-coupon", s.paymentServiceURL), bytes.NewBuffer(couponJSON))
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to create coupon request")
		return "", fmt.Errorf("failed to create coupon request: %w", err)
	}
	couponReqObj.Header.Set("Content-Type", "application/json")

	couponResp, err := s.httpClient.Do(couponReqObj)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Coupon service request failed")
		return "", fmt.Errorf("failed to apply coupon: %w", err)
	}
	defer couponResp.Body.Close()

	if couponResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("coupon service returned status: %d", couponResp.StatusCode)
		span.RecordError(err)
		logger.Error().Err(err).Int("status_code", couponResp.StatusCode).Msg("Coupon service error")
		return "", err
	}

	// Parse the coupon response to get the discounted amount
	var couponRespBody struct {
		DiscountedAmount float64 `json:"discounted_amount"`
	}
	if err := json.NewDecoder(couponResp.Body).Decode(&couponRespBody); err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Failed to decode coupon response")
		return "", fmt.Errorf("failed to decode coupon response: %w", err)
	}

	logger.Info().
		Float64("original_amount", paymentRespBody.Amount).
		Float64("discounted_amount", couponRespBody.DiscountedAmount).
		Str("coupon_code", couponCode).
		Msg("Coupon applied successfully")

	// 5. Create order record in database (simulated)
	logger.Debug().Msg("Creating order record")
	span.AddEvent("Creating order record")
	// In a real implementation, this would save to a database
	// For now, we just log and return the order ID

	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.Float64("order.amount", couponRespBody.DiscountedAmount),
		attribute.String("coupon.code", couponCode),
	)
	logger.Info().
		Str("order_id", orderID).
		Str("user_id", userID).
		Int("items_count", len(items)).
		Float64("final_amount", couponRespBody.DiscountedAmount).
		Msg("Order created successfully")

	return orderID, nil
}

// Helper function to calculate the total amount
func calculateAmount(items []string) int {
	// In a real implementation, this would look up prices from a database
	// For now, we just return a fixed amount per item
	return len(items) * 10
}
