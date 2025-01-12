package service

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"kuroko.com/analystics/internal/config"
	"kuroko.com/analystics/internal/types"
)

func (s *Service) GetPathById(ctx context.Context, id string) (*types.GraphData, error) {
	_id, _ := strconv.ParseUint(id, 10, 32)
	// Open a session
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: config.Neo4jDatabase})
	defer session.Close(ctx)

	// Query the graph data
	query := `
		MATCH (p1:Operation)-[:CALLS{pathId: $pathId}]->(p2:Operation)
		RETURN p1, p2
	`
	data := types.GraphData{}
	operations, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, map[string]any{
			"pathId": _id,
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		return records, nil
	})

	if err != nil {
		return nil, err
	}
	for _, operation := range operations.([]*neo4j.Record) {
		operation1, _ := operation.Get("p1")
		operation2, _ := operation.Get("p2")
		data.Nodes = append(data.Nodes, types.Node{
			ID: operation1.(neo4j.Node).Props["name"].(string),
		})
		data.Nodes = append(data.Nodes, types.Node{
			ID: operation2.(neo4j.Node).Props["name"].(string),
		})
		data.Edges = append(data.Edges, types.Edge{
			Source: operation1.(neo4j.Node).Props["name"].(string),
			Target: operation2.(neo4j.Node).Props["name"].(string),
		})
	}

	return &data, nil

}
