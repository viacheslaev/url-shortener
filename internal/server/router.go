package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

func NewRouter(cfg *config.Config, urlHandler *link.URLHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateShortLink)
	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)

	return mux
}
