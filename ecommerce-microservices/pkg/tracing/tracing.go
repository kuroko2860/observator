package tracing

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes an OTLP exporter, and configures the corresponding trace provider
func InitTracer(serviceName, zipkinURL string) (func(context.Context) error, error) {
	// Create Zipkin exporter
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		return nil, err
	}

	// Configure the trace provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return a function to shutdown the tracer provider
	return tp.Shutdown, nil
}

// Tracer returns a named tracer
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// NewTracedHTTPClient returns an HTTP client that automatically propagates trace context
func NewTracedHTTPClient() *http.Client {
	return &http.Client{
		Transport: &TracingTransport{
			Base: http.DefaultTransport,
		},
		Timeout: 5 * time.Second,
	}
}

// TracingTransport is an http.RoundTripper that automatically injects trace context
type TracingTransport struct {
	Base http.RoundTripper
}

// RoundTrip implements http.RoundTripper
func (t *TracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Inject trace context into the outgoing request
	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))

	// Call the base transport
	return t.Base.RoundTrip(req)
}
