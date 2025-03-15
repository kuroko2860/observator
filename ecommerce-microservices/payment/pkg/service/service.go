package service

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	"kltn/ecommerce-microservices/pkg/tracing"
)

// PaymentService describes the service
type PaymentService interface {
	CalculateMoney(ctx context.Context, orderID string, items []string) (float64, error)
	ApplyCoupon(ctx context.Context, orderID string, couponCode string, amount float64) (float64, error)
}

type basicPaymentService struct{}

// NewBasicPaymentService returns a naive, stateless implementation of PaymentService
func NewBasicPaymentService() PaymentService {
	return &basicPaymentService{}
}

// CalculateMoney implements PaymentService
func (s *basicPaymentService) CalculateMoney(ctx context.Context, orderID string, items []string) (float64, error) {
	// Create a span for the calculate money operation
	tracer := tracing.Tracer("payment-service")
	ctx, span := tracer.Start(ctx, "CalculateMoney")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.Int("items.count", len(items)),
	)

	if len(items) == 0 {
		span.RecordError(errors.New("no items to calculate"))
		return 0, errors.New("no items to calculate")
	}

	// In a real implementation, this would calculate the total based on item prices
	// For simplicity, we'll just use a fixed price per item
	total := float64(len(items)) * 10.0
	
	span.SetAttributes(attribute.Float64("total.amount", total))
	return total, nil
}

// ApplyCoupon implements PaymentService
func (s *basicPaymentService) ApplyCoupon(ctx context.Context, orderID string, couponCode string, amount float64) (float64, error) {
	// Create a span for the apply coupon operation
	tracer := tracing.Tracer("payment-service")
	ctx, span := tracer.Start(ctx, "ApplyCoupon")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.String("coupon.code", couponCode),
		attribute.Float64("original.amount", amount),
	)

	if couponCode == "" {
		return amount, nil
	}

	// In a real implementation, this would validate the coupon and apply the discount
	// For simplicity, we'll just apply a 10% discount for any coupon
	discountedAmount := amount * 0.9
	
	span.SetAttributes(attribute.Float64("discounted.amount", discountedAmount))
	return discountedAmount, nil
}