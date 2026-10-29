import socket

from fastapi import FastAPI

app = FastAPI()


@app.get("/")
def get_ip():
    hostname = socket.gethostname()
    ip = socket.gethostbyname(hostname)
    return {"ip": ip, "message": f"The IP address is: {ip}"}


@app.get("/greet")
def say_hello():
    return {"message": "Hi how are you"}
