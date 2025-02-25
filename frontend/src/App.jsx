import { BrowserRouter, Route, Routes } from "react-router-dom";
import Home from "./page/Home";
import { useDispatch, useSelector } from "react-redux";
import { useEffect } from "react";
import { fetchAllServices } from "./redux/services/slice";
import { Box, CircularProgress } from "@mui/material";
import Error from "./page/Error";
import Layout from "./layout/Layout";
import TopApi from "./page/TopApi";
import ApiStatistic from "./page/ApiStatistic";
import ApiLong from "./page/ApiLong";
import ApiCalled from "./page/ApiCalled";
import Search from "./page/Search";
import ServiceDetail from "./page/ServiceDetail";
import TopService from "./page/TopService";
import LongPath from "./page/LongPath";
import PathSearch from "./page/PathSearch";
import PathDetail from "./page/PathDetail";
import NotFound from "./page/NotFound";

function App() {
  const dispatch = useDispatch();
  const { status, error } = useSelector((state) => state.services);
  useEffect(() => {
    if (status == "idle") {
      dispatch(fetchAllServices());
    }
  }, [dispatch, status]);
  if (status == "loading") {
    return (
      <Box>
        <CircularProgress />
      </Box>
    );
  }
  if (status == "failed") {
    return <Error error={error} />;
  }
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="top-api" element={<TopApi />} />
          <Route path="api-statistic" element={<ApiStatistic />} />
          <Route path="api-long" element={<ApiLong />} />
          <Route path="api-called" element={<ApiCalled />} />
          <Route path="search" element={<Search />} />
          <Route
            path="service-detail/:service_name"
            element={<ServiceDetail />}
          />
          <Route path="top-service" element={<TopService />} />
          <Route path="long-path" element={<LongPath />} />
          <Route path="search-path" element={<PathSearch />} />
          <Route path="path-detail/:path_id" element={<PathDetail />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
