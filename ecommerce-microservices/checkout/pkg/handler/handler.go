package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/checkout/pkg/service"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// CheckoutHandler handles HTTP requests for the checkout service
type CheckoutHandler struct {
	service service.CheckoutService
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(svc service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		service: svc,
	}
}

// RegisterRoutes registers the handler routes with the Echo instance
func (h *CheckoutHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/checkout", h.UserCheckout)
}

// UserCheckout handles the checkout request
func (h *CheckoutHandler) UserCheckout(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()
	
	// Create a span for this handler
	tracer := tracing.Tracer("checkout-handler")
	ctx, span := tracer.Start(ctx, "UserCheckout")
	defer span.End()
	
	// Parse request
	var req struct {
		UserID string   `json:"user_id"`
		Items  []string `json:"items"`
	}
	
	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	// Log request details with trace information
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("user_id", req.UserID).
		Int("items_count", len(req.Items)).
		Logger()
	
	logger.Info().Msg("Processing checkout request")
	
	// Call service
	orderID, err := h.service.UserCheckout(ctx, req.UserID, req.Items)
	if err != nil {
		logger.Error().Err(err).Msg("Checkout failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	logger.Info().Str("order_id", orderID).Msg("Checkout successful")
	
	// Return response
	return c.JSON(http.StatusOK, map[string]string{"order_id": orderID})
}