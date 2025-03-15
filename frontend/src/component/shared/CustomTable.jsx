import { useState, useCallback } from "react";
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
  Box,
  Typography,
  useTheme,
  useMediaQuery,
  Skeleton,
  Tooltip,
} from "@mui/material";

const CustomTable = ({
  headings,
  data = [],
  onRowClick,
  isLoading = false,
  emptyMessage = "No data available",
  className,
}) => {
  const [pg, setPg] = useState(0);
  const [rpg, setRpg] = useState(5);
  const [order, setOrder] = useState("desc");
  const [orderBy, setOrderBy] = useState("count");

  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));
  const isTablet = useMediaQuery(theme.breakpoints.down("md"));

  const handleSortBy = useCallback(
    (name) => () => {
      const isAsc = orderBy === name && order === "asc";
      setOrder(isAsc ? "desc" : "asc");
      setOrderBy(name);
    },
    [order, orderBy]
  );

  const getComparator = useCallback(() => {
    return (a, b) => {
      // Handle string comparison
      if (typeof a[orderBy] === "string") {
        return order === "asc"
          ? a[orderBy].localeCompare(b[orderBy])
          : b[orderBy].localeCompare(a[orderBy]);
      }
      // Handle numeric comparison
      return (a[orderBy] - b[orderBy]) * (order === "desc" ? -1 : 1);
    };
  }, [order, orderBy]);

  const handleChangePage = useCallback((e, newPage) => {
    setPg(newPage);
  }, []);

  const handleChangeRowsPerPage = useCallback((e) => {
    setRpg(parseInt(e.target.value, 10));
    setPg(0);
  }, []);

  // Render loading skeleton
  if (isLoading) {
    return (
      <Box className={`w-full ${className || ""}`}>
        <TableContainer component={Paper} className="rounded-lg shadow-sm">
          <Table size={isMobile ? "small" : "medium"}>
            <TableHead>
              <TableRow className="bg-gray-50">
                {headings.map(({ name, label }) => (
                  <TableCell
                    key={name}
                    className="px-3 md:px-6 py-2 md:py-3 font-semibold"
                  >
                    {label}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {[...Array(5)].map((_, index) => (
                <TableRow key={index}>
                  {headings.map(({ name }) => (
                    <TableCell key={name} className="px-3 md:px-6 py-2 md:py-3">
                      <Skeleton animation="wave" />
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    );
  }

  // Render empty state
  if (!data || data.length === 0) {
    return (
      <Box className={`w-full ${className || ""}`}>
        <Paper className="rounded-lg p-8 text-center shadow-sm">
          <Typography variant="body1" color="textSecondary">
            {emptyMessage}
          </Typography>
        </Paper>
      </Box>
    );
  }

  return (
    <Box className={`w-full ${className || ""}`}>
      <TableContainer
        component={Paper}
        className="rounded-lg shadow-sm overflow-x-auto"
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
              {headings.map(({ sortable, name, label }) =>
                sortable ? (
                  <TableCell
                    key={name}
                    className="px-3 md:px-6 py-2 md:py-3 font-semibold"
                    sx={{ whiteSpace: "nowrap" }}
                  >
                    <TableSortLabel
                      active={orderBy === name}
                      direction={orderBy === name ? order : "desc"}
                      onClick={handleSortBy(name)}
                    >
                      {label}
                    </TableSortLabel>
                  </TableCell>
                ) : (
                  <TableCell
                    key={name}
                    className="px-3 md:px-6 py-2 md:py-3 font-semibold"
                    sx={{ whiteSpace: "nowrap" }}
                  >
                    {label}
                  </TableCell>
                )
              )}
            </TableRow>
          </TableHead>
          <TableBody>
            {data
              .slice()
              .sort(getComparator())
              .slice(pg * rpg, pg * rpg + rpg)
              .map((rowData, index) => (
                <TableRow
                  key={index}
                  className="hover:bg-gray-50 transition-colors duration-150"
                  onClick={() => onRowClick && onRowClick(rowData)}
                  sx={{
                    cursor: onRowClick ? "pointer" : "default",
                    "&:last-child td, &:last-child th": { border: 0 },
                  }}
                >
                  {headings.map(({ name, render }) => (
                    <TableCell
                      key={name}
                      className="px-3 md:px-6 py-2 md:py-3"
                      sx={{
                        whiteSpace: isTablet ? "nowrap" : "normal",
                        maxWidth: isTablet ? "150px" : "300px",
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                      }}
                    >
                      {render ? (
                        render(rowData)
                      ) : (
                        <Tooltip
                          title={String(rowData[name])}
                          arrow
                          placement="top"
                        >
                          <Typography variant="body2" noWrap={isTablet}>
                            {rowData[name]}
                          </Typography>
                        </Tooltip>
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>

      <TablePagination
        component="div"
        count={data.length}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        page={pg}
        rowsPerPage={rpg}
        rowsPerPageOptions={[5, 10, 25]}
        labelRowsPerPage={isMobile ? "Rows:" : "Rows per page:"}
        className="border-t border-gray-200"
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

export default CustomTable;
