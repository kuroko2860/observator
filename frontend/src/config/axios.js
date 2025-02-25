import _axios from "axios";
const axios = _axios.create({ baseURL: "http://localhost:8585/api" });
export default axios;
