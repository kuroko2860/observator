import { Container, Typography } from "@mui/material";

const CustomContainer = ({ title, children }) => {
  return (
    <Container className="bg-white p-4 rounded-lg shadow-md">
      <Typography variant="h5" className="text-2xl font-bold mb-4">
        {title}
      </Typography>
      {children}
    </Container>
  );
};

export default CustomContainer;
