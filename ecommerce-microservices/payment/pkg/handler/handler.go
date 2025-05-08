package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/payment/pkg/service"
)

// PaymentHandler handles HTTP requests for the payment service
type PaymentHandler struct {
	service service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(svc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: svc,
	}
}

// RegisterRoutes registers the handler routes with the Echo instance
func (h *PaymentHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/calculate-money", h.CalculateMoney)
	e.POST("/apply-coupon", h.ApplyCoupon)
}

// CalculateMoney handles the calculate money request
func (h *PaymentHandler) CalculateMoney(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	// tracer := otel.Tracer("payment-handler")
	// ctx, span := tracer.Start(ctx, "CalculateMoney-Handler")
	// defer span.End()

	// Parse request
	var req struct {
		OrderID string   `json:"order_id"`
		Items   []string `json:"items"`
	}

	if err := c.Bind(&req); err != nil {
		// Add error tag to span
		// span.SetAttributes(attribute.Bool("error", true))
		// span.SetAttributes(attribute.String("error.message", err.Error()))
		// span.RecordError(err)

		log.Error().Err(err).Msg("Invalid request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Log request details with trace information
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("order_id", req.OrderID).
		Int("items_count", len(req.Items)).
		Logger()

	logger.Info().Msg("Processing calculate money request")

	// Call service
	amount, err := h.service.CalculateMoney(ctx, req.OrderID, req.Items)
	if err != nil {
		// Add error tag to span
		// span.SetAttributes(attribute.Bool("error", true))
		// span.SetAttributes(attribute.String("error.message", err.Error()))
		// span.RecordError(err)

		logger.Error().Err(err).Msg("Calculate money failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().Float64("amount", amount).Msg("Calculate money completed")

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"amount": amount,
	})
}

// ApplyCoupon handles the apply coupon request
func (h *PaymentHandler) ApplyCoupon(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	// tracer := otel.Tracer("payment-handler")
	// ctx, span := tracer.Start(ctx, "ApplyCoupon-Handler")
	// defer span.End()

	// Parse request
	var req struct {
		OrderID    string  `json:"order_id"`
		CouponCode string  `json:"coupon_code"`
		Amount     float64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		// Add error tag to span
		// span.SetAttributes(attribute.Bool("error", true))
		// span.SetAttributes(attribute.String("error.message", err.Error()))
		// span.RecordError(err)

		log.Error().Err(err).Msg("Invalid request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Log request details with trace information
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("order_id", req.OrderID).
		Str("coupon_code", req.CouponCode).
		Float64("amount", req.Amount).
		Logger()

	logger.Info().Msg("Processing apply coupon request")

	// Call service
	discountedAmount, err := h.service.ApplyCoupon(ctx, req.OrderID, req.CouponCode, req.Amount)
	if err != nil {
		// Add error tag to span
		// span.SetAttributes(attribute.Bool("error", true))
		// span.SetAttributes(attribute.String("error.message", err.Error()))
		// span.RecordError(err)

		logger.Error().Err(err).Msg("Apply coupon failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().
		Float64("original_amount", req.Amount).
		Float64("discounted_amount", discountedAmount).
		Msg("Apply coupon completed")

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"discounted_amount": discountedAmount,
	})
}
