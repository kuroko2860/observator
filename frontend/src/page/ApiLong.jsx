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
  Alert,
  Tooltip,
  Container,
} from "@mui/material";
import { AccessTime, Speed, Warning } from "@mui/icons-material";
import dayjs from "dayjs";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { SubmitButtons } from "../component/shared/Common";
import { ThresholdInput, TimeRangeInput } from "../component/shared/Input";
import useFetchData from "../hook/useFetchData";
import CustomTable from "../component/shared/CustomTable";
import { useState, useCallback, useMemo } from "react";

const TABLE_HEADINGS = [
  { sortable: false, name: "service_name", label: "Service Name" },
  { sortable: false, name: "uri_path", label: "URI Path" },
  { sortable: false, name: "method", label: "Method" },
  {
    sortable: true,
    name: "count",
    label: "Exceed Count",
    render: (row) => (
      <Chip label={row.count} color="warning" size="small" variant="outlined" />
    ),
  },
  {
    sortable: true,
    name: "avg_latency",
    label: "Average Latency",
    render: (row) => (
      <Tooltip title={`${row.avg_latency} ms`} arrow>
        <Box className="flex items-center gap-1">
          <Typography variant="body2">{row.avg_latency}</Typography>
          <Typography variant="caption" color="text.secondary">
            ms
          </Typography>
        </Box>
      </Tooltip>
    ),
  },
];

const FORM_DEFAULT_VALUES = {
  threshold: 1000, // Default threshold of 1000ms (1 second)
  from: null,
  to: null,
};

const ApiLong = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  const [isSearching, setIsSearching] = useState(false);
  const { data, error, loading, fetchData } = useFetchData(
    "/api-statistics/long"
  );
  const methods = useForm({ defaultValues: FORM_DEFAULT_VALUES });

  const handleRowClick = useCallback(
    ({ service_name, uri_path, method }) => {
      navigate(
        `/api-statistics?service_name=${service_name}&uri_path=${uri_path}&method=${method}`
      );
    },
    [navigate]
  );

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

  const transformData = useCallback((apiData) => {
    if (!apiData || !Array.isArray(apiData)) return [];

    return apiData
      .map(({ _id: api, count, avg_latency }) => ({
        ...api,
        count,
        avg_latency: Math.round(avg_latency),
      }))
      .sort((a, b) => b.count - a.count);
  }, []);

  const transformedData = useMemo(() => {
    return transformData(data);
  }, [data, transformData]);

  return (
    <Container maxWidth="lg" className="py-4 px-2 md:px-4">
      <Card className="p-4 md:p-6 rounded-lg shadow-sm">
        <Box className="flex items-center justify-center gap-2 mb-4">
          <Speed color="warning" fontSize={isMobile ? "medium" : "large"} />
          <Typography
            variant={isMobile ? "h6" : "h5"}
            className="font-bold text-center"
          >
            API latency analysis
          </Typography>
        </Box>

        <Divider className="mb-4" />

        <Box className="my-4 flex flex-col gap-4">
          <Typography variant="body1" className="text-gray-700">
            Find APIs that exceed a specified latency threshold
          </Typography>

          <FormProvider {...methods}>
            <form onSubmit={methods.handleSubmit(handleSubmit)}>
              <Grid2 container spacing={2}>
                <Grid2 item xs={12} md={8}>
                  <TimeRangeInput className="w-full" />
                </Grid2>
                <Grid2 item xs={12} md={4}>
                  <ThresholdInput
                    label="Latency threshold (ms)"
                    className="w-full"
                    helperText="APIs with latency above this value will be shown"
                    min={100}
                    step={100}
                  />
                </Grid2>
                <Grid2 item xs={12}>
                  <SubmitButtons
                    isLoading={isSearching}
                    loadingText="Analyzing..."
                    submitText="Analyze Latency"
                    className="w-full md:w-auto"
                  />
                </Grid2>
              </Grid2>
            </form>
          </FormProvider>
        </Box>

        <Box className="mt-6">
          {loading && !isSearching && (
            <Fade in={true}>
              <Box className="flex justify-center py-8">
                <CircularProgress />
              </Box>
            </Fade>
          )}

          {error && (
            <Alert severity="error" className="mb-4" variant="outlined">
              {error.message || error}
            </Alert>
          )}

          {!loading && !error && data && data.length === 0 && (
            <Alert severity="info" className="mb-4" icon={<AccessTime />}>
              No APIs exceeding the latency threshold were found in the selected
              time range.
            </Alert>
          )}

          {!loading && !error && transformedData.length > 0 && (
            <Box>
              <Typography
                variant="subtitle1"
                className="mb-3 font-medium flex items-center gap-1"
              >
                <Warning fontSize="small" color="warning" />
                Found {transformedData.length} APIs exceeding latency threshold
              </Typography>

              <CustomTable
                headings={TABLE_HEADINGS}
                data={transformedData}
                onRowClick={handleRowClick}
                isLoading={isSearching}
                emptyMessage="No APIs exceeding the latency threshold"
                className="w-full"
              />
            </Box>
          )}
        </Box>
      </Card>
    </Container>
  );
};

export default ApiLong;
