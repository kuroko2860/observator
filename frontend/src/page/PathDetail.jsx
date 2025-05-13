import TabContext from "@mui/lab/TabContext";
import TabList from "@mui/lab/TabList";
import TabPanel from "@mui/lab/TabPanel";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  CircularProgress,
  Container,
  Grid2,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import { BarChart } from "@mui/x-charts";
import { ArrowDropDownIcon } from "@mui/x-date-pickers/icons";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import HopDetails from "../component/HopDetail";
import PathTree from "../component/PathTree";
import BarChartCard from "../component/shared/BarChartCard";
import { SubmitButtons } from "../component/shared/Common";
import CustomContainer from "../component/shared/CustomContainer";
import { TimeRangeInput, TimeUnitInput } from "../component/shared/Input";
import StatCard from "../component/shared/StatCard";
import useFetchData from "../hook/useFetchData";
import axios from "../config/axios";

// Helper functions
const getFromTime = (date) =>
  date?.$d.getTime() ||
  dayjs()
    .add(1, "minute")
    .second(0)
    .millisecond(0)
    .subtract(1, "hour")
    .valueOf();
const getEndTime = () =>
  dayjs().add(1, "minute").second(0).millisecond(0).valueOf();

const formatDateKeys = (obj) => {
  return Object.keys(obj || {}).map((key) =>
    new Date(parseInt(key)).toLocaleString()
  );
};

// Format timestamp for display
const formatTimestamp = (timestamp) => {
  if (!timestamp) return "";
  const date = new Date(timestamp);
  const now = new Date();
  const diff = (now.getTime() - date.getTime()) / 1000;

  if (diff < 60) {
    return "now";
  } else if (diff < 3600) {
    return `${Math.floor(diff / 60)} minutes ago`;
  } else if (diff < 86400) {
    return `${Math.floor(diff / 3600)} hours ago`;
  } else if (diff < 604800) {
    return `${Math.floor(diff / 86400)} days ago`;
  } else {
    return date.toLocaleString();
  }
};

