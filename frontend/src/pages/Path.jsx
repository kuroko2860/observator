import { useEffect, useState } from "react";
import ForceGraph2D from "react-force-graph-2d";
import { useParams } from "react-router-dom";
import { CustomDateTimeRangePicker } from "../components/CustomTimePicker";

const GraphVisualizer = () => {
  const [graphData, setGraphData] = useState({ nodes: [], links: [] });
  const { id } = useParams();
  useEffect(() => {
    const fetchGraphData = async () => {
      try {
        const response = await fetch(`http://localhost:8585/api/paths/${id}`);
        const data = await response.json();

        // Prepare data for react-force-graph
        const nodes = [
          ...new Map(data.nodes.map((item) => [item["id"], item])).values(),
        ]; // Remove duplicate nodes
        const links = data.edges.map((edge) => ({
          source: edge.source,
          target: edge.target,
        }));

        setGraphData({ nodes, links });
      } catch (error) {
        console.error("Error fetching graph data:", error);
      }
    };

    fetchGraphData();
  }, []);

  return (
    <div style={{ width: "600px", height: "100vh" }}>
      <CustomDateTimeRangePicker />
      <ForceGraph2D
        graphData={graphData}
        nodeLabel="label"
        nodeAutoColorBy="color"
        nodeCanvasObject={(node, ctx, globalScale) => {
          const label = node.label || node.id; // Use label if available, fallback to id
          const fontSize = 12 / globalScale; // Adjust font size based on zoom level
          ctx.font = `${fontSize}px Sans-Serif`;
          ctx.textAlign = "center";
          ctx.textBaseline = "middle";

          // Draw node circle
          ctx.fillStyle = node.color;
          ctx.beginPath();
          ctx.arc(node.x, node.y, 5, 0, 2 * Math.PI, false); // Node radius is 5
          ctx.fill();

          // Draw text in the center of the node
          ctx.fillStyle = "black"; // Text color
          ctx.fillText(label, node.x, node.y); // Text position at node's center
        }}
        linkDirectionalArrowLength={5}
        linkDirectionalArrowRelPos={1}
        linkWidth={2}
        linkColor={() => "#aaa"}
        onNodeClick={(node) => alert(`Clicked node: ${node.label || node.id}`)}
        onLinkClick={(link) =>
          alert(`Clicked link: ${link.source.id} -> ${link.target.id}`)
        }
      />
    </div>
  );
};

export default GraphVisualizer;
