package middleware

import (
	"context"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
)

type AccountRepository interface {
	FindActiveAccountByEmail(ctx context.Context, email string) (*account.Account, error)
	FindAccountStatusByPublicID(ctx context.Context, publicID string) (*account.AccountStatus, error)
}
