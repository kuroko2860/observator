package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"kltn/ecommerce-microservices/address/pkg/endpoint"
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

	// GET /address/{user_id}
	r.Methods("GET").Path("/address/{user_id}").Handler(otelhttp.NewHandler(
		httptransport.NewServer(
			endpoints.GetAddress,
			decodeGetAddressRequest,
			encodeResponse,
			options...,
		),
		"GET /address/{user_id}",
	))

	return r
}

func decodeGetAddressRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		return nil, http.ErrBodyNotAllowed
	}
	return endpoint.GetAddressRequest{UserID: userID}, nil
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
