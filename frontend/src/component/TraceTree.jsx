import { Box, Container, Typography } from "@mui/material";
import CytoscapeComponent from "react-cytoscapejs";
import OperationLog from "./OperationLog";
import { useState } from "react";

// Styles for cytoscape nodes and edges
const cytoscapeStyles = [
  {
    selector: "node",
    style: {
      label: "data(name)",
      "text-valign": "center",
      "text-halign": "center",
      "background-color": "#0074D9", // Default color
      color: "white",
      "font-size": "12px",
      "text-outline-width": 2,
      "text-outline-color": "#0074D9",
    },
  },
  {
    selector: "node[?error]", // Select nodes with error property
    style: {
      "background-color": "#D32F2F", // Red color for error nodes
      "text-outline-color": "#D32F2F",
    },
  },
  {
    selector: "edge",
    style: {
      width: 2,
      "line-color": "#aaa",
      "curve-style": "bezier",
      "target-arrow-shape": "triangle",
      "font-size": "10px",
      "text-rotation": "autorotate",
      "text-background-opacity": 1,
      "text-background-color": "white",
      "text-background-padding": "3px",
    },
  },
];

// Configuration for cytoscape component
const cytoscapeConfig = {
  minZoom: 0.5,
  maxZoom: 1.5,
  wheelSensitivity: 0.2,
  style: {
    width: "100%",
    height: "600px",
    border: "1px solid black",
  },
  layout: {
    name: "breadthfirst",
  },
};

const TraceTree = ({ traceId, path, spanErrors, spanIds, spanMap }) => {
  const treeData = transformTreeData(path, spanErrors, spanIds);
  const [selectedSpanId, setSelectedSpanId] = useState(null);
  const handleNodeClick = (cy) => {
    cy.on("tap", "node", (event) => _handleNodeClick(event));
  };

  const _handleNodeClick = (event) => {
    const node = event.target;
    const spanId = node.data("spanId");
    if (spanId === selectedSpanId) {
      setSelectedSpanId(null);
      return;
    }
    setSelectedSpanId(spanId);
  };

  return (
    <Container className="h-full flex">
      <Box className={selectedSpanId ? "w-2/3 pr-2" : "w-full"}>
        <CytoscapeComponent
          elements={treeData}
          stylesheet={cytoscapeStyles}
          cytoscapeConfig={{ wheelSensitivity: 0.1 }}
          {...cytoscapeConfig}
          cy={handleNodeClick}
        />
      </Box>
      {selectedSpanId && (
        <Box className="w-1/3 pr-2">
          <Typography variant="h6">
            Service: {spanMap[selectedSpanId].service}
          </Typography>
          <Typography variant="h6">
            Operation: {spanMap[selectedSpanId].operation}
          </Typography>
          <OperationLog traceId={traceId} spanId={selectedSpanId} />
        </Box>
      )}
    </Container>
  );
};

const transformTreeData = (path, spanErrors, spanIds) => {
  const nodes = path.operations.map((node) => ({
    data: { ...node, error: spanErrors[node.id], spanId: spanIds[node.id] },
  }));
  const edges = path.hops.map((edge) => ({ data: edge }));

  return [...nodes, ...edges];
};

export default TraceTree;
