import { Box, Container } from "@mui/material";
import CytoscapeComponent from "react-cytoscapejs";

const PathTree = ({ path, handleLinkClick = () => {} }) => {
  const treeData = transformTreeData(path);

  return (
    <Container className="h-full flex flex-col">
      <Box className="flex-1">
        <CytoscapeComponent
          elements={treeData}
          minZoom={0.5} // Minimum zoom-out level
          maxZoom={1.5} // Maximum zoom-in level
          wheelSensitivity={0.2} // Adjust scroll wheel zoom speed (default is 1)
          stylesheet={[
            {
              selector: "node",
              style: {
                label: "data(name)", // Display the label from data
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
          ]}
          cytoscapeConfig={{ wheelSensitivity: 0.1 }}
          style={{ width: "100%", height: "600px", border: "1px solid black" }}
          layout={{ name: "breadthfirst" }}
          cy={(cy) => cy.on("tap", "edge", (event) => handleLinkClick(event))}
        />
      </Box>
    </Container>
  );
};
const transformTreeData = (path) => {
  const res = [
    ...path.operations.map((node) => {
      return {
        data: node,
      };
    }),
    ...path.hops.map((edge) => {
      return {
        data: edge,
      };
    }),
  ];

  return res;
};

export default PathTree;
