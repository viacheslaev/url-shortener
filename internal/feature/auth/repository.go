package auth

import (
	"context"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
)

type Repository interface {
	FindActiveAccountByEmail(ctx context.Context, email string) (*account.Account, error)
}
