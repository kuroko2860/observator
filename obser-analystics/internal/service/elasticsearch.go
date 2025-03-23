package service

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

var esClient *elasticsearch.Client

// InitElasticsearch initializes the Elasticsearch client
func (s *Service) InitElasticsearch() error {
	// Get Elasticsearch URL from environment variable or use default
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	// Create Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Test the connection
	res, err := esClient.Info()
	if err != nil {
		return fmt.Errorf("error connecting to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	// Check if the response contains an error
	if res.IsError() {
		return fmt.Errorf("Elasticsearch returned an error: %s", res.String())
	}

	return nil
}
