package account

import "context"

type AccountRepository interface {
	CreateAccount(ctx context.Context, email string, passwordHash string) (string, error)
	FindActiveAccountByEmail(ctx context.Context, email string) (*Account, error)
	FindAccountStatusByPublicID(ctx context.Context, publicID string) (*AccountStatus, error)
}
