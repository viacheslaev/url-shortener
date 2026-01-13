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
		INSERT INTO links (code, long_url, expires_at, account_public_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, link.Code, link.LongURL, link.ExpiresAt, link.AccountPublicId)
	return err
}

// GetLongLink returns original URL for the given short code or ErrNotFound if the link does not exist.
func (r *LinkRepository) GetLongLink(ctx context.Context, code string) (link.LongLink, error) {
	const query = `
		SELECT id, long_url, expires_at
		FROM links
		WHERE code = $1
		`
	var longLink link.LongLink
	err := r.db.QueryRowContext(ctx, query, code).Scan(&longLink.Id, &longLink.LongURL, &longLink.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return link.LongLink{}, link.ErrNotFound
	}
	if err != nil {
		return link.LongLink{}, err
	}
	return longLink, nil
}

// GetLinkByCodeAndAccountPublicId returns internal link id only if the link belongs to the given account by account_public_id
func (r *LinkRepository) GetLinkByCodeAndAccountPublicId(ctx context.Context, code string, accountPublicId string) (int64, error) {
	const query = `
		SELECT id
		FROM links
		WHERE code = $1 AND account_public_id = $2
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query, code, accountPublicId).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, link.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DeleteExpiredLinks delete expired links using UTC timezone
func (r *LinkRepository) DeleteExpiredLinks(ctx context.Context) (int64, error) {
	const q = `
        DELETE FROM links
        WHERE expires_at IS NOT NULL
           AND (expires_at AT TIME ZONE 'UTC')::timestamptz <= (NOW() AT TIME ZONE 'UTC')::timestamptz;
    `
	res, err := r.db.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
