const express = require("express");
const { v4: uuidv4 } = require("uuid");
const fs = require("fs");

const app = express();
const PORT = 4111;
const SERVICE_NUMS = 7;
const MIN_SPAN_NUM = 7;
const MAX_SPAN_NUM = 12;

function getRandomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function generateSpans(traceId, baseTimestamp, spanNum) {
  let spans = [];
  let timestamp = baseTimestamp;
  let rootSpan = {
    traceId: traceId,
    parentId: "",
    id: "span1",
    name: `Root operation ${getRandomInt(1, SERVICE_NUMS)}`,
    timestamp: timestamp,
    duration: getRandomInt(50000, 200000),
    localEndpoint: `service-${getRandomInt(1, SERVICE_NUMS)}`,
    tags: { "http.status_code": "200" },
  };
  spans.push(rootSpan);

  let parentSpanIds = ["span1"];

  for (let i = 2; i <= spanNum; i++) {
    let spanId = `span${i}`;
    let parentId = parentSpanIds[getRandomInt(0, parentSpanIds.length - 1)];
    let duration = getRandomInt(10000, 200000);
    let hasError = getRandomInt(0, 30) === 0; // 1/30 chance of an error

    let newSpan = {
      traceId: traceId,
      id: spanId,
      parentId: parentId,
      name: `Operation ${i}`,
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: duration,
      localEndpoint: `service-${getRandomInt(1, SERVICE_NUMS)}`,
      tags: {
        "http.status_code": hasError ? "500" : "200",
        error: hasError ? "true" : "",
      },
    };
    spans.push(newSpan);
    parentSpanIds.push(spanId);
  }
  return spans;
}

function generateSpansStatic(traceId, baseTimestamp, spanNum = 10) {
  let timestamp = baseTimestamp;
  const traceData = [
    {
      traceId: traceId,
      id: "span-1",
      parentId: null,
      name: "service-a-entry-diff",
      localEndpoint: "service-a",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "http.method": "GET",
        "http.path": "/api/entry",
      },
    },
    {
      traceId: traceId,
      id: "span-2",
      parentId: "span-1",
      name: "service-b-handler",
      localEndpoint: "service-b",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "http.method": "POST",
        "http.path": "/api/handle",
      },
    },
    {
      traceId: traceId,
      id: "span-3",
      parentId: "span-2",
      name: "service-c-validate",
      localEndpoint: "service-c",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        validation: "passed",
      },
    },
    {
      traceId: traceId,
      id: "span-4",
      parentId: "span-3",
      name: "service-d-query",
      localEndpoint: "service-d",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "db.query": "SELECT * FROM users",
        "db.row_count": "10",
      },
    },
    {
      traceId: traceId,
      id: "span-5",
      parentId: "span-3",
      name: "service-d-cache-check",
      localEndpoint: "service-d",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "cache.hit": "false",
        error: "not found",
      },
    },
    {
      traceId: traceId,
      id: "span-6",
      parentId: "span-2",
      name: "service-e-logging",
      localEndpoint: "service-e",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "log.level": "INFO",
        "log.message": "Processing started",
      },
    },
    {
      traceId: traceId,
      id: "span-7",
      parentId: "span-4",
      name: "service-c-aggregation",
      localEndpoint: "service-c",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        aggregation: "completed",
      },
    },
    {
      traceId: traceId,
      id: "span-8",
      parentId: "span-5",
      name: "service-d-retry-query",
      localEndpoint: "service-d",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "db.query": "SELECT * FROM users WHERE id = 123",
        error: "Timeout",
      },
    },
    // {
    //   traceId: traceId,
    //   id: "span-9",
    //   parentId: "span-8",
    //   name: "service-e-logging",
    //   localEndpoint: "service-e",
    //   timestamp: timestamp + getRandomInt(1000, 5000),
    //   duration: getRandomInt(10000, 200000),
    //   tags: {
    //     "log.level": "ERROR",
    //     "log.message": "Failed to process",
    //   },
    // },
    {
      traceId: traceId,
      id: "span-10",
      parentId: "span-1",
      name: "service-a-summary",
      localEndpoint: "service-a",
      timestamp: timestamp + getRandomInt(1000, 5000),
      duration: getRandomInt(10000, 200000),
      tags: {
        "http.method": "GET",
        "http.path": "/api/summary",
      },
    },
  ];
  return traceData;
}

const saveLog = (log) => {
  // Save log to file log.txt
  console.log(log);
  fs.appendFile("log.txt", log + "\n", (err) => {
    if (err) {
      console.error(err);
    }
  });
};

app.get("/api/v2/traces", (req, res) => {
  const traceNum = getRandomInt(10, 50);

  const traces = [];
  for (let i = 0; i < traceNum; i++) {
    const baseTimestamp = Date.now();
    const traceId = uuidv4();
    const spanNum = getRandomInt(MIN_SPAN_NUM, MAX_SPAN_NUM);
    // const trace = generateSpans(traceId, baseTimestamp, spanNum);
    const trace = generateSpansStatic(traceId, baseTimestamp);
    console.log(traceId);
    // for (let j = 0; j < spanNum; j++) {
    //   const log = `TraceId: ${traceId}, SpanId: ${trace[j].id}, Input: "input of operation", Output: "output of operation"`;
    //   // saveLog(log);
    // }

    traces.push(trace);
  }
  res.json(traces);
});

app.get("/api/v2/trace/:traceID", (req, res) => {
  const { traceID } = req.params;
  const baseTimestamp = Date.now();

  if (traceID === "abcd1234") {
    const trace = generateSpans(traceID, baseTimestamp);
    return res.json(trace);
  }
  res.status(404).json({ error: "Trace not found" });
});

app.listen(PORT, () => {
  console.log(`Mock Zipkin API running on port ${PORT}`);
});
