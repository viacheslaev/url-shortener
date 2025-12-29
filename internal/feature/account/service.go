package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	repo Repository
}

func NewAccountService(repo Repository) *AccountService {
	return &AccountService{repo: repo}
}

func (service *AccountService) Register(ctx context.Context, email string, regPassword string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	regPassword = strings.TrimSpace(regPassword)

	if !isValidEmail(email) {
		return "", errors.New("invalid email")
	}

	if !isValidRegistrationPassword(regPassword) {
		return "", errors.New("password must be at least 6 characters")
	}

	passwordHash, err := generatePasswordHash(regPassword)
	if err != nil {
		return "", err
	}

	publicId, err := service.repo.CreateAccount(ctx, email, passwordHash)
	if err != nil {
		// todo: refactor, make func transactional
		// Handle unique violation (email already exists).
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return "", ErrEmailAlreadyExists
		}

		return "", err
	}

	return publicId, nil
}

func generatePasswordHash(regPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(regPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password hash failed: %v", err)
		return "", fmt.Errorf("password hash failed")
	}
	return string(hash), nil
}
