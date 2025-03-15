package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/order/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// CreateOrderRequest collects the request parameters for the CreateOrder method
type CreateOrderRequest struct {
	UserID string   `json:"user_id"`
	Items  []string `json:"items"`
}

// CreateOrderResponse collects the response parameters for the CreateOrder method
type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Err     string `json:"error,omitempty"`
}

// MakeCreateOrderEndpoint returns an endpoint that invokes CreateOrder on the service
func MakeCreateOrderEndpoint(s service.OrderService, logger log.Logger) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrderRequest)
		orderID, err := s.CreateOrder(ctx, req.UserID, req.Items)
		if err != nil {
			return CreateOrderResponse{Err: err.Error()}, nil
		}
		return CreateOrderResponse{OrderID: orderID}, nil
	}

	// Add logging middleware
	loggingMiddl := logging.EndpointLoggingMiddleware(logger)(endpoint)

	// Add tracing middleware
	tracingMiddl := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("order-endpoint")
		ctx, span := tracer.Start(ctx, "CreateOrderEndpoint")
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

// Endpoints collects all of the endpoints that compose the order service
type Endpoints struct {
	CreateOrder endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the corresponding method on the provided service
func MakeEndpoints(s service.OrderService, logger log.Logger) Endpoints {
	return Endpoints{
		CreateOrder: MakeCreateOrderEndpoint(s, logger),
	}
}
