import { useState } from "react";
import axios from "../config/axios";

const useFetchData = (url) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const fetchData = async (params) => {
    try {
      setLoading(true);
      const response = await axios.get(url, { params });
      setData(response.data);
      setLoading(false);
    } catch (error) {
      setError(error);
      setLoading(false);
    }
  };
  return { data, loading, error, fetchData };
};
export default useFetchData;
