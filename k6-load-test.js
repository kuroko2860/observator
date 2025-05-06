import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
  vus: 100, // 200 virtual users
  duration: "10m", // chạy 10 giây
  rps: 2000, // target 6000 request/s
};

let data = JSON.stringify({
  user_id: "user1",
  items: ["item1", "item2"],
});

export default function () {
  let res = http.post("http://checkout:8080/checkout", data, {
    headers: {
      "Content-Type": "application/json",
    },
  });
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(0.001);
}
