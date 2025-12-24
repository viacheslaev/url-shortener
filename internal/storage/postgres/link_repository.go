package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Save(ctx context.Context, link link.ShortLink) error {
	const query = `
		INSERT INTO links (code, long_url, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, link.Code, link.LongURL, link.ExpiresAt)
	return err
}

func (r *LinkRepository) GetLongURL(ctx context.Context, code string) (string, error) {
	const query = `
		SELECT long_url
		FROM links
		WHERE code = $1
	`
	var longURL string
	err := r.db.QueryRowContext(ctx, query, code).Scan(&longURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", link.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return longURL, nil
}
