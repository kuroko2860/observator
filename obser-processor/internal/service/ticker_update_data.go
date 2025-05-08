package service

import (
	"fmt"
	"time"
)

func (s *Service) StartTickerUpdateData(timeInterval int) *time.Ticker {
	ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at ", t)
			// go s.UpdateDataStatistic(context.Background())
		}
	}()

	return ticker
}
