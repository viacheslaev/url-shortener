package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/viacheslaev/url-shortener/internal/config"
)

// TokenIssuer issues JWT tokens
type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

func NewTokenIssuer(cfg *config.Config) *TokenIssuer {
	return &TokenIssuer{
		secret: []byte(cfg.JWTSecret),
		ttl:    time.Duration(cfg.JWTAccessTokenTTL) * time.Hour,
		issuer: cfg.ServiceName}
}

func (issuer *TokenIssuer) IssueJWT(userID string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(issuer.ttl)),
		Subject:   userID,
		Issuer:    issuer.issuer,
		Audience:  []string{"url-shortener-api"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(issuer.secret)
}