// Components
const PathStatistics = ({ data, unit }) => (
  <Grid2 container spacing={2}>
    <StatCard title="Count" value={data.count} unit="calls" />
    <StatCard title="Frequency" value={data.frequency} unit={`calls/${unit}`} />
    <BarChartCard title="Frequency distribution" caption={`Calls per ${unit}`}>
      <BarChart
        width={600}
        height={300}
        xAxis={[
          {
            scaleType: "band",
            data: formatDateKeys(data.distribution),
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
);

const TraceTable = ({ traces, onViewTrace }) => (
  <Table>
    <TableHead>
      <TableRow>
        <TableCell>Root</TableCell>
        <TableCell>Start time</TableCell>
        <TableCell>Spans</TableCell>
        <TableCell>Duration</TableCell>
        <TableCell></TableCell>
      </TableRow>
    </TableHead>
    <TableBody>
      {traces.length === 0 ? (
        <TableRow>
          <TableCell colSpan={5} align="center">
            No traces found
          </TableCell>
        </TableRow>
      ) : (
        traces.map((trace) => (
          <TableRow key={trace.id}>
            <TableCell>
              {trace.root_service + " "}
              <small className="text-gray-500">{trace.root_operation}</small>
            </TableCell>
            <TableCell>{formatTimestamp(trace.start_time)}</TableCell>
            <TableCell>{trace.span_num}</TableCell>
            <TableCell>{trace.duration}</TableCell>
            <TableCell>
              <Button
                variant="contained"
                color="primary"
                size="small"
                onClick={() => onViewTrace(trace.trace_id)}
              >
                View
              </Button>
            </TableCell>
          </TableRow>
        ))
      )}
    </TableBody>
  </Table>
);

const PathDetail = () => {
  const navigate = useNavigate();
  const { path_id } = useParams();
  const [activeTab, setActiveTab] = useState("1");
  const [hopID, setHopID] = useState(null);
  const [traces, setTraces] = useState([]);
  const [showHopDetail, setShowHopDetail] = useState(false);
  const [hopParams, setHopParams] = useState();

  const { data, loading, error, fetchData } = useFetchData(`/paths/${path_id}`);

  const methods = useForm({
    defaultValues: {
      unit: "minute",
      from: null,
      to: null,
    },
  });

  useEffect(() => {
    async function _fetchData() {
      console.log("fetching data");
      const params = {
        from: dayjs()
          .add(1, "minute")
          .second(0)
          .millisecond(0)
          .subtract(1, "hour")
          .valueOf(),
        to: dayjs().add(1, "minute").second(0).millisecond(0).valueOf(),
        unit: "minute",
      };

      await fetchData(params);
      const { data: trace_data } = await axios.get(`/paths/${path_id}/traces`, {
        params,
      });
      setTraces(trace_data || []);
    }
    _fetchData();
  }, []);

  const onSubmit = async (formData) => {
    const params = {
      ...formData,
      from: getFromTime(formData.from),
      to: formData.to?.$d.getTime() || getEndTime(),
    };

    await fetchData(params);
    const { data: trace_data } = await axios.get(`/paths/${path_id}/traces`, {
      params,
    });
    setTraces(trace_data || []);
  };

  const handleLinkClick = (event) => {
    const edge = event.target;
    const sourceId = edge.data("source");
    const targetId = edge.data("target");

    setHopID(`${sourceId}_${targetId}_${path_id}`);
    const data = methods.getValues();
    setHopParams({
      from: getFromTime(data.from),
      to: data.to?.$d.getTime() || getEndTime(),
      unit: methods.getValues("unit"),
    });
    setShowHopDetail(true);
  };

  return (
    <Box className="p-4 space-y-4">
      <CustomContainer title={"View path detail"}>
        <Container className="flex items-center justify-center">
          <FormProvider {...methods}>
            <form onSubmit={methods.handleSubmit(onSubmit)}>
              <Grid2 container spacing={2}>
                <TimeRangeInput />
                <TimeUnitInput />
                <SubmitButtons />
              </Grid2>
            </form>
          </FormProvider>
        </Container>
      </CustomContainer>

      {loading && <CircularProgress />}
      {error && <div className="text-red-500">{error}</div>}

      {data && (
        <Box sx={{ width: "100%", typography: "body1" }}>
          <TabContext value={activeTab}>
            <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
              <TabList
                onChange={(_, value) => setActiveTab(value)}
                aria-label="tabs path"
              >
                <Tab label="View path" value="1" />
                <Tab label="View trace" value="2" />
              </TabList>
            </Box>

            <TabPanel value="1" keepMounted>
              <Box className="space-y-4">
                <Accordion defaultExpanded>
                  <AccordionSummary expandIcon={<ArrowDropDownIcon />}>
                    <Typography>View path call tree</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <CustomContainer>
                      <Box className="flex w-full">
                        <Box
                          className={showHopDetail ? "w-1/2 pr-2" : "w-full"}
                        >
                          <PathTree
                            path={data.path_info}
                            handleLinkClick={handleLinkClick}
                          />
                        </Box>
                        {showHopDetail && (
                          <Box className="w-1/2 pl-2">
                            <HopDetails
                              hopID={hopID}
                              params={hopParams}
                              setShowHopDetail={setShowHopDetail}
                              unit={data.unit}
                            />
                          </Box>
                        )}
                      </Box>
                    </CustomContainer>
                  </AccordionDetails>
                </Accordion>

                <Accordion>
                  <AccordionSummary expandIcon={<ArrowDropDownIcon />}>
                    <Typography>View path statistics</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <CustomContainer>
                      <PathStatistics
                        data={data}
                        unit={methods.getValues("unit")}
                      />
                    </CustomContainer>
                  </AccordionDetails>
                </Accordion>
              </Box>
            </TabPanel>

            <TabPanel value="2" keepMounted>
              <TraceTable
                traces={traces}
                onViewTrace={(id) => navigate(`/trace-detail/${id}`)}
              />
            </TabPanel>
          </TabContext>
        </Box>
      )}
    </Box>
  );
};

export default PathDetail;
