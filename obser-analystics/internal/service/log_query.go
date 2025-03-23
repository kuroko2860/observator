package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

// QueryLogsByTraceID retrieves logs with the specified trace ID from Elasticsearch
func (s *Service) QueryLogsByTraceID(ctx context.Context, traceID string, from, size int) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"trace_id": traceID,
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{
				"start_time": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return s.executeElasticsearchQuery(ctx, query)
}

// QueryLogsBySpanID retrieves logs with the specified span ID from Elasticsearch
func (s *Service) QueryLogsBySpanID(ctx context.Context, spanID string, from, size int) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"span_id": spanID,
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{
				"start_time": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return s.executeElasticsearchQuery(ctx, query)
}

// QueryLogsByTraceAndSpanID retrieves logs with both the given trace ID and span ID
func (s *Service) QueryLogsByTraceAndSpanID(ctx context.Context, traceID, spanID string, from, size int) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"trace_id": traceID,
						},
					},
					{
						"match": map[string]interface{}{
							"span_id": spanID,
						},
					},
				},
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{
				"start_time": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return s.executeElasticsearchQuery(ctx, query)
}

// QueryLogsByService retrieves logs from the specified service
func (s *Service) QueryLogsByService(ctx context.Context, serviceName string, from, size int) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"service_name": serviceName,
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{
				"start_time": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return s.executeElasticsearchQuery(ctx, query)
}

// QueryLogsByTimeRange retrieves logs within the specified time range
func (s *Service) QueryLogsByTimeRange(ctx context.Context, startTime, endTime int64, from, size int) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"start_time": map[string]interface{}{
					"gte": startTime,
					"lte": endTime,
				},
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{
				"start_time": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return s.executeElasticsearchQuery(ctx, query)
}

// executeElasticsearchQuery executes a query against Elasticsearch
func (s *Service) executeElasticsearchQuery(ctx context.Context, query map[string]interface{}) ([]map[string]interface{}, error) {
	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	// Create Elasticsearch request
	req := esapi.SearchRequest{
		Index: []string{"microservices-logs-*"},
		Body:  strings.NewReader(string(queryJSON)),
	}

	// Execute the request
	res, err := req.Do(ctx, esClient)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing elasticsearch error response: %w", err)
		}
		return nil, fmt.Errorf("elasticsearch error: %v", e)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing elasticsearch response: %w", err)
	}

	// Extract the hits
	hits, found := result["hits"].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("hits not found in elasticsearch response")
	}

	hitsArray, found := hits["hits"].([]interface{})
	if !found {
		return nil, fmt.Errorf("hits array not found in elasticsearch response")
	}

	// Extract the documents
	var logs []map[string]interface{}
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		logs = append(logs, source)
	}

	return logs, nil
}

// FindHttpLogEntriesByTraceId retrieves all log entries with the given trace ID from MongoDB
func (s *Service) FindHttpLogEntriesByTraceId(ctx context.Context, traceId string) ([]model.HttpLogEntry, error) {
	var logs []model.HttpLogEntry
	err := httpLogEntryCollection.Find(ctx, bson.M{"trace_id": traceId}).All(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// FindHttpLogEntriesBySpanId retrieves all log entries with the given span ID from MongoDB
func (s *Service) FindHttpLogEntriesBySpanId(ctx context.Context, spanId string) ([]model.HttpLogEntry, error) {
	var logs []model.HttpLogEntry
	err := httpLogEntryCollection.Find(ctx, bson.M{"span_id": spanId}).All(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// FindHttpLogEntriesByTraceAndSpanId retrieves all log entries with both the given trace ID and span ID from MongoDB
func (s *Service) FindHttpLogEntriesByTraceAndSpanId(ctx context.Context, traceId, spanId string) ([]model.HttpLogEntry, error) {
	var logs []model.HttpLogEntry
	err := httpLogEntryCollection.Find(ctx, bson.M{
		"trace_id": traceId,
		"span_id":  spanId,
	}).All(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
