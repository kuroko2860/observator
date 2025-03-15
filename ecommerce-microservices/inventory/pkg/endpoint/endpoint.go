package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/inventory/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// UpdateInventoryRequest collects the request parameters for the UpdateInventory method
type UpdateInventoryRequest struct {
	OrderID string   `json:"order_id"`
	Items   []string `json:"items"`
}

// UpdateInventoryResponse collects the response parameters for the UpdateInventory method
type UpdateInventoryResponse struct {
	Err string `json:"error,omitempty"`
}

// VerifyInventoryRequest collects the request parameters for the VerifyInventory method
type VerifyInventoryRequest struct {
	Items []string `json:"items"`
}

// VerifyInventoryResponse collects the response parameters for the VerifyInventory method
type VerifyInventoryResponse struct {
	Available bool   `json:"available"`
	Err       string `json:"error,omitempty"`
}

// MakeUpdateInventoryEndpoint returns an endpoint that invokes UpdateInventory on the service
func MakeUpdateInventoryEndpoint(s service.InventoryService, logger log.Logger) endpoint.Endpoint {
	baseEndpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateInventoryRequest)
		err := s.UpdateInventory(ctx, req.OrderID, req.Items)
		if err != nil {
			return UpdateInventoryResponse{Err: err.Error()}, nil
		}
		return UpdateInventoryResponse{}, nil
	}

	// Add logging middleware
	loggingMiddl := logging.EndpointLoggingMiddleware(logger)(baseEndpoint)

	// Add tracing middleware
	tracingMiddl := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("inventory-endpoint")
		ctx, span := tracer.Start(ctx, "UpdateInventoryEndpoint")
		defer span.End()

		// Extract trace information
		spanContext := trace.SpanContextFromContext(ctx)
		if spanContext.IsValid() {
			logger.Log("trace_id", spanContext.TraceID().String())
		}

		return loggingMiddl(ctx, request)
	}

	return tracingMiddl
}

// MakeVerifyInventoryEndpoint returns an endpoint that invokes VerifyInventory on the service
func MakeVerifyInventoryEndpoint(s service.InventoryService, logger log.Logger) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyInventoryRequest)
		available, err := s.VerifyInventory(ctx, req.Items)
		if err != nil {
			return VerifyInventoryResponse{Err: err.Error()}, nil
		}
		return VerifyInventoryResponse{Available: available}, nil
	}

	// Add logging middleware
	endpoint = logging.EndpointLoggingMiddleware(logger)(endpoint)

	// Add tracing middleware
	endpoint = func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("inventory-endpoint")
		ctx, span := tracer.Start(ctx, "VerifyInventoryEndpoint")
		defer span.End()

		// Extract trace information
		spanContext := trace.SpanContextFromContext(ctx)
		if spanContext.IsValid() {
			logger.Log("trace_id", spanContext.TraceID().String())
		}

		return endpoint(ctx, request)
	}

	return endpoint
}

// Endpoints collects all of the endpoints that compose the inventory service
type Endpoints struct {
	UpdateInventory endpoint.Endpoint
	VerifyInventory endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the corresponding method on the provided service
func MakeEndpoints(s service.InventoryService, logger log.Logger) Endpoints {
	return Endpoints{
		UpdateInventory: MakeUpdateInventoryEndpoint(s, logger),
		VerifyInventory: MakeVerifyInventoryEndpoint(s, logger),
	}
}
