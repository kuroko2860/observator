import {
  Container,
  CircularProgress,
  TextField,
  Typography,
  Box,
  Card,
  Grid2 as Grid,
  useTheme,
  useMediaQuery,
  Fade,
} from "@mui/material";
import useFetchData from "../hook/useFetchData";
import { BarChart } from "@mui/x-charts";
import { useNavigate } from "react-router-dom";
import { useState, useEffect, useCallback } from "react";
import dayjs from "dayjs";
import { CustomForm } from "../component/shared/Common";
import { TimeRangeInput } from "../component/shared/Input";

const TopService = () => {
  const { data, loading, error, fetchData } = useFetchData(
    "/services/top-called"
  );
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  // Fetch data on initial load
  useEffect(() => {
    const defaultStartTime = dayjs()
      .add(1, "minute")
      .second(0)
      .millisecond(0)
      .subtract(1, "hour")
      .valueOf();
    const defaultEndTime = dayjs()
      .add(1, "minute")
      .second(0)
      .millisecond(0)
      .valueOf();

    fetchData({
      limit: 10,
      from: defaultStartTime,
      to: defaultEndTime,
    });
  }, []);

  const handleSubmit = async (formData) => {
    const defaultStartTime = dayjs()
      .add(1, "minute")
      .second(0)
      .millisecond(0)
      .subtract(1, "hour")
      .valueOf();
    const defaultEndTime = dayjs()
      .add(1, "minute")
      .second(0)
      .millisecond(0)
      .valueOf();

    fetchData({
      ...formData,
      limit: 10,
      from: formData.from?.$d.getTime() || defaultStartTime,
      to: formData.to?.$d.getTime() || defaultEndTime,
    });
  };

  return (
    <Container
      className="flex flex-col gap-4 p-4 md:p-6 bg-white shadow-lg rounded-lg"
      maxWidth="lg"
    >
      <Typography
        variant="h5"
        className="text-xl md:text-2xl font-bold text-center"
      >
        View most called services
      </Typography>

      <Card className="p-3 md:p-4 mb-4">
        <CustomForm onSubmit={handleSubmit}>
          <TimeRangeInput />
        </CustomForm>
      </Card>

      <Fade in={loading}>
        <Box
          className={`flex items-center justify-center ${
            loading ? "h-64" : "h-0"
          }`}
        >
          <CircularProgress />
        </Box>
      </Fade>

      {error && (
        <Card className="p-4 bg-red-50 border border-red-200">
          <Typography variant="body1" className="text-red-600">
            {error}
          </Typography>
        </Card>
      )}

      <Fade in={!loading && !!data}>
        <Box>
          {data && <TopServiceCalledChart data={data} isMobile={isMobile} />}
        </Box>
      </Fade>
    </Container>
  );
};

const TopServiceCalledChart = ({ data, isMobile }) => {
  const [limit, setLimit] = useState(Object.keys(data).length);
  const navigate = useNavigate();
  const theme = useTheme();
  const [chartDimensions, setChartDimensions] = useState({
    width: 700,
    height: 350,
  });

  // Update chart dimensions based on window size
  useEffect(() => {
    const updateDimensions = () => {
      const containerWidth =
        document.getElementById("chart-container")?.clientWidth || 700;
      setChartDimensions({
        width: Math.min(containerWidth - 40, 900),
        height: isMobile ? 300 : 400,
      });
    };

    updateDimensions();
    window.addEventListener("resize", updateDimensions);
    return () => window.removeEventListener("resize", updateDimensions);
  }, [isMobile]);

  const sortedData = Object.entries(data)
    .sort((a, b) => b[1] - a[1])
    .slice(0, limit);

  const serviceNames = sortedData.map(([name]) => name);
  const serviceCounts = sortedData.map(([, count]) => count);

  const handleLimitChange = useCallback((event) => {
    const value = Number(event.target.value);
    setLimit(value > 0 ? value : 1);
  }, []);

  const handleBarClick = useCallback(
    (_event, data) => {
      navigate(`/service-detail/${data.axisValue}`);
    },
    [navigate]
  );

  return (
    <Card
      className="p-4 flex flex-col justify-center items-center gap-4"
      elevation={2}
    >
      {/* <Typography variant="h6" className="font-bold text-center">
        Top {Math.min(limit, Object.keys(data).length)} called services
      </Typography> */}

      {/* <Grid container spacing={2} className="flex "> */}
      <Grid item xs={12} sm={6} md={4}>
        <TextField
          type="number"
          required
          value={limit}
          onChange={handleLimitChange}
          label="Number of services to display"
          className="w-full"
          variant="outlined"
          inputProps={{ min: 1, max: Object.keys(data).length }}
          helperText="Adjust to show more or fewer services"
        />
      </Grid>

      <Grid item xs={12} id="chart-container">
        <Box className="mt-4 overflow-x-auto">
          {isMobile ? (
            // Mobile-optimized chart with vertical layout
            <BarChart
              width={chartDimensions.width}
              height={chartDimensions.height}
              layout="vertical"
              margin={{
                left: 120,
                right: 10,
                top: 10,
                bottom: 30,
              }}
              yAxis={[
                {
                  scaleType: "band",
                  data: serviceNames,
                  tickLabelStyle: {
                    fontSize: 12,
                    textAnchor: "end",
                  },
                },
              ]}
              series={[
                {
                  data: serviceCounts,
                  label: "Call Count",
                  color: theme.palette.primary.main,
                },
              ]}
              onAxisClick={handleBarClick}
              tooltip={{ trigger: "item" }}
            />
          ) : (
            // Desktop chart with horizontal layout
            <BarChart
              width={chartDimensions.width}
              height={chartDimensions.height}
              margin={{
                left: 40,
                right: 10,
                top: 10,
                bottom: 70,
              }}
              xAxis={[
                {
                  scaleType: "band",
                  data: serviceNames,
                  tickPlacement: "middle",
                  tickLabelPlacement: "tick",
                  tickLabelStyle: {
                    fontSize: 10,
                  },
                },
              ]}
              series={[
                {
                  data: serviceCounts,
                  label: "Call Count",
                  color: theme.palette.primary.main,
                },
              ]}
              onAxisClick={handleBarClick}
              tooltip={{ trigger: "item" }}
            />
          )}
          <Typography
            variant="caption"
            className="block text-center mt-2 text-gray-600"
          >
            Click on a service name to view details
          </Typography>
        </Box>
      </Grid>
      {/* </Grid> */}
    </Card>
  );
};

export default TopService;
