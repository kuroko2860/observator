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
  Container,
} from "@mui/material";
import { FormProvider, useForm } from "react-hook-form";
import { TimeRangeInput } from "../component/shared/Input";
import CustomContainer from "../component/shared/CustomContainer";
import BarChartCard from "../component/shared/BarChartCard";
import { BarChart } from "@mui/x-charts";
import { useNavigate, useParams } from "react-router-dom";
import useFetchData from "../hook/useFetchData";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
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

  // Fetch data on initial load
  useEffect(() => {
    setPage(0);
    const params = {
      from: dayjs()
        .add(1, "minute")
        .second(0)
        .millisecond(0)
        .subtract(1, "hour")
        .valueOf(),
      to: dayjs().add(1, "minute").second(0).millisecond(0).valueOf(),
    };
    fetchData(params);
  }, []);

  // Handlers
  const handleSubmit = (formData) => {
    setPage(0);
    const params = {
      from:
        formData.from?.$d.getTime() ||
        dayjs()
          .add(1, "minute")
          .second(0)
          .millisecond(0)
          .subtract(1, "hour")
          .valueOf(),
      to:
        formData.to?.$d.getTime() ||
        dayjs().add(1, "minute").second(0).millisecond(0).valueOf(),
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
      title="List of operations"
      caption="Operation call count"
      className="mb-4 "
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
    <CustomContainer
      title="List of HTTP API"
      className="my-4 shadow border-gray-200 border-[1px]"
    >
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
                  <TableCell>{count}</TableCell>
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
    <Box className="flex flex-col gap-4 items-center p-6 max-w-[1200px]">
      <Typography
        variant="h5"
        className="text-xl md:text-2xl font-bold text-center"
      >
        View service{" "}
        <span className="text-[#1976d2] font-bold">{service_name}</span> in
        detail
      </Typography>

      <FormProvider {...methods}>
        <form
          onSubmit={methods.handleSubmit(handleSubmit)}
          className="flex items-center justify-center gap-8"
        >
          <TimeRangeInput />
          <SubmitButtons />
        </form>
      </FormProvider>

      {loading && <CircularProgress className="mx-auto" />}
      {error && <p className="text-red-500">{error}</p>}

      {data && (
        <Container className="mt-4 flex flex-col gap-4">
          {sortedHttpApi ? (
            renderHttpApiTable()
          ) : (
            <Typography variant="h6" className="text-center">
              No HTTP API found
            </Typography>
          )}
          {hasOperations ? (
            renderOperationsChart()
          ) : (
            <Typography variant="h6" className="text-center">
              No operation found
            </Typography>
          )}
        </Container>
      )}
    </Box>
  );
};

export default ServiceDetail;
