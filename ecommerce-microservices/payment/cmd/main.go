package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/payment/pkg/handler"
	"kltn/ecommerce-microservices/payment/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

func main() {
	// Parse command line flags
	var (
		httpAddr  = flag.String("http.addr", ":8082", "HTTP listen address")
		zipkinURL = flag.String("zipkin.url", "http://localhost:9411/api/v2/spans", "Zipkin server URL")
		natsURL   = flag.String("nats.url", "nats://nats:4222", "NATS server URL")
		esURL     = flag.String("es.url", "http://elasticsearch:9200", "Elasticsearch URL")
	)
	flag.Parse()

	// Initialize NATS connection first (before configuring zerolog)
	err := logging.InitNATS(*natsURL)
	if err != nil {
		// Use standard log here since zerolog isn't configured yet
		log.Error().Err(err).Msg("Failed to connect to NATS")
	} else {
		defer logging.CloseNATS()

		// Set up custom zerolog writer that sends logs to NATS
		natsWriter := &NATSLogWriter{ServiceName: "address-service"}
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		multi := zerolog.MultiLevelWriter(consoleWriter, natsWriter)

		// Configure zerolog with the multi-writer
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()
	}

	// Set log level
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Initialize the tracer
	shutdown, err := tracing.InitTracer("payment-service", *zipkinURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize tracer")
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// Initialize NATS connection
	err = logging.InitNATS(*natsURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to NATS")
	} else {
		defer logging.CloseNATS()
	}

	// Initialize Elasticsearch connection
	err = logging.InitElasticsearch(*esURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Elasticsearch")
	} else {
		// Set up NATS to Elasticsearch bridge
		err = logging.SetupNATSToElasticsearchBridge(logging.GetNATSConnection(), "microservices-logs")
		if err != nil {
			log.Error().Err(err).Msg("Failed to set up NATS to Elasticsearch bridge")
		}
	}

	// Create the service
	svc := service.NewPaymentService()

	// Create and register handlers
	paymentHandler := handler.NewPaymentHandler(svc)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware("payment-service"))
	e.Use(createLoggingMiddleware("payment-service"))

	// Register routes
	paymentHandler.RegisterRoutes(e)

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

// NATSLogWriter is a custom zerolog writer that sends logs to NATS
type NATSLogWriter struct {
	ServiceName string
}

// Write implements io.Writer interface
func (w *NATSLogWriter) Write(p []byte) (n int, err error) {
	// Parse the JSON log entry
	var logEvent map[string]interface{}
	if err := json.Unmarshal(p, &logEvent); err != nil {
		return 0, err
	}

	// Create a log entry for NATS
	entry := logging.HttpLogEntry{
		ServiceName:   w.ServiceName,
		URIPath:       "internal", // Indicates this is an internal log, not an HTTP request
		Method:        "LOG",
		StartTime:     time.Now().UnixMilli(),
		StartTimeDate: time.Now().Format(time.RFC3339),
	}

	// Copy relevant fields from the log event
	if level, ok := logEvent["level"].(string); ok {
		entry.URIPath = "internal/" + level // Include log level in the path
	}

	if msg, ok := logEvent["message"].(string); ok {
		entry.ErrorMessage = msg
	}

	// Include caller information if available
	if caller, ok := logEvent["caller"].(string); ok {
		entry.Referer = caller
	}

	// Include trace and span IDs if available
	if traceID, ok := logEvent["trace_id"].(string); ok {
		entry.TraceId = traceID
	}

	if spanID, ok := logEvent["span_id"].(string); ok {
		entry.SpanId = spanID
	}

	// Publish to NATS
	logging.PublishLogEntry(entry)

	return len(p), nil
}
