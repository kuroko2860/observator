package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) GetAllTracesOfPath(ctx context.Context, pathId string) [][]*model.Span {
	var spans []*model.Span
	spanCollection.Find(ctx, bson.M{"pathId": pathId}).All(&spans)
	return [][]*model.Span{spans}
}
