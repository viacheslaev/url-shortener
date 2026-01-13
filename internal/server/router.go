package server

import (
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"github.com/viacheslaev/url-shortener/internal/feature/analytics"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server/middleware"
)

func NewRouter(
	urlHandler *link.URLHandler,
	accRegisterHandler *account.RegisterHandler,
	authHandler *auth.AuthHandler,
	analyticsHandler *analytics.AnalyticsHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	// JWT authorization
	mux.Handle("POST /api/v1/urls", authMiddleware.Authorize(http.HandlerFunc(urlHandler.CreateShortLink)))
	mux.Handle("GET /api/v1/links/{code}/stats", authMiddleware.Authorize(http.HandlerFunc(analyticsHandler.GetStats)))

	// Public
	mux.HandleFunc("GET /{code}", urlHandler.ResolveShortLink)
	mux.HandleFunc("POST /api/v1/auth/register", accRegisterHandler.RegisterAccount)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	return mux
}
