import dayjs from "dayjs";

export const ApiStatisticDefault = {
  service_name: "",
  uri_path: "",
  method: "",
  unit: "",
  from: null,
  to: null,
};
export const DefaultFromTo = {
  to: dayjs().add(1, "minute").second(0).millisecond(0).valueOf(),
  from: dayjs()
    .add(1, "minute")
    .second(0)
    .millisecond(0)
    .subtract(1, "hour")
    .valueOf(),
};
export const DefaultUnit = {
  unit: "minute",
};
