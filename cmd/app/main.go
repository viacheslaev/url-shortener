package main

import (
	"log"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server"
	"github.com/viacheslaev/url-shortener/internal/server/middleware"
)

func main() {
	cfg := config.Load()

	service := link.NewURLService()
	handler := link.NewURLHandler(cfg, service)

	router := middleware.Logging(server.NewRouter(cfg, handler))

	log.Printf("listening on %s\n", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, router))
}
