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
import { TimeRangeInput } from "../component/shared/Input";
import CustomContainer from "../component/shared/CustomContainer";
import BarChartCard from "../component/shared/BarChartCard";
import { BarChart } from "@mui/x-charts";
import { useNavigate, useParams } from "react-router-dom";
import useFetchData from "../hook/useFetchData";
import dayjs from "dayjs";
import { useState } from "react";
import { SubmitButtons } from "../component/shared/Common";

const ServiceDetail = () => {
  // Hooks
  const methods = useForm();
  const navigate = useNavigate();
  const { service_name } = useParams();
  const { data, error, loading, fetchData } = useFetchData(
    `/services/${service_name}`
  );
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);

  // Data processing
  const sortedHttpApi =
    data?.http_api?.sort((a, b) => b.count - a.count) || null;
  const hasOperations = data && Object.keys(data.operations).length > 0;

  // Handlers
  const handleSubmit = (formData) => {
    setPage(0);
    const params = {
      from: formData.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        formData.to?.$d.getTime() ||
        dayjs().startOf("day").add(1, "day").valueOf(),
    };
    fetchData(params);
  };

  const handlePageChange = (_, newPage) => setPage(newPage);

  const handleRowsPerPageChange = (event) => {
    setRowsPerPage(parseInt(event.target.value));
    setPage(0);
  };

  // Chart data preparation
  const operationsData =
    data && Object.entries(data.operations).sort((a, b) => b[1] - a[1]);

  const renderOperationsChart = () => (
    <BarChartCard
      title="Operations"
      caption="Operation call count"
      className="mb-4"
    >
      <BarChart
        layout="horizontal"
        grid={{ vertical: true }}
        width={600}
        height={400}
        yAxis={[
          {
            scaleType: "band",
            tickPlacement: "middle",
            data: operationsData.map(([name]) => name),
            dataKey: "Service name",
          },
        ]}
        series={[
          {
            data: operationsData.map(([, count]) => count),
            label: "Count",
          },
        ]}
      />
    </BarChartCard>
  );

  const renderHttpApiTable = () => (
    <CustomContainer title="HTTP API" className="mt-4">
      <TableContainer component={Paper}>
        <Table aria-label="simple table" className="table-auto min-w-[650px]">
          <TableHead>
            <TableRow>
              <TableCell>STT</TableCell>
              <TableCell>URI Path</TableCell>
              <TableCell>Method</TableCell>
              <TableCell>Count</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sortedHttpApi
              .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
              .map(({ _id, count }, index) => (
                <TableRow
                  key={`${_id.uri_path}-${_id.method}`}
                  className="border-b-0 last:border-b-0"
                  onClick={() =>
                    navigate(
                      `/api-statistics?service_name=${service_name}&uri_path=${_id.uri_path}&method=${_id.method}`
                    )
                  }
                  sx={{ cursor: "pointer" }}
                >
                  <TableCell>{index + 1}</TableCell>
                  <TableCell>{_id.uri_path}</TableCell>
                  <TableCell>{_id.method}</TableCell>
                  <TableCell align="right">{count}</TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        count={sortedHttpApi.length}
        page={page}
        onPageChange={handlePageChange}
        rowsPerPage={rowsPerPage}
        onRowsPerPageChange={handleRowsPerPageChange}
        rowsPerPageOptions={[5, 10, 15]}
        className="mt-4"
      />
    </CustomContainer>
  );

  return (
    <Box className="flex flex-col items-center p-4 max-w-[1200px]">
      <Typography variant="h4" className="text-3xl font-bold mb-4">
        Service Detail: {service_name}
      </Typography>

      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(handleSubmit)}>
          <TimeRangeInput className="mb-4" />
          <SubmitButtons />
        </form>
      </FormProvider>

      {loading && <CircularProgress className="mx-auto" />}
      {error && <p className="text-red-500">{error}</p>}

      {data && (
        <CustomContainer className="mt-4">
          {hasOperations ? (
            renderOperationsChart()
          ) : (
            <Typography variant="h6" className="text-center">
              No operation found
            </Typography>
          )}

          {sortedHttpApi ? (
            renderHttpApiTable()
          ) : (
            <Typography variant="h6" className="text-center">
              No HTTP API found
            </Typography>
          )}
        </CustomContainer>
      )}
    </Box>
  );
};

export default ServiceDetail;
