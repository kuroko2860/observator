package service

import (
	"context"
	"math"

	"kuroko.com/processor/internal/types"
)

func (s *Service) ConvertTraceToGraph(ctx context.Context, trace []*types.SpanResponse) (*types.GraphNode, error) {
	root := &types.GraphNode{}
	mp := make(map[string]*types.GraphNode)

	for _, span := range trace {
		var node *types.GraphNode
		if _node, exists := mp[span.Id]; exists {
			_node.Span = span
			node = _node
		} else {
			node = &types.GraphNode{Span: span}
			mp[span.Id] = node
		}

		if span.ParentId == "" {
			root = node
		} else {
			if parent, exists := mp[span.ParentId]; exists {
				parent.Children = append(parent.Children, node)
			} else {
				parent = &types.GraphNode{}
				parent.Children = append(parent.Children, node)
				mp[span.ParentId] = node
			}
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
