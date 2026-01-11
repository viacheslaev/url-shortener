package auth

import (
	"context"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        AccountRepository
	tokenIssuer *TokenIssuer
}

func NewAuthService(repo AccountRepository, tokenIssuer *TokenIssuer) *AuthService {
	return &AuthService{
		repo:        repo,
		tokenIssuer: tokenIssuer,
	}
}

func (service *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	account, err := service.repo.FindActiveAccountByEmail(ctx, email)
	if err != nil {
		log.Printf("login failed: email=%s not found or account disabled", email)
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)) != nil {
		log.Printf("login failed: account=%s reason=wrong password", account.PublicID)
		return "", ErrInvalidCredentials
	}

	return service.tokenIssuer.IssueJWT(account.PublicID)
}
