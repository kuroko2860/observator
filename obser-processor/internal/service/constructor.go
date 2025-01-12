package service

import (
	"context"
	"math"

	"kuroko.com/processor/internal/types"
)

func (s *Service) ConvertTraceToGraph(ctx context.Context, trace []*types.SpanResponse) (*types.GraphNode, bool) {
	root := &types.GraphNode{}
	var hasError bool = false
	mp := make(map[string]*types.GraphNode)

	for _, span := range trace {
		if s.isSpanError(ctx, span) {
			hasError = true
		}
		var node *types.GraphNode
		if _node, exists := mp[span.Id]; exists {
			_node.Span = span
			node = _node
		} else {
			node = &types.GraphNode{Span: span}
			mp[span.Id] = node
		}

		if span.ParentId != "" {
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
	return root, hasError
}

func (s *Service) constructPathFromGraph(ctx context.Context, root *types.GraphNode, pathId uint32) *types.Path {
	path := &types.Path{
		Id:           pathId,
		LongestChain: s.caculateLongestChain(ctx, root),
		TreeHop:      s.constructTreeHopFromGraph(ctx, root),
	}
	return path
}

func (s *Service) constructTreeHopFromGraph(ctx context.Context, root *types.GraphNode) *types.TreeHop {
	if root == nil {
		return nil
	}
	rootHop := &types.TreeHop{
		OperationName: root.Span.Name,
		ServiceName:   root.Span.LocalEndpoint.ServiceName,
	}
	for _, child := range root.Children {
		rootHop.Children = append(rootHop.Children, s.constructTreeHopFromGraph(ctx, child))
	}
	return rootHop
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
