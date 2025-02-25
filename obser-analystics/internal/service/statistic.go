package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) FindServiceStatisticByDate(ctx context.Context, date string) ([]model.ServiceStatisticObject, error) {
	var res []model.ServiceStatisticObject
	err := serviceStatisticObjectCollection.Find(ctx, bson.M{"date": date}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) FindServiceStatisticByDateAndName(ctx context.Context, date string, svcName string) ([]model.ServiceStatisticObject, error) {
	var res []model.ServiceStatisticObject
	err := serviceStatisticObjectCollection.Find(ctx, bson.M{"date": date, "service_name": svcName}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) FindURIStatisticByDate(ctx context.Context, date string) ([]model.URIStatisticObject, error) {
	var res []model.URIStatisticObject
	err := uriStatisticObjectCollection.Find(ctx, bson.M{"date": date}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) FindURIStatisticByDateAndUri(ctx context.Context, date string, uriPath string) ([]model.URIStatisticObject, error) {
	var res []model.URIStatisticObject
	err := uriStatisticObjectCollection.Find(ctx, bson.M{"date": date, "uri_path": uriPath}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
