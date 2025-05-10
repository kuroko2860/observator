import { Box, Button, Grid2, CircularProgress } from "@mui/material";
import { useEffect } from "react";
import useFetchData from "../hook/useFetchData";
import BarChartCard from "./shared/BarChartCard";
import { BarChart } from "@mui/x-charts";

const HopDetails = ({ hopID, params, setShowHopDetail }) => {
  const { data, loading, error, fetchData } = useFetchData(`/hops/${hopID}`);

  useEffect(() => {
    fetchData(params);
  }, [params]);

  const renderBarChart = (label, title, caption, distributionData) => (
    <BarChartCard
      title={title}
      caption={caption}
      className="col-span-2 lg:col-span-1"
    >
      <BarChart
        width={400}
        height={250}
        className="w-full"
        xAxis={[
          {
            scaleType: "band",
            data: Object.keys(distributionData || {}).map((item) =>
              new Date(parseInt(item)).toLocaleString()
            ),
            label: "Timestamp",
            tickPlacement: "start",
            tickLabelPlacement: "tick",
          },
        ]}
        series={[
          {
            data: Object.values(distributionData || {}),
            label: label,
          },
        ]}
      />
    </BarChartCard>
  );

  if (loading) {
    return (
      <CircularProgress className="animate-spin h-5 w-5 border-b-2 border-gray-900 rounded-full" />
    );
  }

  if (error) {
    return <p className="text-red-500">{error.message}</p>;
  }

  if (!data) return null;

  return (
    <Box className="bg-white p-4 rounded-lg shadow-lg flex flex-wrap gap-2">
      <Grid2
        container
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-[49%]"
      >
        {renderBarChart(
          "Count",
          "Distribution",
          "Hop call distribution",
          data.distribution
        )}
      </Grid2>

      <Grid2
        container
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-[49%]"
      >
        {renderBarChart(
          "% Error",
          "Error distribution",
          "Error distribution",
          data.error_dist
        )}
      </Grid2>
      <Grid2
        container
        className="gap-4 w-full flex flex-col justify-center items-center"
      >
        {renderBarChart(
          "microseconds (us)",
          "Latency distribution",
          "Latency distribution",
          data.latency
        )}
      </Grid2>
      <Button
        type="button"
        onClick={() => setShowHopDetail(false)}
        variant="outlined"
        className="order-2 sm:order-1 w-full sm:w-auto"
        sx={{
          borderRadius: "8px",
          textTransform: "none",
          fontWeight: 500,
        }}
      >
        Close
      </Button>
    </Box>
  );
};

export default HopDetails;
