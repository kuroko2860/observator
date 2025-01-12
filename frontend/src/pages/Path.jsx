import { useEffect, useRef, useState } from "react";
import ForceGraph2D from "react-force-graph-2d";
import { useParams } from "react-router-dom";

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
        // const dagData = {
        //   nodes: [
        //     { id: "A", label: "Node A" },
        //     { id: "B", label: "Node B" },
        //     { id: "C", label: "Node C" },
        //     { id: "D", label: "Node D" },
        //   ],
        //   links: [
        //     { source: "A", target: "B" },
        //     { source: "A", target: "C" },
        //     { source: "B", target: "D" },
        //     { source: "C", target: "D" },
        //   ],
        // };
        // setGraphData(dagData);
      } catch (error) {
        console.error("Error fetching graph data:", error);
      }
    };

    fetchGraphData();
  }, []);
  const fgRef = useRef();

  useEffect(() => {
    if (fgRef.current) {
      fgRef.current.d3Force("charge").strength(-300); // Adjust repulsion strength
    }
  }, []);

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ForceGraph2D
        ref={fgRef}
        graphData={graphData}
        nodeLabel="label"
        nodeAutoColorBy="id"
        linkDirectionalArrowLength={6}
        linkDirectionalArrowRelPos={1} // Position arrow at the end
        linkWidth={2}
        linkDirectionalParticles={1}
        linkDirectionalParticleSpeed={0.005}
        onNodeClick={(node) => alert(`Clicked node: ${node.label}`)}
        onLinkClick={(link) =>
          alert(`Clicked link: ${link.source.id} -> ${link.target.id}`)
        }
      />
    </div>
  );
};

export default GraphVisualizer;
