import {
  Container,
  Typography,
  Box,
  useTheme,
  useMediaQuery,
  Paper,
} from "@mui/material";

const CustomContainer = ({
  title,
  children,
  className,
  maxWidth = "lg",
  elevation = 2,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <Paper
      elevation={elevation}
      className={`rounded-lg overflow-hidden transition-shadow duration-300 hover:shadow-md ${
        className || ""
      }`}
      sx={{
        backgroundColor: theme.palette.background.paper,
      }}
    >
      <Container maxWidth={maxWidth} disableGutters className="flex flex-col">
        {title && (
          <Box
            className="p-3 md:p-4 border-b"
            sx={{
              borderColor: theme.palette.divider,
              backgroundColor: theme.palette.background.default,
            }}
          >
            <Typography
              variant={isMobile ? "h6" : "h5"}
              className="font-bold text-center"
              sx={{
                fontSize: isMobile ? "1.1rem" : "1.5rem",
                lineHeight: 1.3,
              }}
            >
              {title}
            </Typography>
          </Box>
        )}

        <Box className="p-3 md:p-6">{children}</Box>
      </Container>
    </Paper>
  );
};

export default CustomContainer;
