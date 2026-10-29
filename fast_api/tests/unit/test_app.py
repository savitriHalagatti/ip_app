"""Unit tests: call the API in-process with TestClient. No server, no network."""

from fastapi.testclient import TestClient

import app as app_module

client = TestClient(app_module.app)


def fake_ip(monkeypatch, ip="10.1.2.3"):
    """Replace the real IP lookup so the test result is always the same."""
    monkeypatch.setattr(app_module.socket, "gethostname", lambda: "test-host")
    monkeypatch.setattr(app_module.socket, "gethostbyname", lambda name: ip)


def test_root_returns_ip_as_json(monkeypatch):
    fake_ip(monkeypatch)
    response = client.get("/")
    assert response.status_code == 200
    assert response.headers["content-type"] == "application/json"
    assert response.json() == {"ip": "10.1.2.3", "message": "The IP address is: 10.1.2.3"}


def test_root_looks_up_this_machines_hostname(monkeypatch):
    looked_up = []
    monkeypatch.setattr(app_module.socket, "gethostname", lambda: "test-host")
    monkeypatch.setattr(app_module.socket, "gethostbyname", lambda name: looked_up.append(name) or "10.9.9.9")
    client.get("/")
    assert looked_up == ["test-host"]


def test_greet():
    response = client.get("/greet")
    assert response.status_code == 200
    assert response.json() == {"message": "Hi how are you"}


def test_unknown_path_returns_404():
    assert client.get("/does-not-exist").status_code == 404


def test_post_is_not_allowed():
    assert client.post("/").status_code == 405
