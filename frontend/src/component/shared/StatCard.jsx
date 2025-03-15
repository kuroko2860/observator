import {
  Card,
  CardContent,
  Stack,
  Typography,
  Box,
  useTheme,
  useMediaQuery,
  Tooltip,
} from "@mui/material";

const StatCard = ({
  title,
  value,
  unit,
  className,
  tooltip,
  icon,
  trend,
  trendValue,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  // Format value for better display
  const formattedValue =
    typeof value === "number" && !Number.isInteger(value)
      ? value.toFixed(2)
      : value;

  // Determine trend color
  const getTrendColor = () => {
    if (!trend) return "inherit";
    return trend === "up"
      ? theme.palette.success.main
      : theme.palette.error.main;
  };

  return (
    <Card
      className={`transition-all duration-300 hover:shadow-lg ${
        className || ""
      }`}
      sx={{
        borderRadius: "12px",
        backgroundColor: theme.palette.background.paper,
        height: "100%",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <CardContent
        sx={{
          padding: isMobile ? "16px" : "20px",
          flex: 1,
          display: "flex",
          flexDirection: "column",
        }}
      >
        <Box className="flex justify-between items-center mb-3">
          <Typography
            variant={isMobile ? "subtitle1" : "h6"}
            component="h3"
            className="font-semibold text-gray-700"
            sx={{ fontSize: isMobile ? "1rem" : "1.1rem" }}
          >
            {title}
          </Typography>
          {icon && <Box className="text-gray-500">{icon}</Box>}
        </Box>

        <Stack
          direction="row"
          spacing={1}
          alignItems="baseline"
          className="mt-2"
        >
          <Tooltip title={tooltip || ""} arrow placement="top">
            <Typography
              variant={isMobile ? "h5" : "h4"}
              component="span"
              className="font-bold"
              sx={{
                fontSize: isMobile ? "1.5rem" : "2rem",
                lineHeight: 1.2,
              }}
            >
              {formattedValue}
            </Typography>
          </Tooltip>

          <Typography
            variant="body2"
            component="span"
            className="text-gray-500"
            sx={{ fontSize: isMobile ? "0.75rem" : "0.875rem" }}
          >
            {unit}
          </Typography>
        </Stack>

        {trend && (
          <Box className="mt-2 flex items-center">
            <Box
              component="span"
              className="flex items-center"
              sx={{ color: getTrendColor() }}
            >
              {trend === "up" ? "↑" : "↓"}
              <Typography
                variant="caption"
                component="span"
                sx={{
                  marginLeft: "4px",
                  color: getTrendColor(),
                  fontWeight: 500,
                }}
              >
                {trendValue}
              </Typography>
            </Box>
          </Box>
        )}
      </CardContent>
    </Card>
  );
};

export default StatCard;
