package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/types"
)

func (s *Service) ProcessGraph(ctx context.Context, root *types.GraphNode, hasError bool) {
	pathId, err := s.CaculatePathId(ctx, root, 0)
	if err != nil {
		fmt.Printf("caculate path id err: %s\n", err)
		return
	}
	if !s.pathIdExisted(ctx, pathId) {
		path := s.constructPathFromGraph(ctx, root, pathId)
		s.insertPath(ctx, path)
	}
	pathEvent := &types.PathEvent{
		PathId:    pathId,
		TraceId:   root.Span.TraceId,
		Timestamp: root.Span.Timestamp,
		HasError:  hasError,
	}
	s.insertPathEvent(ctx, pathEvent)
	newRoot := &types.GraphNode{
		Span: &types.SpanResponse{
			Name: "",
			LocalEndpoint: struct{ ServiceName string }{
				ServiceName: "",
			},
		},
		Children: make([]*types.GraphNode, 0),
	}
	s.dfs(ctx, newRoot, pathId)
}

func (s *Service) insertPathEvent(ctx context.Context, pathEvent *types.PathEvent) {
	pathEventCollection.InsertOne(ctx, pathEvent)
}

func (s *Service) insertPath(ctx context.Context, path *types.Path) {
	pathCollection.InsertOne(ctx, path)
}

func (s *Service) pathIdExisted(ctx context.Context, pathId uint32) bool {
	var path *types.Path
	pathCollection.Find(ctx, bson.M{"id": pathId}).One(&path)
	return path != nil
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

func (s *Service) dfs(ctx context.Context, root *types.GraphNode, pathId uint32) {
	if root == nil {
		return
	}
	for _, child := range root.Children {
		operation := types.Operation{
			Id:          s.CaculateOperationId(ctx, child),
			Name:        child.Span.Name,
			ServiceName: child.Span.LocalEndpoint.ServiceName,
		}
		if !s.operationExisted(ctx, operation) {
			s.insertOperation(ctx, operation)
		}
		hop := s.generateHop(ctx, root, child, pathId)
		if !s.hopExisted(ctx, hop) {
			s.insertHop(ctx, hop)
		}

		hopEvent := &types.HopEvent{
			HopId:     hop.Id,
			Timestamp: child.Span.Timestamp / 1000, // to milisecond
			Duration:  child.Span.Duration,
			HasError:  s.isSpanError(ctx, child.Span),
		}
		s.insertHopEvent(ctx, hopEvent)
		s.dfs(ctx, child, pathId)
	}
}

func (s *Service) insertHopEvent(ctx context.Context, hopEvent *types.HopEvent) {
	hopEventCollection.InsertOne(ctx, hopEvent)
}

func (s *Service) insertHop(ctx context.Context, hop *types.Hop) {
	hopCollection.InsertOne(ctx, hop)
}

func (s *Service) hopExisted(ctx context.Context, hop *types.Hop) bool {
	var h *types.Hop
	hopCollection.Find(ctx, bson.M{"id": hop.Id}).One(&h)
	return h != nil
}

func (s *Service) insertOperation(ctx context.Context, operation types.Operation) {
	operationCollection.InsertOne(ctx, operation)
}

func (s *Service) operationExisted(ctx context.Context, operation types.Operation) bool {
	var o types.Operation
	operationCollection.Find(ctx, bson.M{"id": operation.Id}).One(&o)
	return o != types.Operation{}
}

func (s *Service) CaculateOperationId(ctx context.Context, child *types.GraphNode) uint32 {
	return HashCode(child.Span.Name) + HashCode(child.Span.LocalEndpoint.ServiceName)
}

func (s *Service) generateHop(ctx context.Context, root *types.GraphNode,
	child *types.GraphNode, pathId uint32) *types.Hop {
	hop := &types.Hop{
		PathId:              pathId,
		CallerOperationName: root.Span.Name,
		CallerServiceName:   root.Span.LocalEndpoint.ServiceName,
		CalledOperationName: child.Span.Name,
		CalledServiceName:   child.Span.LocalEndpoint.ServiceName,
	}
	hop.Id = s.CaculateHopId(ctx, hop)
	return hop
}

func (s *Service) CaculateHopId(ctx context.Context, hop *types.Hop) uint32 {
	return (HashCode(hop.CallerOperationName)+HashCode(hop.CallerServiceName))*31 +
		(HashCode(hop.CalledOperationName)+HashCode(hop.CalledServiceName))*37
}
