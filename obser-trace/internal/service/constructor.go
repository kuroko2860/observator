package service

import (
	"context"
	"math"

	"kuroko.com/processor/internal/types"
)

func (s *Service) ConvertTraceToGraph(ctx context.Context, trace []*types.SpanResponse) (*types.GraphNode, error) {
	if len(trace) == 0 {
		return nil, nil
	}

	// Create a map to store nodes by their span ID
	nodeMap := make(map[string]*types.GraphNode, len(trace))
	var root *types.GraphNode

	// First pass: create all nodes
	for _, span := range trace {
		node, exists := nodeMap[span.ID]
		if !exists {
			node = &types.GraphNode{
				Span:     span,
				Children: []*types.GraphNode{},
			}
			nodeMap[span.ID] = node
		} else {
			node.Span = span
		}

		// Identify the root node (no parent)
		if span.ParentID == "" {
			root = node
		}
	}

	// Second pass: build the tree structure
	for _, span := range trace {
		if span.ParentID != "" {
			parent, exists := nodeMap[span.ParentID]
			if !exists {
				// Create a placeholder parent if it doesn't exist yet
				parent = &types.GraphNode{
					Children: []*types.GraphNode{},
				}
				nodeMap[span.ParentID] = parent
			}
			parent.Children = append(parent.Children, nodeMap[span.ID])
		}
	}

	return root, nil
}

func (s *Service) caculateLongestChain(ctx context.Context, root *types.GraphNode) int {
	if root == nil {
		return 0
	}
	max := 0
	for _, child := range root.Children {
		max = int(math.Max(float64(max), float64(s.caculateLongestChain(ctx, child)+1)))
	}
	return max
}

func (s *Service) caculateLongestErrorChain(ctx context.Context, root *types.GraphNode) int {
	if root == nil {
		return 0
	}
	max := 0
	for _, child := range root.Children {
		max = int(math.Max(float64(max), float64(s.caculateLongestErrorChain(ctx, child)+1)))
	}
	return max
}
