import {
  Button,
  FormControl,
  Grid2 as Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TableContainer,
  TextField,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  TablePagination,
} from "@mui/material";
import {
  Controller,
  FormProvider,
  useForm,
  useFormContext,
} from "react-hook-form";
import { useState } from "react";
import { twMerge } from "tailwind-merge";

const TextInput = ({ name, label, className }) => {
  const methods = useFormContext();
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
        render={({ field }) => (
          <TextField
            {...field}
            fullWidth
            label={label}
            variant="outlined"
            className="p-2 border-2 border-gray-300 rounded-lg"
          />
        )}
      />
    </Grid>
  );
};

const NumberInput = ({ name, label, className }) => {
  const methods = useFormContext();
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
        render={({ field }) => (
          <TextField
            type="number"
            {...field}
            fullWidth
            label={label}
            variant="outlined"
            className="p-2 border-2 border-gray-300 rounded-lg"
          />
        )}
      />
    </Grid>
  );
};
const SelectionInput = ({ labelId, name, label, options, className }) => {
  const methods = useFormContext();
  return (
    <Grid
      item
      xs={12}
      sm={6}
      className={twMerge("flex flex-col gap-2", className)}
    >
      <FormControl fullWidth>
        <InputLabel id={labelId}>{label}</InputLabel>
        <Controller
          name={name}
          control={methods.control}
          render={({ field }) => (
            <Select
              {...field}
              labelId={labelId}
              id={name}
              style={{ minWidth: "160px" }}
              label={label}
              className="p-2 border-2 border-gray-300 rounded-lg"
            >
              {options.map(([value, label]) => (
                <MenuItem key={label} value={value}>
                  {label}
                </MenuItem>
              ))}
            </Select>
          )}
        />
      </FormControl>
    </Grid>
  );
};

const SubmitButtons = ({ className }) => {
  const methods = useFormContext();
  return (
    <Grid item className={twMerge("flex justify-end gap-2", className)}>
      <Button type="submit" variant="contained">
        Submit
      </Button>
      <Button
        type="button"
        onClick={() => methods.reset()}
        variant="outlined"
        className="border-2 border-gray-300 rounded-lg"
      >
        Reset
      </Button>
    </Grid>
  );
};
const CustomForm = ({ onSubmit, defaultValues, children, className }) => {
  const methods = useForm({ defaultValues });
  return (
    <FormProvider {...methods}>
      <form
        onSubmit={methods.handleSubmit(onSubmit)}
        className={twMerge("flex flex-col gap-2", className)}
      >
        {children}
        <SubmitButtons />
      </form>
    </FormProvider>
  );
};

const CustomTablePaging = ({ data, headers, children, className }) => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  return (
    <>
      <TableContainer
        component={Paper}
        className={twMerge("overflow-x-auto", className)}
      >
        <Table>
          <TableHead>
            <TableRow>
              {headers.map((h) => (
                <TableCell key={h} className="p-2 font-bold">
                  {h}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {data.slice(pg * rpg, pg * rpg + rpg).map((s, index) => (
              <TableRow key={index}>
                <TableCell className="p-2">{pg * rpg + index + 1}</TableCell>
                <TableCell className="p-2">{children}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        count={data.length}
        onPageChange={(e, pg) => setPg(pg)}
        onRowsPerPageChange={(e) => {
          setRpg(parseInt(e.target.value), 10);
          setPg(0);
        }}
        page={pg}
        rowsPerPage={rpg}
        rowsPerPageOptions={[5, 10]}
        className="p-2"
      ></TablePagination>
    </>
  );
};

export {
  TextInput,
  NumberInput,
  SelectionInput,
  SubmitButtons,
  CustomForm,
  CustomTablePaging,
};
