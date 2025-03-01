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
import { CustomForm } from "../component/Common";
import { TimeRangeInput } from "../component/Input";
import dayjs from "dayjs";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

function TopApi() {
  const topApiFetcher = useFetchData("/api-statistics/top-called");

  const onSubmit = async (data) => {
    topApiFetcher.fetchData({
      ...data,
      limit: 10,
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    });
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

const TopApiCalledTable = ({ data }) => {
  const [order, setOrder] = useState("desc");
  const [orderBy, setOrderBy] = useState("count");
  const getComparator = () => {
    if (orderBy === "count") {
      return (a, b) => (a.count - b.count) * (order === "desc" ? -1 : 1);
    }
    if (orderBy === "err-count") {
      return (a, b) =>
        (a["err_count"] - b["err_count"]) * (order === "desc" ? -1 : 1);
    }
    if (orderBy === "err-rate") {
      return (a, b) =>
        (((a.err_count / a.count) * 100).toFixed(2) -
          ((b.err_count / b.count) * 100).toFixed(2)) *
        (order === "desc" ? -1 : 1);
    }
  };
  const handleSortByCount = () => {
    const isAsc = orderBy === "count" && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy("count");
  };
  const handleSortByErrCount = () => {
    const isAsc = orderBy === "err-count" && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy("err-count");
  };
  const handleSortByErrRate = () => {
    const isAsc = orderBy === "err-rate" && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy("err-rate");
  };
  const navigate = useNavigate();

  return (
    <Card
      variant="outlined"
      sx={{
        height: "100%",
        flexGrow: 1,
        boxShadow: "0 0 10px rgba(0, 0, 0, 0.1)",
        borderRadius: "0.5rem",
      }}
    >
      <Typography
        variant="h5"
        textAlign="center"
        className="text-gray-900 font-bold"
      >
        Top 10 called APIs
      </Typography>
      <TableContainer component={Paper} className="mt-4">
        <Table
          sx={{ minWidth: 650 }}
          aria-label="simple table"
          className="text-gray-900"
        >
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
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
                <TableSortLabel
                  active={orderBy === "count"}
                  direction={orderBy === "count" ? order : "desc"}
                  onClick={handleSortByCount}
                  className="text-gray-900"
                >
                  Count
                </TableSortLabel>
              </TableCell>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
                <TableSortLabel
                  active={orderBy === "err-count"}
                  direction={orderBy === "err-count" ? order : "desc"}
                  onClick={handleSortByErrCount}
                  className="text-gray-900"
                >
                  Error Count
                </TableSortLabel>
              </TableCell>
              <TableCell className="px-4 py-2 bg-gray-200 border-b border-gray-300 text-right">
                <TableSortLabel
                  active={orderBy === "err-rate"}
                  direction={orderBy === "err-rate" ? order : "desc"}
                  onClick={handleSortByErrRate}
                  className="text-gray-900"
                >
                  Error Rate
                </TableSortLabel>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data
              .slice(0, 10)
              .sort(getComparator())
              .map(
                (
                  { _id: { service_name, uri_path, method }, count, err_count },
                  index
                ) => (
                  <TableRow
                    key={index}
                    onClick={() =>
                      navigate(
                        `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
                      )
                    }
                    sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                    className="hover:bg-gray-100"
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
