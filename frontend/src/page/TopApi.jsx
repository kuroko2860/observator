import {
  Card,
  Container,
  TableCell,
  TableContainer,
  TableSortLabel,
  Typography,
  CircularProgress,
  Table,
  TableBody,
  TableRow,
  TableHead,
  Paper,
} from "@mui/material";
import useFetchData from "../hook/useFetchData";
import { CustomForm } from "../component/shared/Common";
import { TimeRangeInput } from "../component/shared/Input";
import dayjs from "dayjs";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

// Constants
const SORT_ORDERS = {
  COUNT: "count",
  ERROR_COUNT: "err-count",
  ERROR_RATE: "err-rate",
};

const DEFAULT_LIMIT = 10;

// Main component
function TopApi() {
  const topApiFetcher = useFetchData("/api-statistics/top-called");

  const onSubmit = async (data) => {
    const currentDate = dayjs();
    const requestData = {
      ...data,
      limit: DEFAULT_LIMIT,
      from: data.from?.$d.getTime() || currentDate.startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() ||
        currentDate.startOf("day").add(1, "day").valueOf(),
    };
    topApiFetcher.fetchData(requestData);
  };

  return (
    <Container className="flex flex-col gap-4 p-6 bg-white shadow-lg rounded-lg">
      <Typography variant="h5" className="text-gray-900">
        View most called APIs
      </Typography>
      <CustomForm onSubmit={onSubmit} className="flex flex-col gap-4">
        <TimeRangeInput />
      </CustomForm>
      {topApiFetcher.loading && <CircularProgress className="self-center" />}
      {topApiFetcher.error && (
        <div className="text-red-500">{topApiFetcher.error}</div>
      )}
      {topApiFetcher.data && <TopApiCalledTable data={topApiFetcher.data} />}
    </Container>
  );
}

// Table component
const TopApiCalledTable = ({ data }) => {
  const [order, setOrder] = useState("desc");
  const [orderBy, setOrderBy] = useState(SORT_ORDERS.COUNT);
  const navigate = useNavigate();

  const getComparator = () => {
    const sortMultiplier = order === "desc" ? -1 : 1;

    const comparators = {
      [SORT_ORDERS.COUNT]: (a, b) => (a.count - b.count) * sortMultiplier,
      [SORT_ORDERS.ERROR_COUNT]: (a, b) =>
        (a.err_count - b.err_count) * sortMultiplier,
      [SORT_ORDERS.ERROR_RATE]: (a, b) =>
        (a.err_count / a.count - b.err_count / b.count) * sortMultiplier * 100,
    };

    return comparators[orderBy] || comparators[SORT_ORDERS.COUNT];
  };

  const handleSort = (field) => () => {
    const isAsc = orderBy === field && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy(field);
  };

  const renderTableHeader = (label, field) => (
    <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
      <TableSortLabel
        active={orderBy === field}
        direction={orderBy === field ? order : "desc"}
        onClick={handleSort(field)}
        className="text-gray-900"
      >
        {label}
      </TableSortLabel>
    </TableCell>
  );

  const handleRowClick = (service_name, uri_path, method) => {
    navigate(
      `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
    );
  };

  return (
    <Card
      variant="outlined"
      className="h-full flex flex-col gap-4 p-6 bg-white shadow-lg rounded-lg"
    >
      <Typography
        variant="h5"
        textAlign="center"
        className="text-gray-900 font-bold"
      >
        Top 10 called APIs
      </Typography>
      <TableContainer component={Paper} className="mt-4">
        <Table className="min-w-full" aria-label="API statistics table">
          <TableHead>
            <TableRow>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-left">
                STT
              </TableCell>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-left">
                Service
              </TableCell>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
                URI path
              </TableCell>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
                Method
              </TableCell>
              {renderTableHeader("Count", SORT_ORDERS.COUNT)}
              {renderTableHeader("Error Count", SORT_ORDERS.ERROR_COUNT)}
              {renderTableHeader("Error Rate", SORT_ORDERS.ERROR_RATE)}
            </TableRow>
          </TableHead>
          <TableBody>
            {data
              .slice(0, DEFAULT_LIMIT)
              .sort(getComparator())
              .map(
                (
                  { _id: { service_name, uri_path, method }, count, err_count },
                  index
                ) => (
                  <TableRow
                    key={`${service_name}-${uri_path}-${method}`}
                    onClick={() =>
                      handleRowClick(service_name, uri_path, method)
                    }
                    className="hover:bg-gray-100 cursor-pointer"
                  >
                    <TableCell className="px-4 py-2 border-b border-gray-300">
                      {index + 1}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300">
                      {service_name}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300 text-right">
                      {uri_path}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300 text-right">
                      {method}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300 text-right">
                      {count}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300 text-right">
                      {err_count}
                    </TableCell>
                    <TableCell className="px-4 py-2 border-b border-gray-300 text-right">
                      {((err_count / count) * 100).toFixed(2)}%
                    </TableCell>
                  </TableRow>
                )
              )}
          </TableBody>
        </Table>
      </TableContainer>
    </Card>
  );
};

export default TopApi;
