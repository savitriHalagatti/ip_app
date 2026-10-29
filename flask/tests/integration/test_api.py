"""Integration tests: call the real running app (local or Docker container) over HTTP.

Run:  BASE_URL=http://localhost:8000 pytest tests/integration
"""

import ipaddress
import os

import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8000")
PREFIX = "The IP address is: "


def test_root_returns_a_valid_ip():
    response = requests.get(f"{BASE_URL}/", timeout=5)
    assert response.status_code == 200
    assert response.text.startswith(PREFIX)
    ipaddress.ip_address(response.text[len(PREFIX) :])  # raises an error if it is not a real IP address


def test_greet():
    response = requests.get(f"{BASE_URL}/greet", timeout=5)
    assert response.status_code == 200
    assert response.text == "Hi how are you"
