package service

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

// check time nay co user nao dung service nao
func (s *Service) CheckUserFromTo(ctx context.Context, input model.TimeInput, serviceName string) ([]string, error) {
	entries := []model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{
		"start_time": bson.M{
			"$gte": input.StartTime,
			"$lte": input.EndTime,
		},
		"service_name": serviceName,
	}).All(&entries)

	rs := []string{}
	for _, hle := range entries {
		if !contains(rs, hle.UserId) {
			if hle.Method != "GET" { // post, put moi la tao du lieu
				rs = append(rs, hle.UserId)
			}
		}
	}
	return rs, nil
}

func (s *Service) CheckUserFromToWithPath(ctx context.Context, input model.TimeInput, serviceName string, path string) ([]string, error) {
	entries := []model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{
		"start_time": bson.M{
			"$gte": input.StartTime,
			"$lte": input.EndTime,
		},
		"service_name": serviceName,
	}).All(&entries)

	rs := []string{}
	for _, hle := range entries {
		if !contains(rs, hle.UserId) {
			if hle.Method != "GET" { // post, put moi la tao du lieu
				if strings.Contains(hle.URIPath, path) {
					rs = append(rs, hle.UserId)
				}
			}
		}
	}
	return rs, nil
}
