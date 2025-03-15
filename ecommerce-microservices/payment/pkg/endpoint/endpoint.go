package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/trace"

	"kltn/ecommerce-microservices/payment/pkg/service"
	"kltn/ecommerce-microservices/pkg/logging"
	"kltn/ecommerce-microservices/pkg/tracing"
)

// CalculateMoneyRequest collects the request parameters for the CalculateMoney method
type CalculateMoneyRequest struct {
	OrderID string   `json:"order_id"`
	Items   []string `json:"items"`
}

// CalculateMoneyResponse collects the response parameters for the CalculateMoney method
type CalculateMoneyResponse struct {
	Amount float64 `json:"amount"`
	Err    string  `json:"error,omitempty"`
}

// ApplyCouponRequest collects the request parameters for the ApplyCoupon method
type ApplyCouponRequest struct {
	OrderID    string  `json:"order_id"`
	CouponCode string  `json:"coupon_code"`
	Amount     float64 `json:"amount"`
}

// ApplyCouponResponse collects the response parameters for the ApplyCoupon method
type ApplyCouponResponse struct {
	DiscountedAmount float64 `json:"discounted_amount"`
	Err              string  `json:"error,omitempty"`
}

// MakeCalculateMoneyEndpoint returns an endpoint that invokes CalculateMoney on the service
func MakeCalculateMoneyEndpoint(s service.PaymentService, logger log.Logger) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CalculateMoneyRequest)
		amount, err := s.CalculateMoney(ctx, req.OrderID, req.Items)
		if err != nil {
			return CalculateMoneyResponse{Err: err.Error()}, nil
		}
		return CalculateMoneyResponse{Amount: amount}, nil
	}

	// Add logging middleware
	loggingMiddl := logging.EndpointLoggingMiddleware(logger)(endpoint)

	// Add tracing middleware
	tracingMiddl := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("payment-endpoint")
		ctx, span := tracer.Start(ctx, "CalculateMoneyEndpoint")
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

// MakeApplyCouponEndpoint returns an endpoint that invokes ApplyCoupon on the service
func MakeApplyCouponEndpoint(s service.PaymentService, logger log.Logger) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ApplyCouponRequest)
		discountedAmount, err := s.ApplyCoupon(ctx, req.OrderID, req.CouponCode, req.Amount)
		if err != nil {
			return ApplyCouponResponse{Err: err.Error()}, nil
		}
		return ApplyCouponResponse{DiscountedAmount: discountedAmount}, nil
	}

	// Add logging middleware
	logMiddl := logging.EndpointLoggingMiddleware(logger)(endpoint)

	// Add tracing middleware
	tracingMiddl := func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := tracing.Tracer("payment-endpoint")
		ctx, span := tracer.Start(ctx, "ApplyCouponEndpoint")
		defer span.End()

		// Extract trace information
		spanContext := trace.SpanContextFromContext(ctx)
		if spanContext.IsValid() {
			logger.Log("trace_id", spanContext.TraceID().String())
		}

		return logMiddl(ctx, request)
	}

	return tracingMiddl
}

// Endpoints collects all of the endpoints that compose the payment service
type Endpoints struct {
	CalculateMoney endpoint.Endpoint
	ApplyCoupon    endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the corresponding method on the provided service
func MakeEndpoints(s service.PaymentService, logger log.Logger) Endpoints {
	return Endpoints{
		CalculateMoney: MakeCalculateMoneyEndpoint(s, logger),
		ApplyCoupon:    MakeApplyCouponEndpoint(s, logger),
	}
}
