package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/address/pkg/service"
)

// AddressHandler handles HTTP requests for the address service
type AddressHandler struct {
	service service.AddressService
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(svc service.AddressService) *AddressHandler {
	return &AddressHandler{
		service: svc,
	}
}

// RegisterRoutes registers the handler routes with the Echo instance
func (h *AddressHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/address/:user_id", h.GetAddress)
}

// GetAddress handles the get address request
func (h *AddressHandler) GetAddress(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	tracer := otel.Tracer("address-handler")
	ctx, span := tracer.Start(ctx, "GetAddress-Handler")
	defer span.End()

	// Get user ID from path parameter
	userID := c.Param("user_id")

	// Log request details with trace information
	spanContext := trace.SpanContextFromContext(ctx)
	logger := log.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Str("user_id", userID).
		Logger()

	logger.Info().Msg("Processing get address request")

	// Call service
	address, err := h.service.GetAddress(ctx, userID)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Get address failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().
		Str("street", address.Street).
		Str("city", address.City).
		Msg("Address retrieved successfully")

	// Return response
	return c.JSON(http.StatusOK, address)
}
