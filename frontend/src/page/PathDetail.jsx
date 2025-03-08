import { FormProvider, useForm } from "react-hook-form";
import useFetchData from "../hook/useFetchData";
import {
  Accordion,
  Box,
  Grid2,
  Typography,
  CircularProgress,
  AccordionDetails,
  AccordionSummary,
} from "@mui/material";
import { useState } from "react";
import { ArrowDropDownIcon } from "@mui/x-date-pickers/icons";
import { TimeRangeInput, TimeUnitInput } from "../component/shared/Input";
import { SubmitButtons } from "../component/shared/Common";
import CustomContainer from "../component/shared/CustomContainer";
import StatCard from "../component/shared/StatCard";
import BarChartCard from "../component/shared/BarChartCard";
import PathTree from "../component/PathTree";
import { useParams } from "react-router-dom";
import dayjs from "dayjs";
import { BarChart } from "@mui/x-charts";
import HopDetails from "../component/HopDetail";
const PathDetail = () => {
  const { path_id } = useParams();
  const [hopID, setHopID] = useState(null);
  const { data, loading, error, fetchData } = useFetchData(`/paths/${path_id}`);
  const methods = useForm({
    defaultValues: {
      unit: "hour",
      from: null,
      to: null,
    },
  });
  const onSubmit = async (data) => {
    const params = {
      ...data,
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
    };
    await fetchData(params);
  };

  const [showHopDetail, setShowHopDetail] = useState(false);
  const [params, setParams] = useState();
  const handleLinkClick = (event) => {
    const edge = event.target;
    const sourceId = edge.data("source");
    const targetId = edge.data("target");
    setHopID(`${sourceId}_${targetId}_${path_id}`);

    const params = {
      from: data.from?.$d.getTime() || dayjs().startOf("day").valueOf(),
      to:
        data.to?.$d.getTime() || dayjs().startOf("day").add(1, "day").valueOf(),
      unit: methods.getValues("unit"),
    };
    setShowHopDetail(true);
    setParams(params);
  };

  return (
    <Box className="p-4 space-y-4">
      <Typography variant="h4" className="font-bold mb-4">
        Path Detail
      </Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)} className="space-y-4">
          <Grid2 container spacing={2}>
            <TimeRangeInput />
            <TimeUnitInput />
            <SubmitButtons />
          </Grid2>
        </form>
      </FormProvider>
      {loading && <CircularProgress />}
      {error && <div className="text-red-500">{error}</div>}
      {data && (
        <Box className="space-y-4">
          <Accordion defaultExpanded>
            <AccordionSummary expandIcon={<ArrowDropDownIcon />}>
              <Typography>View path call tree</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <CustomContainer>
                <PathTree
                  path={data.path_info}
                  handleLinkClick={handleLinkClick}
                />
              </CustomContainer>
            </AccordionDetails>
          </Accordion>
          <Accordion>
            <AccordionSummary expandIcon={<ArrowDropDownIcon />}>
              <Typography>View path statistics</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <CustomContainer>
                <Grid2 container spacing={2}>
                  <StatCard title={"Count"} value={data.count} unit="calls" />
                  <StatCard
                    title={"Frequency"}
                    value={data.frequency}
                    unit={`calls/${methods.getValues("unit")}`}
                  />
                  <BarChartCard
                    title={"Frequency distribution"}
                    caption={`Calls per ${methods.getValues("unit")}`}
                  >
                    <BarChart
                      width={600}
                      height={300}
                      xAxis={[
                        {
                          scaleType: "band",
                          data: Object.keys(data.distribution || {}).map(
                            (key) => new Date(parseInt(key)).toLocaleString()
                          ),
                          label: "Timestamp",
                          tickPlacement: "start",
                          tickLabelPlacement: "tick",
                        },
                      ]}
                      series={[
                        {
                          data: Object.values(data.distribution || {}),
                          label: "Count",
                        },
                      ]}
                    />
                  </BarChartCard>
                </Grid2>
              </CustomContainer>
            </AccordionDetails>
          </Accordion>
          <Accordion>
            <AccordionSummary expandIcon={<ArrowDropDownIcon />}>
              <Typography>View Error</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <CustomContainer>
                <Grid2 container spacing={2}>
                  <StatCard
                    title={"Count"}
                    value={data.error_count}
                    unit="errors"
                  />
                  <StatCard
                    title={"Error rate"}
                    value={data.error_rate * 100}
                    unit="%"
                  />
                  <BarChartCard
                    title={"Error distribution"}
                    caption={`Errors per ${methods.getValues("unit")}`}
                  >
                    <BarChart
                      width={600}
                      height={300}
                      xAxis={[
                        {
                          scaleType: "band",
                          data: Object.keys(data.error_dist || {}).map((key) =>
                            new Date(parseInt(key)).toLocaleString()
                          ),
                          label: "Timestamp",
                          tickPlacement: "start",
                          tickLabelPlacement: "tick",
                        },
                      ]}
                      series={[
                        {
                          data: Object.values(data.error_dist || {}),
                          label: "Count",
                        },
                      ]}
                    />
                  </BarChartCard>
                </Grid2>
              </CustomContainer>
            </AccordionDetails>
          </Accordion>
          {showHopDetail && (
            <HopDetails
              hopID={hopID}
              params={params}
              setShowHopDetail={setShowHopDetail}
              unit={data.unit}
              className="mt-4"
            />
          )}
        </Box>
      )}
    </Box>
  );
};

export default PathDetail;
