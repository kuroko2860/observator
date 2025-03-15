import {
  Box,
  Button,
  Grid2,
  TableContainer,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Paper,
  CircularProgress,
} from "@mui/material";
import CustomContainer from "./shared/CustomContainer";
import { useEffect } from "react";
import useFetchData from "../hook/useFetchData";
import StatCard from "./shared/StatCard";
import BarChartCard from "./shared/BarChartCard";
import { BarChart } from "@mui/x-charts";

const HopDetails = ({ hopID, params, setShowHopDetail, unit }) => {
  const { data, loading, error, fetchData } = useFetchData(`/hops/${hopID}`);

  useEffect(() => {
    fetchData(params);
  }, [params]);

  const renderBarChart = (data, title, caption, distributionData) => (
    <BarChartCard
      title={title}
      caption={caption}
      className="col-span-2 lg:col-span-1"
    >
      <BarChart
        width={600}
        height={350}
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
            label: "Count",
          },
        ]}
      />
    </BarChartCard>
  );

  const renderLatencyTable = (latencyData) => (
    <TableContainer component={Paper} className="overflow-x-auto">
      <Table className="min-w-full">
        <TableHead>
          <TableRow>
            <TableCell>Key</TableCell>
            <TableCell>Value (microsecond)</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {Object.entries(latencyData || {}).map(([key, value], index) => (
            <TableRow key={index}>
              <TableCell>{key}</TableCell>
              <TableCell>{value}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
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
    <Box className="bg-white p-4 rounded-lg shadow-lg">
      <CustomContainer title="Hop statistic" className="mb-4">
        <Grid2
          container
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
        >
          <StatCard
            title="Count"
            value={data.count}
            unit="calls"
            className="col-span-1"
          />
          <StatCard
            title="Frequency"
            value={data.frequency}
            unit={`calls/${unit}`}
            className="col-span-1"
          />
          {renderBarChart(
            data,
            "Distribution",
            "Hop call distribution",
            data.distribution
          )}
        </Grid2>
      </CustomContainer>

      <CustomContainer title="Error" className="mb-4">
        <Grid2
          container
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
        >
          <StatCard
            title="Error count"
            value={data.error_count}
            unit="calls"
            className="col-span-1"
          />
          <StatCard
            title="Error rate"
            value={data.error_rate * 100}
            unit="%"
            className="col-span-1"
          />
          {renderBarChart(
            data,
            "Error distribution",
            "error distribution",
            data.error_dist
          )}
        </Grid2>
      </CustomContainer>

      <CustomContainer title="Latency" className="mb-4">
        {renderLatencyTable(data.latency)}
      </CustomContainer>

      <Button
        type="button"
        onClick={() => setShowHopDetail(false)}
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
      >
        Close
      </Button>
    </Box>
  );
};

export default HopDetails;
