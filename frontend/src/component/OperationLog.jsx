import { useState, useEffect } from "react";
import axios from "../config/axios";
import { Box, Typography, Paper, CircularProgress } from "@mui/material";
import { styled } from "@mui/material/styles";

const LogContainer = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
  maxHeight: "500px",
  overflowY: "auto",
  fontFamily: "monospace",
  whiteSpace: "pre-wrap",
  backgroundColor: "#f5f5f5",
  marginTop: theme.spacing(2),
}));

const LogEntry = styled(Box)(({ theme, severity }) => ({
  padding: theme.spacing(1),
  marginBottom: theme.spacing(1),
  borderRadius: theme.shape.borderRadius,
  backgroundColor:
    severity === "error"
      ? "#ffebee"
      : severity === "warn"
      ? "#fff8e1"
      : "#e8f5e9",
  borderLeft: `4px solid ${
    severity === "error"
      ? "#f44336"
      : severity === "warn"
      ? "#ff9800"
      : "#4caf50"
  }`,
}));

const OperationLog = ({ traceId, spanId }) => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchLogs = async () => {
      if (!traceId && !spanId) {
        setLogs([]);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        let url;

        // Determine which API endpoint to use based on the provided IDs
        if (traceId && spanId) {
          url = `/logs/elasticsearch/trace/${traceId}/span/${spanId}`;
        } else if (traceId) {
          url = `/logs/elasticsearch/trace/${traceId}`;
        } else if (spanId) {
          url = `/logs/elasticsearch/span/${spanId}`;
        }

        const response = await axios.get(url);
        setLogs(response.data || []);
      } catch (err) {
        console.error("Error fetching logs:", err);
        setError(`Failed to fetch logs: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchLogs();
  }, [traceId, spanId]);

  const getSeverityFromLog = (log) => {
    // Determine log severity based on status code or log level
    if (log.status_code >= 500 || log.uri_path?.includes("internal/error")) {
      return "error";
    } else if (
      log.status_code >= 400 ||
      log.uri_path?.includes("internal/warn")
    ) {
      return "warn";
    }
    return "info";
  };

  const formatLogEntry = (log) => {
    const timestamp = new Date(log.start_time).toISOString();
    const service = log.service_name || "unknown";
    const method = log.method || "";
    const path = log.uri_path || "";
    const status = log.status_code || "";
    const duration = log.duration ? `${log.duration}ms` : "";
    const message = log.error_message || "Request completed";

    return (
      <LogEntry
        key={log.start_time + Math.random()}
        severity={getSeverityFromLog(log)}
      >
        <Typography variant="body2" component="div">
          <strong>Time:</strong> {timestamp}
        </Typography>
        <Typography variant="body2" component="div">
          <strong>Service:</strong> {service} | <strong>Method:</strong>{" "}
          {method} | <strong>Path:</strong> {path}
        </Typography>
        {status && (
          <Typography variant="body2" component="div">
            <strong>Status:</strong> {status} | <strong>Duration:</strong>{" "}
            {duration}
          </Typography>
        )}
        <Typography variant="body2" component="div">
          <strong>Message:</strong> {message}
        </Typography>
        {log.trace_id && (
          <Typography variant="body2" component="div">
            <strong>Trace ID:</strong> {log.trace_id}
          </Typography>
        )}
        {log.span_id && (
          <Typography variant="body2" component="div">
            <strong>Span ID:</strong> {log.span_id}
          </Typography>
        )}
      </LogEntry>
    );
  };

  return (
    <Box sx={{ width: "100%" }}>
      <Typography variant="h6" gutterBottom>
        Operation Logs
      </Typography>

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Typography color="error" sx={{ my: 2 }}>
          {error}
        </Typography>
      ) : logs?.length === 0 ? (
        <Typography sx={{ my: 2 }}>
          {traceId || spanId
            ? "No logs found for the specified trace or span ID."
            : "Please provide a trace ID or span ID to view logs."}
        </Typography>
      ) : (
        <LogContainer>{logs?.map((log) => formatLogEntry(log))}</LogContainer>
      )}
    </Box>
  );
};

export default OperationLog;
