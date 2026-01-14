package account

import (
	"context"
	"errors"
	"testing"
)

type mockAccountRepo struct {
	createAccountFunc func(ctx context.Context, email string, passwordHash string) (string, error)
}

func (m *mockAccountRepo) CreateAccount(ctx context.Context, email, passwordHash string) (string, error) {
	if m.createAccountFunc != nil {
		return m.createAccountFunc(ctx, email, passwordHash)
	}
	return "", errors.New("not implemented")
}

func TestAccountService_Register_NormalizesAndCreatesAccount(t *testing.T) {
	repo := &mockAccountRepo{createAccountFunc: func(ctx context.Context, email string, passwordHash string) (string, error) {
		if email != "user@example.com" {
			t.Fatalf("expected normalized email user@example.com, got %q", email)
		}
		if passwordHash == "" {
			t.Fatalf("expected password hash to be set")
		}
		return "public-id-1", nil
	}}

	svc := NewAccountService(repo)

	publicID, err := svc.RegisterAccount(context.Background(), "USER@Example.COM ", "  secret1 ")

	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if publicID != "public-id-1" {
		t.Fatalf("unexpected public id: %s", publicID)
	}
}

func TestAccountService_Register_InvalidEmail(t *testing.T) {
	svc := NewAccountService(&mockAccountRepo{})

	_, err := svc.RegisterAccount(context.Background(), "not-an-email", "secret1")

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAccountService_Register_WeakPassword(t *testing.T) {
	svc := NewAccountService(&mockAccountRepo{})

	_, err := svc.RegisterAccount(context.Background(), "user@example.com", "123")

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAccountService_Register_EmailAlreadyExists(t *testing.T) {
	repo := &mockAccountRepo{createAccountFunc: func(ctx context.Context, email string, passwordHash string) (string, error) {
		return "", ErrEmailAlreadyExists
	}}
	svc := NewAccountService(repo)

	_, err := svc.RegisterAccount(context.Background(), "user@example.com", "secret1")

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}
