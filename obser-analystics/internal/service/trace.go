package service

import "context"

func (s *Service) GetAllSpans(ctx context.Context) []string {
	return []string{"svc 1", "svc 2"}
}
