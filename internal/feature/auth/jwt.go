package auth

import (
	"fmt"
	"log"
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

// VerifyJWT validates and verifies JWT and returns trusted RegisteredClaims.
func VerifyJWT(tokenString, secret, issuer, audience string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %w", ErrUnauthorized)
			}
			return []byte(secret), nil
		},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		log.Printf("JWT parse error: %s", err)
		return nil, fmt.Errorf("jwt validation failed: %w", ErrUnauthorized)
	}

	return claims, nil
}
