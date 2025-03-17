package model

type ApiStatistic struct {
	ServiceName string
	URIPath     string
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
		URIPath     string
		Method      string
		ServiceName string
	}
	Count int
}

type CalledApiResponse struct {
	Id struct {
		URIPath     string
		Method      string
		ServiceName string
	}
	Count int
}
type HopDetail struct {
	HopInfo      *Hop          `json:"hop_info"`
	Count        int           `json:"count"`
	Frequency    float32       `json:"frequency"`
	Distribution map[int64]int `json:"distribution"`
	ErrorCount   int           `json:"error_count"`
	ErrorRate    float32       `json:"error_rate"`
	ErrorDist    map[int64]int `json:"error_dist"`
	Latency      map[int64]int `json:"latency"`
}

type ServiceDetail struct {
	Operations any `json:"operations"`
	HttpApi    any `json:"http_api"`
}
type PathDetail struct {
	PathInfo     *Path         `json:"path_info"`
	Count        int           `json:"count"`
	Frequency    float32       `json:"frequency"`
	Distribution map[int64]int `json:"distribution"`
	ErrorCount   int           `json:"error_count"`
	ErrorRate    float32       `json:"error_rate"`
	ErrorDist    map[int64]int `json:"error_dist"`
}

type GroupResult struct {
	ID struct {
		HopId uint32
	}
	Count int
}

// ServiceOperation represents a service-operation pair
type ServiceOperation struct {
	ID        int    `json:"id"`
	Service   string `json:"service"`
	Operation string `json:"operation"`
}

type RequestPayload struct {
	Pairs []ServiceOperation `json:"pairs"`
}

type PathResponse struct {
	TotalCount int    `json:"total_count"`
	Paths      []Path `json:"paths"`
}

type Path struct {
	ID                string          `json:"id" bson:"_id"`
	PathID            uint32          `json:"path_id" bson:"path_id"`
	CreatedAt         int64           `json:"created_at" bson:"created_at"`
	LongestChain      int             `json:"longest_chain" bson:"longest_chain"`
	LongestErrorChain int             `json:"longest_error_chain" bson:"longest_error_chain"`
	Operations        []PathOperation `json:"operations" bson:"operations"`
	Hops              []PathHop       `json:"hops" bson:"hops"`
}

type PathOperation struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Service string `json:"service" bson:"service"`
}

type PathHop struct {
	ID     string `json:"id" bson:"_id"`
	Source string `json:"source" bson:"source"`
	Target string `json:"target" bson:"target"`
}

// ResultResponse is the API response structure
type ResultResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type TraceSummaryResponse struct {
	RootService   string `json:"root_service"`
	RootOperation string `json:"root_operation"`
	TraceId       string `json:"trace_id"`
	StartTime     int64  `json:"start_time"`
	SpanNum       int    `json:"span_num"`
	Duration      int    `json:"duration"`
}
type TraceResponse struct {
	Spans      []*Span         `json:"spans"`
	Path       *Path           `json:"path"`
	SpanErrors map[string]bool `json:"span_errors"`
}
