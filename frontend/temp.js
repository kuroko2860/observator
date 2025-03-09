// Sample trace data sets - in a real app, these could be loaded based on active tab
const traceDataSets = {
  timeline: {
    duration: 112819, // Total duration in ms
    services: [
      {
        name: "WEB CLUSTER",
        spans: [{ id: "12.819", name: "GET", startTime: 0, duration: 112819 }],
      },
      {
        name: "QUICKIE SERVICE",
        spans: [
          {
            id: "1.571",
            name: "gimme_stuff",
            startTime: 1000,
            duration: 5000,
          },
        ],
      },
      {
        name: "WEB SERVER",
        spans: [
          { id: "109.358", name: "GET", startTime: 10000, duration: 99000 },
        ],
      },
      {
        name: "SOME SERVICE",
        spans: [
          { id: "4.346", name: "get", startTime: 20000, duration: 10000 },
        ],
      },
      {
        name: "MEMCACHED",
        spans: [
          { id: "0.279", name: "Get", startTime: 22000, duration: 2000 },
          { id: "0.319", name: "Incr", startTime: 26000, duration: 1500 },
        ],
      },
      {
        name: "BIG ASS SERVICE",
        spans: [
          {
            id: "97.155",
            name: "getThemData",
            startTime: 30000,
            duration: 80000,
          },
        ],
      },
      {
        name: "SERVICE",
        spans: [
          {
            id: "5.718",
            name: "getStuff",
            startTime: 35000,
            duration: 10000,
          },
        ],
      },
      {
        name: "DATAR",
        spans: [
          {
            id: "84.537",
            name: "findThings",
            startTime: 40000,
            duration: 70000,
          },
        ],
      },
      {
        name: "THINGIE",
        spans: [
          { id: "1.435", name: "getMoar", startTime: 95000, duration: 5000 },
        ],
      },
      {
        name: "OTHER DATA SERVICE",
        spans: [
          { id: "11.723", name: "get", startTime: 100000, duration: 10000 },
        ],
      },
      {
        name: "FINAL DATA SERVICE",
        spans: [
          { id: "8.902", name: "get", startTime: 106000, duration: 5000 },
        ],
      },
      {
        name: "MEMCACHED",
        spans: [
          { id: "0.386", name: "Get", startTime: 107000, duration: 1000 },
          { id: "0.322", name: "Get", startTime: 108000, duration: 1000 },
          { id: "0.332", name: "Get", startTime: 109000, duration: 1000 },
        ],
      },
    ],
  },
  latency: {
    duration: 112819,
    services: [
      {
        name: "WEB CLUSTER",
        spans: [
          {
            id: "12.819",
            name: "GET",
            startTime: 0,
            duration: 112819,
            latency: "high",
          },
        ],
      },
      {
        name: "BIG ASS SERVICE",
        spans: [
          {
            id: "97.155",
            name: "getThemData",
            startTime: 30000,
            duration: 80000,
            latency: "high",
          },
        ],
      },
      {
        name: "DATAR",
        spans: [
          {
            id: "84.537",
            name: "findThings",
            startTime: 40000,
            duration: 70000,
            latency: "high",
          },
        ],
      },
      {
        name: "OTHER DATA SERVICE",
        spans: [
          {
            id: "11.723",
            name: "get",
            startTime: 100000,
            duration: 10000,
            latency: "medium",
          },
        ],
      },
      {
        name: "SERVICE",
        spans: [
          {
            id: "5.718",
            name: "getStuff",
            startTime: 35000,
            duration: 10000,
            latency: "medium",
          },
        ],
      },
      {
        name: "QUICKIE SERVICE",
        spans: [
          {
            id: "1.571",
            name: "gimme_stuff",
            startTime: 1000,
            duration: 5000,
            latency: "low",
          },
        ],
      },
      {
        name: "MEMCACHED",
        spans: [
          {
            id: "0.279",
            name: "Get",
            startTime: 22000,
            duration: 2000,
            latency: "low",
          },
          {
            id: "0.319",
            name: "Incr",
            startTime: 26000,
            duration: 1500,
            latency: "low",
          },
        ],
      },
    ],
  },
  errors: {
    duration: 112819,
    services: [
      {
        name: "WEB CLUSTER",
        spans: [
          {
            id: "12.819",
            name: "GET",
            startTime: 0,
            duration: 112819,
            status: "ok",
          },
        ],
      },
      {
        name: "DATAR",
        spans: [
          {
            id: "84.537",
            name: "findThings",
            startTime: 40000,
            duration: 70000,
            status: "error",
          },
        ],
      },
      {
        name: "MEMCACHED",
        spans: [
          {
            id: "0.319",
            name: "Incr",
            startTime: 26000,
            duration: 1500,
            status: "warning",
          },
        ],
      },
      {
        name: "FINAL DATA SERVICE",
        spans: [
          {
            id: "8.902",
            name: "get",
            startTime: 106000,
            duration: 5000,
            status: "warning",
          },
        ],
      },
    ],
  },
  dependencies: {
    duration: 112819,
    services: [
      {
        name: "WEB CLUSTER",
        spans: [
          {
            id: "12.819",
            name: "GET",
            startTime: 0,
            duration: 112819,
            dependencies: ["WEB SERVER"],
          },
        ],
      },
      {
        name: "WEB SERVER",
        spans: [
          {
            id: "109.358",
            name: "GET",
            startTime: 10000,
            duration: 99000,
            dependencies: ["SOME SERVICE", "BIG ASS SERVICE"],
          },
        ],
      },
      {
        name: "BIG ASS SERVICE",
        spans: [
          {
            id: "97.155",
            name: "getThemData",
            startTime: 30000,
            duration: 80000,
            dependencies: ["DATAR"],
          },
        ],
      },
      {
        name: "DATAR",
        spans: [
          {
            id: "84.537",
            name: "findThings",
            startTime: 40000,
            duration: 70000,
            dependencies: ["THINGIE"],
          },
        ],
      },
      {
        name: "THINGIE",
        spans: [
          {
            id: "1.435",
            name: "getMoar",
            startTime: 95000,
            duration: 5000,
            dependencies: ["OTHER DATA SERVICE"],
          },
        ],
      },
      {
        name: "OTHER DATA SERVICE",
        spans: [
          {
            id: "11.723",
            name: "get",
            startTime: 100000,
            duration: 10000,
            dependencies: ["FINAL DATA SERVICE"],
          },
        ],
      },
      {
        name: "FINAL DATA SERVICE",
        spans: [
          {
            id: "8.902",
            name: "get",
            startTime: 106000,
            duration: 5000,
            dependencies: ["MEMCACHED"],
          },
        ],
      },
      {
        name: "MEMCACHED",
        spans: [
          { id: "0.386", name: "Get", startTime: 107000, duration: 1000 },
        ],
      },
    ],
  },
};
