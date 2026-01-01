package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

// NewRouter wires HTTP routes.
//
// Requires JWT:
//   - POST /api/v1/urls
//
// Public:
//   - GET /{code}     (redirect)
//   - POST /api/v1/auth/register
//   - POST /api/v1/auth/login
func NewRouter(
	urlHandler *link.URLHandler,
	accRegisterHandler *account.RegisterHandler,
	authHandler *auth.AuthHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateShortLink)

	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)
	mux.HandleFunc("POST /api/v1/auth/register", accRegisterHandler.RegisterAccount)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	return mux
}
