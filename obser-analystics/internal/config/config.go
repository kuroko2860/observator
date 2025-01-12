package config

const (
	// ProjectName is the name of the project
	ProjectName = "obser-analystics"

	// ProjectVersion is the version of the project
	ProjectVersion = "1.0.0"

	MONGO_URI      = "mongodb://localhost:27017"
	MONGO_DATABASE = "kltn"
	INTERVAL       = 600 // second
	LOOKBACK       = 60000
	LIMIT          = 99999

	TRACE_HOST = "localhost"
	TRACE_PORT = 4111

	Neo4jURI      = "bolt://localhost:7687" // Update with your Neo4j URI
	Neo4jUsername = "neo4j"
	Neo4jPassword = "123456789"
	Neo4jDatabase = "neo4j"
)
