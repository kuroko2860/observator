import { useState, useEffect, useMemo } from "react";
import {
  Box,
  Typography,
  Tabs,
  Tab,
  Paper,
  Divider,
  Chip,
  CircularProgress,
} from "@mui/material";
import { useParams } from "react-router-dom";
import axios from "../config/axios";
import PathTree from "../component/PathTree";

const TraceDetail = () => {
  const { id } = useParams();
  const [trace, setTrace] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedSpan, setSelectedSpan] = useState(null);
  const [viewMode, setViewMode] = useState("timeline");
  const [collapsedSpans, setCollapsedSpans] = useState({});

  useEffect(() => {
    const fetchTrace = async () => {
      try {
        setLoading(true);

        const { data } = await axios.get(`/traces/${id}`, {
          params: {
            from: 1703664522100,
            to: 1703664522600,
          },
        });

        setTrace(data);
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };

    fetchTrace();
  }, [id]);

  // Process trace data for visualization
  const processedData = useMemo(() => {
    if (!trace || !trace.spans || trace.spans.length === 0) return null;

    // Find the minimum timestamp (start of the trace)
    const minTimestamp = Math.min(
      ...trace.spans.map((span) => span.timestamp * 1000)
    );

    // Calculate total duration - find the max end time of any span
    const maxEndTime = Math.max(
      ...trace.spans.map((span) => span.timestamp * 1000 + span.duration)
    );
    const totalDuration = maxEndTime - minTimestamp;

    // Build span hierarchy
    const spanMap = {};
    trace.spans.forEach((span) => {
      spanMap[span.id] = {
        ...span,
        relativeStartTime: span.timestamp * 1000 - minTimestamp,
        relativeEndTime: span.timestamp * 1000 + span.duration - minTimestamp,
        children: [],
      };
    });

    // Connect parents and children
    const rootSpans = [];
    trace.spans.forEach((span) => {
      if (span.parent_id && spanMap[span.parent_id]) {
        spanMap[span.parent_id].children.push(spanMap[span.id]);
      } else {
        rootSpans.push(spanMap[span.id]);
      }
    });

    // Group spans by service
    const serviceMap = {};
    trace.spans.forEach((span) => {
      if (!serviceMap[span.service]) {
        serviceMap[span.service] = [];
      }
      serviceMap[span.service].push(spanMap[span.id]);
    });

    return {
      totalDuration,
      services: Object.entries(serviceMap).map(([name, spans]) => ({
        name,
        spans: spans.sort((a, b) => a.relativeStartTime - b.relativeStartTime),
      })),
      minTimestamp,
      maxEndTime,
      path: trace.path, // Include path information from your model
      rootSpans,
      spanMap,
    };
  }, [trace]);

  // Toggle span collapse
  const toggleSpanCollapse = (spanId, e) => {
    if (e) e.stopPropagation();
    setCollapsedSpans((prev) => ({
      ...prev,
      [spanId]: !prev[spanId],
    }));
  };

  // Check if a span is visible based on collapsed state
  const isSpanVisible = (span) => {
    if (!span.parent_id) return true;

    // Check if any parent is collapsed
    let currentParentId = span.parent_id;
    while (currentParentId) {
      if (collapsedSpans[currentParentId]) {
        return false;
      }
      const parentSpan = processedData?.spanMap[currentParentId];
      currentParentId = parentSpan?.parent_id;
    }
    return true;
  };

  // Generate color map for services
  const serviceColorMap = useMemo(() => {
    if (!processedData) return {};

    const colorPalette = [
      "#1976D2", // blue
      "#388E3C", // green
      "#D32F2F", // red
      "#7B1FA2", // purple
      "#FFA000", // amber
      "#0097A7", // cyan
      "#C2185B", // pink
      "#5D4037", // brown
      "#455A64", // blue-grey
      "#F57C00", // orange
    ];

    const colorMap = {};
    processedData.services.forEach((service, index) => {
      colorMap[service.name] = colorPalette[index % colorPalette.length];
    });

    return colorMap;
  }, [processedData]);

  // Calculate position and width for a span
  const calculateSpanStyle = (span) => {
    if (!processedData) return {};

    const leftPercent =
      (span.relativeStartTime / processedData.totalDuration) * 100;
    const widthPercent = (span.duration / processedData.totalDuration) * 100;

    return {
      left: `${leftPercent}%`,
      width: `${Math.max(widthPercent, 0.5)}%`, // Ensure minimum width for visibility
    };
  };

  // Format timestamp for display
  const formatTimestamp = (timestamp) => {
    if (!timestamp) return "";
    const date = new Date(timestamp);
    const now = new Date();
    const diff = (now.getTime() - date.getTime()) / 1000;

    if (diff < 60) {
      return "now";
    } else if (diff < 3600) {
      return `${Math.floor(diff / 60)} minutes ago`;
    } else if (diff < 86400) {
      return `${Math.floor(diff / 3600)} hours ago`;
    } else if (diff < 604800) {
      return `${Math.floor(diff / 86400)} days ago`;
    } else {
      return date.toLocaleString();
    }
  };

  // Format duration for display
  const formatDuration = (duration) => {
    if (!duration) return "0ms";
    if (duration < 1000) return `${duration}μs`;
    if (duration < 1000000) return `${(duration / 1000).toFixed(2)}ms`;
    return `${(duration / 1000000).toFixed(2)}s`;
  };

  // Handle span click
  const handleSpanClick = (span) => {
    setSelectedSpan(span);
  };

  // Generate time markers for the timeline
  const timeMarkers = useMemo(() => {
    if (!processedData) return [];

    const markers = [];
    const totalDurationMs = processedData.totalDuration / 1000; // Convert to ms
    const step = Math.max(Math.ceil(totalDurationMs / 6), 1); // Divide into ~6 segments, minimum 1ms

    for (let i = 0; i <= totalDurationMs; i += step) {
      markers.push({
        position: ((i * 1000) / processedData.totalDuration) * 100,
        label: `${i}ms`,
      });
    }

    return markers;
  }, [processedData]);

  // Helper function to render hierarchical spans
  const renderHierarchicalSpans = (spans, depth = 0) => {
    if (!spans) return null;

    return spans.flatMap((span) => {
      const hasChildren = span.children && span.children.length > 0;
      const isCollapsed = collapsedSpans[span.id];

      const rows = [
        <tr
          key={span.id}
          className={`hover:bg-gray-50 cursor-pointer ${
            selectedSpan?.id === span.id ? "bg-blue-50" : ""
          }`}
          onClick={(e) => {
            if (e.target.closest(".collapse-toggle")) {
              // Don't select span when clicking the collapse toggle
              return;
            }
            handleSpanClick(span);
          }}
        >
          <td className="p-2 border-t border-gray-200">
            <Box className="flex items-center">
              <Box style={{ width: `${depth * 20}px` }} />
              {hasChildren && (
                <Box
                  className="collapse-toggle mr-2 cursor-pointer text-gray-500 hover:text-gray-800"
                  onClick={() => toggleSpanCollapse(span.id)}
                >
                  {isCollapsed ? "▶" : "▼"}
                </Box>
              )}
              <Box className="flex-1">
                <Box
                  className="w-2 h-2 rounded-full inline-block mr-2"
                  style={{
                    backgroundColor: span.error
                      ? "#F44336"
                      : serviceColorMap[span.service] || "#1976D2",
                  }}
                />
                {span.operation}
              </Box>
            </Box>
          </td>
          <td className="p-2 border-t border-gray-200">{span.service}</td>
          <td className="p-2 border-t border-gray-200 text-right">
            {formatDuration(span.duration)}
          </td>
          <td className="p-2 border-t border-gray-200 text-right">
            {formatTimestamp(span.timestamp)}
          </td>
          <td className="p-2 border-t border-gray-200 text-center">
            {span.error ? (
              <Chip
                label="Error"
                size="small"
                className="bg-red-100 text-red-800"
              />
            ) : (
              <Chip
                label="OK"
                size="small"
                className="bg-green-100 text-green-800"
              />
            )}
          </td>
        </tr>,
      ];

      // Add children if not collapsed
      if (hasChildren && !isCollapsed) {
        rows.push(...renderHierarchicalSpans(span.children, depth + 1));
      }

      return rows;
    });
  };

  if (loading) {
    return (
      <Box className="flex justify-center p-8">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Typography color="error" className="p-4">
        Error loading trace: {error}
      </Typography>
    );
  }

  if (!processedData) {
    return <Typography className="p-4">No trace data available</Typography>;
  }

  return (
    <Box className="bg-white p-4 rounded-lg shadow-md">
      {/* Trace header */}
      <Box className="mb-4">
        <Typography variant="h5" className="font-bold">
          Trace: {id}
        </Typography>
        <Box className="flex gap-4 mt-2">
          <Typography>
            Duration: {formatDuration(processedData.totalDuration)}
          </Typography>
          <Typography>Services: {processedData.services.length}</Typography>
          <Typography>Total Spans: {trace.spans.length}</Typography>
          {processedData.path && (
            <Typography>Path ID: {processedData.path.path_id}</Typography>
          )}
        </Box>
      </Box>

      {/* View mode tabs */}
      <Tabs
        value={viewMode}
        onChange={(_, newValue) => setViewMode(newValue)}
        className="mb-4"
      >
        <Tab value="timeline" label="Timeline" />
        <Tab value="hierarchy" label="Hierarchy" />
        {processedData.path && <Tab value="graph" label="Graph" />}
      </Tabs>

      {/* Timeline view */}
      {viewMode === "timeline" && (
        <Box className="border rounded-lg overflow-hidden">
          {/* Time markers */}
          <Box className="flex border-b border-gray-200 pl-48 pr-4 relative h-8">
            {timeMarkers.map((marker, idx) => (
              <Box
                key={idx}
                className="absolute text-xs text-gray-500"
                style={{ left: `calc(${marker.position}% + 48px)` }}
              >
                {marker.label}
              </Box>
            ))}
          </Box>

          {/* Flatten the hierarchy for display */}
          {(() => {
            // Create a flat list of all spans with their hierarchy info
            const allSpans = [];

            // Helper function to recursively add spans to the flat list
            const addSpansToList = (spans, depth = 0) => {
              if (!spans) return;

              spans.forEach((span) => {
                // Only add visible spans
                if (isSpanVisible(span)) {
                  allSpans.push({
                    ...span,
                    depth,
                  });
                }

                // Add children if not collapsed
                if (
                  span.children &&
                  span.children.length > 0 &&
                  !collapsedSpans[span.id]
                ) {
                  addSpansToList(span.children, depth + 1);
                }
              });
            };

            // Start with root spans
            addSpansToList(processedData.rootSpans);

            // Render each span in its own row
            return allSpans.map((span, idx) => {
              const hasChildren = span.children && span.children.length > 0;
              const isCollapsed = collapsedSpans[span.id];

              return (
                <Box key={`span-${span.id}-${idx}`} className="flex">
                  <Box
                    className="w-48 p-2 border-b border-gray-200 flex items-center"
                    style={{
                      backgroundColor: serviceColorMap[span.service],
                      color: "white",
                    }}
                  >
                    {span.service}
                  </Box>
                  <Box className="flex-grow relative border-b border-gray-200 h-10">
                    <Box
                      className="absolute rounded text-xs text-white px-1 overflow-hidden whitespace-nowrap cursor-pointer hover:opacity-80 flex items-center"
                      style={{
                        ...calculateSpanStyle(span),
                        backgroundColor: span.error
                          ? "#F44336"
                          : serviceColorMap[span.service],
                        height: "24px",
                        top: "3px",
                        marginLeft: `${span.depth * 20}px`, // Indent based on hierarchy
                      }}
                      onClick={() => handleSpanClick(span)}
                      title={`${span.operation} (${formatDuration(
                        span.duration
                      )})`}
                    >
                      {hasChildren && (
                        <Box
                          className="mr-1 text-white cursor-pointer"
                          onClick={(e) => {
                            e.stopPropagation();
                            toggleSpanCollapse(span.id);
                          }}
                        >
                          {isCollapsed ? "▶" : "▼"}
                        </Box>
                      )}
                      {span.operation}
                    </Box>
                  </Box>
                </Box>
              );
            });
          })()}
        </Box>
      )}

      {/* Hierarchical view */}
      {viewMode === "hierarchy" && (
        <Box className="border rounded-lg overflow-auto">
          <table className="min-w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-2 text-left">Span</th>
                <th className="p-2 text-left">Service</th>
                <th className="p-2 text-right">Duration</th>
                <th className="p-2 text-right">Start Time</th>
                <th className="p-2 text-center">Status</th>
              </tr>
            </thead>
            <tbody>{renderHierarchicalSpans(processedData.rootSpans)}</tbody>
          </table>
        </Box>
      )}

      {/* Path Info view */}
      {viewMode === "graph" && processedData.path && (
        <Box className="border rounded-lg p-4">
          <PathTree path={processedData.path} />
        </Box>
      )}

      {/* Selected span details */}
      {selectedSpan && (
        <Paper className="mt-4 p-4">
          <Typography variant="h6" className="font-bold mb-2">
            Span Details
          </Typography>
          <Box className="grid grid-cols-2 gap-4">
            <Box>
              <Typography variant="subtitle2" className="text-gray-500">
                ID
              </Typography>
              <Typography>{selectedSpan.id}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" className="text-gray-500">
                Operation
              </Typography>
              <Typography>{selectedSpan.operation}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" className="text-gray-500">
                Service
              </Typography>
              <Typography>{selectedSpan.service || "unknown"}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" className="text-gray-500">
                Duration
              </Typography>
              <Typography>{formatDuration(selectedSpan.duration)}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" className="text-gray-500">
                Start Time
              </Typography>
              <Typography>{formatTimestamp(selectedSpan.timestamp)}</Typography>
            </Box>
            {selectedSpan.parent_id && (
              <Box>
                <Typography variant="subtitle2" className="text-gray-500">
                  Parent ID
                </Typography>
                <Typography>{selectedSpan.parent_id}</Typography>
              </Box>
            )}
          </Box>

          {/* Tags */}
          {selectedSpan.tags && Object.keys(selectedSpan.tags).length > 0 && (
            <>
              <Divider className="my-3" />
              <Typography variant="subtitle1" className="font-bold mb-2">
                Tags
              </Typography>
              <Box className="flex flex-wrap gap-2">
                {Object.entries(selectedSpan.tags).map(([key, value], idx) => (
                  <Chip
                    key={idx}
                    label={`${key}: ${value}`}
                    size="small"
                    className="bg-gray-100"
                  />
                ))}
              </Box>
            </>
          )}

          {/* Annotations */}
          {selectedSpan.annotations && selectedSpan.annotations.length > 0 && (
            <>
              <Divider className="my-3" />
              <Typography variant="subtitle1" className="font-bold mb-2">
                Annotations
              </Typography>
              <Box>
                {selectedSpan.annotations.map((anno, idx) => (
                  <Box key={idx} className="mb-2">
                    <Typography variant="body2" className="text-gray-500">
                      {formatTimestamp(anno.timestamp)}
                    </Typography>
                    <Typography>{anno.value}</Typography>
                  </Box>
                ))}
              </Box>
            </>
          )}
        </Paper>
      )}
    </Box>
  );
};

export default TraceDetail;
