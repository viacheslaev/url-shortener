package middleware

import (
	"log"
	"net/http"
)

type logWriter struct {
	http.ResponseWriter
	status int
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logWriter := &logWriter{ResponseWriter: w}

		log.Printf("request %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(logWriter, r)

		log.Printf("response status %d", logWriter.status)
	})
}

func (lw *logWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}
