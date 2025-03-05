package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiniu/qmgo"
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
	root, err := s.ConvertTraceToGraph(ctx, trace)
	if err != nil {
		fmt.Printf("Error when converting trace: %s\n", err.Error())
		return err
	}
	pathId, _ := s.CaculatePathId(ctx, root, 0)
	if !s.IsPathExist(ctx, pathId) {
		s.InsertEntityFromGraph(ctx, root, pathId)
		s.InsertPath(ctx, root, pathId)

		pathIds[pathId] = true
		pathIdCollection.InsertOne(ctx, bson.M{"_id": pathId})
	} else {
		fmt.Printf("path %d already exists\n", pathId)
	}
	s.ProcessGraph(ctx, root, pathId)

	for _, sr := range trace {
		span := convertSrToSpan(sr)
		span.PathID = pathId
		spanCollection.InsertOne(ctx, span)
	}

	return err
}

func (s *Service) InsertPath(ctx context.Context, root *types.GraphNode, pathId uint32) {

}

// insert operation, hop
func (s *Service) InsertEntityFromGraph(ctx context.Context, root *types.GraphNode, pathId uint32) {
	if root == nil {
		return
	}
	opID := generateOperationID(root.Span)
	var op types.Operation
	err := operationCollection.Find(ctx, bson.M{"_id": opID}).One(&op)
	if err != nil {
		if qmgo.IsErrNoDocuments(err) {
			op = types.Operation{
				ID:      opID,
				Name:    root.Span.Name,
				Service: root.Span.LocalEndpoint,
			}
			operationCollection.InsertOne(ctx, &op)
		}
	}

	for _, child := range root.Children {
		hopID := generateHopID(root, child, pathId)
		var hop types.Hop
		err := hopCollection.Find(ctx, bson.M{"_id": hopID}).One(&hop)
		if err != nil {
			if qmgo.IsErrNoDocuments(err) {
				hop = types.Hop{
					ID:              hopID,
					PathID:          pathId,
					CallerService:   root.Span.LocalEndpoint,
					CallerOperation: root.Span.Name,
					CalledService:   child.Span.LocalEndpoint,
					CalledOperation: child.Span.Name,
				}
				hopCollection.InsertOne(ctx, &hop)
			}
		}
		s.InsertEntityFromGraph(ctx, child, pathId)
	}
}
func (s *Service) IsPathExist(ctx context.Context, pathId uint32) bool {
	return pathIds[pathId]
}
func convertSrToSpan(sr *types.SpanResponse) *types.Span {
	var span types.Span
	span.ID = sr.Id
	span.TraceID = sr.TraceId
	span.Service = sr.LocalEndpoint
	span.Operation = sr.Name
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

func generateOperationID(sr *types.SpanResponse) string {
	return strings.ToUpper(sr.LocalEndpoint + "_" + sr.Name)
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
