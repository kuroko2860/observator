package service

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"kuroko.com/processor/internal/types"
)

func (s *Service) ReceiveNATSMsg(m *nats.Msg) error {
	msgCount.Inc()
	entry := types.HttpLogEntry{}
	err := json.Unmarshal(m.Data, &entry)
	if err != nil {
		return err
	}
	s.ProcessHttpLogEntry(m.Subject, entry)
	return nil
}

func (s *Service) ProcessHttpLogEntry(key string, entry types.HttpLogEntry) error {
	if entry.URIPath == "/-/ready" || entry.URIPath == "/metrics" {
		return nil
	}
	// bo service goi service, su dung tracer de track
	// if !strings.Contains(entry.Host, "abc.vn") {
	// 	return nil
	// }
	s.CreateHttpLogEntry(context.Background(), &entry)
	return nil
}
