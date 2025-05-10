import {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  TablePagination,
  Paper,
} from "@mui/material";
import CustomContainer from "../component/shared/CustomContainer";
import { useSelector } from "react-redux";
import { getAllServices } from "../redux/services/selector";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
function Home() {
  const services = useSelector(getAllServices);
  const navigate = useNavigate();
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);

  return (
    <CustomContainer
      title={"Services list"}
      className="bg-white shadow-md rounded p-4"
    >
      {/* <CheckInOutSelector /> */}
      <TableContainer component={Paper} className="overflow-x-auto">
        <Table className="min-w-full">
          <TableHead>
            <TableRow>
              <TableCell className="px-6 py-3 border-b border-gray-200 bg-gray-50 text-center text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                STT
              </TableCell>
              <TableCell className="px-6 py-3 border-b border-gray-200 bg-gray-50 text-center text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                Service Name
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {services.slice(pg * rpg, pg * rpg + rpg).map((s, index) => (
              <TableRow
                key={index}
                className="hover:bg-gray-100 cursor-pointer"
                onClick={() => navigate(`/service-detail/${s}`)}
              >
                <TableCell className="px-6 py-4 whitespace-no-wrap border-b border-gray-200">
                  {pg * rpg + index + 1}
                </TableCell>
                <TableCell className="px-6 py-4 whitespace-no-wrap border-b border-gray-200">
                  {s}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        count={services.length}
        onPageChange={(e, pg) => setPg(pg)}
        onRowsPerPageChange={(e) => {
          setRpg(parseInt(e.target.value), 10);
          setPg(0);
        }}
        page={pg}
        rowsPerPage={rpg}
        rowsPerPageOptions={[5, 10, 25, 50]}
        className="mt-4"
      />
    </CustomContainer>
  );
}

export default Home;
