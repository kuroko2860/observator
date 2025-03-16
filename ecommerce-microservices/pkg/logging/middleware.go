package logging

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type HttpLogEntry struct {
	ServiceName   string `json:"service_name" bson:"service_name"`
	URIPath       string `json:"uri_path" bson:"uri_path"`
	Referer       string `json:"referer" bson:"referer"`
	UserId        string `json:"user_id" bson:"user_id"`
	Username      string `json:"username" bson:"username"`
	StartTime     int64  `json:"start_time" bson:"start_time"`
	Method        string `json:"method" bson:"method"`
	StartTimeDate string `json:"start_time_date" bson:"start_time_date"`
	Host          string `json:"host" bson:"host"`

	Protocol  string `json:"protocol" bson:"protocol"`
	RemoteIP  string `json:"remote_ip" bson:"remote_ip"`
	RequestId string `json:"request_id" bson:"request_id"`
	// TraceId      string `json:"trace_id" bson:"trace_id"`
	// PathId       uint32 `json:"path_id" bson:"path_id"`
	UserAgent    string `json:"user_agent" bson:"user_agent"`
	Duration     int64  `json:"duration" bson:"duration"`
	ResquestSize string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize int64  `json:"response_size" bson:"response_size"`
	StatusCode   int    `json:"status_code" bson:"status_code"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
}

// natsConn is the shared NATS connection used for publishing logs
var natsConn *nats.Conn

// InitNATS initializes the connection to NATS server
func InitNATS(natsURL string) error {
	var err error
	natsConn, err = nats.Connect(natsURL)
	return err
}

// CloseNATS closes the connection to NATS server
func CloseNATS() {
	if natsConn != nil {
		natsConn.Close()
	}
}

// HTTPMiddleware returns a handler that logs the HTTP requests
func HTTPMiddleware(logger log.Logger, serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()
			r.Header.Set("X-Request-ID", uuid.New().String())

			// Wrap the response writer to capture the status code and response size
			wrw := &responseWriter{w: w, status: http.StatusOK, responseSize: 0}

			next.ServeHTTP(wrw, r)

			duration := time.Since(begin)

			// Create HttpLogEntry
			entry := HttpLogEntry{
				ServiceName:   serviceName,
				URIPath:       r.URL.Path,
				Referer:       r.Referer(),
				UserId:        r.Header.Get("User-ID"),
				Method:        r.Method,
				StartTime:     begin.UnixMilli(),
				StartTimeDate: begin.Format(time.RFC3339),
				Host:          r.Host,
				Protocol:      r.Proto,
				RemoteIP:      r.RemoteAddr,
				RequestId:     r.Header.Get("X-Request-ID"),
				UserAgent:     r.UserAgent(),
				Duration:      duration.Milliseconds(),
				StatusCode:    wrw.status,
				ResquestSize:  r.Header.Get("Content-Length"),
				ResponseSize:  wrw.responseSize,
			}
			logger.Log("entry", entry)

			// Publish to NATS if connection is available
			if natsConn != nil && natsConn.IsConnected() {
				entryJSON, err := json.Marshal(entry)
				if err == nil {
					// Publish to "logs" subject
					err = natsConn.Publish("logs", entryJSON)
					if err != nil {
						logger.Log("msg", "Failed to publish log to NATS", "error", err)
					} else {
						logger.Log("msg", "Published log to NATS")
					}

				} else {
					logger.Log("msg", "Failed to marshal log entry", "error", err)
				}
			}

			// Still log to the console
			logger.Log(
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrw.status,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"took", duration,
			)
		})
	}
}

// EndpointLoggingMiddleware logs the duration of each request
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log(
					"transport_error", err,
					"took", time.Since(begin),
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// ServerFinalizer is a go-kit HTTP server finalizer that logs the response status
func ServerFinalizer(logger log.Logger) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		logger.Log(
			"method", r.Method,
			"path", r.URL.Path,
			"status", code,
		)
	}
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code and response size
type responseWriter struct {
	w           http.ResponseWriter
	status      int
	responseSize int64
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.w.Write(b)
	rw.responseSize += int64(n)
	return n, err
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}
