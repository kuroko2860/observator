import {
  Box,
  CircularProgress,
  TableCell,
  Typography,
  Table,
  TableBody,
  TableRow,
  TableContainer,
  TableHead,
  TablePagination,
  Paper,
} from "@mui/material";
import { FormProvider, useForm } from "react-hook-form";
import { TimeRangeInput } from "../component/Input";
import CustomContainer from "../component/CustomContainer";
import BarChartCard from "../component/BarChartCard";
import { BarChart } from "@mui/x-charts";
import { useNavigate, useParams } from "react-router-dom";
import useFetchData from "../hook/useFetchData";
import dayjs from "dayjs";
import { useState } from "react";
import { SubmitButtons } from "../component/Common";
const ServiceDetail = () => {
  const methods = useForm();
  const navigate = useNavigate();
  const { service_name } = useParams();
  const { data, error, loading, fetchData } = useFetchData(
    `/services/${service_name}`
  );
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);

  const onSubmit = (data) => {
    setPg(0);
    const params = {
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    };
    fetchData(params);
  };
  let sortedData = null;
  if (data?.http_api) {
    sortedData = data.http_api.sort((a, b) => b.count - a.count);
  }
  return (
    <Box>
      <Typography variant="h4">Service Detail: {service_name}</Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <TimeRangeInput />
          <SubmitButtons />
        </form>
      </FormProvider>
      {loading && <CircularProgress />}
      {error && <p>{error}</p>}
      {data && (
        <CustomContainer>
          {Object.keys(data.operations).length > 0 ? (
            <BarChartCard title={"Operations"} caption={"Operation call count"}>
              <BarChart
                layout="horizontal"
                grid={{ vertical: true }}
                width={600}
                height={400}
                yAxis={[
                  {
                    scaleType: "band",
                    tickPlacement: "middle",
                    data: Object.entries(data.operations)
                      .sort((a, b) => b[1] - a[1])
                      .map((op) => op[0]),
                    dataKey: "Service name",
                  },
                ]}
                series={[
                  {
                    data: Object.entries(data.operations)
                      .sort((a, b) => b[1] - a[1])
                      .map((op) => op[1]),
                    label: "Count",
                  },
                ]}
              />
            </BarChartCard>
          ) : (
            <Typography variant="h6">No operation found</Typography>
          )}
          {sortedData ? (
            <CustomContainer title={"HTTP API"}>
              <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                  <TableHead>
                    <TableRow>
                      <TableCell>STT</TableCell>
                      <TableCell>Endpoint</TableCell>
                      <TableCell>Method</TableCell>
                      <TableCell>Count</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {sortedData
                      .slice(pg * rpg, pg * rpg + rpg)
                      .map(({ _id, count }, index) => (
                        <TableRow
                          key={index}
                          sx={{
                            "&:last-child td, &:last-child th": { border: 0 },
                          }}
                          onClick={() =>
                            navigate(
                              `/api-statistics?service_name=${service_name}&endpoint=${_id.endpoint}&method=${_id.method}`
                            )
                          }
                        >
                          <TableCell>{index + 1}</TableCell>
                          <TableCell>{_id.endpoint}</TableCell>
                          <TableCell>{_id.method}</TableCell>
                          <TableCell align="right">{count}</TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
              </TableContainer>
              <TablePagination
                count={sortedData.length}
                page={pg}
                onPageChange={(e, page) => setPg(page)}
                rowsPerPage={rpg}
                onRowsPerPageChange={(e) => {
                  setRpg(e.target.value);
                  setPg(0);
                }}
                rowsPerPageOptions={[5, 10, 15]}
              />
            </CustomContainer>
          ) : (
            <Typography variant="h6">No HTTP API found</Typography>
          )}
        </CustomContainer>
      )}
    </Box>
  );
};

export default ServiceDetail;
