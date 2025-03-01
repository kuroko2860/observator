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
    <Box className="flex">
      <CssBaseline />
      <AppBar className="fixed top-0 left-0 z-10 sm:w-[calc(100%-240px)] sm:ml-240">
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
                  <NavLink className={"custom-navlink"} to="/api-statistics">
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
                <MenuItem to="/api-statistics" text="API statistic" />
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
                {/* <MenuItem to="/long-path" text="Long path" /> */}
                <MenuItem to="/search-path" text="Search path" />
                <MenuItem to="/path-detail/2298794235" text="Path detail" />
              </List>
            </Popup>
          </List>
        </Box>
      </Drawer>

      <Box
        component="main"
        className="bg-lightGrey-100 dark:bg-darkGrey-900 flex-1 h-full overflow-auto"
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
