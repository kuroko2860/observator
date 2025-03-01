import { Box, CircularProgress, Grid2, Typography } from "@mui/material";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { SubmitButtons, TextInput } from "../component/shared/Common";
import CustomContainer from "../component/shared/CustomContainer";
import CustomTable from "../component/shared/CustomTable";
import {
  EndpointInput,
  MethodInput,
  ServiceNameInput,
  TimeRangeInput,
} from "../component/shared/Input";
import axios from "../config/axios";
import useFetchData from "../hook/useFetchData";

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

  const headings = [
    { sortable: false, name: "service_name", label: "Service Name" },
    { sortable: false, name: "uri_path", label: "URI Path" },
    { sortable: false, name: "method", label: "Method" },
    { sortable: false, name: "caller", label: "Caller" },
    { sortable: true, name: "count", label: "Count" },
    { sortable: true, name: "err-count", label: "Error Count" },
    { sortable: true, name: "err-rate", label: "Error Rate" },
  ];

  const [endpoints, setEndpoints] = useState([]);
  const serviceName = methods.watch("service_name");
  const onSubmit = async (data) => {
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
      methods.resetField("uri_path");
      methods.resetField("method");
      fetchEndpointsFromService(serviceName);
    }
  }, [methods, serviceName]);

  const filterData = ({ service_name, uri_path, method }) => {
    return (
      service_name.includes(methods.watch("service_name")) &&
      uri_path.includes(methods.watch("uri_path")) &&
      method.includes(methods.watch("method"))
    );
  };
  const transform = (data) => {
    return data
      .map(({ api, caller, count, err_count }) => ({
        ...api,
        caller,
        count,
        err_count,
        err_rate: ((err_count / count) * 100).toFixed(2),
      }))
      .filter(filterData);
  };

  const handleRowClick = (rowData) => {
    const { service_name, uri_path, method } = rowData;
    navigate(
      `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
    );
  };

  return (
    <Box className="flex flex-col items-center p-4 gap-4 bg-white rounded-lg shadow-md">
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
          <CustomTable
            headings={headings}
            data={transform(data)}
            onRowClick={handleRowClick}
          />
        </CustomContainer>
      )}
    </Box>
  );
}

export default ApiCalled;
