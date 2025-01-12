import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import DashboardLayout from "./layout/DashboardLayout";
import Home from "./pages/Home";
import Trace from "./pages/Trace";
import Trending from "./pages/Trending";
import GraphVisualizer from "./pages/Path";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<DashboardLayout />}>
          <Route index element={<Home />} />
          <Route path="traces" element={<Trace />} />
          <Route path="trending" element={<Trending />} />
          <Route path="path/:id" element={<GraphVisualizer />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
