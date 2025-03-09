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
  // let timestamp = baseTimestamp;
  const traceData = [
    {
      id: "span-root",
      traceId: traceId,
      parentId: null,
      name: "GET /api/products",
      localEndpoint: "frontend",
      timestamp: 1703664522000,
      duration: 450000,
      tags: {
        "http.method": "GET",
        "http.path": "/api/products",
        component: "web",
      },
    },
    {
      id: "span-gateway1",
      traceId: traceId,
      parentId: "span-root",
      name: "route-request",
      localEndpoint: "api-gateway",
      timestamp: 1703664522050,
      duration: 400000,

      tags: {
        "gateway.route": "/products",
        "request.id": "req-123",
      },
    },
    {
      id: "span-auth1",
      traceId: traceId,
      parentId: "span-gateway1",
      name: "authenticate",
      localEndpoint: "auth-service",
      timestamp: 1703664522100,
      duration: 150000,

      tags: {
        "auth.type": "jwt",
        "user.id": "user-456",
      },
    },
    {
      id: "span-product1",
      traceId: traceId,
      parentId: "span-gateway1",
      name: "get-products",
      localEndpoint: "product-service",
      timestamp: 1703664522150,
      duration: 200000,

      tags: {
        "product.category": "electronics",
        "query.limit": "50",
      },
    },
    {
      id: "span-cache1",
      traceId: traceId,
      parentId: "span-product1",
      name: "check-cache",
      localEndpoint: "cache-service",
      timestamp: 1703664522200,
      duration: 50000,

      tags: {
        "cache.key": "products:electronics",
        "cache.hit": "false",
      },
    },
    {
      id: "span-db1",
      traceId: traceId,
      parentId: "span-product1",
      name: "query-database",
      localEndpoint: "database",
      timestamp: 1703664522250,
      duration: 100000,

      tags: {
        "db.query": "SELECT * FROM products",
        "db.rows": "42",
      },
    },
    {
      id: "span-inventory1",
      traceId: traceId,
      parentId: "span-product1",
      name: "check-inventory",
      localEndpoint: "inventory-service",
      timestamp: 1703664522300,
      duration: 75000,

      tags: {
        "inventory.check": "bulk",
        "products.count": "42",
      },
    },
    {
      id: "span-pricing1",
      traceId: traceId,
      parentId: "span-product1",
      name: "calculate-prices",
      localEndpoint: "pricing-service",
      timestamp: 1703664522350,
      duration: 80000,

      tags: {
        "pricing.type": "dynamic",
        currency: "USD",
      },
    },
    {
      id: "span-db2",
      traceId: traceId,
      parentId: "span-pricing1",
      name: "get-price-rules",
      localEndpoint: "database",
      timestamp: 1703664522400,
      duration: 30000,

      tags: {
        "db.query": "SELECT * FROM price_rules",
        "db.table": "price_rules",
      },
    },
    {
      id: "span-notification1",
      traceId: traceId,
      parentId: "span-gateway1",
      name: "send-notification",
      localEndpoint: "notification-service",
      timestamp: 1703664522450,
      duration: 60000,

      tags: {
        "notification.type": "product-view",
        channel: "analytics",
      },
    },
    {
      id: "span-queue1",
      traceId: traceId,
      parentId: "span-notification1",
      name: "queue-message",
      localEndpoint: "message-queue",
      timestamp: 1703664522500,
      duration: 40000,

      tags: {
        "queue.name": "analytics",
        "message.size": "1kb",
      },
    },
    {
      id: "span-analytics1",
      traceId: traceId,
      parentId: "span-queue1",
      name: "process-analytics",
      localEndpoint: "analytics-service",
      timestamp: 1703664522550,
      duration: 25000,

      tags: {
        "analytics.event": "product-view",
        "processing.type": "async",
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
