import {
  Box,
  CircularProgress,
  Grid2,
  Typography,
  Card,
  useTheme,
  useMediaQuery,
  Fade,
  Divider,
  Chip,
} from "@mui/material";
import dayjs from "dayjs";
import { useEffect, useState, useMemo, useCallback } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { SubmitButtons, TextInput } from "../component/shared/Common";
import CustomTable from "../component/shared/CustomTable";
import {
  EndpointInput,
  MethodInput,
  ServiceNameInput,
  TimeRangeInput,
} from "../component/shared/Input";
import axios from "../config/axios";
import useFetchData from "../hook/useFetchData";
import { Api, ErrorOutline } from "@mui/icons-material";

// Table headings configuration
const TABLE_HEADINGS = [
  { sortable: false, name: "service_name", label: "Service Name" },
  { sortable: false, name: "uri_path", label: "URI Path" },
  { sortable: false, name: "method", label: "Method" },
  { sortable: false, name: "username", label: "Caller" },
  { sortable: true, name: "count", label: "Count" },
  { sortable: true, name: "err-count", label: "Error Count" },
  {
    sortable: true,
    name: "err-rate",
    label: "Error Rate",
    render: (row) => (
      <Chip
        label={`${row["err-rate"]}%`}
        size="small"
        color={parseFloat(row["err-rate"]) > 5 ? "error" : "success"}
        variant="outlined"
      />
    ),
  },
];

// Default form values
const DEFAULT_FORM_VALUES = {
  username: "",
  from: null,
  to: null,
  service_name: "",
  uri_path: "",
  method: "",
};

function ApiUsage() {
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  const { data, loading, error, fetchData } = useFetchData(
    "/api-statistics/user-called"
  );

  // Form setup
  const methods = useForm({
    defaultValues: DEFAULT_FORM_VALUES,
  });

  const [endpoints, setEndpoints] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  const serviceName = methods.watch("service_name");

  // Handle form submission
  const handleSubmit = useCallback(
    async (formData) => {
      setIsSearching(true);
      try {
        const params = {
          ...formData,
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
        await fetchData(params);
      } finally {
        setIsSearching(false);
      }
    },
    [fetchData]
  );

  // Fetch endpoints when service name changes
  const fetchEndpointsFromService = useCallback(async (service) => {
    try {
      const res = await axios.get(`/services/${service}/endpoints`);
      setEndpoints(["", ...res.data]);
    } catch (error) {
      console.error("Error fetching endpoints:", error);
    }
  }, []);

  useEffect(() => {
    if (serviceName) {
      methods.setValue("uri_path", "");
      methods.setValue("method", "");
      fetchEndpointsFromService(serviceName);
    }
  }, [serviceName, methods, fetchEndpointsFromService]);

  // Filter data based on form values
  const filterData = useCallback(
    (item) => {
      const { service_name, uri_path, method } = item;
      const formServiceName = methods.watch("service_name");
      const formUriPath = methods.watch("uri_path");
      const formMethod = methods.watch("method");

      return (
        service_name.includes(formServiceName) &&
        uri_path.includes(formUriPath) &&
        method.includes(formMethod)
      );
    },
    [methods]
  );

  // Transform API data for display
  const transformedData = useMemo(() => {
    if (!data) return [];

    return data
      .map(({ _id: api, count, err_count }) => ({
        ...api,
        count,
        err_count,
        "err-rate": ((err_count / count) * 100).toFixed(2),
      }))
      .filter(filterData);
  }, [data, filterData]);

  // Handle row click to navigate to API statistics
  const handleRowClick = useCallback(
    (rowData) => {
      const { service_name, uri_path, method } = rowData;
      navigate(
        `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
      );
    },
    [navigate]
  );

  return (
    <Box className="flex flex-col gap-4 p-3 md:p-6">
      <Card className="p-4 md:p-6 rounded-lg shadow-sm">
        <Typography
          variant={isMobile ? "h6" : "h5"}
          className="flex items-center justify-center gap-2 mb-4 font-bold"
        >
          <Api color="primary" />
          API usage analytics
        </Typography>
        <Box className="mt-2">
          <Divider />
        </Box>

        <FormProvider {...methods}>
          <form
            onSubmit={methods.handleSubmit(handleSubmit)}
            className="flex flex-col gap-3 mt-8"
          >
            <Grid2 container spacing={2}>
              <Grid2 item xs={12} md={6}>
                <ServiceNameInput
                  className="mb-2"
                  helperText="Select the service to analyze"
                />
              </Grid2>
              <Grid2 item xs={12} md={6}>
                <EndpointInput
                  endpoints={endpoints}
                  className="mb-2"
                  helperText="Select the API endpoint"
                />
              </Grid2>
              <Grid2 item xs={12} sm={6} md={3}>
                <MethodInput className="mb-2" helperText="HTTP method" />
              </Grid2>
              <Grid2 item xs={12} sm={6} md={3}>
                <TextInput
                  name="username"
                  label="Username"
                  className="mb-2 w-full"
                />
              </Grid2>
              <Grid2 item xs={12} md={6}>
                <TimeRangeInput className="mb-2" />
              </Grid2>
            </Grid2>
            <SubmitButtons
              isLoading={isSearching}
              loadingText="Analyzing..."
              submitText="Analyze API"
              className={"justify-start"}
            />
          </form>
        </FormProvider>

        {/* Loading and Error States */}
        <Fade in={loading && !isSearching}>
          <Box className="flex justify-center my-8">
            <CircularProgress />
          </Box>
        </Fade>

        {error && (
          <Box className="p-4 bg-red-50 border border-red-200 rounded-md my-4">
            <Typography className="text-red-600 flex items-center gap-2">
              <ErrorOutline fontSize="small" />
              {error}
            </Typography>
          </Box>
        )}

        {/* Results Table */}
        {data && !loading && (
          <Box>
            <CustomTable
              headings={TABLE_HEADINGS}
              data={transformedData}
              onRowClick={handleRowClick}
              isLoading={isSearching}
              emptyMessage="No API calls match your search criteria"
              className="w-full"
            />
          </Box>
        )}
      </Card>
    </Box>
  );
}

export default ApiUsage;
