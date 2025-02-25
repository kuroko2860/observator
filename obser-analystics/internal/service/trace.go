package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/analystics/internal/model"
)

func (s *Service) GetAllTracesOfPath(ctx context.Context, pathId uint32) (map[string][]*model.Span, error) {
	var spans []*model.Span
	err := spanCollection.Find(ctx, bson.M{"path_id": pathId}).All(&spans)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]*model.Span)
	for _, span := range spans {
		if _, ok := res[span.TraceId]; !ok {
			res[span.TraceId] = make([]*model.Span, 0)
		}
		res[span.TraceId] = append(res[span.TraceId], span)
	}
	return res, nil
}

func (s *Service) GetTraceById(ctx context.Context, traceId string) ([]*model.Span, error) {
	var spans []*model.Span
	err := spanCollection.Find(ctx, bson.M{"trace_id": traceId}).All(&spans)
	if err != nil {
		return nil, err
	}
	return spans, nil
}

func (s *Service) GetPathByOperation(ctx context.Context, serviceName, operation string) ([]*model.Path, error) {
	return nil, nil
}
func (s *Service) GetHopStatistic(ctx context.Context, callerSvc, callerOp, calledSvc, calledOp string) ([]*model.HopStatistic, error) {

	return nil, nil
}
