import { Card, CardContent, Stack } from "@mui/material";

const StatCard = ({ title, value, unit }) => {
  return (
    <Card className="shadow-md rounded-lg p-4">
      <CardContent className="text-lg font-bold">{title}</CardContent>
      <Stack direction="column" className="space-y-2">
        <Stack className="flex items-center space-x-2">
          <h4 className="text-3xl">{value}</h4>
          <p className="text-sm">{unit}</p>
        </Stack>
      </Stack>
    </Card>
  );
};

export default StatCard;
