package types

type Time struct {
	ID   string `json:"id" bson:"_id"`
	Time int64  `json:"time" bson:"time"`
}

type Span struct {
	ID      string `json:"id" bson:"_id"`
	PathID  uint32 `json:"path_id" bson:"path_id"`
	TraceID string `json:"trace_id" bson:"trace_id"`

	Service   string `json:"service" bson:"service"`
	Operation string `json:"operation" bson:"operation"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Duration  int    `json:"duration" bson:"duration"`
	Error     string `json:"error" bson:"error"`
}

type PathEvent struct {
	ID        string `json:"id" bson:"_id"`
	PathID    uint32 `json:"path_id" bson:"path_id"`
	TraceID   string `json:"trace_id" bson:"trace_id"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

type HopEvent struct {
	ID        string `json:"id" bson:"_id"`
	HopID     string `json:"hop_id" bson:"hop_id"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Duration  int    `json:"duration" bson:"duration"`
	HasError  bool   `json:"has_error" bson:"has_error"`
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

type Hop struct {
	ID              string `json:"id" bson:"_id"`
	PathID          uint32 `json:"path_id" bson:"path_id"`
	CallerOperation string `json:"caller_operation" bson:"caller_operation"`
	CallerService   string `json:"caller_service" bson:"caller_service"`
	CalledOperation string `json:"called_operation" bson:"called_operation"`
	CalledService   string `json:"called_service" bson:"called_service"`
}

type Operation struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Service string `json:"service" bson:"service"`
}
