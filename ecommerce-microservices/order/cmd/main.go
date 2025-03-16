package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/order/pkg/handler"
	"kltn/ecommerce-microservices/order/pkg/service"

	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	// Parse command line flags
	var (
		httpAddr            = flag.String("http.addr", ":8081", "HTTP listen address")
		paymentServiceURL   = flag.String("payment.url", "http://localhost:8082", "Payment service URL")
		inventoryServiceURL = flag.String("inventory.url", "localhost:8083", "Inventory service gRPC address")
		addressServiceURL   = flag.String("address.url", "http://localhost:8084", "Address service URL")
		zipkinURL           = flag.String("zipkin.url", "http://localhost:9411/api/v2/spans", "Zipkin server URL")
		natsURL             = flag.String("nats.url", "nats://nats:4222", "NATS server URL")
	)
	flag.Parse()

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Timestamp().Logger()

	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Initialize the tracer
	shutdown, err := tracing.InitTracer("order-service", *zipkinURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize tracer")
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// Create an instrumented HTTP client for service-to-service communication
	httpClient := tracing.NewTracedHTTPClient()

	// Initialize NATS connection
	err = logging.InitNATS(*natsURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to NATS")
	} else {
		defer logging.CloseNATS()
	}

	// Ensure inventory service URL doesn't have http:// prefix for gRPC
	inventoryURL := *inventoryServiceURL
	if strings.HasPrefix(inventoryURL, "http://") {
		inventoryURL = strings.TrimPrefix(inventoryURL, "http://")
	}

	// Update the route registration section in main.go

	// Create the service
	svc := service.NewOrderService(*paymentServiceURL, inventoryURL, *addressServiceURL, httpClient)

	// Create and register handlers
	orderHandler := handler.NewOrderHandler(svc)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware("order-service"))
	e.Use(createLoggingMiddleware("order-service"))

	// Register routes
	orderHandler.RegisterRoutes(e)

	// Start server in a goroutine
	go func() {
		log.Info().Str("transport", "HTTP").Str("addr", *httpAddr).Msg("Starting server")
		if err := e.Start(*httpAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

// createLoggingMiddleware creates a middleware that logs requests with trace and span IDs
func createLoggingMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			// Extract trace context
			spanContext := trace.SpanContextFromContext(req.Context())
			traceID := ""
			spanID := ""
			if spanContext.IsValid() {
				traceID = spanContext.TraceID().String()
				spanID = spanContext.SpanID().String()
			}

			// Process request
			err := next(c)

			// Log request details
			duration := time.Since(start)

			// Create log entry with zerolog
			logger := log.With().
				Str("service", serviceName).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Str("trace_id", traceID).
				Str("span_id", spanID).
				Str("remote_addr", req.RemoteAddr).
				Str("user_agent", req.UserAgent()).
				Dur("duration", duration).
				Logger()

			// Log based on status code
			if res.Status >= 500 {
				logger.Error().Err(err).Msg("Server error")
			} else if res.Status >= 400 {
				logger.Warn().Err(err).Msg("Client error")
			} else {
				logger.Info().Msg("Request completed")
			}

			// Create HttpLogEntry for NATS
			entry := logging.HttpLogEntry{
				ServiceName:   serviceName,
				URIPath:       req.URL.Path,
				Referer:       req.Referer(),
				UserId:        req.Header.Get("User-ID"),
				Method:        req.Method,
				StartTime:     start.UnixMilli(),
				StartTimeDate: start.Format(time.RFC3339),
				Host:          req.Host,
				Protocol:      req.Proto,
				RemoteIP:      req.RemoteAddr,
				RequestId:     c.Response().Header().Get(echo.HeaderXRequestID),
				TraceId:       traceID,
				SpanId:        spanID,
				UserAgent:     req.UserAgent(),
				Duration:      duration.Milliseconds(),
				StatusCode:    res.Status,
				ResquestSize:  req.Header.Get("Content-Length"),
				ResponseSize:  int64(res.Size),
			}

			// Publish to NATS
			logging.PublishLogEntry(entry)

			return err
		}
	}
}
