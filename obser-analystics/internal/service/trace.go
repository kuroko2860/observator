package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/types"
)

func (s *Service) GetAllTracesOfPath(ctx context.Context, pathId string) [][]*types.Span {
	var spans []*types.Span
	spanCollection.Find(ctx, bson.M{"pathId": pathId}).All(&spans)
	return []string{"svc 1", "svc 2"}
}
