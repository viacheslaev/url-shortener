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
	linkHandler *link.LinkHandler,
	accRegisterHandler *account.RegisterHandler,
	authHandler *auth.AuthHandler,
	analyticsHandler *analytics.AnalyticsHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	// Swagger (public)
	mux.Handle("/swagger/", SwaggerHandler())

	// Auth (public)
	mux.HandleFunc("POST /api/v1/auth/register", accRegisterHandler.RegisterAccount)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Protected
	mux.Handle("POST /api/v1/urls", authMiddleware.Authorize(http.HandlerFunc(linkHandler.CreateShortLink)))
	mux.Handle("GET /api/v1/links/{code}/stats", authMiddleware.Authorize(http.HandlerFunc(analyticsHandler.GetStats)))

	// Public redirect
	mux.HandleFunc("GET /{code}", linkHandler.ResolveShortLink)

	return mux
}
