package config

const (
	// ProjectName is the name of the project
	ProjectName = "obser-processor"

	// ProjectVersion is the version of the project
	ProjectVersion = "1.0.0"

	MONGO_URI  = "mongodb://localhost:27017"
	TRACE_HOST = "localhost"

	MONGO_DATABASE = "kltn"
	INTERVAL       = 60 // second
	LOOKBACK       = 60 // second
	LIMIT          = 99999

	TRACE_PORT = 9411
	NATS_URL   = "nats://localhost:4222"
)
