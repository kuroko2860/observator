// Mock data generator for Zipkin traces
export const generateMockTrace = (traceId) => {
  // Generate a random timestamp within the last hour
  const baseTimestamp = Date.now() - Math.floor(Math.random() * 3600000);

  // Services in the trace
  const services = [
    "frontend",
    "api-gateway",
    "auth-service",
    "user-service",
    "product-service",
    "database",
  ];

  // Generate a unique path ID
  const pathID = `path-${Math.random().toString(36).substring(2, 10)}`;

  // Create root span (usually the frontend or entry point)
  const rootSpan = {
    id: `span-${Math.random().toString(36).substring(2, 10)}`,
    traceId: traceId,
    parentId: null,
    operation: "GET /products",
    service: services[0],
    timestamp: baseTimestamp,
    duration: 350000, // 350ms total
    error: false,
    tags: {
      "http.method": "GET",
      "http.path": "/products",
      component: "web",
    },
  };

  // Create child spans
  const spans = [rootSpan];

  // API Gateway span
  const gatewaySpan = createChildSpan(rootSpan, {
    operation: "route-request",
    service: services[1],
    duration: 320000,
    tags: {
      component: "proxy",
      "http.path": "/api/products",
    },
  });
  spans.push(gatewaySpan);

  // Auth service span
  const authSpan = createChildSpan(gatewaySpan, {
    operation: "validate-token",
    service: services[2],
    duration: 45000,
    tags: {
      "auth.type": "jwt",
      "user.id": "user-123",
    },
  });
  spans.push(authSpan);

  // User service span
  const userSpan = createChildSpan(gatewaySpan, {
    operation: "get-user-preferences",
    service: services[3],
    duration: 65000,
    tags: {
      "user.id": "user-123",
    },
  });
  spans.push(userSpan);

  // Product service span (with error)
  const productSpan = createChildSpan(gatewaySpan, {
    operation: "list-products",
    service: services[4],
    duration: 180000,
    error: Math.random() > 0.7, // 30% chance of error
    tags: {
      "result.count": "42",
      "query.filter": "category=electronics",
    },
  });
  spans.push(productSpan);

  // Database spans
  const dbSpan1 = createChildSpan(userSpan, {
    operation: "SELECT users",
    service: services[5],
    duration: 30000,
    tags: {
      "db.statement": "SELECT * FROM users WHERE id = ?",
      "db.type": "mysql",
    },
  });
  spans.push(dbSpan1);

  const dbSpan2 = createChildSpan(productSpan, {
    operation: "SELECT products",
    service: services[5],
    duration: 120000,
    error: productSpan.error, // Inherit error state from parent
    tags: {
      "db.statement": "SELECT * FROM products WHERE category = ? LIMIT 50",
      "db.type": "mysql",
    },
  });
  spans.push(dbSpan2);

  // Add some random small spans
  for (let i = 0; i < 5; i++) {
    const parentSpan = spans[Math.floor(Math.random() * spans.length)];
    const randomService = services[Math.floor(Math.random() * services.length)];
    const randomDuration = Math.floor(Math.random() * 50000) + 5000;

    spans.push(
      createChildSpan(parentSpan, {
        operation: `operation-${i}`,
        service: randomService,
        duration: randomDuration,
        tags: {
          "custom.tag": `value-${i}`,
        },
      })
    );
  }

  // Create path data
  const operations = services.map((service, index) => ({
    id: `op-${index}`,
    service: service,
    name: `Operation in ${service}`,
  }));

  const hops = [];
  for (let i = 0; i < services.length - 1; i++) {
    hops.push({
      id: `hop-${i}`,
      source: `op-${i}`,
      target: `op-${i + 1}`,
    });
  }

  return {
    spans,
    path: {
      pathID,
      createdAt: new Date().toISOString(),
      longestChain: 4,
      longestErrorChain: productSpan.error ? 2 : 0,
      operations,
      hops,
    },
  };
};

// Helper function to create a child span
function createChildSpan(parentSpan, options) {
  const offset = Math.floor(Math.random() * 20000); // Random offset up to 20ms

  return {
    id: `span-${Math.random().toString(36).substring(2, 10)}`,
    traceId: parentSpan.traceId,
    parentId: parentSpan.id,
    operation: options.operation,
    service: options.service,
    timestamp: parentSpan.timestamp + offset,
    duration: options.duration || 50000,
    error: options.error || false,
    tags: options.tags || {},
  };
}

// Generate a mock trace with a specific ID
export const getMockTrace = (traceId) => {
  return generateMockTrace(traceId);
};

