package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
)

const port = "8000"

// interfaceAddrs is a variable so unit tests can replace the real network cards.
var interfaceAddrs = net.InterfaceAddrs

// getIP returns the first non-loopback IPv4 address of this machine (container).
func getIP() string {
	addrs, err := interfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "127.0.0.1"
}

func writeJSON(w http.ResponseWriter, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

// newRouter sets up the URLs. Tests use it without starting a real server.
func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		ip := getIP()
		writeJSON(w, map[string]string{"ip": ip, "message": "The IP address is: " + ip})
	})

	mux.HandleFunc("GET /greet", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"message": "Hi how are you"})
	})

	return mux
}

func main() {
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, newRouter()))
}
