package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"kltn/ecommerce-microservices/payment/pkg/endpoint"
	"kltn/ecommerce-microservices/pkg/logging"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on predefined paths
func NewHTTPHandler(endpoints endpoint.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerFinalizer(logging.ServerFinalizer(logger)),
	}

	// POST /calculate-money
	r.Methods("POST").Path("/calculate-money").Handler(otelhttp.NewHandler(
		httptransport.NewServer(
			endpoints.CalculateMoney,
			decodeCalculateMoneyRequest,
			encodeResponse,
			options...,
		),
		"POST /calculate-money",
	))

	// POST /apply-coupon
	r.Methods("POST").Path("/apply-coupon").Handler(otelhttp.NewHandler(
		httptransport.NewServer(
			endpoints.ApplyCoupon,
			decodeApplyCouponRequest,
			encodeResponse,
			options...,
		),
		"POST /apply-coupon",
	))

	return r
}

func decodeCalculateMoneyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CalculateMoneyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeApplyCouponRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.ApplyCouponRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}