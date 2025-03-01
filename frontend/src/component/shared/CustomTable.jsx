import { useState } from "react";
import {
  Table,
  TableBody,
  TableHead,
  TableRow,
  TableSortLabel,
  Paper,
  TableCell,
  TableContainer,
  TablePagination,
} from "@mui/material";

const CustomTable = ({ headings, data, onRowClick }) => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);

  const [order, setOrder] = useState("desc");
  const [orderBy, setOrderBy] = useState("count");

  const handleSortBy = (name) => {
    const isAsc = orderBy === name && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy(name);
  };
  const getComparator = () => {
    return (a, b) => (a[orderBy] - b[orderBy]) * (order === "desc" ? -1 : 1);
  };

  return (
    <>
      <TableContainer component={Paper} className="rounded-lg">
        <Table>
          <TableHead>
            <TableRow>
              {headings.map(({ sortable, name, label }) =>
                sortable ? (
                  <TableCell key={name} className="px-6 py-3">
                    <TableSortLabel
                      active={orderBy === name}
                      direction={orderBy === name ? order : "desc"}
                      onClick={handleSortBy(name)}
                    >
                      {label}
                    </TableSortLabel>
                  </TableCell>
                ) : (
                  <TableCell key={name} className="px-6 py-3">
                    {label}
                  </TableCell>
                )
              )}
            </TableRow>
          </TableHead>
          <TableBody>
            {data
              .sort(getComparator())
              .slice(pg * rpg, pg * rpg + rpg)
              .map((rowData, index) => (
                <TableRow
                  key={index}
                  className="hover:bg-gray-100"
                  onClick={() => onRowClick(rowData)}
                >
                  {headings.map(({ name }) => (
                    <TableCell key={name} className="px-6 py-3">
                      {rowData[name]}
                    </TableCell>
                  ))}
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
        rowsPerPageOptions={[5, 10, 25]}
      />
    </>
  );
};

export default CustomTable;
