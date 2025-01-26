import {
  Button,
  FormControl,
  Grid2,
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

const TextInput = ({ name, label }) => {
  const methods = useFormContext();
  return (
    <Grid2 item="true" xs={12} sm={6}>
      <Controller
        name={name}
        control={methods.control}
        render={({ field }) => (
          <TextField {...field} fullWidth label={label} variant="outlined" />
        )}
      />
    </Grid2>
  );
};

const NumberInput = ({ name, label }) => {
  const methods = useFormContext();
  return (
    <Grid2 item="true" xs={12} sm={6}>
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
          />
        )}
      />
    </Grid2>
  );
};
const SelectionInput = ({ labelId, name, label, options }) => {
  const methods = useFormContext();
  return (
    <Grid2 item="true" xs={12} sm={6}>
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
    </Grid2>
  );
};

const SubmitButtons = () => {
  const methods = useFormContext();
  return (
    <Grid2 item="true">
      <Button type="submit">Submit</Button>
      <Button type="button" onClick={() => methods.reset()} variant="outlined">
        Reset
      </Button>
    </Grid2>
  );
};
const CustomForm = ({ onSubmit, defaultValues, children }) => {
  const methods = useForm({ defaultValues });
  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(onSubmit)}>
        {children}
        <SubmitButtons />
      </form>
    </FormProvider>
  );
};

const CustomTablePaging = ({ data, headers, children }) => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  return (
    <>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              {headers.map((h) => (
                <TableCell key={h}>{h}</TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {data.slice(pg * rpg, pg * rpg + rpg).map((s, index) => (
              <TableRow key={index}>
                <TableCell>{pg * rpg + index + 1}</TableCell>
                <TableCell>{children}</TableCell>
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
