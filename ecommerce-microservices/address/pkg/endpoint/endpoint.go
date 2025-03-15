package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/address/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// GetAddressRequest collects the request parameters for the GetAddress method
type GetAddressRequest struct {
	UserID string `json:"user_id"`
}

// GetAddressResponse collects the response parameters for the GetAddress method
type GetAddressResponse struct {
	Address service.Address `json:"address,omitempty"`
	Err     string          `json:"error,omitempty"`
}

// MakeGetAddressEndpoint returns an endpoint that invokes GetAddress on the service
func MakeGetAddressEndpoint(s service.AddressService, logger log.Logger) endpoint.Endpoint {
	baseEndpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAddressRequest)
		address, err := s.GetAddress(ctx, req.UserID)
		if err != nil {
			return GetAddressResponse{Err: err.Error()}, nil
		}
		return GetAddressResponse{Address: address}, nil
	}

	// Add logging middleware
	loggingMiddleware := logging.EndpointLoggingMiddleware(logger)(baseEndpoint)

	// Add tracing middleware
	tracingMiddleware := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("address-endpoint")
		ctx, span := tracer.Start(ctx, "GetAddressEndpoint")
		defer span.End()

		// Extract trace information
		spanContext := trace.SpanContextFromContext(ctx)
		if spanContext.IsValid() {
			logger.Log("trace_id", spanContext.TraceID().String())
		}

		return loggingMiddleware(ctx, request)
	}

	return tracingMiddleware
}

// Endpoints collects all of the endpoints that compose the address service
type Endpoints struct {
	GetAddress endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the corresponding method on the provided service
func MakeEndpoints(s service.AddressService, logger log.Logger) Endpoints {
	return Endpoints{
		GetAddress: MakeGetAddressEndpoint(s, logger),
	}
}
