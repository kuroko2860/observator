package types

type AlertGetObject struct {
	ID          string       `json:"id" bson:"id"`
	URIPath     string       `json:"uri_path" bson:"uri_path"`
	Referer     string       `json:"referer" bson:"referer"`
	ServiceName string       `json:"service_name" bson:"service_name"`
	Ignore      bool         `json:"ignore" bson:"ignore"`
	Entry       HttpLogEntry `json:"entry" bson:"entry"`
}

type URIObject struct {
	ID          string `json:"id" bson:"id"`
	ServiceName string `json:"service_name" bson:"service_name"`
	Method      string `json:"method" bson:"method"`
	URIPath     string `json:"uri_path" bson:"uri_path"`
}
type ServiceObject struct {
	ServiceName string `json:"service_name" bson:"service_name"`
}

type ServiceStatisticObject struct {
	ID          string        `json:"id" bson:"id"`
	Date        string        `json:"date" bson:"date"`
	ServiceName string        `json:"service_name" bson:"service_name"`
	Statistic   map[int]int64 `json:"statistic" bson:"statistic"`
}

type URIStatisticObject struct {
	ID          string        `json:"id" bson:"id"`
	Date        string        `json:"date" bson:"date"`
	URIPath     string        `json:"uri_path" bson:"uri_path"`
	ServiceName string        `json:"service_name" bson:"service_name"`
	Method      string        `json:"method" bson:"method"`
	Statistic   map[int]int64 `json:"statistic" bson:"statistic"`
}

type StatisticDone struct {
	// date format yyyyMMdd
	Date string `json:"_id" bson:"_id"`
}

type HttpLogEntry struct {
	ServiceName   string `json:"service_name" bson:"service_name"`
	URIPath       string `json:"uri_path" bson:"uri_path"`
	Referer       string `json:"referer" bson:"referer"`
	UserId        string `json:"user_id" bson:"user_id"`
	StartTime     int64  `json:"start_time" bson:"start_time"`
	Method        string `json:"method" bson:"method"`
	StartTimeDate string `json:"start_time_date" bson:"start_time_date"`
	Host          string `json:"host" bson:"host"`

	Protocol  string `json:"protocol" bson:"protocol"`
	RemoteIP  string `json:"remote_ip" bson:"remote_ip"`
	RequestId string `json:"request_id" bson:"request_id"`
	// TraceId      string `json:"trace_id" bson:"trace_id"`
	// PathId       uint32 `json:"path_id" bson:"path_id"`
	UserAgent    string `json:"user_agent" bson:"user_agent"`
	Duration     int64  `json:"duration" bson:"duration"`
	ResquestSize string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize int64  `json:"response_size" bson:"response_size"`
	StatusCode   int    `json:"status_code" bson:"status_code"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
}
