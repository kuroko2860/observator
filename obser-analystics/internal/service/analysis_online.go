package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) CheckOnlineUser(ctx context.Context, timeInput model.TimeInput) ([]string, error) {
	entries := []model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{
		"uri_path": "/admin/sessions/refresh",
		"start_time": bson.M{
			"$gte": timeInput.StartTime,
			"$lte": timeInput.EndTime,
		},
	}).All(&entries)
	rs := []string{}
	for _, hle := range entries {
		if !contains(rs, hle.UserId) {
			rs = append(rs, hle.UserId)
		}
	}
	return rs, nil
}

func (s *Service) CheckOnlineTime(ctx context.Context, input model.TimeInput, userId string) ([]model.OnlineTimeOutput, error) {
	rs := []model.OnlineTimeOutput{}

	entries := []model.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{
		"uri_path": "/admin/sessions/refresh",
		"user_id":  userId,
		"start_time": bson.M{
			"$gte": input.StartTime,
			"$lte": input.EndTime,
		},
	}).Sort("start_time").All(&entries)

	start := false
	out := model.OnlineTimeOutput{
		UserId:    userId,
		Username:  "",
		StartTime: 0,
		EndTime:   0,
	}
	var anchor int64
	for i := 0; i < len(entries); i++ {
		hle := entries[i]
		if !start { // bat dau moi
			start = true
			out.StartTime = hle.StartTime
			anchor = hle.StartTime
		} else {
			if hle.StartTime-anchor < (360 + 30) { // 6 phut + delta
				anchor = hle.StartTime // set lai anchor
			} else { // ket thuc
				out.EndTime = hle.StartTime
				rs = append(rs, out)

				start = false
				out = model.OnlineTimeOutput{
					UserId:    userId,
					Username:  "",
					StartTime: 0,
					EndTime:   0,
				}
			}

			if i == len(entries)-1 {
				out.EndTime = hle.StartTime
				rs = append(rs, out)
			}
		}
	}
	return rs, nil
}
