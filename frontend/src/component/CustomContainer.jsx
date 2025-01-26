import { Container, Typography } from "@mui/material";

const CustomContainer = ({ title, children }) => {
  return (
    <Container className="custom-container">
      <Typography variant="h5">{title}</Typography>
      {children}
    </Container>
  );
};

export default CustomContainer;
