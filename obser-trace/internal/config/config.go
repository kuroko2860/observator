package config

const (
	// ProjectName is the name of the project
	ProjectName = "obser-processor"

	// ProjectVersion is the version of the project
	ProjectVersion = "1.0.0"

	// MONGO_URI  = "mongodb://mongo-db:27017"
	// TRACE_HOST = "mock-zipkin"
	// Neo4jURI   = "bolt://neo4j-db:7687"
	MONGO_URI  = "mongodb://localhost:27017"
	TRACE_HOST = "localhost"
	Neo4jURI   = "bolt://localhost:7687"

	MONGO_DATABASE = "kltn"
	INTERVAL       = 600 // second
	LOOKBACK       = 60  // second
	LIMIT          = 99999

	TRACE_PORT = 9411

	Neo4jUsername = "neo4j"
	Neo4jPassword = "123456789"
	Neo4jDatabase = "neo4j"
)
