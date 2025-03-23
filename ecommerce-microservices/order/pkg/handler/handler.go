package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/order/pkg/service"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// OrderHandler handles HTTP requests for the order service
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: svc,
	}
}

// RegisterRoutes registers the handler routes with the Echo instance
func (h *OrderHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/orders", h.CreateOrder)
}

// CreateOrder handles the create order request
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	tracer := tracing.Tracer("order-handler")
	ctx, span := tracer.Start(ctx, "CreateOrder-handler")
	defer span.End()

	// Parse request
	var req struct {
		UserID string   `json:"user_id"`
		Items  []string `json:"items"`
	}

	if err := c.Bind(&req); err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
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

	logger.Info().Msg("Processing order creation request")

	// Call service
	orderID, err := h.service.CreateOrder(ctx, req.UserID, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)
		
		logger.Error().Err(err).Msg("Order creation failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().Str("order_id", orderID).Msg("Order created successfully")

	// Return response
	return c.JSON(http.StatusOK, map[string]string{"order_id": orderID})
}
