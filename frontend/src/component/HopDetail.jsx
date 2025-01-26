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
} from "@mui/material";
import CustomContainer from "./CustomContainer";
import { useEffect } from "react";
import useFetchData from "../hook/useFetchData";
import CircularProgess from "./CircularProgess";
import StatCard from "./StatCard";
import BarChartCard from "./BarChartCard";
import { BarChart } from "@mui/x-charts";

const HopDetails = ({ params, setShowHopDetail, unit }) => {
  const { data, loading, error, fetchData } = useFetchData(
    "/path-analystic/hop"
  );
  useEffect(() => {
    fetchData(params);
  }, [fetchData, params]);

  if (loading) {
    return <CircularProgess />;
  }
  if (error) {
    return <p>Error: {error.message}</p>;
  }
  return (
    data && (
      <Box>
        <CustomContainer title={"Hop info"}>
          <Grid2 container>
            <StatCard
              title="Caller service"
              value={data.HopInfo.caller_service_name}
            />
            <StatCard
              title="Caller operation"
              value={data.HopInfo.caller_operation_name}
            />
            <StatCard
              title="Called service"
              value={data.HopInfo.called_service_name}
            />
            <StatCard
              title="Called operation"
              value={data.HopInfo.called_operation_name}
            />
          </Grid2>
        </CustomContainer>

        <CustomContainer title={"Hop statistic"}>
          <Grid2 container>
            <StatCard title="Count" value={data.Count} unit="calls" />
            <StatCard
              title={"Frequency"}
              value={data.Frequency}
              unit={`calls/${unit}`}
            />
            <BarChartCard
              title={"Distribution"}
              caption={"Hop call distribution"}
            >
              <BarChart
                width={600}
                height={350}
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
        <CustomContainer title={"Error"}>
          <Grid2 container>
            <StatCard
              title={"Error count"}
              value={data.ErrorCount}
              unit={"calls"}
            />
            <StatCard
              title={"Error rate"}
              value={data.ErrorRate * 100}
              unit={"%"}
            />
            <BarChartCard
              title={"Error distribution"}
              caption={"error distribution"}
            >
              <BarChart
                width={600}
                height={350}
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
        <CustomContainer title={"Latency"}>
          <TableContainer component={Paper}>
            <Table>
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
        <Button type="button" onClick={() => setShowHopDetail(false)}>
          Close
        </Button>
      </Box>
    )
  );
};

export default HopDetails;
