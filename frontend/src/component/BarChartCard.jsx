import { Card, CardContent, Stack, Typography } from "@mui/material";

function BarChartCard({ title, caption, children }) {
  return (
    <Card
      variant="outlined"
      className="h-full flex flex-col rounded-lg shadow-md"
      sx={{ height: "100%", flexGrow: 1 }}
    >
      <CardContent className="flex items-center justify-center p-4">
        <Typography
          variant="h5"
          className="text-center font-bold text-gray-800"
        >
          {title}
        </Typography>
      </CardContent>
      <Stack
        direction="column"
        className="flex flex-col items-center justify-between gap-1 p-4"
      >
        {children}
        <Typography variant="caption" className="text-center text-gray-500">
          {caption}
        </Typography>
      </Stack>
    </Card>
  );
}

export default BarChartCard;
