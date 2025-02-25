package service

import (
	"hash/fnv"

	"kuroko.com/processor/internal/types"
)

func HashCode(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

func (s *Service) isSpanError(span *types.SpanResponse) bool {
	return span.Tags["error"] != ""
}
