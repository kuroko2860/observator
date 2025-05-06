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
	// uriId := http_log_entry.ServiceName + "*" + http_log_entry.Method + "*" + http_log_entry.URIPath
	// _, err = uriObjectCollection.UpsertId(ctx, uriId, types.URIObject{
	// 	ID:          uriId,
	// 	ServiceName: http_log_entry.ServiceName,
	// 	Method:      http_log_entry.Method,
	// 	URIPath:     http_log_entry.URIPath,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = svcObjectCollection.UpsertId(ctx, http_log_entry.ServiceName, types.ServiceObject{
	// 	ServiceName: http_log_entry.ServiceName,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	return result.InsertedID, nil
}
