const express = require("express");
const os = require("os");

const port = 8000;

function getIp() {
  const nets = os.networkInterfaces();
  for (const name of Object.keys(nets)) {
    for (const net of nets[name]) {
      if (net.family === "IPv4" && !net.internal) {
        return net.address;
      }
    }
  }
  return "127.0.0.1";
}

const app = express();

app.get("/", (_req, res) => {
  const ip = getIp();
  res.json({ ip, message: `The IP address is: ${ip}` });
});

app.get("/greet", (_req, res) => {
  res.json({ message: "Hi how are you" });
});

// Start the server only when run directly (node app.js), not when tests import the app
if (require.main === module) {
  app.listen(port, "0.0.0.0", () => {
    console.log(`Listening on port ${port}`);
  });
}

module.exports = { app, getIp };
