package dto

import "kuroko.com/analystics/internal/model"

type ApiStatistic struct {
	ServiceName string
	Endpoint    string
	Method      string
	From        int64
	To          int64
	Unit        string

	Count         int
	Frequency     float32
	Distribution  map[int64]int
	ErrorCount    int
	ErrorRate     float32
	ErrorDist     map[int]int // err code with count
	ErrorDistTime map[int64]int
	Latency       map[string]int // max min p50 p99
}
type LongApiResponse struct {
	Id struct {
		Endpoint    string
		Method      string
		ServiceName string
	}
	Count int
}

type CalledApiResponse struct {
	Id struct {
		Endpoint    string
		Method      string
		ServiceName string
	}
	Count int
}
type HopDetail struct {
	HopInfo      *model.Hop
	Count        int
	Frequency    float32
	Distribution map[int64]int
	ErrorCount   int
	ErrorRate    float32
	ErrorDist    map[int64]int  // err code with count
	Latency      map[string]int // max min p50 p99
}

type ServiceDetail struct {
	Operations any
	HttpApi    any
}
type PathDetail struct {
	PathInfo     *model.Path
	Count        int
	Frequency    float32
	Distribution map[int64]int
	ErrorCount   int
	ErrorRate    float32
	ErrorDist    map[int64]int
}
