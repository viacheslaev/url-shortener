package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/viacheslaev/url-shortener/internal/config"
)

func testConfig() *config.Config {
	return &config.Config{
		JWTAccessTokenTTL: 24,
		ServiceName:       "url-shortener",
		JWTSecret:         "secret",
	}
}

func TestIssueAndVerifyJWT_OK(t *testing.T) {
	conf := testConfig()
	tokenIssuer := NewTokenIssuer(conf)

	token, err := tokenIssuer.IssueJWT("user-123")
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	claims, err := VerifyJWT(token, conf.JWTSecret, conf.ServiceName, "url-shortener-api")
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Fatalf("expected subject user-123, got %s", claims.Subject)
	}

	if claims.Issuer != conf.ServiceName {
		t.Fatalf("wrong tokenIssuer: %s", claims.Issuer)
	}
}

func TestVerifyJWT_InvalidSecret(t *testing.T) {
	conf := testConfig()
	issuer := NewTokenIssuer(conf)
	token, _ := issuer.IssueJWT("user-1")

	_, err := VerifyJWT(token, "wrong-secret", conf.ServiceName, "url-shortener-api")
	if err == nil {
		t.Fatal("expected ErrUnauthorized")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatal("expected ErrUnauthorized")
	}
}

func TestVerifyJWT_InvalidIssuer(t *testing.T) {
	conf := testConfig()
	issuer := NewTokenIssuer(conf)
	token, _ := issuer.IssueJWT("user-1")

	_, err := VerifyJWT(token, conf.JWTSecret, "wrong-issuer", "url-shortener-api")
	if err == nil {
		t.Fatal("expected ErrUnauthorized")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatal("expected ErrUnauthorized")
	}
}

func TestVerifyJWT_InvalidAudience(t *testing.T) {
	conf := testConfig()
	issuer := NewTokenIssuer(conf)
	token, _ := issuer.IssueJWT("user-1")

	_, err := VerifyJWT(token, conf.JWTSecret, conf.ServiceName, "wrong-api")

	if err == nil {
		t.Fatal("expected ErrUnauthorized")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatal("expected ErrUnauthorized")
	}
}

func TestVerifyJWT_Expired(t *testing.T) {
	cfg := testConfig()
	cfg.JWTAccessTokenTTL = -1 // expired ttl

	issuer := NewTokenIssuer(cfg)
	token, _ := issuer.IssueJWT("user-1")

	if _, err := VerifyJWT(token, cfg.JWTSecret, cfg.ServiceName, "url-shortener-api"); err == nil {
		t.Fatal("expected expired token error")
	}
}

func TestVerifyJWT_RejectsInvalidSigningMethod(t *testing.T) {
	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		Issuer:    testConfig().ServiceName,
		Audience:  []string{"url-shortener-api"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token, _ := tk.SignedString(jwt.UnsafeAllowNoneSignatureType)

	if _, err := VerifyJWT(token, testConfig().JWTSecret, testConfig().ServiceName, "url-shortener-api"); err == nil {
		t.Fatal("expected none-alg rejection")
	}
}
