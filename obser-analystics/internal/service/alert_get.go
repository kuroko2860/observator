package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) FindAllAlertGet(ctx context.Context) ([]model.AlertGetObject, error) {
	rs := []model.AlertGetObject{}
	alertGetCollection.Find(ctx, bson.M{"ignore": false}).All(&rs)
	return rs, nil
}

func (s *Service) IgnoreAlertGet(ctx context.Context, id string) error {
	rs := model.AlertGetObject{}
	alertGetCollection.Find(ctx, bson.M{"id": id}).One(&rs)
	rs.Ignore = true

	err := alertGetCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": rs})
	if err != nil {
		return err
	}
	return nil
}
