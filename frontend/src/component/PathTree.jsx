import { Box, Container } from "@mui/material";
import CytoscapeComponent from "react-cytoscapejs";

// Styles for cytoscape nodes and edges
const cytoscapeStyles = [
  {
    selector: "node",
    style: {
      label: "data(name)",
      "text-valign": "center",
      "text-halign": "center",
      "background-color": "#0074D9",
      color: "white",
      "font-size": "12px",
      "text-outline-width": 2,
      "text-outline-color": "#0074D9",
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

const PathTree = ({ path, handleLinkClick = () => {} }) => {
  const treeData = transformTreeData(path);

  const handleEdgeClick = (cy) => {
    cy.on("tap", "edge", (event) => handleLinkClick(event));
  };

  return (
    <Container className="h-full flex flex-col">
      <Box className="flex-1">
        <CytoscapeComponent
          elements={treeData}
          stylesheet={cytoscapeStyles}
          cytoscapeConfig={{ wheelSensitivity: 0.1 }}
          {...cytoscapeConfig}
          cy={handleEdgeClick}
        />
      </Box>
    </Container>
  );
};

const transformTreeData = (path) => {
  const nodes = path.operations.map((node) => ({ data: node }));
  const edges = path.hops.map((edge) => ({ data: edge }));

  return [...nodes, ...edges];
};

export default PathTree;
