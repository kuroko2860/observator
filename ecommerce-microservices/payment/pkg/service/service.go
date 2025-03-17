package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/pkg/tracing"
)

// PaymentService describes the service
type PaymentService interface {
	CalculateMoney(ctx context.Context, orderID string, items []string) (float64, error)
	ApplyCoupon(ctx context.Context, orderID string, couponCode string, amount float64) (float64, error)
}

type paymentService struct{}

// NewPaymentService returns a new implementation of PaymentService
func NewPaymentService() PaymentService {
	return &paymentService{}
}

// CalculateMoney implements PaymentService
func (s *paymentService) CalculateMoney(ctx context.Context, orderID string, items []string) (float64, error) {
	// Create a span for the calculate money operation
	tracer := tracing.Tracer("payment-service")
	ctx, span := tracer.Start(ctx, "CalculateMoney-Service")
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

	logger.Debug().Msg("Starting payment calculation")

	if len(items) == 0 {
		err := errors.New("no items to calculate")
		span.RecordError(err)
		logger.Error().Err(err).Msg("Payment calculation failed")
		return 0, err
	}

	// In a real implementation, this would calculate the total based on item prices
	// For simplicity, we'll just use a fixed price per item
	total := float64(len(items)) * 10.0

	span.SetAttributes(attribute.Float64("total.amount", total))
	logger.Info().Float64("amount", total).Msg("Payment calculated successfully")
	return total, nil
}

// ApplyCoupon implements PaymentService
func (s *paymentService) ApplyCoupon(ctx context.Context, orderID string, couponCode string, amount float64) (float64, error) {
	// Create a span for the apply coupon operation
	tracer := tracing.Tracer("payment-service")
	ctx, span := tracer.Start(ctx, "ApplyCoupon-Service")
	defer span.End()

	// Extract trace context for logging
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("order_id", orderID).
		Str("coupon_code", couponCode).
		Float64("original_amount", amount).
		Logger()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.String("coupon.code", couponCode),
		attribute.Float64("original.amount", amount),
	)

	logger.Debug().Msg("Applying coupon")

	if couponCode == "" {
		logger.Info().Msg("No coupon code provided, returning original amount")
		return amount, nil
	}

	// In a real implementation, this would validate the coupon and apply the discount
	// For simplicity, we'll just apply a 10% discount for any coupon
	discountedAmount := amount * 0.9

	span.SetAttributes(attribute.Float64("discounted.amount", discountedAmount))
	logger.Info().
		Float64("original_amount", amount).
		Float64("discounted_amount", discountedAmount).
		Msg("Coupon applied successfully")

	return discountedAmount, nil
}
