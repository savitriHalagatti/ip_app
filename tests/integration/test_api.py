"""Integration tests: call the real running app (local or Docker container) over HTTP.

Run:  BASE_URL=http://localhost:8000 pytest tests/integration
"""

import ipaddress
import os

import httpx

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8000")


def test_root_returns_a_valid_ip():
    response = httpx.get(f"{BASE_URL}/", timeout=5)
    assert response.status_code == 200
    body = response.json()
    ipaddress.ip_address(body["ip"])  # raises an error if it is not a real IP address
    assert body["message"] == f"The IP address is: {body['ip']}"


def test_greet():
    response = httpx.get(f"{BASE_URL}/greet", timeout=5)
    assert response.status_code == 200
    assert response.json() == {"message": "Hi how are you"}


def test_swagger_docs_are_served():
    assert httpx.get(f"{BASE_URL}/docs", timeout=5).status_code == 200
