package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) FindAllHttpLogEntry(ctx context.Context) ([]model.HttpLogEntry, error) {
	rs := []model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{}).All(&rs)
	return rs, nil
}
func (s *Service) FindHttpLogEntryById(ctx context.Context, id string) (model.HttpLogEntry, error) {
	rs := model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{"_id": id}).One(&rs)
	return rs, nil
}
