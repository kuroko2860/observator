import { Typography } from "@mui/material";

function NotFound() {
  return (
    <div className="d-flex justify-content-center align-items-center h-100 flex-column">
      <Typography variant="h1" component="h1" gutterBottom>
        404
      </Typography>
      <Typography variant="h5" component="h2" gutterBottom>
        Page Not Found
      </Typography>
    </div>
  );
}

export default NotFound;
