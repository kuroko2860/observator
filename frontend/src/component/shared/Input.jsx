import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import {
  Grid2,
  useTheme,
  useMediaQuery,
  Box,
  Typography,
  Tooltip,
} from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import { NumberInput, SelectionInput } from "./Common";
import { useSelector } from "react-redux";
import { getAllServices } from "../../redux/services/selector";
import { Info } from "@mui/icons-material";

const TimeRangeInput = ({ className }) => {
  const methods = useFormContext();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Grid2
      container
      spacing={2}
      className={`flex flex-col md:flex-row gap-2 ${className || ""}`}
    >
      <Grid2 item xs={12} sm={6}>
        <Controller
          name="from"
          control={methods.control}
          render={({ field, fieldState: { error } }) => (
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DateTimePicker
                {...field}
                views={["year", "month", "day", "hours", "minutes", "seconds"]}
                format="YYYY-MM-DD HH:mm:ss"
                label="From time"
                className="w-full"
                slotProps={{
                  textField: {
                    variant: "outlined",
                    error: !!error,
                    helperText: error?.message,
                    sx: {
                      "& .MuiInputBase-root": {
                        borderRadius: "8px",
                      },
                      "& .MuiOutlinedInput-notchedOutline": {
                        borderWidth: "1px",
                      },
                      "& .MuiInputLabel-root": {
                        fontSize: isMobile ? "0.875rem" : "1rem",
                      },
                    },
                  },
                  mobilePaper: {
                    sx: {
                      "& .MuiPickersToolbar-root": {
                        padding: isMobile ? "8px" : "16px",
                      },
                    },
                  },
                }}
              />
            </LocalizationProvider>
          )}
        />
      </Grid2>
      <Grid2 item xs={12} sm={6}>
        <Controller
          name="to"
          control={methods.control}
          render={({ field, fieldState: { error } }) => (
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DateTimePicker
                {...field}
                views={["year", "month", "day", "hours", "minutes", "seconds"]}
                format="YYYY-MM-DD HH:mm:ss"
                label="To time"
                className="w-full"
                slotProps={{
                  textField: {
                    variant: "outlined",
                    error: !!error,
                    helperText: error?.message,
                    sx: {
                      "& .MuiInputBase-root": {
                        borderRadius: "8px",
                      },
                      "& .MuiOutlinedInput-notchedOutline": {
                        borderWidth: "1px",
                      },
                      "& .MuiInputLabel-root": {
                        fontSize: isMobile ? "0.875rem" : "1rem",
                      },
                    },
                  },
                  mobilePaper: {
                    sx: {
                      "& .MuiPickersToolbar-root": {
                        padding: isMobile ? "8px" : "16px",
                      },
                    },
                  },
                }}
              />
            </LocalizationProvider>
          )}
        />
      </Grid2>
    </Grid2>
  );
};

const ServiceNameInput = ({
  label,
  className,
  required = false,
  helperText,
}) => {
  const services = useSelector(getAllServices);
  const options = services.map((s) => [s, s]);

  return (
    <Box className={className || ""}>
      <Box className="flex items-center gap-1 mb-1">
        {label && (
          <Typography variant="body2" className="text-gray-700 font-medium">
            {label || "Service name"}
            {required && <span className="text-red-500 ml-1">*</span>}
          </Typography>
        )}
        {helperText && (
          <Tooltip title={helperText} arrow placement="top">
            <Info fontSize="small" className="text-gray-400 cursor-help" />
          </Tooltip>
        )}
      </Box>
      <SelectionInput
        labelId={"service-name"}
        name={"service_name"}
        label={label || "Service name"}
        options={options}
        required={required}
        className="w-full"
      />
    </Box>
  );
};

const EndpointInput = ({
  endpoints,
  className,
  required = false,
  helperText,
}) => {
  const options = endpoints.map((e) => [e, e]);

  return (
    <Box className={className || ""}>
      <Box className="flex items-center gap-1 mb-1">
        {helperText && (
          <Tooltip title={helperText} arrow placement="top">
            <Info
              fontSize="small"
              className="text-gray-400 cursor-help ml-auto"
            />
          </Tooltip>
        )}
      </Box>
      <SelectionInput
        labelId={"uri_path"}
        name={"uri_path"}
        label={"URI Path"}
        options={options}
        required={required}
        className="w-full"
      />
    </Box>
  );
};

const ThresholdInput = ({
  label,
  className,
  min = 0,
  max,
  step = 1,
  required = false,
  helperText,
}) => {
  return (
    <Box className={className || ""}>
      <NumberInput
        name={"threshold"}
        label={label}
        min={min}
        max={max}
        step={step}
        required={required}
        helperText={helperText}
        className="w-full"
      />
    </Box>
  );
};

const MethodInput = ({ className, required = false, helperText }) => {
  const options = [
    ["GET", "GET"],
    ["POST", "POST"],
    ["PUT", "PUT"],
    ["PATCH", "PATCH"],
    ["DELETE", "DELETE"],
  ];

  return (
    <Box className={className || ""}>
      <SelectionInput
        labelId={"method"}
        name={"method"}
        label={"Method"}
        options={options}
        required={required}
        helperText={helperText}
        className="w-full"
      />
    </Box>
  );
};

const TimeUnitInput = ({ className, required = false, helperText }) => {
  const options = [
    ["second", "Second"],
    ["minute", "Minute"],
    ["hour", "Hour"],
    ["day", "Day"],
  ];

  return (
    <Box className={className || ""}>
      <SelectionInput
        labelId={"time-unit"}
        label={"Time unit"}
        name={"unit"}
        options={options}
        required={required}
        helperText={helperText}
        className="w-full"
      />
    </Box>
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
