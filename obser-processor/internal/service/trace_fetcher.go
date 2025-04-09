package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/types"
)

var pathIds map[uint32]bool

func (s *Service) ProcessTrace(ctx context.Context, trace []*types.SpanResponse) error {
	root, err := s.ConvertTraceToGraph(ctx, trace)
	if err != nil {
		fmt.Printf("Error when converting trace: %s\n", err.Error())
		return err
	}
	pathId, err := s.CaculatePathId(ctx, root, 0)
	if err != nil {
		fmt.Printf("Error when caculating path id: %s\n", err.Error())
		return err
	}
	if !s.IsPathExist(ctx, pathId) {
		s.InsertEntityFromGraph(ctx, root, pathId)
		s.InsertPath(ctx, root, pathId)

		pathIds[pathId] = true
		pathIdCollection.InsertOne(ctx, bson.M{"_id": pathId})
	}
	s.ProcessGraph(ctx, root, pathId)

	for _, sr := range trace {
		span := convertSrToSpan(sr)
		span.PathID = pathId
		spanCollection.InsertOne(ctx, span)
	}

	return err
}

func convertSrToSpan(sr *types.SpanResponse) *types.Span {
	var span types.Span
	_, hasErrorTag := sr.Tags["error"]
	_, hasErrorMessageTag := sr.Tags["error.message"]
	span.ID = sr.ID
	span.TraceID = sr.TraceID
	span.Service = sr.LocalEndpoint.ServiceName
	span.Operation = sr.Name
	span.Timestamp = sr.Timestamp
	span.Duration = sr.Duration
	span.Error = sr.Tags["error"] + sr.Tags["error.message"]
	span.HasError = hasErrorTag || hasErrorMessageTag
	span.ParentID = sr.ParentID
	return &span
}

func (s *Service) InsertPath(ctx context.Context, root *types.GraphNode, pathId uint32) {
	path := types.Path{
		ID:                uuid.NewString(),
		PathID:            pathId,
		CreatedAt:         root.Span.Timestamp / 1000,
		Operations:        []types.PathOperation{},
		Hops:              []types.PathHop{},
		LongestChain:      s.caculateLongestChain(ctx, root),
		LongestErrorChain: s.caculateLongestErrorChain(ctx, root), // todo
	}

	// Use a map to track processed operations to avoid duplicates
	nodeMap := make(map[string]types.PathOperation)

	// Use a separate map to track processed nodes to ensure we visit each node exactly once
	processedNodes := make(map[string]bool)

	var processNode func(node *types.GraphNode)

	processNode = func(node *types.GraphNode) {
		// Skip if node is nil or has no span
		if node == nil || node.Span == nil {
			return
		}

		// Skip if we've already processed this node
		if processedNodes[node.Span.ID] {
			return
		}

		// Mark this node as processed
		processedNodes[node.Span.ID] = true

		// Create operation ID
		id := strings.ToUpper(node.Span.LocalEndpoint.ServiceName + "_" + node.Span.Name)

		// Add operation if it doesn't exist
		if _, exists := nodeMap[id]; !exists {
			nodeOp := types.PathOperation{
				ID:      id,
				Name:    node.Span.Name,
				Service: node.Span.LocalEndpoint.ServiceName,
			}
			nodeMap[id] = nodeOp
			path.Operations = append(path.Operations, nodeOp)
		}

		// Process all children
		for _, child := range node.Children {
			if child == nil || child.Span == nil {
				continue
			}

			// Create hop from this node to child
			targetId := strings.ToUpper(child.Span.LocalEndpoint.ServiceName + "_" + child.Span.Name)
			hop := types.PathHop{
				ID:     uuid.NewString(),
				Source: id,
				Target: targetId,
			}
			path.Hops = append(path.Hops, hop)

			// Process child node
			processNode(child)
		}
	}

	// Start processing from the root
	processNode(root)

	fmt.Printf("Path created with %d operations and %d hops\n", len(path.Operations), len(path.Hops))
	pathCollection.InsertOne(ctx, &path)
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
				Service: root.Span.LocalEndpoint.ServiceName,
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
					CallerService:   root.Span.LocalEndpoint.ServiceName,
					CallerOperation: root.Span.Name,
					CalledService:   child.Span.LocalEndpoint.ServiceName,
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

func generateOperationID(sr *types.SpanResponse) string {
	return strings.ToUpper(sr.LocalEndpoint.ServiceName + "_" + sr.Name)
}
