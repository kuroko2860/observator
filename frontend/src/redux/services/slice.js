import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "../../config/axios";

export const fetchAllServices = createAsyncThunk("/services", async () => {
  const res = await axios.get("services");
  return res.data;
});
export const servicesSlice = createSlice({
  name: "services",
  unitialState: {
    services: [],
    status: "idle",
    error: null,
  },
  reducer: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchAllServices.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchAllServices.fulfilled, (state, action) => {
        state.status = "succeeded";
        state.services = action.payload;
      })
      .addCase(fetchAllServices.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.error.message;
      });
  },
});

export default servicesSlice.reducer;
