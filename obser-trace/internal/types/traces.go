package types

type Time struct {
	ID   string `json:"_id" bson:"_id"`
	Time int64  `json:"time" bson:"time"`
}

type Span struct {
	Id      string `json:"id" bson:"id"`
	PathId  uint32 `json:"path_id" bson:"path_id"`
	TraceId string `json:"trace_id" bson:"trace_id"`

	ServiceName   string `json:"service_name" bson:"service_name"`
	OperationName string `json:"operation_name" bson:"operation_name"`
	Timestamp     int64  `json:"timestamp" bson:"timestamp"`
	Duration      int    `json:"duration" bson:"duration"`
	Error         string `json:"error" bson:"error"`
}

type PathEvent struct {
	PathId  uint32 `json:"path_id" bson:"path_id"`
	TraceId string `json:"trace_id" bson:"trace_id"`

	LongestChain int   `json:"longest_chain" bson:"longest_chain"`
	Timestamp    int64 `json:"timestamp" bson:"timestamp"`
	HasError     bool  `json:"has_error" bson:"has_error"`
	HasLoop      bool  `json:"has_loop" bson:"has_loop"`
}

type HopEvent struct {
	PathId              uint32 `json:"path_id" bson:"path_id"`
	CallerOperationName string `json:"caller_operation_name" bson:"caller_operation_name"`
	CallerServiceName   string `json:"caller_service_name" bson:"caller_service_name"`
	CalledOperationName string `json:"called_operation_name" bson:"called_operation_name"`
	CalledServiceName   string `json:"called_service_name" bson:"called_service_name"`
	Timestamp           int64  `json:"timestamp" bson:"timestamp"`
	Duration            int    `json:"duration" bson:"duration"`
	HasError            bool   `json:"has_error" bson:"has_error"`
}
