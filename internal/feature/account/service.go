package account

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (service *AccountService) RegisterAccount(ctx context.Context, email string, regPassword string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	regPassword = strings.TrimSpace(regPassword)

	if !isValidEmail(email) {
		return "", fmt.Errorf("email validation error: %w", ErrEmailInvalidFormat)
	}

	if !isValidRegistrationPassword(regPassword) {
		return "", fmt.Errorf("password validation error: %w", ErrPasswordInvalidFormat)
	}

	passwordHash, err := generatePasswordHash(regPassword)
	if err != nil {
		return "", fmt.Errorf("password hashing failed: %w", err)
	}

	publicID, err := service.repo.CreateAccount(ctx, email, passwordHash)
	if err != nil {
		return "", fmt.Errorf("account creation failed: %w", err)
	}

	return publicID, nil
}

func generatePasswordHash(regPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(regPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password hash failed: %v", err)
		return "", fmt.Errorf("password hash failed")
	}
	return string(hash), nil
}
