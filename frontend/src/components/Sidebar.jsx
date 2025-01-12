import { NavLink } from "react-router-dom";
import HomeIcon from "@mui/icons-material/Home";
import TimelineIcon from "@mui/icons-material/Timeline";
import SettingsIcon from "@mui/icons-material/Settings";

const Sidebar = () => {
  return (
    <div className="h-screen w-64 bg-gray-800 text-white flex flex-col">
      <div className="p-4 text-lg font-bold border-b border-gray-700">
        APM Dashboard
      </div>
      <nav className="flex-1 p-4 space-y-4">
        <NavLink
          to="/"
          className={({ isActive }) =>
            `flex items-center gap-3 p-2 rounded-md ${
              isActive ? "bg-gray-700" : "hover:bg-gray-700"
            }`
          }
        >
          <HomeIcon />
          Home
        </NavLink>
        <NavLink
          to="/traces"
          className={({ isActive }) =>
            `flex items-center gap-3 p-2 rounded-md ${
              isActive ? "bg-gray-700" : "hover:bg-gray-700"
            }`
          }
        >
          <TimelineIcon />
          Traces
        </NavLink>
        <NavLink
          to="/trending"
          className={({ isActive }) =>
            `flex items-center gap-3 p-2 rounded-md ${
              isActive ? "bg-gray-700" : "hover:bg-gray-700"
            }`
          }
        >
          <SettingsIcon />
          Settings
        </NavLink>
      </nav>
    </div>
  );
};

export default Sidebar;
