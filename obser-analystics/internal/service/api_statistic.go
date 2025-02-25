package service

import (
	"context"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kuroko.com/analystics/internal/model"
)

type Log struct {
	StartTime  int64
	StatusCode int
	Duration   int64
}

func (s *Service) GetApiStatisticService(ctx context.Context, serviceName, endpoint, method, _from, _to, unit string) (*model.ApiStatistic, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	interval := ParseUnitToInterval(unit)

	res := &model.ApiStatistic{
		ServiceName: serviceName,
		Endpoint:    endpoint,
		Method:      method,
		From:        from,
		To:          to,
		Unit:        unit,
	}
	filter := bson.M{
		"service_name": serviceName,
		"uri_path":     endpoint,
		"method":       method,
		"start_time": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	selected := bson.M{
		"_id":         0,
		"start_time":  1,
		"status_code": 1,
		"duration":    1,
	}
	var logs []*Log
	err := httpLogEntryCollection.Find(ctx, filter).Select(selected).All(&logs)
	if err != nil {
		return nil, err
	}
	count := len(logs)
	if count == 0 {
		return nil, nil
	}
	res.Count = int(count)
	freq := float32(count) * float32(interval) / float32(to-from)
	res.Frequency = freq
	res.Distribution = s.GetDistributionApiUsageService(ctx, logs, from, to, interval)
	errCount, errDist, errTimeDist := s.GetApiErrorService(ctx, logs, from, to, interval)
	res.ErrorCount = errCount
	res.ErrorRate = float32(errCount) / float32(count)
	res.ErrorDist = errDist
	res.ErrorDistTime = errTimeDist
	res.Latency = s.GetLatencyService(ctx, logs, from, to, interval)
	return res, nil
}

func (s *Service) GetLatencyService(ctx context.Context, logs []*Log, from, to int64, unit int64) map[string]int {
	var latencies []int
	for _, log := range logs {
		latencies = append(latencies, int(log.Duration))
	}
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	count := len(latencies)
	var sum int = 0
	for _, l := range latencies {
		sum += l
	}
	res := map[string]int{
		"max": latencies[count-1],
		"min": latencies[0],
		"avg": sum / count,
		"p50": latencies[count/2],
		"p95": latencies[int(float32(count)*float32(0.95))],
		"p99": latencies[int(float32(count)*float32(0.99))],
	}
	return res
}

func (s *Service) GetApiErrorService(ctx context.Context, logs []*Log, from, to, unit int64) (int, map[int]int, map[int64]int) {
	count := 0
	var dist = map[int]int{}
	for _, log := range logs {
		if log.StatusCode >= 400 && log.StartTime <= 600 {
			dist[log.StatusCode]++
			count++
		}
	}
	var res = map[int64]int{}
	_from := (from / unit) * unit
	_to := (to / unit) * unit
	for i := _from; i <= _to; i += unit {
		res[i] = 0
	}
	for _, log := range logs {
		if log.StatusCode >= 400 && log.StartTime <= 600 {
			key := (log.StartTime / unit) * unit
			res[key]++
		}
	}
	return count, dist, res
}

func (s *Service) GetDistributionApiUsageService(ctx context.Context, logs []*Log, from, to, unit int64) map[int64]int {
	var res = map[int64]int{}
	_from := (from / unit) * unit
	_to := (to / unit) * unit
	for i := _from; i <= _to; i += unit {
		res[i] = 0
	}
	for _, log := range logs {
		key := (log.StartTime / unit) * unit
		res[key]++
	}
	return res
}

func (s *Service) GetLongApiService(ctx context.Context, from, to, threshold string) ([]bson.M, error) {
	fromInt, toInt := ParseFromToStringToInt(from, to)
	thresholdNumber, _ := strconv.ParseInt(threshold, 10, 32)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"start_time": bson.M{
					"$gte": fromInt,
					"$lte": toInt,
				},
				"duration": bson.M{
					"$gte": thresholdNumber,
				},
			},
		},
	}
	groupStg := bson.D{
		{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"service_name": "$service_name",
					"endpoint":     "$endpoint",
					"method":       "$method",
				},
				"count": bson.M{
					"$sum": 1,
				},
				"avg_latency": bson.M{
					"$avg": "$duration",
				},
			},
		},
	}
	pipeline := mongo.Pipeline{
		matchStg, groupStg,
	}
	var result = []bson.M{}
	err := httpLogEntryCollection.Aggregate(ctx, pipeline).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetCalledApiService(ctx context.Context, from, to, username string) ([]bson.M, error) {
	fromInt, toInt := ParseFromToStringToInt(from, to)
	match := bson.M{
		"start_time": bson.M{
			"$gte": fromInt,
			"$lte": toInt,
		},
	}
	if len(username) != 0 {
		match["username"] = username
	}
	matchStg := bson.D{
		{
			Key: "$match", Value: match,
		},
	}
	projectStg := bson.D{
		{
			Key: "$project", Value: bson.M{
				"service_name": 1,
				"endpoint":     1,
				"method":       1,
				"username":     1,
				"is_error": bson.M{
					"$cond": bson.A{
						bson.M{"$gte": bson.A{"$status_code", 400}},
						1,
						0,
					},
				},
			},
		},
	}
	groupStg := bson.D{
		{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"service_name": "$service_name",
				"endpoint":     "$endpoint",
				"method":       "$method",
				"username":     "$username",
			},
			"count": bson.M{
				"$sum": 1,
			},
			"err_count": bson.M{
				"$sum": "$is_error",
			},
		}},
	}
	pipeline := mongo.Pipeline{
		matchStg, projectStg, groupStg,
	}
	var result = []bson.M{}
	err := httpLogEntryCollection.Aggregate(ctx, pipeline).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetTopCalledApi(ctx context.Context, _from, _to, _limit string) ([]bson.M, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"start_time": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
	}
	projectStg := bson.D{
		{
			Key: "$project", Value: bson.M{
				"service_name": 1,
				"endpoint":     1,
				"method":       1,
				"is_error": bson.M{
					"$cond": bson.A{
						bson.M{"$gte": bson.A{"$status_code", 400}},
						1,
						0,
					},
				},
			},
		},
	}
	groupStg := bson.D{
		{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"service_name": "$service_name",
				"endpoint":     "$endpoint",
				"method":       "$method",
			},
			"count": bson.M{
				"$sum": 1,
			},
			"err_count": bson.M{
				"$sum": "$is_error",
			},
		}},
	}
	sortStg := bson.D{
		{
			Key: "$sort", Value: bson.M{
				"count": -1,
			},
		},
	}
	pipeline := mongo.Pipeline{
		matchStg, projectStg, groupStg, sortStg,
	}
	var result = []bson.M{}
	err := httpLogEntryCollection.Aggregate(ctx, pipeline).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetHttpApiByService(ctx context.Context, _from, _to, service_name string) ([]bson.M, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"service_name": service_name,
				"start_time": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
	}
	groupStg := bson.D{
		{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"endpoint": "$endpoint",
					"method":   "$method",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	pipeline := mongo.Pipeline{
		matchStg, groupStg,
	}
	var result = []bson.M{}
	err := httpLogEntryCollection.Aggregate(ctx, pipeline).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
