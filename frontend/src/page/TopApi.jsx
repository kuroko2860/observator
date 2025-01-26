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
    <Container>
      <Typography variant="h5">View most called APIs</Typography>
      <CustomForm onSubmit={onSubmit}>
        <TimeRangeInput />
      </CustomForm>
      {topApiFetcher.loading && <CircularProgress />}
      {topApiFetcher.error && <div>{topApiFetcher.error}</div>}
      {topApiFetcher.data && (
        <TopApiCalledTable data={topApiFetcher.data}></TopApiCalledTable>
      )}
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
    <Card variant="outlined" sx={{ height: "100%", flexGrow: 1 }}>
      <Typography variant="h5" textAlign="center">
        Top 10 called APIs
      </Typography>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>STT</TableCell>
              <TableCell>Service</TableCell>
              <TableCell align="right">Endpoint</TableCell>
              <TableCell align="right">Method</TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={orderBy === "count"}
                  direction={orderBy === "count" ? order : "desc"}
                  onClick={handleSortByCount}
                >
                  Count
                </TableSortLabel>
              </TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={orderBy === "err-count"}
                  direction={orderBy === "err-count" ? order : "desc"}
                  onClick={handleSortByErrCount}
                >
                  Error Count
                </TableSortLabel>
              </TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={orderBy === "err-rate"}
                  direction={orderBy === "err-rate" ? order : "desc"}
                  onClick={handleSortByErrRate}
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
                  { _id: { service_name, endpoint, method }, count, err_count },
                  index
                ) => (
                  <TableRow
                    key={index}
                    onClick={() =>
                      navigate(
                        `/api-statistics?service_name=${service_name}&endpoint=${endpoint}&method=${method}`
                      )
                    }
                    sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                  >
                    <TableCell>{index + 1}</TableCell>
                    <TableCell>{service_name}</TableCell>
                    <TableCell align="right">{endpoint}</TableCell>
                    <TableCell align="right">{method}</TableCell>
                    <TableCell align="right">{count}</TableCell>
                    <TableCell align="right">{err_count}</TableCell>
                    <TableCell align="right">
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
