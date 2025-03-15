import { Card, CardContent, Stack, Typography, useTheme, useMediaQuery } from "@mui/material";
import { memo } from "react";

function BarChartCard({ title, caption, children, className }) {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  return (
    <Card
      variant="outlined"
      className={`rounded-lg shadow-md overflow-hidden ${className || ''}`}
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        transition: "all 0.3s ease",
        "&:hover": {
          boxShadow: theme.shadows[4],
        },
      }}
    >
      <CardContent 
        className="flex items-center justify-center p-3 md:p-4"
        sx={{ 
          borderBottom: `1px solid ${theme.palette.divider}`,
          backgroundColor: theme.palette.background.default,
        }}
      >
        <Typography
          variant={isMobile ? "h6" : "h5"}
          className="text-center font-bold text-gray-800"
        >
          {title}
        </Typography>
      </CardContent>
      
      <Stack
        direction="column"
        className="flex-grow flex flex-col items-center justify-between p-3 md:p-4"
        sx={{ 
          minHeight: isMobile ? "250px" : "300px",
          overflow: "auto"
        }}
      >
        <div className="w-full h-full flex items-center justify-center">
          {children}
        </div>
        
        {caption && (
          <Typography 
            variant="caption" 
            className="text-center text-gray-500 mt-2 px-2"
            sx={{ fontSize: isMobile ? '0.7rem' : '0.75rem' }}
          >
            {caption}
          </Typography>
        )}
      </Stack>
    </Card>
  );
}

export default memo(BarChartCard);