export const sampleTrace = {
  spans: [
    {
      id: "span-root",
      traceId: "trace-xyz789",
      parentId: null,
      operation: "GET /api/products",
      service: "frontend",
      timestamp: 1703664522000,
      duration: 450000,
      error: false,
      tags: {
        "http.method": "GET",
        "http.path": "/api/products",
        component: "web",
      },
    },
    {
      id: "span-gateway1",
      traceId: "trace-xyz789",
      parentId: "span-root",
      operation: "route-request",
      service: "api-gateway",
      timestamp: 1703664522050,
      duration: 400000,
      error: false,
      tags: {
        "gateway.route": "/products",
        "request.id": "req-123",
      },
    },
    {
      id: "span-auth1",
      traceId: "trace-xyz789",
      parentId: "span-gateway1",
      operation: "authenticate",
      service: "auth-service",
      timestamp: 1703664522100,
      duration: 150000,
      error: false,
      tags: {
        "auth.type": "jwt",
        "user.id": "user-456",
      },
    },
    {
      id: "span-product1",
      traceId: "trace-xyz789",
      parentId: "span-gateway1",
      operation: "get-products",
      service: "product-service",
      timestamp: 1703664522150,
      duration: 200000,
      error: false,
      tags: {
        "product.category": "electronics",
        "query.limit": "50",
      },
    },
    {
      id: "span-cache1",
      traceId: "trace-xyz789",
      parentId: "span-product1",
      operation: "check-cache",
      service: "cache-service",
      timestamp: 1703664522200,
      duration: 50000,
      error: false,
      tags: {
        "cache.key": "products:electronics",
        "cache.hit": "false",
      },
    },
    {
      id: "span-db1",
      traceId: "trace-xyz789",
      parentId: "span-product1",
      operation: "query-database",
      service: "database",
      timestamp: 1703664522250,
      duration: 100000,
      error: false,
      tags: {
        "db.query": "SELECT * FROM products",
        "db.rows": "42",
      },
    },
    {
      id: "span-inventory1",
      traceId: "trace-xyz789",
      parentId: "span-product1",
      operation: "check-inventory",
      service: "inventory-service",
      timestamp: 1703664522300,
      duration: 75000,
      error: false,
      tags: {
        "inventory.check": "bulk",
        "products.count": "42",
      },
    },
    {
      id: "span-pricing1",
      traceId: "trace-xyz789",
      parentId: "span-product1",
      operation: "calculate-prices",
      service: "pricing-service",
      timestamp: 1703664522350,
      duration: 80000,
      error: false,
      tags: {
        "pricing.type": "dynamic",
        currency: "USD",
      },
    },
    {
      id: "span-db2",
      traceId: "trace-xyz789",
      parentId: "span-pricing1",
      operation: "get-price-rules",
      service: "database",
      timestamp: 1703664522400,
      duration: 30000,
      error: false,
      tags: {
        "db.query": "SELECT * FROM price_rules",
        "db.table": "price_rules",
      },
    },
    {
      id: "span-notification1",
      traceId: "trace-xyz789",
      parentId: "span-gateway1",
      operation: "send-notification",
      service: "notification-service",
      timestamp: 1703664522450,
      duration: 60000,
      error: false,
      tags: {
        "notification.type": "product-view",
        channel: "analytics",
      },
    },
    {
      id: "span-queue1",
      traceId: "trace-xyz789",
      parentId: "span-notification1",
      operation: "queue-message",
      service: "message-queue",
      timestamp: 1703664522500,
      duration: 40000,
      error: false,
      tags: {
        "queue.name": "analytics",
        "message.size": "1kb",
      },
    },
    {
      id: "span-analytics1",
      traceId: "trace-xyz789",
      parentId: "span-queue1",
      operation: "process-analytics",
      service: "analytics-service",
      timestamp: 1703664522550,
      duration: 25000,
      error: false,
      tags: {
        "analytics.event": "product-view",
        "processing.type": "async",
      },
    },
  ],
  path: {
    pathID: "path-123abc",
    createdAt: "2023-12-27T12:15:22.000Z",
    longestChain: 4,
    longestErrorChain: 0,
    operations: [
      {
        id: "op-1",
        service: "frontend",
        name: "Frontend Request",
      },
      {
        id: "op-2",
        service: "api-gateway",
        name: "API Gateway",
      },
      {
        id: "op-3",
        service: "product-service",
        name: "Product Service",
      },
      {
        id: "op-4",
        service: "database",
        name: "Database Operations",
      },
    ],
    hops: [
      {
        id: "hop-1",
        source: "op-1",
        target: "op-2",
      },
      {
        id: "hop-2",
        source: "op-2",
        target: "op-3",
      },
      {
        id: "hop-3",
        source: "op-3",
        target: "op-4",
      },
    ],
  },
};
