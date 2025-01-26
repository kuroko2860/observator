import {
  Box,
  Button,
  Grid2,
  TableContainer,
  CircularProgress,
  TablePagination,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from "@mui/material";
import { FormProvider, useForm } from "react-hook-form";
import { ServiceNameInput, ThresholdInput } from "../component/Input";
import { SubmitButtons } from "../component/Common";
import CustomContainer from "../component/CustomContainer";
import { Link } from "react-router-dom";
import { useState } from "react";

import useFetchData from "../hook/useFetchData";

const LongPath = () => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  const { data, error, loading, fetchData } = useFetchData(
    "/path-analysistic/long"
  );
  const methods = useForm({
    defaultValues: {
      threshold: null,
      service_name: "",
    },
  });
  const onSubmit = async (data) => {
    setPg(0);
    const params = {
      ...data,
    };
    await fetchData(params);
  };
  let sortedData = null;
  if (data) {
    sortedData = data
      .sort((a, b) => b.longest_chain - a.longest_chain)
      .filter((item) =>
        item.tree_hop.service_name.includes(methods.getValues("service_name"))
      );
  }
  return (
    <Box>
      <Typography variant="h4">Longest Path</Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2}>
            <ThresholdInput name="threshold" />
            <ServiceNameInput label={"Entry service"} />
            <SubmitButtons />
          </Grid2>
        </form>
      </FormProvider>
      <CustomContainer>
        {loading && <CircularProgress />}
        {error && <Typography>{error}</Typography>}
        {sortedData && sortedData.length > 0 ? (
          <>
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Entry Service</TableCell>
                    <TableCell>Longest Chain</TableCell>
                    <TableCell>Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {sortedData
                    .slice(pg * rpg, pg * rpg + rpg)
                    .map((path, index) => (
                      <TableRow
                        key={index}
                        sx={{
                          "&:last-child td, &:last-child th": { border: 0 },
                        }}
                      >
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
                setRpg(parseInt(e.target.value), 10);
                setPg(0);
              }}
              page={pg}
              rowsPerPage={rpg}
              rowsPerPageOptions={[5, 10, 25]}
            />
          </>
        ) : (
          <Typography variant="h5">{loading || "No data"}</Typography>
        )}
      </CustomContainer>
    </Box>
  );
};

export default LongPath;
