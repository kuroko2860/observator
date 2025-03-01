import { useState } from "react";
import dayjs from "dayjs";
import { Container, Box } from "@mui/material";
import HopDetails from "./HopDetail";
import CytoscapeComponent from "react-cytoscapejs";

const PathTree = ({ pathTree, from, to, unit }) => {
  const treeData = transformTreeData(pathTree);
  const [showHopDetail, setShowHopDetail] = useState(false);
  const [params, setParams] = useState();
  const handleLinkClick = (event, cy) => {
    const edge = event.target;
    const sourceId = edge.data("source");
    const targetId = edge.data("target");

    // Find the source and target nodes
    const source = cy.getElementById(sourceId);
    const target = cy.getElementById(targetId);

    // Get labels from nodes
    const callerOp = source.data("name");
    const callerSvc = source.data("service");
    const calledOp = target.data("name");
    const calledSvc = target.data("service");

    const params = {
      from: from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to: to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
      unit: unit || "hour",
      caller_svc: callerSvc,
      caller_op: callerOp,
      called_svc: calledSvc,
      called_op: calledOp,
    };
    setShowHopDetail(true);
    setParams(params);
  };

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
          cy={(cy) =>
            cy.on("tap", "edge", (event) => handleLinkClick(event, cy))
          }
        />
      </Box>
      {showHopDetail && (
        <HopDetails
          params={params}
          setShowHopDetail={setShowHopDetail}
          unit={unit}
          className="mt-4"
        />
      )}
    </Container>
  );
};
const transformTreeData = (pathTree) => {
  const res = [
    ...pathTree.nodes.map((node) => {
      return {
        data: {
          id: node.id,
          name: node.operation,
          service: node.service,
        },
      };
    }),
    ...pathTree.edges.map((edge) => {
      return {
        data: {
          id: edge.id,
          source: edge.source,
          target: edge.target,
        },
      };
    }),
  ];

  return res;
};

export default PathTree;
