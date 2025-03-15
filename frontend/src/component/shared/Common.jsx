import {
  Box,
  Button,
  CircularProgress,
  FormControl,
  Grid2 as Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { useCallback, useState } from "react";
import {
  Controller,
  FormProvider,
  useForm,
  useFormContext,
} from "react-hook-form";
import { twMerge } from "tailwind-merge";

function TextInput({
  name,
  label,
  className,
  helperText,
  required = false,
  disabled = false,
}) {
  const methods = useFormContext();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Grid
      item
      xs={12}
      sm={6}
      className={twMerge("flex flex-col gap-2", className)}
    >
      <Controller
        name={name}
        control={methods.control}
        render={({ field, fieldState: { error } }) => (
          <TextField
            {...field}
            fullWidth
            label={label}
            variant="outlined"
            required={required}
            disabled={disabled}
            error={!!error}
            helperText={error ? error.message : helperText}
            className="transition-all duration-200"
            sx={{
              "& .MuiInputBase-root": {
                borderRadius: "8px",
              },
              "& .MuiOutlinedInput-notchedOutline": {
                borderWidth: "1px",
              },
              "& .MuiInputLabel-root": {
                fontSize: isMobile ? "0.875rem" : "1rem",
              },
            }}
          />
        )}
      />
    </Grid>
  );
}

const NumberInput = ({
  name,
  label,
  className,
  min,
  max,
  step = 1,
  helperText,
  required = false,
  disabled = false,
}) => {
  const methods = useFormContext();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Grid
      item
      xs={12}
      sm={6}
      className={twMerge("flex flex-col gap-2", className)}
    >
      <Controller
        name={name}
        control={methods.control}
        render={({ field, fieldState: { error } }) => (
          <TextField
            type="number"
            {...field}
            fullWidth
            label={label}
            variant="outlined"
            required={required}
            disabled={disabled}
            error={!!error}
            helperText={error ? error.message : helperText}
            inputProps={{
              min: min,
              max: max,
              step: step,
            }}
            className="transition-all duration-200"
            sx={{
              "& .MuiInputBase-root": {
                borderRadius: "8px",
              },
              "& .MuiOutlinedInput-notchedOutline": {
                borderWidth: "1px",
              },
              "& .MuiInputLabel-root": {
                fontSize: isMobile ? "0.875rem" : "1rem",
              },
            }}
          />
        )}
      />
    </Grid>
  );
};

const SelectionInput = ({
  labelId,
  name,
  label,
  options,
  className,
  helperText,
  required = false,
  disabled = false,
}) => {
  const methods = useFormContext();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Grid
      item
      xs={12}
      sm={6}
      className={twMerge("flex flex-col gap-2", className)}
    >
      <FormControl
        fullWidth
        required={required}
        disabled={disabled}
        error={!!methods.formState.errors[name]}
      >
        <InputLabel
          id={labelId}
          sx={{
            fontSize: isMobile ? "0.875rem" : "1rem",
          }}
        >
          {label}
        </InputLabel>
        <Controller
          name={name}
          control={methods.control}
          render={({ field }) => (
            <Select
              {...field}
              labelId={labelId}
              id={name}
              label={label}
              className="transition-all duration-200"
              sx={{
                borderRadius: "8px",
                "& .MuiOutlinedInput-notchedOutline": {
                  borderWidth: "1px",
                },
              }}
            >
              {options.map(([value, label]) => (
                <MenuItem key={value} value={value}>
                  {label}
                </MenuItem>
              ))}
            </Select>
          )}
        />
        {(methods.formState.errors[name] || helperText) && (
          <Box
            className="text-xs mt-1 ml-3"
            sx={{
              color: methods.formState.errors[name]
                ? theme.palette.error.main
                : theme.palette.text.secondary,
            }}
          >
            {methods.formState.errors[name]?.message || helperText}
          </Box>
        )}
      </FormControl>
    </Grid>
  );
};

