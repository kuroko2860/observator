package model

type Operation struct {
	Id          uint32
	Name        string
	ServiceName string
}
type Hop struct {
	Id                  uint32
	PathId              uint32
	CallerOperationName string
	CallerServiceName   string
	CalledOperationName string
	CalledServiceName   string
}
type Path struct {
	Id           uint32
	LongestChain int
	TreeHop      *TreeHop
}
type TreeHop struct {
	OperationName string
	ServiceName   string
	Children      []*TreeHop
}
type GraphNode struct {
	Span     *SpanResponse
	Children []*GraphNode
}
type LogMiddlewareEvent struct {
	ServiceName  string
	Method       string
	URI          string
	URIPath      string
	Protocol     string
	Host         string
	RemoteIP     string
	RequestId    string
	TraceId      string
	PathId       uint32
	UserId       string
	Referer      string
	UserAgent    string
	StartTime    int64
	Duration     int64
	ResquestSize string
	ResponseSize int64
	StatusCode   int
	ErrorMessage string
}
type SpanResponse struct {
	TraceId       string
	Id            string
	ParentId      string
	King          string
	LocalEndpoint struct {
		ServiceName string
	}
	RemoteEndpoint struct {
		ServiceName string
	}
	Name      string
	Timestamp int64
	Duration  int
	Tags      map[string]string
}
type PathEvent struct {
	PathId    uint32
	Timestamp int64
	HasError  bool
}
type HopEvent struct {
	HopId               uint32
	CalledOperationName string
	CalledServiceName   string
	Timestamp           int64
	Duration            int
	HasError            bool
}

type Span struct {
	Id            string
	TraceId       string
	ServiceName   string
	OperationName string
	Timestamp     int64
	Duration      int
	Error         string
	PathId        uint32
}

type Node struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type GraphData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
