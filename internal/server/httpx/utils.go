package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	// If behind a proxy, X-Forwarded-For may exist.
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return normalizeIP(parsed)
		}
	}

	// RemoteAddr is usually "IP:port" (IPv4) or "[IPv6]:port".
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if parsed := net.ParseIP(host); parsed != nil {
			return normalizeIP(parsed)
		}
		return host
	}

	if parsed := net.ParseIP(r.RemoteAddr); parsed != nil {
		return normalizeIP(parsed)
	}

	return r.RemoteAddr
}

func WriteErr(w http.ResponseWriter, status int, msg string) {
	if err := writeJSON(w, status, map[string]string{"error": msg}); err != nil {
		log.Printf("failed to write JSON error response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func WriteResponse(w http.ResponseWriter, status int, v any) {
	if err := writeJSON(w, status, v); err != nil {
		log.Printf("failed to write JSON response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_, err := w.Write(buf.Bytes())
	return err
}

func normalizeIP(ip net.IP) string {
	// Map IPv6 loopback to IPv4 loopback (for local/dev convenience)
	if ip.IsLoopback() {
		return "127.0.0.1"
	}

	// If it's an IPv4-mapped address or normal IPv4, return dotted form
	if v4 := ip.To4(); v4 != nil {
		return v4.String()
	}

	// Otherwise keep IPv6 string
	return ip.String()
}
