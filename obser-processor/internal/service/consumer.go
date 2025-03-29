package service

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	v1 "go.opentelemetry.io/proto/otlp/common/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
	"kuroko.com/processor/internal/types"
)

var (
	httpAddr   = flag.String("http.addr", ":8085", "HTTP listen address")
	natsURL    = flag.String("nats.url", "nats://localhost:4222", "NATS server URL")
	natsSubj   = flag.String("nats.subject", "otel-spans", "NATS subject for spans")
	bufferTime = flag.Duration("buffer.time", 30*time.Second, "Time to buffer spans before processing")
)

// TraceStore stores spans by trace ID
type TraceStore struct {
	mu     sync.RWMutex
	traces map[string][]*tracepb.Span
	times  map[string]time.Time
}

// NewTraceStore creates a new trace store
func NewTraceStore() *TraceStore {
	return &TraceStore{
		traces: make(map[string][]*tracepb.Span),
		times:  make(map[string]time.Time),
	}
}

// AddSpan adds a span to the store
func (ts *TraceStore) AddSpan(span *tracepb.Span) {
	traceID := fmt.Sprintf("%x", span.TraceId)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.traces[traceID] = append(ts.traces[traceID], span)
	ts.times[traceID] = time.Now()
}

// GetExpiredTraces returns traces that have not been updated for the given duration
func (ts *TraceStore) GetExpiredTraces(d time.Duration) map[string][]*tracepb.Span {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	expired := make(map[string][]*tracepb.Span)
	now := time.Now()

	for traceID, lastUpdate := range ts.times {
		if now.Sub(lastUpdate) > d {
			expired[traceID] = ts.traces[traceID]
			delete(ts.traces, traceID)
			delete(ts.times, traceID)
		}
	}

	return expired
}

// ConvertSpanToDocument converts a span to an Elasticsearch document
func convertSpanToSpanResponse(span *tracepb.Span) *types.SpanResponse {
	traceID := fmt.Sprintf("%x", span.TraceId)
	spanID := fmt.Sprintf("%x", span.SpanId)

	var parentSpanID string
	if len(span.ParentSpanId) > 0 {
		parentSpanID = fmt.Sprintf("%x", span.ParentSpanId)
	}

	startTime := int64(span.StartTimeUnixNano)
	endTime := int64(span.EndTimeUnixNano)
	duration := int((endTime - startTime)) // Convert to microseconds

	// Convert attributes to map
	attributes := make(map[string]interface{})
	for _, attr := range span.Attributes {
		attributes[attr.Key] = convertAnyValue(attr.Value)
	}

	// Convert events
	events := make([]map[string]interface{}, 0, len(span.Events))
	for _, event := range span.Events {
		eventMap := map[string]interface{}{
			"name":      event.Name,
			"timestamp": event.TimeUnixNano,
		}

		eventAttrs := make(map[string]interface{})
		for _, attr := range event.Attributes {
			eventAttrs[attr.Key] = convertAnyValue(attr.Value)
		}

		if len(eventAttrs) > 0 {
			eventMap["attributes"] = eventAttrs
		}

		events = append(events, eventMap)
	}

	// Convert links
	links := make([]map[string]interface{}, 0, len(span.Links))
	for _, link := range span.Links {
		linkMap := map[string]interface{}{
			"trace_id": fmt.Sprintf("%x", link.TraceId),
			"span_id":  fmt.Sprintf("%x", link.SpanId),
		}

		linkAttrs := make(map[string]interface{})
		for _, attr := range link.Attributes {
			linkAttrs[attr.Key] = convertAnyValue(attr.Value)
		}

		if len(linkAttrs) > 0 {
			linkMap["attributes"] = linkAttrs
		}

		links = append(links, linkMap)
	}

	return &types.SpanResponse{
		TraceID:   traceID,
		ID:        spanID,
		ParentID:  parentSpanID,
		Name:      span.Name,
		Kind:      span.Kind.String(),
		Timestamp: startTime / 1000,
		Duration:  duration / 1000,
		LocalEndpoint: types.SpanEndpoint{
			ServiceName: attributes["service.name"].(string),
		},
		Tags: convertAttrributes(attributes),
	}
}
func convertAttrributes(attributes map[string]interface{}) map[string]string {
	// Convert attributes to map
	attributesMap := make(map[string]string)
	for key, value := range attributes {
		attributesMap[key] = fmt.Sprintf("%v", value)
	}
	return attributesMap
}

