package types

type SpanResponse struct {
	TraceID       string            `json:"traceId"`
	ID            string            `json:"id"`
	ParentID      string            `json:"parentId"`
	Kind          string            `json:"kind"`
	LocalEndpoint SpanEndpoint      `json:"localEndpoint"`
	Name          string            `json:"name"`
	Timestamp     int64             `json:"timestamp"` // microseconds
	Duration      int               `json:"duration"`  // microseconds
	Events        []map[string]any  `json:"events"`
	Links         []map[string]any  `json:"links"`
	Tags          map[string]string `json:"tags"`
}
type SpanEndpoint struct {
	ServiceName string `json:"serviceName"`
	IPv4        string `json:"ipv4"`
	Port        int    `json:"port"`
}

type GraphNode struct {
	Span     *SpanResponse
	Children []*GraphNode
}
