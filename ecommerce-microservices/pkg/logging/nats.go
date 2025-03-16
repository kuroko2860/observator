package logging

import (
	"encoding/json"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

var (
	natsConn *nats.Conn
	natsMu   sync.Mutex
)

// InitNATS initializes the connection to NATS server
func InitNATS(natsURL string) error {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn != nil && natsConn.IsConnected() {
		return nil
	}

	var err error
	natsConn, err = nats.Connect(natsURL,
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Error().Err(err).Msg("Disconnected from NATS")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info().Msg("Reconnected to NATS")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Info().Msg("NATS connection closed")
		}),
	)

	if err != nil {
		return err
	}

	log.Info().Str("url", natsURL).Msg("Connected to NATS")
	return nil
}

// CloseNATS closes the connection to NATS server
func CloseNATS() {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn != nil {
		natsConn.Close()
		natsConn = nil
		log.Info().Msg("NATS connection closed")
	}
}

// HttpLogEntry represents a log entry for HTTP requests
type HttpLogEntry struct {
	ServiceName   string `json:"service_name" bson:"service_name"`
	URIPath       string `json:"uri_path" bson:"uri_path"`
	Referer       string `json:"referer" bson:"referer"`
	UserId        string `json:"user_id" bson:"user_id"`
	Method        string `json:"method" bson:"method"`
	StartTime     int64  `json:"start_time" bson:"start_time"`
	StartTimeDate string `json:"start_time_date" bson:"start_time_date"`
	Host          string `json:"host" bson:"host"`
	Protocol      string `json:"protocol" bson:"protocol"`
	RemoteIP      string `json:"remote_ip" bson:"remote_ip"`
	RequestId     string `json:"request_id" bson:"request_id"`
	TraceId       string `json:"trace_id" bson:"trace_id"`
	SpanId        string `json:"span_id" bson:"span_id"`
	UserAgent     string `json:"user_agent" bson:"user_agent"`
	Duration      int64  `json:"duration" bson:"duration"`
	ResquestSize  string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize  int64  `json:"response_size" bson:"response_size"`
	StatusCode    int    `json:"status_code" bson:"status_code"`
	ErrorMessage  string `json:"error_message,omitempty" bson:"error_message,omitempty"`
}

// PublishLogEntry publishes a log entry to NATS
func PublishLogEntry(entry HttpLogEntry) {
	natsMu.Lock()
	defer natsMu.Unlock()

	if natsConn == nil || !natsConn.IsConnected() {
		log.Warn().Msg("Cannot publish log entry: NATS not connected")
		return
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log entry")
		return
	}

	err = natsConn.Publish("logs", entryJSON)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish log to NATS")
	}
}