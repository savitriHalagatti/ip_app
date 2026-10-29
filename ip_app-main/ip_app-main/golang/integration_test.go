//go:build integration

// Integration tests: call the real running app (local or Docker container) over HTTP.
// Run:  BASE_URL=http://localhost:8000 go test -tags=integration -v ./...
package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
)

func baseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8000"
}

func TestLiveRootReturnsAValidIP(t *testing.T) {
	resp, err := http.Get(baseURL() + "/")
	if err != nil {
		t.Fatalf("app not reachable: %v", err)
	}
	defer resp.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if resp.StatusCode != http.StatusOK || net.ParseIP(body["ip"]) == nil {
		t.Errorf("status %d, body %v", resp.StatusCode, body)
	}
	if body["message"] != "The IP address is: "+body["ip"] {
		t.Errorf("message = %q", body["message"])
	}
}

func TestLiveGreet(t *testing.T) {
	resp, err := http.Get(baseURL() + "/greet")
	if err != nil {
		t.Fatalf("app not reachable: %v", err)
	}
	defer resp.Body.Close()
	text, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || string(text) != "{\"message\":\"Hi how are you\"}\n" {
		t.Errorf("got %d %q", resp.StatusCode, text)
	}
}
