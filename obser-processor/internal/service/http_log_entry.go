package service

import (
	"context"
	"time"

	"kuroko.com/processor/internal/types"
)

func (s *Service) CreateHttpLogEntry(ctx context.Context, http_log_entry *types.HttpLogEntry) (any, error) {
	http_log_entry.StartTimeDate = time.Unix(http_log_entry.StartTime, 0).Local().Format("20030628")

	result, err := httpLogEntryCollection.InsertOne(ctx, http_log_entry)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
