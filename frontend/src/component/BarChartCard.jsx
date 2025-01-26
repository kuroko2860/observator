import { Card, CardContent, Stack, Typography } from "@mui/material";

function BarChartCard({ title, caption, children }) {
  return (
    <Card variant="outlined" sx={{ height: "100%", flexGrow: 1 }}>
      <CardContent>
        <Typography variant="h5" textAlign="center">
          {title}
        </Typography>
      </CardContent>
      <Stack
        direction="column"
        sx={{
          justifyContent: "space-between",
          alignItems: "center",
          flexGrow: 1,
          gap: 1,
        }}
      >
        {children}
        <Typography
          variant="caption"
          sx={{ color: "text.secondary", textAlign: "center" }}
        >
          {caption}
        </Typography>
      </Stack>
    </Card>
  );
}

export default BarChartCard;
