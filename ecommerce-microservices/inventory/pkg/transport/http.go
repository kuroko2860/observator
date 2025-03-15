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

	"kltn/ecommerce-microservices/inventory/pkg/endpoint"
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

	// POST /update-inventory
	r.Methods("POST").Path("/update-inventory").Handler(otelhttp.NewHandler(
		httptransport.NewServer(
			endpoints.UpdateInventory,
			decodeUpdateInventoryRequest,
			encodeResponse,
			options...,
		),
		"POST /update-inventory",
	))

	// POST /verify-inventory
	r.Methods("POST").Path("/verify-inventory").Handler(otelhttp.NewHandler(
		httptransport.NewServer(
			endpoints.VerifyInventory,
			decodeVerifyInventoryRequest,
			encodeResponse,
			options...,
		),
		"POST /verify-inventory",
	))

	return r
}

func decodeUpdateInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.UpdateInventoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeVerifyInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.VerifyInventoryRequest
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