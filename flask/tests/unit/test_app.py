"""Unit tests: call the app in-process with Flask's test client. No server, no network."""

import pytest

import app as app_module


@pytest.fixture
def client():
    app_module.app.config["TESTING"] = True
    return app_module.app.test_client()


def fake_ip(monkeypatch, ip="10.1.2.3"):
    """Replace the real IP lookup so the test result is always the same."""
    monkeypatch.setattr(app_module.socket, "gethostname", lambda: "test-host")
    monkeypatch.setattr(app_module.socket, "gethostbyname", lambda name: ip)


def test_root_returns_ip_as_text(client, monkeypatch):
    fake_ip(monkeypatch)
    response = client.get("/")
    assert response.status_code == 200
    assert response.mimetype == "text/html"
    assert response.get_data(as_text=True) == "The IP address is: 10.1.2.3"


def test_root_looks_up_this_machines_hostname(client, monkeypatch):
    looked_up = []
    monkeypatch.setattr(app_module.socket, "gethostname", lambda: "test-host")
    monkeypatch.setattr(app_module.socket, "gethostbyname", lambda name: looked_up.append(name) or "10.9.9.9")
    client.get("/")
    assert looked_up == ["test-host"]


def test_greet(client):
    response = client.get("/greet")
    assert response.status_code == 200
    assert response.get_data(as_text=True) == "Hi how are you"


def test_unknown_path_returns_404(client):
    assert client.get("/does-not-exist").status_code == 404


def test_post_is_not_allowed(client):
    assert client.post("/").status_code == 405
