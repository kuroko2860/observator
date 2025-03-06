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
	HopInfo      *Hop
	Count        int
	Frequency    float32
	Distribution map[int64]int
	ErrorCount   int
	ErrorRate    float32
	ErrorDist    map[int64]int  // err code with count
	Latency      map[string]int // max min p50 p99
}

type ServiceDetail struct {
	Operations any `json:"operations"`
	HttpApi    any `json:"http_api"`
}
type PathDetail struct {
	PathInfo     *GraphData
	Count        int
	Frequency    float32
	Distribution map[int64]int
	ErrorCount   int
	ErrorRate    float32
	ErrorDist    map[int64]int
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