const SubmitButtons = ({
  className,
  isSubmitting = false,
  resetLabel = "Reset",
  submitLabel = "Submit",
}) => {
  const methods = useFormContext();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Grid
      item
      xs={12}
      className={twMerge(
        "flex flex-col sm:flex-row justify-end gap-2 mt-4",
        className
      )}
    >
      <Button
        type="button"
        onClick={() => methods.reset()}
        variant="outlined"
        disabled={isSubmitting}
        className="order-2 sm:order-1 w-full sm:w-auto"
        sx={{
          borderRadius: "8px",
          padding: isMobile ? "8px 16px" : "10px 24px",
          textTransform: "none",
          fontWeight: 500,
        }}
      >
        {resetLabel}
      </Button>
      <Button
        type="submit"
        variant="contained"
        disabled={isSubmitting}
        className="order-1 sm:order-2 w-full sm:w-auto"
        sx={{
          borderRadius: "8px",
          padding: isMobile ? "8px 16px" : "10px 24px",
          textTransform: "none",
          fontWeight: 500,
        }}
      >
        {isSubmitting ? (
          <CircularProgress size={24} color="inherit" className="mr-2" />
        ) : null}
        {submitLabel}
      </Button>
    </Grid>
  );
};

const CustomForm = ({
  onSubmit,
  defaultValues,
  children,
  className,
  isSubmitting = false,
  resetLabel,
  submitLabel,
  showSubmitButtons = true,
}) => {
  const methods = useForm({ defaultValues });

  const handleSubmit = useCallback(
    async (data) => {
      try {
        await onSubmit(data);
      } catch (error) {
        console.error("Form submission error:", error);
      }
    },
    [onSubmit]
  );

  return (
    <FormProvider {...methods}>
      <form
        onSubmit={methods.handleSubmit(handleSubmit)}
        className={twMerge("flex flex-col gap-3", className)}
        noValidate
      >
        <Grid container spacing={2}>
          {children}
          {showSubmitButtons && (
            <SubmitButtons
              isSubmitting={isSubmitting}
              resetLabel={resetLabel}
              submitLabel={submitLabel}
            />
          )}
        </Grid>
      </form>
    </FormProvider>
  );
};

const CustomTablePaging = ({
  data,
  headers,
  renderRow,
  className,
  rowsPerPageOptions = [5, 10, 25],
}) => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(rowsPerPageOptions[0]);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  const handleChangePage = useCallback((e, newPage) => {
    setPg(newPage);
  }, []);

  const handleChangeRowsPerPage = useCallback((e) => {
    setRpg(parseInt(e.target.value, 10));
    setPg(0);
  }, []);

  if (!data || data.length === 0) {
    return (
      <Box className="flex justify-center items-center p-8 text-gray-500">
        No data available
      </Box>
    );
  }

  return (
    <Box className="w-full">
      <TableContainer
        component={Paper}
        className={twMerge("overflow-x-auto rounded-lg shadow-sm", className)}
        sx={{
          "&::-webkit-scrollbar": {
            height: "8px",
          },
          "&::-webkit-scrollbar-thumb": {
            backgroundColor: theme.palette.grey[300],
            borderRadius: "4px",
          },
        }}
      >
        <Table size={isMobile ? "small" : "medium"}>
          <TableHead>
            <TableRow className="bg-gray-50">
              {headers.map((header, index) => (
                <TableCell
                  key={index}
                  className="font-semibold"
                  sx={{
                    whiteSpace: "nowrap",
                    padding: isMobile ? "8px" : "16px",
                  }}
                >
                  {header}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {data
              .slice(pg * rpg, pg * rpg + rpg)
              .map((row, rowIndex) => renderRow(row, pg * rpg + rowIndex))}
          </TableBody>
        </Table>
      </TableContainer>

      <TablePagination
        component="div"
        count={data.length}
        page={pg}
        onPageChange={handleChangePage}
        rowsPerPage={rpg}
        onRowsPerPageChange={handleChangeRowsPerPage}
        rowsPerPageOptions={rowsPerPageOptions}
        className="border-t border-gray-200"
        labelRowsPerPage={isMobile ? "Rows:" : "Rows per page:"}
        sx={{
          ".MuiTablePagination-selectLabel, .MuiTablePagination-displayedRows":
            {
              fontSize: isMobile ? "0.75rem" : "0.875rem",
            },
        }}
      />
    </Box>
  );
};

export {
  CustomForm,
  CustomTablePaging,
  NumberInput,
  SelectionInput,
  SubmitButtons,
  TextInput,
};
