// Integration tests: call the real running app (local or Docker container) over HTTP.
// Run:  BASE_URL=http://localhost:8000 npm run test:integration
const net = require("net");

const BASE_URL = process.env.BASE_URL || "http://localhost:8000";

test("GET / returns a valid IP", async () => {
  const response = await fetch(`${BASE_URL}/`);
  expect(response.status).toBe(200);
  const body = await response.json();
  expect(net.isIP(body.ip)).not.toBe(0); // 0 means "not an IP address"
  expect(body.message).toBe(`The IP address is: ${body.ip}`);
});

test("GET /greet", async () => {
  const response = await fetch(`${BASE_URL}/greet`);
  expect(response.status).toBe(200);
  expect(await response.json()).toEqual({ message: "Hi how are you" });
});
