package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	findActiveAccountByEmailFunc func(ctx context.Context, email string) (*account.Account, error)
}

func (m *mockAuthRepo) FindActiveAccountByEmail(ctx context.Context, email string) (*account.Account, error) {
	if m.findActiveAccountByEmailFunc != nil {
		return m.findActiveAccountByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func TestAuthService_Login_IssuesJWT(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}

	repo := &mockAuthRepo{findActiveAccountByEmailFunc: func(ctx context.Context, email string) (*account.Account, error) {
		if email != "user@example.com" {
			t.Fatalf("expected normalized email user@example.com, got %q", email)
		}
		return &account.Account{PublicID: "acc-1", PasswordHash: string(hash)}, nil
	}}

	issuer := &TokenIssuer{secret: []byte("secret"), ttl: time.Hour, issuer: "svc"}
	svc := NewAuthService(repo, issuer)

	tok, err := svc.Login(context.Background(), "  USER@Example.COM ", " pass123 ")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("expected valid jwt, err=%v", err)
	}
	if claims.Subject != "acc-1" {
		t.Fatalf("expected subject acc-1, got %q", claims.Subject)
	}
}

func TestAuthService_Login_InvalidCredentials_WhenAccountNotFound(t *testing.T) {
	repo := &mockAuthRepo{findActiveAccountByEmailFunc: func(ctx context.Context, email string) (*account.Account, error) {
		return nil, errors.New("not found")
	}}
	issuer := &TokenIssuer{secret: []byte("secret"), ttl: time.Hour, issuer: "svc"}
	svc := NewAuthService(repo, issuer)

	_, err := svc.Login(context.Background(), "user@example.com", "pass123")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_InvalidCredentials_WhenPasswordWrong(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	repo := &mockAuthRepo{findActiveAccountByEmailFunc: func(ctx context.Context, email string) (*account.Account, error) {
		return &account.Account{PublicID: "acc-1", PasswordHash: string(hash)}, nil
	}}
	issuer := &TokenIssuer{secret: []byte("secret"), ttl: time.Hour, issuer: "svc"}
	svc := NewAuthService(repo, issuer)

	_, err = svc.Login(context.Background(), "user@example.com", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
