package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/types"
)

func (s *Service) ProcessGraph(ctx context.Context, root *types.GraphNode, pathId uint32) {
	pathEvent := &types.PathEvent{
		ID:        uuid.New().String(),
		PathID:    pathId,
		TraceID:   root.Span.TraceID,
		Timestamp: root.Span.Timestamp / 1000,
	}
	pathEventCollection.InsertOne(ctx, pathEvent)
	newRoot := &types.GraphNode{
		Span: &types.SpanResponse{
			Name: "root",
			LocalEndpoint: types.SpanEndpoint{
				ServiceName: "root",
			},
		},
		Children: make([]*types.GraphNode, 0),
	}
	newRoot.Children = append(newRoot.Children, root)
	s.dfs(ctx, newRoot, pathId)
}

func (s *Service) CaculatePathId(ctx context.Context, root *types.GraphNode, level int) (uint32, error) {
	hash := HashCode(root.Span.Name) + HashCode(root.Span.LocalEndpoint.ServiceName) + uint32(level)*31
	for _, child := range root.Children {
		h, err := s.CaculatePathId(ctx, child, level+1)
		if err != nil {
			return 0, err
		}
		hash += h
	}
	return hash, nil
}

// insert hop and hop event
func (s *Service) dfs(ctx context.Context, root *types.GraphNode, pathId uint32) {
	if root == nil {
		return
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
					CallerService:   root.Span.LocalEndpoint.ServiceName,
					CallerOperation: root.Span.Name,
					CalledService:   child.Span.LocalEndpoint.ServiceName,
					CalledOperation: child.Span.Name,
				}
				hopCollection.InsertOne(ctx, &hop)
			}
		}
		hopEvent := &types.HopEvent{
			ID:        uuid.NewString(),
			HopID:     hopID,
			Timestamp: child.Span.Timestamp / 1000,
			Duration:  child.Span.Duration,
			HasError:  s.isSpanError(child.Span),
		}
		hopEventCollection.InsertOne(ctx, hopEvent)
		s.dfs(ctx, child, pathId)
	}
}

func generateHopID(parent, child *types.GraphNode, pathId uint32) string {
	return strings.ToUpper(parent.Span.LocalEndpoint.ServiceName + "_" + parent.Span.Name + "_" + child.Span.LocalEndpoint.ServiceName + "_" + child.Span.Name + "_" + strconv.Itoa(int(pathId)))
}
