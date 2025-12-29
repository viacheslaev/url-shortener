package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

func NewRouter(urlHandler *link.URLHandler, accRegisterHandler *account.RegisterHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateShortLink)
	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)

	mux.HandleFunc("POST /api/v1/auth/register", accRegisterHandler.RegisterAccount)

	return mux
}
