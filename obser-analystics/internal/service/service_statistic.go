package service

import (
	"context"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) GetAllOperationsFromService(ctx context.Context, serviceName string) ([]string, error) {
	var res []string
	err := operationCollection.Find(ctx, bson.M{"service": serviceName}).Distinct("name", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetAllOperationsCountFromService(ctx context.Context, serviceName string, _from, _to string) (map[string]int, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	var res = map[string]int{}
	var spans []model.Span
	err := spanCollection.Find(ctx, bson.M{
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
		"service": serviceName,
	}).All(&spans)
	if err != nil {
		return nil, err
	}
	for _, span := range spans {
		if _, ok := res[span.Operation]; !ok {
			res[span.Operation] = 0
		}
		res[span.Operation] += 1
	}
	return res, nil
}

func (s *Service) GetAllServices(ctx context.Context) ([]string, error) {
	var res []string
	err := operationCollection.Find(ctx, bson.M{}).Distinct("service", &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetServiceDetailService(ctx context.Context, serviceName, from, to string) (*model.ServiceDetail, error) {
	ops, err := s.GetAllOperationsCountFromService(ctx, serviceName, from, to)
	if err != nil {
		return nil, err
	}
	httpApi, err := s.GetHttpServiceApiService(ctx, serviceName, from, to)
	if err != nil {
		return nil, err
	}
	res := &model.ServiceDetail{
		Operations: ops,
		HttpApi:    httpApi,
	}
	return res, nil
}

func (s *Service) GetHttpServiceApiService(ctx context.Context, serviceName, _from, _to string) (any, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"start_time": bson.M{
					"$gte": from,
					"$lte": to,
				},
				"service_name": serviceName,
			},
		},
	}
	groupStg := bson.D{
		{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"uri_path": "$uri_path",
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
	var apis []bson.M
	err := httpLogEntryCollection.Aggregate(ctx, pipeline).All(&apis)
	if err != nil {
		return nil, err
	}
	return apis, nil
}

func (s *Service) GetServiceEndpointService(ctx context.Context, serviceName string) ([]string, error) {
	var res []string
	err := httpLogEntryCollection.Find(ctx, bson.M{"service_name": serviceName}).Select(bson.M{
		"_id":      0,
		"uri_path": 1,
	}).Distinct("uri_path", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetTopCalledService(ctx context.Context, _from, _to, _limit string) (map[string]int, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	limit, _ := strconv.Atoi(_limit)
	var spans []model.Span
	err := spanCollection.Find(ctx, bson.M{
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}).All(&spans)
	if err != nil {
		return nil, err
	}
	var res = map[string]int{}
	for _, span := range spans {
		if _, ok := res[span.Service]; !ok {
			res[span.Service] = 0
		}
		res[span.Service] += 1
	}
	// sort map by value
	var keys []string
	for k := range res {
		keys = append(keys, k)
	}
	var newRes = map[string]int{}
	sort.Slice(keys, func(i, j int) bool {
		return res[keys[i]] > res[keys[j]]
	})
	if len(keys) > limit {
		keys = keys[:limit]
	}
	for _, key := range keys {
		newRes[key] = res[key]
	}
	return newRes, nil
}
