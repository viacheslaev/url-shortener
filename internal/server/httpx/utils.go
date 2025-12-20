package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

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
