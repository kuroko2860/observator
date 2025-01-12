const express = require("express");
const app = express();

// Port for the server
const PORT = 4111;

// Sample trace data
const traceData = [
  {
    trace_id: "trace-12345",
    id: "span-1",
    parent_id: null,
    name: "service-a-entry",
    local_endpoint: { serviceName: "service-a" },
    timestamp: 1615550000000,
    duration: 5000,
    tags: {
      "http.method": "GET",
      "http.path": "/api/entry",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-2",
    parent_id: "span-1",
    name: "service-b-handler",
    local_endpoint: { serviceName: "service-b" },
    timestamp: 1615550001000,
    duration: 4000,
    tags: {
      "http.method": "POST",
      "http.path": "/api/handle",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-3",
    parent_id: "span-2",
    name: "service-c-validate",
    local_endpoint: { serviceName: "service-c" },
    timestamp: 1615550002000,
    duration: 2000,
    tags: {
      validation: "passed",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-4",
    parent_id: "span-3",
    name: "service-d-query",
    local_endpoint: { serviceName: "service-d" },
    timestamp: 1615550003000,
    duration: 1000,
    tags: {
      "db.query": "SELECT * FROM users",
      "db.row_count": "10",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-5",
    parent_id: "span-3",
    name: "service-d-cache-check",
    local_endpoint: { serviceName: "service-d" },
    timestamp: 1615550003100,
    duration: 500,
    tags: {
      "cache.hit": "false",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-6",
    parent_id: "span-2",
    name: "service-e-logging",
    local_endpoint: { serviceName: "service-e" },
    timestamp: 1615550004000,
    duration: 800,
    tags: {
      "log.level": "INFO",
      "log.message": "Processing started",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-7",
    parent_id: "span-4",
    name: "service-c-aggregation",
    local_endpoint: { serviceName: "service-c" },
    timestamp: 1615550005000,
    duration: 1200,
    tags: {
      aggregation: "completed",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-8",
    parent_id: "span-5",
    name: "service-d-retry-query",
    local_endpoint: { serviceName: "service-d" },
    timestamp: 1615550005200,
    duration: 900,
    tags: {
      "db.query": "SELECT * FROM users WHERE id = 123",
      error: "Timeout",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-9",
    parent_id: "span-6",
    name: "service-e-logging",
    local_endpoint: { serviceName: "service-e" },
    timestamp: 1615550006100,
    duration: 400,
    tags: {
      "log.level": "ERROR",
      "log.message": "Failed to process",
    },
  },
  {
    trace_id: "trace-12345",
    id: "span-10",
    parent_id: "span-1",
    name: "service-a-summary",
    local_endpoint: { serviceName: "service-a" },
    timestamp: 1615550007000,
    duration: 1500,
    tags: {
      "http.method": "GET",
      "http.path": "/api/summary",
    },
  },
];

// REST API to fetch trace data
app.get("/api/trace", (req, res) => {
  res.json([traceData]);
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
