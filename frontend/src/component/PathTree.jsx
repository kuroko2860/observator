import { useState } from "react";
import dayjs from "dayjs";
import { Box, Container } from "@mui/material";
import Tree from "react-d3-tree";
import CustomContainer from "./CustomContainer";
import HopDetails from "./HopDetail";

const PathTree = ({ pathTree, from, to, unit }) => {
  const treeData = transformTreeData(pathTree);
  const [showHopDetail, setShowHopDetail] = useState(false);
  const [params, setParams] = useState();
  const handleLinkClick = (source, target) => {
    console.log(source, target);
    const params = {
      from: from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to: to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
      unit: unit || "hour",
      caller_svc: source.data.name,
      caller_op: source.data.attributes.op,
      called_svc: target.data.name,
      called_op: target.data.attributes.op,
    };
    setShowHopDetail(true);
    setParams(params);
  };

  return (
    <Container>
      <Box style={{ width: "100%", height: "70vh" }}>
        <Tree
          data={treeData}
          onLinkClick={handleLinkClick}
          pathClassFunc={() => "custom-link"}
        />
      </Box>
      {showHopDetail && (
        <CustomContainer>
          <HopDetails
            params={params}
            setShowHopDetail={setShowHopDetail}
            unit={unit}
          />
        </CustomContainer>
      )}
    </Container>
  );
};
const transformTreeData = (pathTree) => {
  const MAX_LENGTH = 7;
  const res = {
    name:
      pathTree.service_name.length < MAX_LENGTH
        ? pathTree.service_name
        : pathTree.service_name.slice(0, MAX_LENGTH) + "...",
    attributes: {
      op:
        pathTree.operation_name.length < MAX_LENGTH
          ? pathTree.operation_name
          : pathTree.operation_name.slice(0, MAX_LENGTH) + "...",
    },
    children: [],
  };
  if (pathTree.children) {
    pathTree.children.forEach((child) => {
      res.children.push(transformTreeData(child));
    });
  }
  return res;
};

export default PathTree;
