package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/config"
	"kuroko.com/processor/internal/types"
)

var pathIds map[uint32]bool

func (s *Service) FetchTraces(ctx context.Context) error {
	fmt.Println("Start Fetching traces...")
	lookback := config.LOOKBACK
	limit := config.LIMIT
	endtime := s.FindEndTime(ctx)
	locktime := s.FindLockTime(ctx)

	if endtime != locktime {
		return nil
	}

	end := time.Now().Unix()
	// s.UpdateLockTime(ctx, end)
	s.FetchTracesFromTimeRange(ctx, endtime, end, int64(lookback), int(limit))
	// s.UpdateEndTime(ctx, end)
	return nil
}

func (s *Service) FetchTracesFromTimeRange(ctx context.Context, endOld, endNew, lookback int64, limit int) error {
	if endOld >= endNew {
		return errors.New("from >= to")
	}
	LOOK_BACK_IN_MILIES := lookback * 1000
	END_OLD_IN_MILIES := endOld * 1000
	END_NEW_IN_MILIES := endNew * 1000

	anchor := END_OLD_IN_MILIES
	for {
		from := anchor
		anchor += LOOK_BACK_IN_MILIES
		if anchor > END_NEW_IN_MILIES {
			anchor = END_NEW_IN_MILIES
		}

		to := anchor
		traces, err := s.FetchTracesFromTo(ctx, from, to, limit)
		if err != nil {
			fmt.Println(err)
		}
		for _, trace := range traces {
			s.ProcessTrace(ctx, trace)
		}
		if anchor == END_NEW_IN_MILIES {
			return nil
		}
		time.Sleep(1 * time.Microsecond)
	}
}

func (s *Service) ProcessTrace(ctx context.Context, trace []*types.SpanResponse) error {
	root, errPath := s.ConvertTraceToGraph(ctx, trace)
	pathId, _ := s.CaculatePathId(ctx, root, 0)
	s.ProcessGraph(ctx, root, errPath)
	var err error
	if !s.IsPathExist(ctx, pathId) {
		session := s.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: config.Neo4jDatabase})
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			for _, span := range trace {
				// Create or merge the node for the current service and operation
				query := `
				MERGE (s:Service {name: $service})
				MERGE (o:Operation {name: $operation, service: $service})
				MERGE (s)-[:PERFORMS]->(o)
			`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"service":   span.LocalEndpoint,
					"operation": span.Name,
				})
				if err != nil {
					return nil, err
				}

				// If the span has a parent, create a dependency edge
				if span.ParentId != "" {
					parentSpan := findParentSpan(trace, span.ParentId)
					if parentSpan != nil {
						query = `
						MATCH (p:Operation {name: $parentOperation, service: $parentService})
						MATCH (c:Operation {name: $childOperation, service: $childService})
						MERGE (p)-[:CALLS {pathId: $pathId} ]->(c)
					`
						_, err = tx.Run(ctx, query, map[string]interface{}{
							"parentOperation": parentSpan.Name,
							"parentService":   parentSpan.LocalEndpoint,
							"childOperation":  span.Name,
							"childService":    span.LocalEndpoint,
							"pathId":          pathId,
						})
						if err != nil {
							return nil, err
						}
					}
				}
			}
			return nil, nil
		})
		pathIds[pathId] = true
		pathIdCollection.InsertOne(ctx, bson.M{"_id": pathId})
	} else {
		fmt.Printf("path %d already exists\n", pathId)
	}

	for _, sr := range trace {
		span := convertSrToSpan(sr)
		span.PathId = pathId
		spanCollection.InsertOne(ctx, span)
	}

	return err
}

func (s *Service) IsPathExist(ctx context.Context, pathId uint32) bool {
	return pathIds[pathId]
}
func convertSrToSpan(sr *types.SpanResponse) *types.Span {
	var span types.Span
	span.Id = sr.Id
	span.TraceId = sr.TraceId
	span.ServiceName = sr.LocalEndpoint
	span.OperationName = sr.Name
	span.Timestamp = sr.Timestamp
	span.Duration = sr.Duration
	span.Error = sr.Tags["error"]
	return &span
}
func findParentSpan(spans []*types.SpanResponse, parentID string) *types.SpanResponse {
	for _, span := range spans {
		if span.Id == parentID {
			return span
		}
	}
	return nil
}

func (s *Service) FetchTracesFromTo(ctx context.Context, from int64, to int64, limit int) ([][]*types.SpanResponse, error) {
	if from >= to {
		return nil, errors.New("from >= to")
	}

	end := to
	lookback := to - from
	url := fmt.Sprintf("http://%s:%d/api/v2/traces?lookback=%d&endTs=%d&limit=%d", config.TRACE_HOST, config.TRACE_PORT, lookback, end, limit)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var traces [][]*types.SpanResponse
	if err := json.Unmarshal(resData, &traces); err != nil {
		return nil, err
	}
	return traces, nil
}