// Convert AnyValue to Go type
func convertAnyValue(value *v1.AnyValue) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.Value.(type) {
	case *v1.AnyValue_StringValue:
		return v.StringValue
	case *v1.AnyValue_BoolValue:
		return v.BoolValue
	case *v1.AnyValue_IntValue:
		return v.IntValue
	case *v1.AnyValue_DoubleValue:
		return v.DoubleValue
	case *v1.AnyValue_ArrayValue:
		result := make([]interface{}, 0, len(v.ArrayValue.Values))
		for _, val := range v.ArrayValue.Values {
			result = append(result, convertAnyValue(val))
		}
		return result
	case *v1.AnyValue_KvlistValue:
		result := make(map[string]interface{})
		for _, kv := range v.KvlistValue.Values {
			result[kv.Key] = convertAnyValue(kv.Value)
		}
		return result
	default:
		return nil
	}
}

func (s *Service) StartProcessTrace(nc *nats.Conn) {
	flag.Parse()

	// Create trace store
	store := NewTraceStore()

	// Subscribe to NATS subject
	sub, err := nc.Subscribe(*natsSubj, func(msg *nats.Msg) {
		// Unmarshal protobuf message
		var tracesData tracepb.TracesData
		if err := proto.Unmarshal(msg.Data, &tracesData); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}

		// Process spans
		for _, rs := range tracesData.ResourceSpans {
			// Extract service name from resource
			serviceName := "unknown"
			for _, attr := range rs.Resource.Attributes {
				if attr.Key == "service.name" {
					if sv := attr.Value.GetStringValue(); sv != "" {
						serviceName = sv
					}
					break
				}
			}
			log.Println("Service name:", serviceName)

			// Process spans in each scope
			for _, ss := range rs.ScopeSpans {
				for _, span := range ss.Spans {
					span.Attributes = append(span.Attributes,
						&v1.KeyValue{
							Key: "service.name",
							Value: &v1.AnyValue{
								Value: &v1.AnyValue_StringValue{
									StringValue: serviceName,
								},
							},
						},
					)
					store.AddSpan(span)
				}
			}
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS: %v", err)
	}
	defer sub.Unsubscribe()

	// Start background processor
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(*bufferTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Get expired traces
				expiredTraces := store.GetExpiredTraces(*bufferTime)

				// Process each trace
				for traceID, _spans := range expiredTraces {
					log.Printf("Processing trace %s with %d spans", traceID, len(_spans))
					spans := make([]*types.SpanResponse, 0, len(_spans))
					for _, _span := range _spans {
						span := convertSpanToSpanResponse(_span)
						spans = append(spans, span)
					}
					if err := s.ProcessTrace(ctx, spans); err != nil {
						log.Printf("Failed to process trace %s: %v", traceID, err)
					}

				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// HTTP server for health checks and metrics
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		store.mu.RLock()
		traceCount := len(store.traces)
		store.mu.RUnlock()

		fmt.Fprintf(w, "# HELP traces_count Number of traces in memory\n")
		fmt.Fprintf(w, "# TYPE traces_count gauge\n")
		fmt.Fprintf(w, "traces_count %d\n", traceCount)
	})

	// Start HTTP server
	server := &http.Server{
		Addr: *httpAddr,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for signal to shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	log.Println("Shutting down...")

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Process remaining traces
	expiredTraces := store.GetExpiredTraces(0)

	// Process each trace
	for traceID, spans := range expiredTraces {
		log.Printf("Processing remaining trace %s with %d spans", traceID, len(spans))
	}

	log.Println("Shutdown complete")
}
