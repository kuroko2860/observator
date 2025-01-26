import {
  Container,
  CircularProgress,
  InputLabel,
  TextField,
  Typography,
} from "@mui/material";
import useFetchData from "../hook/useFetchData";
import BarChartCard from "../component/BarChartCard";
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
    <Container>
      <Typography variant="h5">View most called services</Typography>
      <CustomForm onSubmit={onSubmit}>
        <TimeRangeInput />
      </CustomForm>
      {topSvcFetcher.loading && <CircularProgress />}
      {topSvcFetcher.error && <div>{topSvcFetcher.error}</div>}
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
    <BarChartCard
      caption={`Top ${Math.min(
        limit,
        Object.keys(data).length
      )} services called`}
    >
      <InputLabel>Limit</InputLabel>
      <TextField
        type="number"
        required
        defaultValue={10}
        onChange={(e) => setLimit(e.target.value)}
      />
      <BarChart
        width={700}
        height={350}
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
    </BarChartCard>
  );
};

export default TopService;
