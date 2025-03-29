package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/inventory/pkg/service"
)

// InventoryHandler handles HTTP requests for the inventory service
type InventoryHandler struct {
	service service.InventoryService
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: svc,
	}
}

// RegisterRoutes registers the handler routes with the Echo instance
func (h *InventoryHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/inventory/verify", h.VerifyInventory)
	e.POST("/inventory/update", h.UpdateInventory)
}

// VerifyInventory handles the verify inventory request
func (h *InventoryHandler) VerifyInventory(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	tracer := otel.Tracer("inventory-handler")
	ctx, span := tracer.Start(ctx, "VerifyInventory")
	defer span.End()

	// Parse request
	var req struct {
		Items []string `json:"items"`
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
		Int("items_count", len(req.Items)).
		Logger()

	logger.Info().Msg("Processing inventory verification request")

	// Call service
	available, err := h.service.VerifyInventory(ctx, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Inventory verification failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().Bool("available", available).Msg("Inventory verification completed")

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"available": available,
	})
}

// UpdateInventory handles the update inventory request
func (h *InventoryHandler) UpdateInventory(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	tracer := otel.Tracer("inventory-handler")
	ctx, span := tracer.Start(ctx, "UpdateInventory")
	defer span.End()

	// Parse request
	var req struct {
		OrderID string   `json:"order_id"`
		Items   []string `json:"items"`
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
		Str("order_id", req.OrderID).
		Int("items_count", len(req.Items)).
		Logger()

	logger.Info().Msg("Processing inventory update request")

	// Call service
	err := h.service.UpdateInventory(ctx, req.OrderID, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Inventory update failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().Str("order_id", req.OrderID).Msg("Inventory updated successfully")

	// Return response
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
