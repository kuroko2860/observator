package types

// import (
// 	"time"
// )

// // Span represents a trace span
// type Span struct {
// 	TraceID       string                 `bson:"trace_id"`
// 	SpanID        string                 `bson:"span_id"`
// 	ParentSpanID  string                 `bson:"parent_span_id"`
// 	Name          string                 `bson:"name"`
// 	ServiceName   string                 `bson:"service_name"`
// 	Kind          int                    `bson:"kind"`
// 	StartTime     time.Time              `bson:"start_time"`
// 	EndTime       time.Time              `bson:"end_time"`
// 	Duration      int64                  `bson:"duration"`
// 	Attributes    map[string]interface{} `bson:"attributes"`
// 	Status        int                    `bson:"status"`
// 	StatusMessage string                 `bson:"status_message"`
// 	Events        []SpanEvent            `bson:"events"`
// 	Links         []SpanLink             `bson:"links"`
// }

// // SpanEvent represents an event within a span
// type SpanEvent struct {
// 	Name       string                 `bson:"name"`
// 	Timestamp  time.Time              `bson:"timestamp"`
// 	Attributes map[string]interface{} `bson:"attributes"`
// }

// // SpanLink represents a link to another span
// type SpanLink struct {
// 	TraceID    string                 `bson:"trace_id"`
// 	SpanID     string                 `bson:"span_id"`
// 	Attributes map[string]interface{} `bson:"attributes"`
// }
