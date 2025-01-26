package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kuroko.com/analystics/internal/model"
	"kuroko.com/analystics/internal/model/dto"
)

func (s *Service) GetAllOperationsFromService(ctx context.Context, serviceName string) ([]string, error) {
	var res []string
	err := operationCollection.Find(ctx, bson.M{"service_name": serviceName}).Select(bson.M{"_id": 0, "name": 1}).Distinct("name", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GroupResult struct {
	ID struct {
		HopId uint32
	}
	Count int
}

func (s *Service) GetAllOperationsCountFromService(ctx context.Context, serviceName string, _from, _to string) (map[string]int, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	var res = map[string]int{}
	var ops []model.Hop
	err := hopCollection.Find(ctx, bson.M{
		"called_service_name": serviceName,
	}).Select(bson.M{"id": 1, "called_operation_name": 1, "_id": 0}).All(&ops)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"timestamp": bson.M{
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
					"hop_id": "$hop_id",
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
	var hopevent []*GroupResult
	hopEventCollection.Aggregate(ctx, pipeline).All(&hopevent)
	if err != nil {
		return nil, err
	}
	for _, op := range ops {
		for _, e := range hopevent {
			if e.ID.HopId == op.Id {
				res[op.CalledOperationName] += e.Count
			}
		}
	}
	return res, nil
}

func (s *Service) GetAllServices(ctx context.Context) ([]string, error) {
	var res []string
	err := operationCollection.Find(ctx, bson.M{}).Select(bson.M{"_id": 0, "service_name": 1}).Distinct("service_name", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetServiceDetailService(ctx context.Context, serviceName, from, to string) (*dto.ServiceDetail, error) {
	ops, err := s.GetAllOperationsCountFromService(ctx, serviceName, from, to)
	if err != nil {
		return nil, err
	}
	httpApi, err := s.GetHttpServiceApiService(ctx, serviceName, from, to)
	if err != nil {
		return nil, err
	}
	res := &dto.ServiceDetail{
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
	var apis []bson.M
	err := logMiddlewareCollection.Aggregate(ctx, pipeline).All(&apis)
	if err != nil {
		return nil, err
	}
	return apis, nil
}

func (s *Service) GetServiceEndpointService(ctx context.Context, serviceName string) ([]string, error) {
	var res []string
	err := logMiddlewareCollection.Find(ctx, bson.M{"service_name": serviceName}).Select(bson.M{
		"_id":      0,
		"endpoint": 1,
	}).Distinct("endpoint", &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetTopCalledService(ctx context.Context, _from, _to, _limit string) (map[string]int, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	matchStg := bson.D{
		{
			Key: "$match", Value: bson.M{
				"timestamp": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
	}
	groupStg := bson.D{
		{
			Key: "$group", Value: bson.M{
				"_id": "$hop_id",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	pipeline := mongo.Pipeline{
		matchStg, groupStg,
	}
	var _result []struct {
		Id    uint32
		Count int
	}
	err := hopEventCollection.Aggregate(ctx, pipeline).All(&_result)
	if err != nil {
		return nil, err
	}
	var resMap = map[uint32]int{}
	for _, r := range _result {
		resMap[r.Id] += r.Count
	}
	var hops []struct {
		CalledServiceName string
		ID                uint32
	}
	hopCollection.Find(ctx, bson.M{}).Select(bson.M{"called_service_name": 1, "id": 1, "_id": 0}).All(&hops)
	var result = map[string]int{}
	for _, hop := range hops {
		result[hop.CalledServiceName] += resMap[hop.ID]
	}
	return result, nil
}
