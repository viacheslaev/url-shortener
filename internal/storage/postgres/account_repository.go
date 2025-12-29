package postgres

import (
	"context"
	"database/sql"

	"github.com/viacheslaev/url-shortener/internal/feature/account"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, email string, passwordHash string) (string, error) {
	const q = `
		INSERT INTO accounts (email, password_hash)
		VALUES ($1, $2)
		RETURNING public_id
	`
	var publicId string
	err := r.db.QueryRowContext(ctx, q, email, passwordHash).Scan(&publicId)
	if err != nil {
		return "", err
	}
	return publicId, nil
}

func (r *AccountRepository) FindActiveAccountByEmail(ctx context.Context, email string) (*account.Account, error) {
	const q = `
		SELECT public_id, password_hash
		FROM accounts
		WHERE email = $1
		AND deleted_at IS NULL
		AND is_active = true
     
	`
	var a account.Account
	err := r.db.QueryRowContext(ctx, q, email).Scan(&a.PublicID, &a.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
