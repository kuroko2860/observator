import { Box, CircularProgress, Grid2, Typography } from "@mui/material";
import dayjs from "dayjs";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { SubmitButtons } from "../component/shared/Common";
import { ThresholdInput, TimeRangeInput } from "../component/shared/Input";
import useFetchData from "../hook/useFetchData";

import CustomContainer from "../component/shared/CustomContainer";
import CustomTable from "../component/shared/CustomTable";

const ApiLong = () => {
  const navigate = useNavigate();
  const { data, error, loading, fetchData } = useFetchData(
    "/api-statistics/long"
  );
  const methods = useForm({
    defaultValues: {
      threshold: null,
      from: null,
      to: null,
    },
  });

  const headings = [
    { sortable: false, name: "service_name", label: "Service Name" },
    { sortable: false, name: "uri_path", label: "URI Path" },
    { sortable: false, name: "method", label: "Method" },
    { sortable: false, name: "count", label: "Exceed count" },
    { sortable: false, name: "avg_latency", label: "Average Latency" },
  ];

  const handleRowClick = (rowData) => {
    const { service_name, uri_path, method } = rowData;
    navigate(
      `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
    );
  };

  const onSubmit = async (data) => {
    const params = {
      ...data,
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    };
    await fetchData(params);
  };

  const transform = (data) => {
    return data
      .map(({ _id: api, count, avg_latency }) => ({
        ...api,
        count,
        avg_latency,
      }))
      .sort((a, b) => b.count - a.count);
  };
  return (
    <Box className="flex flex-col items-center gap-2 p-2">
      <Typography variant="h5" className="text-2xl font-bold">
        View API exceed latency threshold
      </Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2}>
            <TimeRangeInput className="w-full" />
            <ThresholdInput
              label={"Latency (miliseconds)"}
              className="w-full"
            />
            <SubmitButtons className="w-full" />
          </Grid2>
        </form>
      </FormProvider>
      {loading && <CircularProgress className="m-auto" />}
      {error && <div className="text-red-600">{error.message}</div>}
      {data ? (
        <CustomContainer>
          <CustomTable
            headings={headings}
            data={transform(data)}
            onRowClick={handleRowClick}
          />
        </CustomContainer>
      ) : (
        <Typography variant="h5" className="text-lg font-bold">
          No data
        </Typography>
      )}
    </Box>
  );
};

export default ApiLong;
