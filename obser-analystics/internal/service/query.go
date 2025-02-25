package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) FindService(ctx context.Context) ([]model.ServiceObject, error) {
	rs := []model.ServiceObject{}
	svcObjectCollection.Find(ctx, bson.M{}).All(&rs)
	return rs, nil
}

func (s *Service) FindURI(ctx context.Context) ([]model.URIObject, error) {
	rs := []model.URIObject{}
	uriObjectCollection.Find(ctx, bson.M{}).All(&rs)
	return rs, nil
}
