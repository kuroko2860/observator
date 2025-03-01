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
import CustomContainer from "./CustomContainer";
import { useEffect } from "react";
import useFetchData from "../hook/useFetchData";
import StatCard from "./StatCard";
import BarChartCard from "./BarChartCard";
import { BarChart } from "@mui/x-charts";

const HopDetails = ({ params, setShowHopDetail, unit }) => {
  const { data, loading, error, fetchData } = useFetchData("/hop-detail");
  useEffect(() => {
    fetchData(params);
  }, [params]);

  if (loading) {
    return (
      <CircularProgress className="animate-spin h-5 w-5 border-b-2 border-gray-900 rounded-full" />
    );
  }
  if (error) {
    return <p className="text-red-500">{error.message}</p>;
  }
  return (
    data && (
      <Box className="bg-white p-4 rounded-lg shadow-lg">
        {/* <CustomContainer title={"Hop info"} className="mb-4">
          <Grid2
            container
            className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
          >
            <StatCard
              title="Caller service"
              value={data.HopInfo.caller_service_name}
              className="col-span-1"
            />
            <StatCard
              title="Caller operation"
              value={data.HopInfo.caller_operation_name}
              className="col-span-1"
            />
            <StatCard
              title="Called service"
              value={data.HopInfo.called_service_name}
              className="col-span-1"
            />
            <StatCard
              title="Called operation"
              value={data.HopInfo.called_operation_name}
              className="col-span-1"
            />
          </Grid2>
        </CustomContainer> */}

        <CustomContainer title={"Hop statistic"} className="mb-4">
          <Grid2
            container
            className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
          >
            <StatCard
              title="Count"
              value={data.Count}
              unit="calls"
              className="col-span-1"
            />
            <StatCard
              title={"Frequency"}
              value={data.Frequency}
              unit={`calls/${unit}`}
              className="col-span-1"
            />
            <BarChartCard
              title={"Distribution"}
              caption={"Hop call distribution"}
              className="col-span-2 lg:col-span-1"
            >
              <BarChart
                width={600}
                height={350}
                className="w-full"
                xAxis={[
                  {
                    scaleType: "band",
                    data: Object.keys(data.Distribution || {}).map((item) =>
                      new Date(parseInt(item)).toLocaleString()
                    ),
                    label: "Timestamp",
                    tickPlacement: "start",
                    tickLabelPlacement: "tick",
                  },
                ]}
                series={[
                  {
                    data: Object.values(data.Distribution || {}),
                    label: "Count",
                  },
                ]}
              />
            </BarChartCard>
          </Grid2>
        </CustomContainer>
        <CustomContainer title={"Error"} className="mb-4">
          <Grid2
            container
            className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
          >
            <StatCard
              title={"Error count"}
              value={data.ErrorCount}
              unit={"calls"}
              className="col-span-1"
            />
            <StatCard
              title={"Error rate"}
              value={data.ErrorRate * 100}
              unit={"%"}
              className="col-span-1"
            />
            <BarChartCard
              title={"Error distribution"}
              caption={"error distribution"}
              className="col-span-2 lg:col-span-1"
            >
              <BarChart
                width={600}
                height={350}
                className="w-full"
                xAxis={[
                  {
                    scaleType: "band",
                    data: Object.keys(data.ErrorDist || {}).map((item) =>
                      new Date(parseInt(item)).toLocaleString()
                    ),
                    label: "Timestamp",
                    tickPlacement: "start",
                    tickLabelPlacement: "tick",
                  },
                ]}
                series={[
                  {
                    data: Object.values(data.ErrorDist || {}),
                    label: "Count",
                  },
                ]}
              />
            </BarChartCard>
          </Grid2>
        </CustomContainer>
        <CustomContainer title={"Latency"} className="mb-4">
          <TableContainer component={Paper} className="overflow-x-auto">
            <Table className="min-w-full">
              <TableHead>
                <TableRow>
                  <TableCell>Key</TableCell>
                  <TableCell>Value (microsecond)</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {Object.entries(data.Latency || {}).map(
                  ([key, value], index) => (
                    <TableRow key={index}>
                      <TableCell>{key}</TableCell>
                      <TableCell>{value}</TableCell>
                    </TableRow>
                  )
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CustomContainer>
        <Button
          type="button"
          onClick={() => setShowHopDetail(false)}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Close
        </Button>
      </Box>
    )
  );
};

export default HopDetails;
