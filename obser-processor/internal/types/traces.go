package types

type Time struct {
	ID   string `json:"_id" bson:"_id"`
	Time int64  `json:"time" bson:"time"`
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

type SpanResponse struct {
	TraceId       string `json:"trace_id"`
	Id            string `json:"id"`
	ParentId      string `json:"parent_id"`
	Kind          string `json:"kind"`
	LocalEndpoint struct {
		ServiceName string
	} `json:"local_endpoint"`
	RemoteEndpoint struct {
		ServiceName string
	} `json:"remote_endpoint"`
	Name      string            `json:"name"`
	Timestamp int64             `json:"timestamp"`
	Duration  int               `json:"duration"`
	Tags      map[string]string `json:"tags"`
}

type GraphNode struct {
	Span     *SpanResponse
	Children []*GraphNode
}

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
	HasLoop      bool
	TreeHop      *TreeHop
}

type TreeHop struct {
	OperationName string
	ServiceName   string
	Children      []*TreeHop
}

type PathEvent struct {
	PathId    uint32
	TraceId   string
	Timestamp int64
	HasError  bool
}

type HopEvent struct {
	HopId     uint32
	Timestamp int64
	Duration  int
	HasError  bool
}
