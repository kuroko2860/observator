curl -X PUT "localhost:9200/_template/microservices_logs" -H 'Content-Type: application/json' -d'
{
  "index_patterns": ["microservices-logs-*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "@timestamp": { "type": "date" },
      "service_name": { "type": "keyword" },
      "uri_path": { "type": "keyword" },
      "method": { "type": "keyword" },
      "status_code": { "type": "integer" },
      "duration": { "type": "long" },
      "trace_id": { "type": "keyword" },
      "span_id": { "type": "keyword" },
      "user_id": { "type": "keyword" },
      "remote_ip": { "type": "ip" },
      "user_agent": { "type": "text" }
    }
  }
}'