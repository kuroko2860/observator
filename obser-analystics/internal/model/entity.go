package model

type HopStatistic struct{}

type ServiceStatisticObject struct {
	ID          string
	Date        string
	ServiceName string
	Statistic   map[int]int64
}

type URIStatisticObject struct {
	ID          string
	Date        string
	URIPath     string
	ServiceName string
	Method      string
	Statistic   map[int]int64
}

type StatisticDone struct {
	// date format yyyyMMdd
	Date string
}

type TimeInput struct {
	StartTime int64
	EndTime   int64
}
type OnlineTimeOutput struct {
	StartTime int64
	EndTime   int64
	UserId    string
	Username  string
}
type URIObject struct {
	ID          string
	ServiceName string
	Method      string
	URIPath     string
}
type ServiceObject struct {
	ServiceName string
}

type Hop struct {
	ID              string `json:"id" bson:"_id"`
	CallerService   string `json:"caller_service" bson:"caller_service"`
	CallerOperation string `json:"caller_operation" bson:"caller_operation"`
	CalledService   string `json:"called_service" bson:"called_service"`
	CalledOperation string `json:"called_operation" bson:"called_operation"`
}

type LogMiddlewareEvent struct {
	ServiceName string `json:"service_name" bson:"service_name"`
	Method      string `json:"method" bson:"method"`
	URI         string `json:"uri" bson:"uri"`
	URIPath     string `json:"uri_path" bson:"uri_path"`
	Protocol    string `json:"protocol" bson:"protocol"`
	Host        string `json:"host" bson:"host"`
	RemoteIP    string `json:"remote_ip" bson:"remote_ip"`
	RequestId   string `json:"request_id" bson:"request_id"`
	// TraceId      string `json:"trace_id" bson:"trace_id"`
	PathId       uint32 `json:"path_id" bson:"path_id"`
	UserId       string `json:"user_id" bson:"user_id"`
	Referer      string `json:"referer" bson:"referer"`
	UserAgent    string `json:"user_agent" bson:"user_agent"`
	StartTime    int64  `json:"start_time" bson:"start_time"`
	Duration     int64  `json:"duration" bson:"duration"`
	ResquestSize string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize int64  `json:"response_size" bson:"response_size"`
	StatusCode   int    `json:"status_code" bson:"status_code"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
}

type PathEvent struct {
	PathID    uint32 `json:"path_id" bson:"path_id"`
	TraceID   string `json:"trace_id" bson:"trace_id"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	HasError  bool   `json:"has_error" bson:"has_error"`
}

type HopEvent struct {
	ID        string `json:"id" bson:"_id"`
	HopID     string `json:"hop_id" bson:"hop_id"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Duration  int    `json:"duration" bson:"duration"`
	HasError  bool   `json:"has_error" bson:"has_error"`
}

type Span struct {
	ID        string `json:"id" bson:"_id"`
	TraceID   string `json:"trace_id" bson:"trace_id"`
	ParentID  string `json:"parent_id" bson:"parent_id"`
	Service   string `json:"service" bson:"service"`
	Operation string `json:"operation" bson:"operation"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Duration  int    `json:"duration" bson:"duration"`
	Error     string `json:"error" bson:"error"`
	PathID    uint32 `json:"path_id" bson:"path_id"`
	HasError  bool   `json:"has_error" bson:"has_error"`
}

type AlertGetObject struct {
	ID          string       `json:"id" bson:"id"`
	URI         string       `json:"uri" bson:"uri"`
	Referer     string       `json:"referer" bson:"referer"`
	ServiceName string       `json:"service_name" bson:"service_name"`
	Ignore      bool         `json:"ignore" bson:"ignore"`
	Entry       HttpLogEntry `json:"entry" bson:"entry"`
}
type HttpLogEntry struct {
	ServiceName   string `json:"service_name" bson:"service_name"`
	URI           string `json:"uri" bson:"uri"`
	URIPath       string `json:"uri_path" bson:"uri_path"`
	Referer       string `json:"referer" bson:"referer"`
	UserId        string `json:"user_id" bson:"user_id"`
	StartTime     int64  `json:"start_time" bson:"start_time"`
	Method        string `json:"method" bson:"method"`
	StartTimeDate string `json:"start_time_date" bson:"start_time_date"`
	Host          string `json:"host" bson:"host"`
}

type Node struct {
	ID        string `json:"id"`
	Service   string `json:"service"`
	Operation string `json:"operation"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}

type GraphData struct {
	PathId int64  `json:"path_id"`
	Nodes  []Node `json:"nodes"`
	Edges  []Edge `json:"edges"`
}
