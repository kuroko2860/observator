package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/types"
)

func (s *Service) FindEndTime(ctx context.Context) int64 {
	var endTime types.Time
	err := timeCollection.Find(ctx, bson.M{"_id": "end_time"}).One(&endTime)
	if err != nil {
		fmt.Println(err)
		return time.Now().Add(0 * time.Hour).Unix()
	}
	return endTime.Time
}

func (s *Service) UpdateEndTime(ctx context.Context, endTime int64) error {
	_, err := timeCollection.Upsert(ctx, bson.M{"_id": "end_time"}, bson.M{"time": endTime})
	return err
}

func (s *Service) FindLockTime(ctx context.Context) int64 {
	var lockTime types.Time
	err := timeCollection.Find(ctx, bson.M{"_id": "lock_time"}).One(&lockTime)
	if err != nil {
		fmt.Println(err)
		return time.Now().Add(0 * time.Hour).Unix()
	}
	return lockTime.Time
}

func (s *Service) UpdateLockTime(ctx context.Context, lockTime int64) error {
	_, err := timeCollection.Upsert(ctx, bson.M{"_id": "lock_time"}, bson.M{"time": lockTime})
	return err
}
