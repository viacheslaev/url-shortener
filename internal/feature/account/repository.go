package account

import "context"

type AccountRepository interface {
	CreateAccount(ctx context.Context, email string, passwordHash string) (string, error)
}
