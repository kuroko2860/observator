import { Card, CardContent, Stack } from "@mui/material";

const StatCard = ({ title, value, unit }) => {
  return (
    <Card>
      <CardContent>{title}</CardContent>
      <Stack direction="column">
        <Stack>
          <h4>{value}</h4>
          <p>{unit}</p>
        </Stack>
      </Stack>
    </Card>
  );
};

export default StatCard;
