package account

import "time"

type Account struct {
	ID           int64
	PublicID     string // UUID public (API / JWT)
	Email        string
	PasswordHash string

	IsActive  bool
	DeletedAt *time.Time
	CreatedAt time.Time
}
