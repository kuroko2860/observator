import {
  AppBar,
  Box,
  CssBaseline,
  Drawer,
  List,
  ListItem,
  Toolbar,
  Typography,
} from "@mui/material";
import { NavLink, useLocation, Outlet } from "react-router-dom";
import Popup from "reactjs-popup";

const drawerWidth = 240;
const pathMap = {
  "/": "Services",
  "/top-service": "Services",
  "/top-api": "API",
  "/api-statistic": "API",
  "/api-long": "API",
  "/api-called": "API",
  "/long-path": "Path",
  "/search-path": "Path",
};

function Layout() {
  const location = useLocation();
  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <Typography variant="h6" noWrap component="div">
            {pathMap[location.pathname] || "Dashboard"}
          </Typography>
        </Toolbar>
      </AppBar>

      <Drawer
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          "& .MuiDrawer-paper": {
            width: drawerWidth,
            boxSizing: "border-box",
          },
        }}
        variant="permanent"
        anchor="left"
      >
        <Toolbar />
        <Box sx={{ overflow: "auto" }}>
          <List component="nav">
            <Popup
              trigger={
                <ListItem>
                  <NavLink className={"custom-navlink"} to="">
                    Services
                  </NavLink>
                </ListItem>
              }
              position="right top"
              on={"hover"}
              closeOnDocumentClick
              mouseLeaveDelay={50}
              mouseEnterDelay={0}
              arrow={false}
            >
              <List className="menu">
                <MenuItem to="/" text="Services list" />
                <MenuItem to="/top-service" text="Top service" />
              </List>
            </Popup>
            <Popup
              trigger={
                <ListItem>
                  <NavLink className={"custom-navlink"} to="/api-statistic">
                    API
                  </NavLink>
                </ListItem>
              }
              position="right top"
              on={"hover"}
              closeOnDocumentClick
              mouseLeaveDelay={50}
              mouseEnterDelay={0}
              arrow={false}
            >
              <List className="menu">
                <MenuItem to="/top-api" text="Top API" />
                <MenuItem to="/api-statistic" text="API statistic" />
                <MenuItem to="/api-long" text="Long API" />
                <MenuItem to="/api-called" text="API called" />
              </List>
            </Popup>
            <Popup
              trigger={
                <ListItem>
                  <NavLink className={"custom-navlink"} to="/search-path">
                    Paths
                  </NavLink>
                </ListItem>
              }
              position="right top"
              on={"hover"}
              closeOnDocumentClick
              mouseLeaveDelay={50}
              mouseEnterDelay={0}
              arrow={false}
            >
              <List className="menu">
                <MenuItem to="/long-path" text="Long path" />
                <MenuItem to="/search-path" text="Search path" />
              </List>
            </Popup>
          </List>
        </Box>
      </Drawer>

      <Box
        component="main"
        sx={{
          backgroundColor: (theme) =>
            theme.palette.mode === "light"
              ? theme.palette.grey[100]
              : theme.palette.grey[900],
          flexGrow: 1,
          height: "100vh",
          overflow: "auto",
        }}
      >
        <Toolbar />
        <Box>
          <Outlet />
        </Box>
      </Box>
    </Box>
  );
}
const MenuItem = ({ to, text }) => (
  <ListItem button>
    <NavLink className={"custom-navlink"} to={to}>
      {text}
    </NavLink>
  </ListItem>
);

export default Layout;
