package logging

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

var (
	natsConn *nats.Conn
	natsMu   sync.Mutex
)

// InitNATS initializes the connection to NATS server
func InitNATS(natsURL string) error {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn != nil && natsConn.IsConnected() {
		return nil
	}

	var err error
	natsConn, err = nats.Connect(natsURL,
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Error().Err(err).Msg("Disconnected from NATS")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info().Msg("Reconnected to NATS")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Info().Msg("NATS connection closed")
		}),
	)

	if err != nil {
		return err
	}

	log.Info().Str("url", natsURL).Msg("Connected to NATS")
	return nil
}

// CloseNATS closes the connection to NATS server
func CloseNATS() {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn != nil {
		natsConn.Close()
		natsConn = nil
		log.Info().Msg("NATS connection closed")
	}
}

// HttpLogEntry represents a log entry for HTTP requests
type HttpLogEntry struct {
	ServiceName   string `json:"service_name" bson:"service_name"`
	URIPath       string `json:"uri_path" bson:"uri_path"`
	Referer       string `json:"referer" bson:"referer"`
	UserId        string `json:"user_id" bson:"user_id"`
	Method        string `json:"method" bson:"method"`
	StartTime     int64  `json:"start_time" bson:"start_time"`
	StartTimeDate string `json:"start_time_date" bson:"start_time_date"`
	Host          string `json:"host" bson:"host"`
	Protocol      string `json:"protocol" bson:"protocol"`
	RemoteIP      string `json:"remote_ip" bson:"remote_ip"`
	RequestId     string `json:"request_id" bson:"request_id"`
	TraceId       string `json:"trace_id" bson:"trace_id"`
	SpanId        string `json:"span_id" bson:"span_id"`
	UserAgent     string `json:"user_agent" bson:"user_agent"`
	Duration      int64  `json:"duration" bson:"duration"`
	ResquestSize  string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize  int64  `json:"response_size" bson:"response_size"`
	StatusCode    int    `json:"status_code" bson:"status_code"`
	ErrorMessage  string `json:"error_message,omitempty" bson:"error_message,omitempty"`
}

// PublishLogEntry publishes a log entry to NATS
func PublishHttpRequestLogEntry(entry HttpLogEntry) {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn == nil || !natsConn.IsConnected() {
		log.Warn().Msg("Cannot publish http log entry: NATS not connected")
		return
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal http log entry")
		return
	}

	err = natsConn.Publish("logs.http", entryJSON)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish http log to NATS")
	}
}

type LogEntry struct {
	ServiceName string `json:"service_name"`
	Message     string `json:"message"`
	Level       string `json:"level"`
	Caller      string `json:"caller"`
	TraceID     string `json:"trace_id"`
	SpanID      string `json:"span_id"`
	StartTime   int64  `json:"start_time"`
}

func PublishLogEntry(entry LogEntry) {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn == nil || !natsConn.IsConnected() {
		log.Warn().Msg("Cannot publish log entry: NATS not connected")
		return
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log entry")
		return
	}

	err = natsConn.Publish("log.internal", entryJSON)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish log to NATS")
	}
}

// Add this function to expose the NATS connection
func GetNATSConnection() *nats.Conn {
	natsMu.Lock()
	defer natsMu.Unlock()
	return natsConn
}

// CreateLoggingMiddleware creates a middleware that logs requests with trace and span IDs
func CreateLoggingMiddleware(serviceName string) echo.MiddlewareFunc {
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
				Str("trace_id", traceID).
				Str("span_id", spanID).
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
			entry := HttpLogEntry{
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
				Duration:      duration.Microseconds(),
				StatusCode:    res.Status,
				ResquestSize:  req.Header.Get("Content-Length"),
				ResponseSize:  int64(res.Size),
			}

			// Publish to NATS
			PublishHttpRequestLogEntry(entry)

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
	entry := LogEntry{
		StartTime:   time.Now().UnixMilli(),
		ServiceName: w.ServiceName,
	}

	// Copy relevant fields from the log event
	if level, ok := logEvent["level"].(string); ok {
		entry.Level = level
	}

	if msg, ok := logEvent["message"].(string); ok {
		entry.Message = msg
	}

	// Include caller information if available
	if caller, ok := logEvent["caller"].(string); ok {
		entry.Caller = caller
	}

	// Include trace and span IDs if available
	if traceID, ok := logEvent["trace_id"].(string); ok {
		entry.TraceID = traceID
	}

	if spanID, ok := logEvent["span_id"].(string); ok {
		entry.SpanID = spanID
	}

	// Publish to NATS
	PublishLogEntry(entry)

	return len(p), nil
}
