package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

func NewRouter(svc *link.URLService) http.Handler {
	mux := http.NewServeMux()
	urlHandler := link.NewURLHandler(svc)

	mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateShortLink)
	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)

	return mux
}
