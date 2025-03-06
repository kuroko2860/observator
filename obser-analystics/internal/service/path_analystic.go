package service

import (
	"context"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) GetAllPathsFromOperations(ctx context.Context, pairs []model.ServiceOperation) (*model.PathResponse, error) {
	var paths []model.Path
	var conditions []bson.M
	for _, pair := range pairs {
		conditions = append(conditions, bson.M{"operations": bson.M{"$elemMatch": bson.M{"service": pair.Service, "name": pair.Operation}}})
	}

	filter := bson.M{"$and": conditions}

	err := pathCollection.Find(ctx, filter).All(&paths)
	if err != nil {
		return nil, err
	}
	return &model.PathResponse{Paths: paths, TotalCount: len(paths)}, nil
}

func (s *Service) GetPathDetailById(ctx context.Context, _pathId string, _from, _to, unit string) (*model.PathDetail, error) {
	pathId, _ := strconv.ParseUint(_pathId, 10, 32)
	from, to := ParseFromToStringToInt(_from, _to)
	interval := ParseUnitToInterval(unit)

	res := &model.PathDetail{}
	pathInfo := &model.GraphData{PathId: int64(pathId)}

	res.PathInfo = pathInfo
	filter := bson.M{
		"path_id": pathId,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	var pathEvents []*model.PathEvent
	err := pathEventCollection.Find(ctx, filter).All(&pathEvents)
	if err != nil {
		return nil, err
	}
	if len(pathEvents) == 0 {
		return res, nil
	}
	res.Count, res.ErrorCount, res.Distribution, res.ErrorDist = buildPathEventDistribution(pathEvents, from, to, interval)
	res.Frequency = float32(res.Count) * float32(interval) / float32(to-from)
	res.ErrorRate = float32(res.ErrorCount) / float32(res.Count)
	return res, nil
}

func buildPathEventDistribution(pathEvents []*model.PathEvent, from, to, interval int64) (count, errCount int, pathDist, errDist map[int64]int) {
	pathDist = map[int64]int{}
	errDist = map[int64]int{}
	_from := (from / interval) * interval
	_to := (to / interval) * interval
	for i := _from; i <= _to; i += interval {
		pathDist[i] = 0
		errDist[i] = 0
	}
	for _, e := range pathEvents {
		count++
		key := (e.Timestamp / interval) * interval
		pathDist[key]++
		if e.HasError {
			errCount++
			errDist[key]++
		}
	}
	return count, errCount, pathDist, errDist
}

func (s *Service) GetHopDetailById(ctx context.Context, callerSvc, callerOp, calledSvc, calledOp, _from, _to, unit string) (*model.HopDetail, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	interval := ParseUnitToInterval(unit)

	res := &model.HopDetail{}
	var hopInfo = &model.Hop{
		CallerServiceName:   callerSvc,
		CalledServiceName:   calledSvc,
		CallerOperationName: callerOp,
		CalledOperationName: calledOp,
	}
	res.HopInfo = hopInfo

	filter := bson.M{
		"caller_service":   callerSvc,
		"caller_operation": callerOp,
		"called_service":   calledSvc,
		"called_operation": calledOp,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	var hopEvents []*model.HopEvent
	err := hopEventCollection.Find(ctx, filter).All(&hopEvents)
	if err != nil {
		return nil, err
	}
	if len(hopEvents) == 0 {
		return res, nil
	}
	res.Count, res.ErrorCount, res.Distribution, res.ErrorDist, res.Latency = buildHopEventDistribution(hopEvents, from, to, interval)
	res.Frequency = float32(res.Count) * float32(interval) / float32(to-from)
	res.ErrorRate = float32(res.ErrorCount) / float32(res.Count)
	return res, nil
}

func buildHopEventDistribution(hopEvents []*model.HopEvent, from, to, interval int64) (count, errCount int, hopDist, errDist map[int64]int, latency map[string]int) {
	hopDist = map[int64]int{}
	errDist = map[int64]int{}
	_from := (from / interval) * interval
	_to := (to / interval) * interval
	for i := _from; i <= _to; i += interval {
		hopDist[i] = 0
		errDist[i] = 0
	}
	sort.Slice(hopEvents, func(i, j int) bool {
		return hopEvents[i].Duration < hopEvents[j].Duration
	})
	var sum int = 0
	for _, e := range hopEvents {
		count++
		sum += int(e.Duration)
		key := (e.Timestamp / interval) * interval
		hopDist[key]++
		if e.HasError {
			errCount++
			errDist[key]++
		}
	}
	latency = map[string]int{
		"max": int(hopEvents[count-1].Duration),
		"min": int(hopEvents[0].Duration),
		"avg": sum / count,
		"p50": int(hopEvents[count/2].Duration),
		"p95": int(hopEvents[int(float32(count)*float32(0.95))].Duration),
		"p99": int(hopEvents[int(float32(count)*float32(0.99))].Duration),
	}
	return count, errCount, hopDist, errDist, latency
}

func (s *Service) GetLongPath(ctx context.Context, thresholdStr string) ([]*model.GraphData, error) {
	threshold, _ := strconv.ParseInt(thresholdStr, 10, 32)
	var res = []*model.GraphData{}
	err := pathEventCollection.Find(ctx, bson.M{
		"longest_chain": bson.M{
			"$gte": threshold,
		},
	}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
