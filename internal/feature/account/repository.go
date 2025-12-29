package account

import "context"

type Repository interface {
	CreateAccount(ctx context.Context, email string, passwordHash string) (string, error)
	FindActiveAccountByEmail(ctx context.Context, email string) (*Account, error)
}
