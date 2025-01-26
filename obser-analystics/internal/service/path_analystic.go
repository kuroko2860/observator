package service

import (
	"context"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
	"kuroko.com/analystics/internal/model/dto"
)

func (s *Service) GetAllPathFromHop(ctx context.Context, callerSvc, callerOp, calledSvc, calledOp string) ([]*model.Path, error) {
	var results []*model.Hop
	filter := bson.M{}
	if len(calledOp) > 0 {
		filter["called_operation_name"] = calledOp
	}
	if len(calledSvc) > 0 {
		filter["called_service_name"] = calledSvc
	}
	if len(callerOp) > 0 {
		filter["caller_operation_name"] = callerOp
	}
	if len(callerSvc) > 0 {
		filter["caller_service_name"] = callerSvc
	}
	err := hopCollection.Find(ctx, filter).All(&results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return []*model.Path{}, nil
	}
	var pathIdArr []uint32
	for _, res := range results {
		pathIdArr = append(pathIdArr, res.PathId)
	}
	var paths []*model.Path
	err = pathCollection.Find(ctx, bson.M{
		"id": bson.M{
			"$in": pathIdArr,
		},
	}).All(&paths)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func (s *Service) GetPathDetailById(ctx context.Context, _pathId string, _from, _to, unit string) (*dto.PathDetail, error) {
	pathId, _ := strconv.ParseInt(_pathId, 10, 32)
	from, to := ParseFromToStringToInt(_from, _to)
	interval := ParseUnitToInterval(unit)

	res := &dto.PathDetail{}
	var pathInfo *model.Path

	err := pathCollection.Find(ctx, bson.M{
		"id": pathId,
	}).One(&pathInfo)
	if err != nil {
		return nil, err
	}
	res.PathInfo = pathInfo
	filter := bson.M{
		"path_id": pathId,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	var pathEvents []*model.PathEvent
	err = pathEventCollection.Find(ctx, filter).All(&pathEvents)
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

func (s *Service) GetHopId(ctx context.Context, caller_svc, caller_op, called_svc, called_op string) (uint32, error) {
	var hop *model.Hop
	err := hopCollection.Find(ctx, bson.M{
		"caller_operation_name": caller_op,
		"caller_service_name":   caller_svc,
		"called_operation_name": called_op,
		"called_service_name":   called_svc,
	}).One(&hop)
	if err != nil {
		return 0, err
	}
	return hop.Id, nil
}

func (s *Service) GetHopDetailById(ctx context.Context, hopId uint32, _from, _to, unit string) (*dto.HopDetail, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	interval := ParseUnitToInterval(unit)

	res := &dto.HopDetail{}
	var hopInfo *model.Hop

	err := hopCollection.Find(ctx, bson.M{
		"id": hopId,
	}).One(&hopInfo)
	if err != nil {
		return nil, err
	}
	res.HopInfo = hopInfo
	filter := bson.M{
		"hop_id": hopId,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	var hopEvents []*model.HopEvent
	err = hopEventCollection.Find(ctx, filter).All(&hopEvents)
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

func (s *Service) GetLongPath(ctx context.Context, thresholdStr string) ([]*model.Path, error) {
	threshold, _ := strconv.ParseInt(thresholdStr, 10, 32)
	var res = []*model.Path{}
	err := pathCollection.Find(ctx, bson.M{
		"longest_chain": bson.M{
			"$gte": threshold,
		},
	}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
