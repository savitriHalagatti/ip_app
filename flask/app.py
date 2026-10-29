import socket

from flask import Flask

app = Flask(__name__)


@app.route("/")
def get_ip():
    hostname = socket.gethostname()
    ip = socket.gethostbyname(hostname)
    return f"The IP address is: {ip}"


@app.route("/greet")
def say_hello():
    return "Hi how are you"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000)
