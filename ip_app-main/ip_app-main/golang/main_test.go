// Unit tests: call the router in-process with httptest. No server, no network.
package main

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeNetworkCards replaces the real network cards for one test.
func fakeNetworkCards(t *testing.T, ips ...string) {
	t.Helper()
	original := interfaceAddrs
	t.Cleanup(func() { interfaceAddrs = original })
	interfaceAddrs = func() ([]net.Addr, error) {
		var addrs []net.Addr
		for _, ip := range ips {
			addrs = append(addrs, &net.IPNet{IP: net.ParseIP(ip), Mask: net.CIDRMask(24, 32)})
		}
		return addrs, nil
	}
}

func get(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestRootReturnsIPAsJSON(t *testing.T) {
	fakeNetworkCards(t, "127.0.0.1", "10.1.2.3")
	rec := get(t, http.MethodGet, "/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body["ip"] != "10.1.2.3" || body["message"] != "The IP address is: 10.1.2.3" {
		t.Errorf("body = %v", body)
	}
}

func TestGetIPSkipsLoopbackAndIPv6(t *testing.T) {
	fakeNetworkCards(t, "127.0.0.1", "fe80::1", "192.168.1.20")
	if got := getIP(); got != "192.168.1.20" {
		t.Errorf("getIP() = %q, want 192.168.1.20", got)
	}
}

func TestGetIPFallsBackWhenNoNetwork(t *testing.T) {
	fakeNetworkCards(t, "127.0.0.1")
	if got := getIP(); got != "127.0.0.1" {
		t.Errorf("getIP() = %q, want 127.0.0.1", got)
	}
}

func TestGetIPFallsBackOnError(t *testing.T) {
	original := interfaceAddrs
	t.Cleanup(func() { interfaceAddrs = original })
	interfaceAddrs = func() ([]net.Addr, error) { return nil, errors.New("no network") }
	if got := getIP(); got != "127.0.0.1" {
		t.Errorf("getIP() = %q, want 127.0.0.1", got)
	}
}

func TestGreet(t *testing.T) {
	rec := get(t, http.MethodGet, "/greet")
	if rec.Code != http.StatusOK || rec.Body.String() != "{\"message\":\"Hi how are you\"}\n" {
		t.Errorf("got %d %q", rec.Code, rec.Body.String())
	}
}

func TestUnknownPathReturns404(t *testing.T) {
	if rec := get(t, http.MethodGet, "/does-not-exist"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestPostIsNotAllowed(t *testing.T) {
	if rec := get(t, http.MethodPost, "/"); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}
