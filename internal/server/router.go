package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server/middleware"
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
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/urls", authMiddleware.Authorize(http.HandlerFunc(urlHandler.CreateShortLink)))

	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)
	mux.HandleFunc("POST /api/v1/auth/register", accRegisterHandler.RegisterAccount)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	return mux
}
