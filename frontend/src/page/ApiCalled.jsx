import { useNavigate } from "react-router-dom";
import useFetchData from "../hook/useFetchData";
import { FormProvider, useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import dayjs from "dayjs";
import axios from "../config/axios";
import {
  EndpointInput,
  MethodInput,
  ServiceNameInput,
  TimeRangeInput,
} from "../component/Input";
import {
  Box,
  Grid2,
  Table,
  TableBody,
  TableHead,
  TableRow,
  TableSortLabel,
  Paper,
  TableCell,
  TableContainer,
  TablePagination,
  Typography,
  CircularProgress,
} from "@mui/material";
import CustomContainer from "../component/CustomContainer";
import { SubmitButtons, TextInput } from "../component/Common";

function ApiCalled() {
  const navigate = useNavigate();
  const { data, loading, error, fetchData } = useFetchData(
    "/api-statistics/called"
  );
  const methods = useForm({
    defaultValues: {
      username: "",
      from: null,
      to: null,
      service_name: "",
      uri_path: "",
      method: "",
    },
  });
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  const [endpoints, setEndpoints] = useState([]);
  const serviceName = methods.watch("service_name");
  const onSubmit = async (data) => {
    setPg(0);
    const params = {
      ...data,
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    };
    await fetchData(params);
  };
  const fetchEndpointsFromService = async (service) => {
    try {
      const res = await axios.get(`/services/${service}/endpoints`);
      setEndpoints(["", ...res.data]);
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    if (serviceName) {
      setPg(0);
      methods.setValue("uri_path", "");
      methods.setValue("method", "");
      fetchEndpointsFromService(serviceName);
    }
  }, [serviceName]);
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
  const filterData = ({ _id: api }) => {
    const { service_name, uri_path, method } = api;
    return (
      service_name.includes(methods.watch("service_name")) &&
      uri_path.includes(methods.watch("uri_path")) &&
      method.includes(methods.watch("method"))
    );
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        padding: "1rem",
        gap: "1rem",
      }}
      className="bg-white rounded-lg shadow-md p-4"
    >
      <Typography variant="h5" className="text-2xl font-bold">
        View called API by user
      </Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2} className="grid grid-cols-1 gap-4">
            <TimeRangeInput />
            <TextInput name="username" label="Username" />
            <ServiceNameInput />
            <EndpointInput endpoints={endpoints} />
            <MethodInput />

            <SubmitButtons />
          </Grid2>
        </form>
      </FormProvider>
      {loading && <CircularProgress className="mx-auto" />}
      {error && <Typography className="text-red-500">{error}</Typography>}
      {data && (
        <CustomContainer className="overflow-x-auto">
          <TableContainer component={Paper} className="rounded-lg">
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell className="px-6 py-3">Service name</TableCell>
                  <TableCell className="px-6 py-3">URI Path</TableCell>
                  <TableCell className="px-6 py-3">Method</TableCell>
                  <TableCell className="px-6 py-3">User</TableCell>
                  <TableCell key={"count"} className="px-6 py-3">
                    <TableSortLabel
                      active={orderBy === "count"}
                      direction={orderBy === "count" ? order : "desc"}
                      onClick={handleSortByCount}
                    >
                      Count
                    </TableSortLabel>
                  </TableCell>
                  <TableCell key={"err-count"} className="px-6 py-3">
                    <TableSortLabel
                      active={orderBy === "err-count"}
                      direction={orderBy === "err-count" ? order : "desc"}
                      onClick={handleSortByErrCount}
                    >
                      Error count
                    </TableSortLabel>
                  </TableCell>
                  <TableCell key={"err-rate"} className="px-6 py-3">
                    <TableSortLabel
                      active={orderBy === "err-rate"}
                      direction={orderBy === "err-rate" ? order : "desc"}
                      onClick={handleSortByErrRate}
                    >
                      Error rate
                    </TableSortLabel>
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {data
                  .filter(filterData)
                  .sort(getComparator())
                  .slice(pg * rpg, pg * rpg + rpg)
                  .map(({ _id: api, count, err_count }, index) => (
                    <TableRow
                      key={index}
                      className="hover:bg-gray-100"
                      onClick={() =>
                        navigate(
                          `/api-statistics?service_name=${api.service_name}&uri_path=${api.uri_path}&method=${api.method}`
                        )
                      }
                    >
                      <TableCell className="px-6 py-3">
                        {api.service_name}
                      </TableCell>
                      <TableCell className="px-6 py-3">
                        {api.uri_path}
                      </TableCell>
                      <TableCell className="px-6 py-3">{api.method}</TableCell>
                      <TableCell className="px-6 py-3">
                        {api.username}
                      </TableCell>
                      <TableCell className="px-6 py-3">{count}</TableCell>
                      <TableCell className="px-6 py-3">{err_count}</TableCell>
                      <TableCell className="px-6 py-3">
                        {((err_count / count) * 100).toFixed(2)}%
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
          <TablePagination
            count={data.filter(filterData).length}
            onPageChange={(e, pg) => setPg(pg)}
            onRowsPerPageChange={(e) => {
              setRpg(parseInt(e.target.value), 10);
              setPg(0);
            }}
            page={pg}
            rowsPerPage={rpg}
            rowsPerPageOptions={[5, 10, 25]}
          />
        </CustomContainer>
      )}
    </Box>
  );
}

export default ApiCalled;
