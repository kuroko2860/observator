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
    <Box>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2}>
            <SelectionInput
              name={"caller_svc"}
              label="Caller service"
              labelId="caller-svc"
              options={svcOpt}
            />
            <SelectionInput
              name={"caller_op"}
              label="Caller operation"
              labelId="caller-op"
              options={callerOpOpt}
            />
            <SelectionInput
              name={"called_svc"}
              label="Called service"
              labelId="called-svc"
              options={svcOpt}
            />
            <SelectionInput
              name={"called_op"}
              label="Called operation"
              labelId="called-op"
              options={calledOpOpt}
            />
            <SubmitButtons />
          </Grid2>
        </form>
      </FormProvider>
      <CustomContainer>
        {loading && <CircularProgress />}
        {error && <p>Error: {error.message}</p>}
        {sortedData && sortedData.length > 0 ? (
          <>
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Entry service</TableCell>
                    <TableCell>Longest chain</TableCell>
                    <TableCell>Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {sortedData
                    .slice(pg * rpg, pg * rpg + rpg)
                    .map((path, index) => (
                      <TableRow key={index}>
                        <TableCell>{path.tree_hop.service_name}</TableCell>
                        <TableCell>{path.longest_chain}</TableCell>
                        <TableCell>
                          <Button>
                            <Link to={`/path-detail/${path.id}`}>Detail</Link>
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
            />
          </>
        ) : (
          <p>{loading || "No data"}</p>
        )}
      </CustomContainer>
    </Box>
  );
};

export default PathSearch;
