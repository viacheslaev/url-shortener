package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type ctxKey string

const accountIDKey ctxKey = "account_id"

type AuthMiddleware struct {
	jwtSecret string
	issuer    string
	audience  string
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: cfg.JWTSecret,
		issuer:    cfg.ServiceName,
		audience:  cfg.ServiceName + "-" + "api",
	}
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.WriteErr(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenString := parts[1]
		claims, err := parseJWT(tokenString, m.jwtSecret, m.issuer, m.audience)
		if err != nil {
			httpx.WriteErr(w, http.StatusUnauthorized, auth.ErrUnauthorized.Error())
			return
		}

		ctx := context.WithValue(r.Context(), accountIDKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseJWT(tokenString string, secret string, issuer string, audience string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrUnauthorized
			}
			return []byte(secret), nil
		},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		log.Printf("JWT parse error: %s", err)
		return nil, auth.ErrUnauthorized
	}

	return claims, nil
}
