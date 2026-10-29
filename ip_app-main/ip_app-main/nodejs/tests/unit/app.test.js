// Unit tests: call the app in-process with supertest. No server, no network.
const os = require("os");
const request = require("supertest");
const { app, getIp } = require("../../app");

afterEach(() => jest.restoreAllMocks());

// Replace the real network cards so the test result is always the same
function fakeNetworkCards(cards) {
  jest.spyOn(os, "networkInterfaces").mockReturnValue(cards);
}

test("GET / returns the IP as JSON", async () => {
  fakeNetworkCards({ eth0: [{ family: "IPv4", internal: false, address: "10.1.2.3" }] });
  const response = await request(app).get("/");
  expect(response.status).toBe(200);
  expect(response.headers["content-type"]).toMatch(/application\/json/);
  expect(response.body).toEqual({ ip: "10.1.2.3", message: "The IP address is: 10.1.2.3" });
});

test("getIp skips loopback and IPv6 addresses", () => {
  fakeNetworkCards({
    lo: [{ family: "IPv4", internal: true, address: "127.0.0.1" }],
    eth0: [
      { family: "IPv6", internal: false, address: "fe80::1" },
      { family: "IPv4", internal: false, address: "192.168.1.20" },
    ],
  });
  expect(getIp()).toBe("192.168.1.20");
});

test("getIp falls back to 127.0.0.1 when there is no network", () => {
  fakeNetworkCards({ lo: [{ family: "IPv4", internal: true, address: "127.0.0.1" }] });
  expect(getIp()).toBe("127.0.0.1");
});

test("GET /greet", async () => {
  const response = await request(app).get("/greet");
  expect(response.status).toBe(200);
  expect(response.body).toEqual({ message: "Hi how are you" });
});

test("unknown path returns 404", async () => {
  expect((await request(app).get("/does-not-exist")).status).toBe(404);
});
