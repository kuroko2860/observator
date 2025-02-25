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
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: 2,
        padding: 2,
      }}
    >
      <Typography variant="h4" className="text-2xl font-bold">
        Longest Path
      </Typography>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2} className="grid grid-cols-2 gap-4">
            <ThresholdInput name="threshold" />
            <ServiceNameInput label={"Entry service"} />
            <SubmitButtons className="col-span-2" />
          </Grid2>
        </form>
      </FormProvider>
      <CustomContainer className="bg-white rounded-md shadow-md p-4">
        {loading && <CircularProgress className="mx-auto" />}
        {error && <Typography className="text-red-500">{error}</Typography>}
        {sortedData && sortedData.length > 0 ? (
          <>
            <TableContainer component={Paper} className="overflow-auto">
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell className="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                      Entry Service
                    </TableCell>
                    <TableCell className="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                      Longest Chain
                    </TableCell>
                    <TableCell className="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                      Action
                    </TableCell>
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
                        className="hover:bg-gray-100"
                      >
                        <TableCell className="px-6 py-4 whitespace-no-wrap border-b border-gray-200">
                          {path.tree_hop.service_name}
                        </TableCell>
                        <TableCell className="px-6 py-4 whitespace-no-wrap border-b border-gray-200">
                          {path.longest_chain}
                        </TableCell>
                        <TableCell className="px-6 py-4 whitespace-no-wrap border-b border-gray-200">
                          <Button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
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
              className="mt-4"
            />
          </>
        ) : (
          <Typography variant="h5" className="text-gray-500 text-center">
            {loading || "No data"}
          </Typography>
        )}
      </CustomContainer>
    </Box>
  );
};

export default LongPath;
