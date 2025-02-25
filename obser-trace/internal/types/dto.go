package types

type SpanResponse struct {
	TraceId        string            `json:"traceId"`
	Id             string            `json:"id"`
	ParentId       string            `json:"parentId"`
	Kind           string            `json:"kind"`
	LocalEndpoint  string            `json:"localEndpoint"`
	RemoteEndpoint string            `json:"remoteEndpoint"`
	Name           string            `json:"name"`
	Timestamp      int64             `json:"timestamp"`
	Duration       int               `json:"duration"`
	Tags           map[string]string `json:"tags"`
}

type GraphNode struct {
	Span     *SpanResponse
	Children []*GraphNode
}
