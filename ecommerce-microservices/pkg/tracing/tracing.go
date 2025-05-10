package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// InitTracer creates a new trace provider instance and registers it as global tracer provider.
func InitTracer(serviceName, natsURL string) (func(context.Context) error, error) {
	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create NATS exporter
	exporter, err := NewNatsExporter(natsURL, "traces.service")
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS exporter: %w", err)
	}

	// Create trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(time.Second),
			sdktrace.WithMaxExportBatchSize(100),
		),
	)

	// Set global trace provider
	otel.SetTracerProvider(traceProvider)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return shutdown function
	return func(ctx context.Context) error {
		return traceProvider.Shutdown(ctx)
	}, nil
}
