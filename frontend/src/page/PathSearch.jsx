import {
  Box,
  Button,
  Grid2,
  TableContainer,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  TablePagination,
  Paper,
  CircularProgress,
} from "@mui/material";
import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { SelectionInput, SubmitButtons } from "../component/Common";
import CustomContainer from "../component/CustomContainer";
import useFetchData from "../hook/useFetchData";
import { useSelector } from "react-redux";
import { getAllServices } from "../redux/services/selector";
import { Link } from "react-router-dom";
import axios from "../config/axios";
const PathSearch = () => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  const { data, loading, error, fetchData } = useFetchData("/path-analystic");
  const [calledOpOpt, setCalledOpOpt] = useState([]);
  const [callerOpOpt, setCallerOpOpt] = useState([]);

  const services = useSelector(getAllServices);

  const methods = useForm({
    defaultValues: {
      called_svc: "",
      called_op: "",
      caller_svc: "",
      caller_op: "",
    },
  });
  const onSubmit = (data) => {
    console.log(data);
    setPg(0);
    fetchData(data);
  };
  const svcOpt = services.map((s) => [s, s]);
  const calledSvc = methods.watch("called_svc");
  const callerSvc = methods.watch("caller_svc");

  const fetchCallerOperation = async (callerSvc) => {
    try {
      const res = await axios.get(`/services/${callerSvc}/operations`);
      setCallerOpOpt(res.data.map((op) => [op, op]));
    } catch (error) {
      console.log(error);
    }
  };

  const fetchCalledOperation = async (calledSvc) => {
    try {
      const res = await axios.get(`/services/${calledSvc}/operations`);
      setCalledOpOpt(res.data.map((op) => [op, op]));
    } catch (error) {
      console.log(error);
    }
  };
  useEffect(() => {
    if (calledSvc) {
      fetchCalledOperation(calledSvc);
    }
    if (callerSvc) {
      fetchCallerOperation(callerSvc);
    }
  }, [calledSvc, callerSvc]);
  let sortedData = null;
  if (data) {
    sortedData = data.sort((a, b) => b.longest_chain - a.longest_chain);
  }
  return (
    <Box
      className="flex flex-col gap-4 p-4"
      sx={{
        "& .MuiFormControl-root": {
          width: "100%",
        },
      }}
    >
      <FormProvider {...methods}>
        <form
          onSubmit={methods.handleSubmit(onSubmit)}
          className="flex flex-col gap-4"
        >
          <Grid2 container spacing={2}>
            <Grid2 item xs={12} sm={6}>
              <SelectionInput
                name={"caller_svc"}
                label="Caller service"
                labelId="caller-svc"
                options={svcOpt}
                className="w-full"
              />
            </Grid2>
            <Grid2 item xs={12} sm={6}>
              <SelectionInput
                name={"caller_op"}
                label="Caller operation"
                labelId="caller-op"
                options={callerOpOpt}
                className="w-full"
              />
            </Grid2>
            <Grid2 item xs={12} sm={6}>
              <SelectionInput
                name={"called_svc"}
                label="Called service"
                labelId="called-svc"
                options={svcOpt}
                className="w-full"
              />
            </Grid2>
            <Grid2 item xs={12} sm={6}>
              <SelectionInput
                name={"called_op"}
                label="Called operation"
                labelId="called-op"
                options={calledOpOpt}
                className="w-full"
              />
            </Grid2>
            <Grid2 item xs={12}>
              <SubmitButtons />
            </Grid2>
          </Grid2>
        </form>
      </FormProvider>
      <CustomContainer
        className="flex flex-col gap-4 p-4"
        sx={{
          "& .MuiPaper-root": {
            padding: 0,
          },
        }}
      >
        {loading && <CircularProgress className="m-auto" />}
        {error && <p className="text-red-500">{error.message}</p>}
        {sortedData && sortedData.length > 0 ? (
          <>
            <TableContainer component={Paper} className="overflow-x-auto">
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell className="px-4 py-2">Entry service</TableCell>
                    <TableCell className="px-4 py-2">Longest chain</TableCell>
                    <TableCell className="px-4 py-2">Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {sortedData
                    .slice(pg * rpg, pg * rpg + rpg)
                    .map((path, index) => (
                      <TableRow key={index}>
                        <TableCell className="px-4 py-2">
                          {path.tree_hop.service_name}
                        </TableCell>
                        <TableCell className="px-4 py-2">
                          {path.longest_chain}
                        </TableCell>
                        <TableCell className="px-4 py-2">
                          <Button
                            component={Link}
                            to={`/path-detail/${path.id}`}
                            className="text-blue-500 hover:text-blue-700"
                          >
                            Detail
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                </TableBody>
              </Table>
            </TableContainer>
            <TablePagination
              count={sortedData.length}
              onPageChange={(e, pg) => setPg(pg)}
              onRowsPerPageChange={(e) => {
                setRpg(parseInt(e.target.value, 10));
                setPg(0);
              }}
              page={pg}
              rowsPerPage={rpg}
              rowsPerPageOptions={[5, 10, 25]}
              className="flex justify-center"
            />
          </>
        ) : (
          <p className="text-center">{loading || "No data"}</p>
        )}
      </CustomContainer>
    </Box>
  );
};

export default PathSearch;
