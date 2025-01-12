package service

import (
	"fmt"
	"time"
)

func (s *Service) StartTickerFetchTraceData(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at ", t)
			}
		}
	}()

}
