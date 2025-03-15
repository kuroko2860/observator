import {
  AccessTime,
  Api,
  BarChart as BarChartIcon,
  ErrorOutline,
  ExpandMore,
  Refresh,
  Speed,
} from "@mui/icons-material";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Card,
  Chip,
  CircularProgress,
  Divider,
  Fade,
  Grid2,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { BarChart } from "@mui/x-charts";
import dayjs from "dayjs";
import { useCallback, useEffect, useMemo, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { useLocation } from "react-router-dom";
import BarChartCard from "../component/shared/BarChartCard";
import { SubmitButtons } from "../component/shared/Common";
import CustomContainer from "../component/shared/CustomContainer";
import {
  EndpointInput,
  MethodInput,
  ServiceNameInput,
  TimeRangeInput,
  TimeUnitInput,
} from "../component/shared/Input";
import StatCard from "../component/shared/StatCard";
import axios from "../config/axios";
import { ApiStatisticDefault } from "../config/default";
import useFetchData from "../hook/useFetchData";

const ApiStatistic = ({ defaultValue = ApiStatisticDefault }) => {
  const { search } = useLocation();
  const params = new URLSearchParams(search);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));
  const isTablet = useMediaQuery(theme.breakpoints.down("md"));

  const [isSearching, setIsSearching] = useState(false);
  const [chartDimensions, setChartDimensions] = useState({
    width: 500,
    height: 300,
  });

  const methods = useForm({
    defaultValues: {
      service_name: params.get("service_name") || "",
      uri_path: params.get("uri_path") || "",
      method: params.get("method") || "",
      unit: "hour",
      from: null,
      to: null,
    },
  });

  const apiFetcher = useFetchData("/api-statistics");
  const [endpoints, setEndpoints] = useState([]);
  const serviceName = methods.watch("service_name");
  const timeUnit = methods.watch("unit");

  // Update chart dimensions based on window size
  useEffect(() => {
    const updateDimensions = () => {
      const containerWidth =
        document.getElementById("chart-container")?.clientWidth || 500;
      setChartDimensions({
        width: Math.min(containerWidth - 40, 800),
        height: isMobile ? 250 : isTablet ? 300 : 350,
      });
    };

    updateDimensions();
    window.addEventListener("resize", updateDimensions);
    return () => window.removeEventListener("resize", updateDimensions);
  }, [isMobile, isTablet]);

  const fetchEndpointsFromService = useCallback(async (service) => {
    try {
      const res = await axios.get(`/services/${service}/endpoints`);
      setEndpoints(["", ...res.data]);
    } catch (error) {
      console.error("Failed to fetch endpoints:", error);
    }
  }, []);

  useEffect(() => {
    if (serviceName) {
      fetchEndpointsFromService(serviceName);
    }
  }, [serviceName, fetchEndpointsFromService]);

  const onSubmit = useCallback(
    async (data) => {
      setIsSearching(true);
      try {
        const currentDate = dayjs();
        const requestData = {
          ...defaultValue,
          ...data,
          from: data.from?.$d.getTime() || currentDate.startOf("day").valueOf(),
          to:
            data.to?.$d.getTime() ||
            currentDate.startOf("day").add(1, "day").valueOf(),
        };
        await apiFetcher.fetchData(requestData);
      } finally {
        setIsSearching(false);
      }
    },
    [apiFetcher, defaultValue]
  );

  const getChartData = useCallback(() => {
    if (!apiFetcher.data) return { labels: [], counts: [], errCounts: [] };

    const labels = Object.keys(apiFetcher.data?.Distribution || {}).map(
      (item) => {
        const date = new Date(parseInt(item));
        return isMobile ? date.toLocaleTimeString() : date.toLocaleString();
      }
    );
    const counts = Object.values(apiFetcher.data?.Distribution || {});
    const errCounts = Object.values(apiFetcher.data?.ErrorDistTime || {});

    return { labels, counts, errCounts };
  }, [apiFetcher.data, isMobile]);

  const { labels, counts } = useMemo(() => getChartData(), [getChartData]);

  const renderForm = () => (
    <Card className="p-4 md:p-6 rounded-lg shadow-sm mb-4">
      <Box className="flex items-center gap-2 mb-4">
        <Api color="primary" fontSize={isMobile ? "medium" : "large"} />
        <Typography variant={isMobile ? "h6" : "h5"} className="font-bold">
          API Statistics
        </Typography>
      </Box>

      <Divider className="mb-4" />

      <FormProvider {...methods}>
        <form
          onSubmit={methods.handleSubmit(onSubmit)}
          className="flex flex-col gap-3"
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
              <TimeUnitInput
                className="mb-2"
                helperText="Time unit for aggregation"
              />
            </Grid2>
            <Grid2 item xs={12} md={6}>
              <TimeRangeInput className="mb-2" />
            </Grid2>
            <Grid2 item xs={12}>
              <SubmitButtons
                isLoading={isSearching}
                loadingText="Analyzing..."
                submitText="Analyze API"
                className="w-full md:w-auto"
              />
              {apiFetcher.data && (
                <IconButton
                  color="primary"
                  className="ml-2"
                  onClick={() => methods.handleSubmit(onSubmit)()}
                  title="Refresh data"
                >
                  <Refresh />
                </IconButton>
              )}
            </Grid2>
          </Grid2>
        </form>
      </FormProvider>
    </Card>
  );

  const renderLatencyTable = () => (
    <TableContainer component={Paper} className="overflow-x-auto">
      <Table size={isMobile ? "small" : "medium"}>
        <TableHead>
          <TableRow className="bg-gray-50">
            <TableCell className="font-semibold">Metric</TableCell>
            <TableCell className="font-semibold">Microseconds</TableCell>
            <TableCell className="font-semibold">Milliseconds</TableCell>
            <TableCell className="font-semibold">Seconds</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {Object.entries(apiFetcher.data.Latency).map(
            ([key, value], index) => {
              const msValue = (value / 1000).toFixed(2);
              const secValue = (value / 1000000).toFixed(4);
              return (
                <TableRow key={index} hover>
                  <TableCell className="font-medium">
                    {key.charAt(0).toUpperCase() + key.slice(1)}
                  </TableCell>
                  <TableCell>{value.toLocaleString()}</TableCell>
                  <TableCell>
                    <Chip
                      label={msValue}
                      size="small"
                      color={
                        key === "p99" || key === "max" ? "warning" : "default"
                      }
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell>{secValue}</TableCell>
                </TableRow>
              );
            }
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );

  const renderContent = () => {
    if (apiFetcher.loading && !isSearching) {
      return (
        <Fade in={true}>
          <Box className="flex justify-center py-8">
            <CircularProgress />
          </Box>
        </Fade>
      );
    }

    if (apiFetcher.error) {
      return (
        <Alert
          severity="error"
          className="mb-4"
          variant="outlined"
          icon={<ErrorOutline />}
        >
          {apiFetcher.error}
        </Alert>
      );
    }

    if (!apiFetcher.data) {
      return (
        <Alert severity="info" className="mb-4" icon={<Api />}>
          Select API parameters and click &quot;Analyze API&quot; to view
          statistics
        </Alert>
      );
    }

    return (
      <Box className="flex flex-col gap-4" id="chart-container">
        <Accordion defaultExpanded>
          <AccordionSummary
            expandIcon={<ExpandMore />}
            className="bg-gray-50 hover:bg-gray-100"
          >
            <Box className="flex items-center gap-2">
              <BarChartIcon color="primary" />
              <Typography
                variant={isMobile ? "subtitle1" : "h6"}
                className="font-semibold"
              >
                Overview
              </Typography>
            </Box>
          </AccordionSummary>
          <AccordionDetails className="bg-white p-0">
            <CustomContainer title="API Usage Statistics" className="p-0">
              <Grid2 container spacing={3} className="p-4">
                <Grid2 item xs={12} sm={6} md={4}>
                  <StatCard
                    title="Total Calls"
                    value={apiFetcher.data.Count}
                    unit="requests"
                    icon={<Api fontSize="small" color="primary" />}
                  />
                </Grid2>
                <Grid2 item xs={12} sm={6} md={4}>
                  <StatCard
                    title="Call Frequency"
                    value={apiFetcher.data.Frequency}
                    unit={`calls/${timeUnit}`}
                    icon={<AccessTime fontSize="small" color="primary" />}
                  />
                </Grid2>
                <Grid2 item xs={12} md={4}>
                  <StatCard
                    title="Error Rate"
                    value={(apiFetcher.data.ErrorRate * 100).toFixed(2)}
                    unit="%"
                    icon={
                      <ErrorOutline
                        fontSize="small"
                        color={
                          apiFetcher.data.ErrorRate > 0.05 ? "error" : "success"
                        }
                      />
                    }
                    trend={apiFetcher.data.ErrorRate > 0.05 ? "up" : "down"}
                  />
                </Grid2>
                <Grid2 item xs={12}>
                  <BarChartCard
                    title="Request Distribution"
                    caption={`API calls per ${timeUnit}`}
                  >
                    <BarChart
                      width={chartDimensions.width}
                      height={chartDimensions.height}
                      margin={{
                        top: 20,
                        right: 20,
                        bottom: isMobile ? 80 : 40,
                        left: 40,
                      }}
                      xAxis={[
                        {
                          scaleType: "band",
                          data: labels,
                          label: "Time",
                          tickPlacement: "start",
                          tickLabelPlacement: "tick",
                          labelRotation: isMobile ? -45 : 0,
                        },
                      ]}
                      series={[
                        {
                          data: counts,
                          label: "Requests",
                          color: theme.palette.primary.main,
                        },
                      ]}
                    />
                  </BarChartCard>
                </Grid2>
              </Grid2>
            </CustomContainer>
          </AccordionDetails>
        </Accordion>

        <Accordion>
          <AccordionSummary
            expandIcon={<ExpandMore />}
            className="bg-gray-50 hover:bg-gray-100"
          >
            <Box className="flex items-center gap-2">
              <Speed color="warning" />
              <Typography
                variant={isMobile ? "subtitle1" : "h6"}
                className="font-semibold"
              >
                Latency Metrics
              </Typography>
            </Box>
          </AccordionSummary>
          <AccordionDetails className="bg-white p-0">
            <CustomContainer title="Response Time Analysis" className="p-0">
              <Box className="p-4">{renderLatencyTable()}</Box>
            </CustomContainer>
          </AccordionDetails>
        </Accordion>

        <Accordion>
          <AccordionSummary
            expandIcon={<ExpandMore />}
            className="bg-gray-50 hover:bg-gray-100"
          >
            <Box className="flex items-center gap-2">
              <ErrorOutline color="error" />
              <Typography
                variant={isMobile ? "subtitle1" : "h6"}
                className="font-semibold"
              >
                Error Analysis
              </Typography>
              {apiFetcher.data.ErrorCount > 0 && (
                <Chip
                  label={`${apiFetcher.data.ErrorCount} errors`}
                  size="small"
                  color="error"
                  variant="outlined"
                  className="ml-2"
                />
              )}
            </Box>
          </AccordionSummary>
          <AccordionDetails className="bg-white p-0">
            <CustomContainer title="Error Statistics" className="p-0">
              <Grid2 container spacing={3} className="p-4">
                <Grid2 item xs={12} sm={6} md={4}>
                  <StatCard
                    title="Error Count"
                    value={apiFetcher.data.ErrorCount}
                    unit="errors"
                    icon={<ErrorOutline fontSize="small" color="error" />}
                  />
                </Grid2>
                <Grid2 item xs={12} sm={6} md={4}>
                  <StatCard
                    title="Error Rate"
                    value={(apiFetcher.data.ErrorRate * 100).toFixed(2)}
                    unit="%"
                    icon={<ErrorOutline fontSize="small" color="error" />}
                  />
                </Grid2>
                <BarChartCard title="" caption="Errors by status code">
                  <BarChart
                    width={chartDimensions.width}
                    height={chartDimensions.height}
                    margin={{
                      top: 20,
                      right: 20,
                      bottom: isMobile ? 80 : 40,
                      left: 40,
                    }}
                    xAxis={[
                      {
                        scaleType: "band",
                        data: Object.keys(apiFetcher.data.ErrorDist || {}),
                        labels: "Status Code",
                      },
                    ]}
                    series={[
                      {
                        data: Object.values(apiFetcher.data.ErrorDist || {}),
                        label: "Count",
                      },
                    ]}
                  />
                </BarChartCard>
              </Grid2>
            </CustomContainer>
          </AccordionDetails>
        </Accordion>
      </Box>
    );
  };

  return (
    <Box className="flex flex-col items-center gap-2 p-2">
      <Typography variant="h5" className="text-2xl font-bold">
        API Statistic
      </Typography>
      {renderForm()}
      {renderContent()}
    </Box>
  );
};

export default ApiStatistic;
