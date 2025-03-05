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

// RequestPayload represents the body of the request
type RequestPayload struct {
	Pairs []ServiceOperation `json:"pairs"`
}

// PathResponse represents a path in the graph
type PathResponse struct {
	Nodes []OperationNode `json:"nodes"`
	Path  string          `json:"path"`
}

// OperationNode represents an operation node from Neo4j
type OperationNode struct {
	Name    string `json:"name"`
	Service string `json:"service"`
}

// ResultResponse is the API response structure
type ResultResponse struct {
	Success bool           `json:"success"`
	Paths   []PathResponse `json:"paths"`
	Error   string         `json:"error,omitempty"`
}
