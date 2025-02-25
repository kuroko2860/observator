package service

import (
	"context"
	"fmt"

	"kuroko.com/processor/internal/types"
)

func (s *Service) ProcessGraph(ctx context.Context, root *types.GraphNode, errPath bool) {
	pathId, err := s.CaculatePathId(ctx, root, 0)
	if err != nil {
		fmt.Printf("caculate path id err: %s\n", err)
		return
	}

	pathEvent := &types.PathEvent{
		PathId:  pathId,
		TraceId: root.Span.TraceId,

		LongestChain: s.caculateLongestChain(ctx, root),
		Timestamp:    root.Span.Timestamp,
		HasError:     errPath,
	}
	pathEventCollection.InsertOne(ctx, pathEvent)
	newRoot := &types.GraphNode{
		Span: &types.SpanResponse{
			Name:          "",
			LocalEndpoint: "",
		},
		Children: make([]*types.GraphNode, 0),
	}
	s.dfs(ctx, newRoot, pathId)
}

func (s *Service) CaculatePathId(ctx context.Context, root *types.GraphNode, level int) (uint32, error) {
	hash := HashCode(root.Span.Name) + HashCode(root.Span.LocalEndpoint) + uint32(level)*31
	for _, child := range root.Children {
		h, err := s.CaculatePathId(ctx, child, level+1)
		if err != nil {
			return 0, err
		}
		hash += h
	}
	return hash, nil
}

func (s *Service) dfs(ctx context.Context, root *types.GraphNode, pathId uint32) {
	if root == nil {
		return
	}
	for _, child := range root.Children {
		hopEvent := &types.HopEvent{
			PathId:              pathId,
			CallerServiceName:   root.Span.LocalEndpoint,
			CallerOperationName: root.Span.Name,
			CalledServiceName:   child.Span.LocalEndpoint,
			CalledOperationName: child.Span.Name,

			Timestamp: child.Span.Timestamp / 1000, // to milisecond
			Duration:  child.Span.Duration,
			HasError:  s.isSpanError(child.Span),
		}
		hopEventCollection.InsertOne(ctx, hopEvent)
		s.dfs(ctx, child, pathId)
	}
}
