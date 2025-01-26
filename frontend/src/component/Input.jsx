import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { Grid2 } from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import { NumberInput, SelectionInput } from "./Common";
import { useSelector } from "react-redux";
import { getAllServices } from "../redux/services/selector";

const TimeRangeInput = () => {
  const methods = useFormContext();
  return (
    <Grid2>
      <Grid2>
        <Controller
          name="from"
          control={methods.control}
          render={({ field }) => (
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DateTimePicker
                {...field}
                views={["year", "month", "day", "hours", "minutes", "seconds"]}
                format="YYYY-MM-DD HH-mm-ss"
                label="From time"
              />
            </LocalizationProvider>
          )}
        />
      </Grid2>
      <Grid2>
        <Controller
          name="to"
          control={methods.control}
          render={({ field }) => (
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DateTimePicker
                {...field}
                views={["year", "month", "day", "hours", "minutes", "seconds"]}
                format="YYYY-MM-DD HH-mm-ss"
                label="to time"
              />
            </LocalizationProvider>
          )}
        />
      </Grid2>
    </Grid2>
  );
};

const ServiceNameInput = ({ label }) => {
  const methods = useFormContext();
  const services = useSelector(getAllServices);
  const options = services.map((s) => [s, s]);
  return (
    <SelectionInput
      control={methods.control}
      label={label || "service name"}
      labelId={"service-name"}
      name={"service_name"}
      options={options}
    />
  );
};

const EndpointInput = ({ endpoints }) => {
  const methods = useFormContext();
  const options = endpoints.map((e) => [e, e]);
  return (
    <SelectionInput
      control={methods.control}
      label={"endpoint"}
      labelId={"endpoint"}
      name={"endpoint"}
      options={options}
    />
  );
};
const ThresholdInput = ({ label }) => {
  const methods = useFormContext();
  return (
    <NumberInput control={methods.control} name={"threshold"} label={label} />
  );
};
const MethodInput = () => {
  const methods = useFormContext();
  const options = [
    ["GET", "GET"],
    ["POST", "POST"],
    ["PUT", "PUT"],
    ["PATCH", "PATCH"],
    ["DELETE", "DELETE"],
  ];
  return (
    <SelectionInput
      control={methods.control}
      labelId={"method"}
      name={"method"}
      label={"Method"}
      options={options}
    />
  );
};
const TimeUnitInput = () => {
  const methods = useFormContext();
  const options = [
    ["second", "second"],
    ["minute", "minute"],
    ["hour", "hour"],
    ["day", "day"],
  ];
  return (
    <SelectionInput
      control={methods.control}
      labelId={"time-unit"}
      label={"Time unit"}
      name={"unit"}
      options={options}
    />
  );
};

export {
  TimeRangeInput,
  ServiceNameInput,
  EndpointInput,
  ThresholdInput,
  MethodInput,
  TimeUnitInput,
};
