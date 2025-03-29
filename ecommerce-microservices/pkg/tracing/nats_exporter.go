package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	cv1 "go.opentelemetry.io/proto/otlp/common/v1"
	rv1 "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
)

// NatsExporter exports spans to NATS
type NatsExporter struct {
	conn    *nats.Conn
	subject string
}

// NewNatsExporter creates a new NATS exporter
func NewNatsExporter(url, subject string) (*NatsExporter, error) {
	// Connect to NATS server with retry options
	opts := []nats.Option{
		nats.Name("OpenTelemetry Exporter"),
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(10),
	}

	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &NatsExporter{
		conn:    conn,
		subject: subject,
	}, nil
}

// ExportSpans exports spans to NATS
func (e *NatsExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	// Convert spans to OTLP format
	otlpSpans := convertToOTLP(spans)
	fmt.Println("Exporting spans to NATS:", len(spans))

	// Serialize to protobuf
	data, err := proto.Marshal(otlpSpans)
	if err != nil {
		return fmt.Errorf("failed to marshal spans: %w", err)
	}

	// Publish to NATS
	err = e.conn.Publish(e.subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	return nil
}

// Shutdown closes the NATS connection
func (e *NatsExporter) Shutdown(ctx context.Context) error {
	e.conn.Close()
	return nil
}

// convertToOTLP converts spans to OTLP format
func convertToOTLP(spans []trace.ReadOnlySpan) *tracepb.TracesData {
	resourceSpans := make([]*tracepb.ResourceSpans, 0)

	// Group spans by resource
	resourceMap := make(map[string][]*tracepb.Span)

	for _, span := range spans {
		// Convert span to OTLP format
		traceId := span.SpanContext().TraceID()
		spanId := span.SpanContext().SpanID()
		parentId := span.Parent().SpanID()

		otlpSpan := &tracepb.Span{
			TraceId:           traceId[:],
			SpanId:            spanId[:],
			Name:              span.Name(),
			Kind:              tracepb.Span_SpanKind(span.SpanKind()),
			StartTimeUnixNano: uint64(span.StartTime().UnixNano()),
			EndTimeUnixNano:   uint64(span.EndTime().UnixNano()),
			Status: &tracepb.Status{
				Code:    tracepb.Status_StatusCode(span.Status().Code),
				Message: span.Status().Description,
			},
		}

		// Add parent span ID if available
		if span.Parent().SpanID().IsValid() {
			otlpSpan.ParentSpanId = parentId[:]
		}

		// Add attributes
		attributes := make([]*cv1.KeyValue, 0)
		for _, k := range span.Attributes() {
			attributes = append(attributes, &cv1.KeyValue{
				Key:   string(k.Key),
				Value: attributeValueToOTLP(k.Value),
			})
		}
		otlpSpan.Attributes = attributes

		// Get resource key for grouping - IMPORTANT CHANGE HERE
		// Extract service name from resource attributes, not span attributes
		serviceName := "unknown"
		for _, attr := range span.Resource().Attributes() {
			if attr.Key == "service.name" {
				serviceName = attr.Value.AsString()
				break
			}
		}

		resourceMap[serviceName] = append(resourceMap[serviceName], otlpSpan)
	}

	// Create resource spans
	for resourceKey, otlpSpans := range resourceMap {
		scopeSpans := []*tracepb.ScopeSpans{
			{
				Spans: otlpSpans,
			},
		}

		resourceSpans = append(resourceSpans, &tracepb.ResourceSpans{
			Resource: &rv1.Resource{
				Attributes: []*cv1.KeyValue{
					{
						Key:   "service.name",
						Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: resourceKey}},
					},
				},
			},
			ScopeSpans: scopeSpans,
		})
	}

	return &tracepb.TracesData{
		ResourceSpans: resourceSpans,
	}
}

// attributeValueToOTLP converts an attribute value to OTLP format
func attributeValueToOTLP(v attribute.Value) *cv1.AnyValue {
	switch v.Type() {
	case attribute.STRING:
		return &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: v.AsString()}}
	case attribute.BOOL:
		return &cv1.AnyValue{Value: &cv1.AnyValue_BoolValue{BoolValue: v.AsBool()}}
	case attribute.INT64:
		return &cv1.AnyValue{Value: &cv1.AnyValue_IntValue{IntValue: v.AsInt64()}}
	case attribute.FLOAT64:
		return &cv1.AnyValue{Value: &cv1.AnyValue_DoubleValue{DoubleValue: v.AsFloat64()}}
	default:
		return &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: v.AsString()}}
	}
}
