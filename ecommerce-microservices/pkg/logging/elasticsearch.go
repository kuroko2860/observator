package logging

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
)

var (
	esClient *elastic.Client
	esMu     sync.Mutex
)

// InitElasticsearch initializes the connection to Elasticsearch
func InitElasticsearch(esURL string) error {
	esMu.Lock()
	defer esMu.Unlock()

	if esClient != nil {
		return nil
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return err
	}

	// Check if the connection is successful
	_, _, err = client.Ping(esURL).Do(context.Background())
	if err != nil {
		return err
	}

	esClient = client
	log.Info().Str("url", esURL).Msg("Connected to Elasticsearch")
	return nil
}

// CloseElasticsearch closes the connection to Elasticsearch
func CloseElasticsearch() {
	esMu.Lock()
	defer esMu.Unlock()

	esClient = nil
	log.Info().Msg("Elasticsearch connection closed")
}

// SetupNATSToElasticsearchBridge sets up a subscriber that forwards logs from NATS to Elasticsearch
func SetupNATSToElasticsearchBridge(natsConn *nats.Conn, indexPrefix string) error {
	if natsConn == nil || !natsConn.IsConnected() {
		return ErrNATSNotConnected
	}

	if esClient == nil {
		return ErrElasticsearchNotConnected
	}

	// Subscribe to the logs subject
	_, err := natsConn.Subscribe("logs", func(msg *nats.Msg) {
		var entry HttpLogEntry
		if err := json.Unmarshal(msg.Data, &entry); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal log entry")
			return
		}

		// Create index name with date suffix
		indexName := indexPrefix + "-" + time.Now().Format("2006.01.02")

		// Index the document
		_, err := esClient.Index().
			Index(indexName).
			BodyJson(entry).
			Do(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("Failed to index log entry")
		}
	})

	if err != nil {
		return err
	}

	log.Info().Msg("NATS to Elasticsearch bridge established")
	return nil
}

// Errors
var (
	ErrNATSNotConnected          = nats.ErrConnectionClosed
	ErrElasticsearchNotConnected = elastic.ErrNoClient
)