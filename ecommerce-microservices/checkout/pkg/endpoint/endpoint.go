package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/checkout/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// UserCheckoutRequest collects the request parameters for the UserCheckout method
type UserCheckoutRequest struct {
	UserID string   `json:"user_id"`
	Items  []string `json:"items"`
}

// UserCheckoutResponse collects the response parameters for the UserCheckout method
type UserCheckoutResponse struct {
	OrderID string `json:"order_id"`
	Err     string `json:"error,omitempty"`
}

// MakeUserCheckoutEndpoint returns an endpoint that invokes UserCheckout on the service
func MakeUserCheckoutEndpoint(s service.CheckoutService, logger log.Logger) endpoint.Endpoint {
	baseEndpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserCheckoutRequest)
		orderID, err := s.UserCheckout(ctx, req.UserID, req.Items)
		if err != nil {
			return UserCheckoutResponse{Err: err.Error()}, nil
		}
		return UserCheckoutResponse{OrderID: orderID}, nil
	}

	// Add logging middleware
	loggingMiddleware := logging.EndpointLoggingMiddleware(logger)(baseEndpoint)

	// Add tracing middleware
	tracingMiddleware := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("checkout-endpoint")
		ctx, span := tracer.Start(ctx, "UserCheckoutEndpoint")
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

// Endpoints collects all of the endpoints that compose the checkout service
type Endpoints struct {
	UserCheckout endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the corresponding method on the provided service
func MakeEndpoints(s service.CheckoutService, logger log.Logger) Endpoints {
	return Endpoints{
		UserCheckout: MakeUserCheckoutEndpoint(s, logger),
	}
}
