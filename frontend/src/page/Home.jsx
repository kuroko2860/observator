import {
  Container,
  TableContainer,
  Typography,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  TablePagination,
  Paper,
} from "@mui/material";
import CustomContainer from "../component/CustomContainer";
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
    <Container>
      <Typography variant="h4">Services</Typography>
      <CustomContainer>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>STT</TableCell>
                <TableCell>Service Name</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {services.slice(pg * rpg, pg * rpg + rpg).map((s, index) => (
                <TableRow
                  key={index}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                  onClick={() => navigate(`/service-detail/${s}`)}
                >
                  <TableCell>{pg * rpg + index + 1}</TableCell>
                  <TableCell>{s}</TableCell>
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
        />
      </CustomContainer>
    </Container>
  );
}

export default Home;
