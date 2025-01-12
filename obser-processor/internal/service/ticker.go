package service

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) StartTickerFetchTraceData(interval int) {
	go s.FetchTraces(context.Background())
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at ", t)
				go s.FetchTraces(context.Background())
			}
		}
	}()

}
