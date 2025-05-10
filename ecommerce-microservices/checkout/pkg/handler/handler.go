package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/checkout/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
)

var httpDurationMs *prometheus.HistogramVec

func init() {
	fmt.Println("Initializing metrics...")
	httpDurationMs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_milliseconds",
		Help:    "Histogram of HTTP request latency in milliseconds.",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1ms, 2ms, 4ms, ... ~500ms
	}, []string{"path", "method", "status"})

	prometheus.MustRegister(httpDurationMs)
}

func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		durationMs := float64(time.Since(start).Milliseconds())

		status := http.StatusOK
		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				status = he.Code
			} else {
				status = http.StatusInternalServerError
			}
		} else {
			status = c.Response().Status
		}

		httpDurationMs.WithLabelValues(c.Path(), c.Request().Method, http.StatusText(status)).Observe(durationMs)
		return err
	}
}

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
	e.POST("/checkout", h.UserCheckout, otelecho.Middleware("checkout-service"), MetricsMiddleware, logging.CreateLoggingMiddleware("checkout-service"))
}

// UserCheckout handles the checkout request
func (h *CheckoutHandler) UserCheckout(c echo.Context) error {
	// Extract trace context
	ctx := c.Request().Context()

	// Create a span for this handler
	tracer := otel.Tracer("checkout-handler")
	ctx, span := tracer.Start(ctx, "UserCheckout-handler")
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

	logger.Info().Msg("Processing checkout request")

	// Call service
	orderID, err := h.service.UserCheckout(ctx, req.UserID, req.Items)
	if err != nil {
		// Add error tag to span
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.RecordError(err)

		logger.Error().Err(err).Msg("Checkout failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logger.Info().Str("order_id", orderID).Msg("Checkout successful")

	// Return response
	return c.JSON(http.StatusOK, map[string]string{"order_id": orderID})
}
