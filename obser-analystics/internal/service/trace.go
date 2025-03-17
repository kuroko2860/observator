package service

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) GetAllTracesOfPath(ctx context.Context, pathId uint32, _from, _to string) ([]*model.TraceSummaryResponse, error) {
	from, to := ParseFromToStringToInt(_from, _to)
	var pe []*model.PathEvent
	err := pathEventCollection.Find(ctx, bson.M{
		"path_id":   pathId,
		"timestamp": bson.M{"$gte": from, "$lte": to}},
	).Limit(10).All(&pe)
	if err != nil {
		return nil, err
	}
	if len(pe) == 0 {
		return nil, nil
	}

	var peids []string
	for _, p := range pe {
		peids = append(peids, p.TraceID)
	}

	var spans []*model.Span
	err = spanCollection.Find(ctx, bson.M{
		"trace_id": bson.M{"$in": peids},
	}).All(&spans)

	if err != nil {
		return nil, err
	}
	res := make(map[string]*model.TraceSummaryResponse)
	for _, span := range spans {
		if _, ok := res[span.TraceID]; !ok {
			res[span.TraceID] = &model.TraceSummaryResponse{SpanNum: 0, TraceId: span.TraceID}
		}
		res[span.TraceID].SpanNum += 1
		if span.ParentID == "" {
			res[span.TraceID].RootService = span.Service
			res[span.TraceID].RootOperation = span.Operation
			res[span.TraceID].StartTime = span.Timestamp
			res[span.TraceID].Duration = span.Duration
		}
	}
	response := []*model.TraceSummaryResponse{}
	for _, v := range res {
		response = append(response, v)
	}
	return response, nil
}

func (s *Service) GetTraceById(ctx context.Context, traceId string) (*model.TraceResponse, error) {
	var trace = &model.TraceResponse{}
	var spans []*model.Span
	err := spanCollection.Find(ctx, bson.M{"trace_id": traceId}).All(&spans)
	if err != nil {
		return nil, err
	}
	trace.Spans = spans
	spanErrMap := make(map[string]bool)
	for _, span := range spans {
		if span.Error != "" {
			spanErrMap[strings.ToUpper(span.Service+"_"+span.Operation)] = true
		}
	}
	trace.SpanErrors = spanErrMap

	var path *model.Path
	err = pathCollection.Find(ctx, bson.M{"path_id": spans[0].PathID}).One(&path)
	if err != nil {
		return nil, err
	}
	trace.Path = path
	return trace, nil
}

func (s *Service) GetPathByOperation(ctx context.Context, serviceName, operation string) ([]*model.Path, error) {
	return nil, nil
}
func (s *Service) GetHopStatistic(ctx context.Context, callerSvc, callerOp, calledSvc, calledOp string) ([]*model.HopStatistic, error) {

	return nil, nil
}
