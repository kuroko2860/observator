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
