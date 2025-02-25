import {
  Container,
  CircularProgress,
  TextField,
  Typography,
  Box,
  Card,
  Grid2 as Grid,
} from "@mui/material";
import useFetchData from "../hook/useFetchData";
import { BarChart } from "@mui/x-charts";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import dayjs from "dayjs";
import { CustomForm } from "../component/Common";
import { TimeRangeInput } from "../component/Input";

const TopService = () => {
  const topSvcFetcher = useFetchData("/services/top-called");
  const onSubmit = async (data) => {
    topSvcFetcher.fetchData({
      ...data,
      limit: 10,
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    });
  };

  return (
    <Container
      className="flex flex-col gap-4"
      sx={{
        "& .MuiTypography-h5": {
          fontWeight: 600,
          marginBottom: 2,
        },
      }}
    >
      <Typography variant="h5">View most called services</Typography>
      <CustomForm onSubmit={onSubmit} className="flex flex-col gap-2">
        <TimeRangeInput />
      </CustomForm>
      {topSvcFetcher.loading && (
        <Box className="flex items-center justify-center">
          <CircularProgress />
        </Box>
      )}
      {topSvcFetcher.error && (
        <Typography variant="body1" className="text-red-600">
          {topSvcFetcher.error}
        </Typography>
      )}
      {topSvcFetcher.data && (
        <TopServiceCalledChart
          data={topSvcFetcher.data}
        ></TopServiceCalledChart>
      )}
    </Container>
  );
};

const TopServiceCalledChart = ({ data }) => {
  const [limit, setLimit] = useState(10);
  const _data = Object.entries(data)
    .sort((a, b) => b[1] - a[1])
    .slice(0, limit);
  const navigate = useNavigate();
  const xa = _data.map((d) => d[0]);
  const sr = _data.map((d) => d[1]);
  return (
    <Card className="p-4">
      <Typography variant="h6" className="font-bold">
        Top {Math.min(limit, Object.keys(data).length)} services called
      </Typography>
      <Grid container spacing={2} className="mt-2">
        <Grid item xs={12} sm={6}>
          <TextField
            type="number"
            required
            defaultValue={10}
            onChange={(e) => setLimit(e.target.value)}
            label="Limit"
            className="w-full"
            variant="outlined"
          />
        </Grid>
        <Grid item xs={12}>
          <BarChart
            width={700}
            height={350}
            className="mt-4"
            xAxis={[
              {
                scaleType: "band",
                data: xa,
                tickPlacement: "middle",
                tickLabelPlacement: "tick",
              },
            ]}
            series={[
              {
                data: sr,
                label: "count",
              },
            ]}
            onAxisClick={(_event, data) =>
              navigate(`/service-detail/${data.axisValue}`)
            }
          />
        </Grid>
      </Grid>
    </Card>
  );
};

export default TopService;
