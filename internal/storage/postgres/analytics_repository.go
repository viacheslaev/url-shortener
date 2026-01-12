package postgres

import (
	"context"
	"database/sql"

	"github.com/viacheslaev/url-shortener/internal/feature/analytics"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) InsertClick(ctx context.Context, c analytics.Click) error {
	const q = `
		INSERT INTO link_clicks (link_id, ip_address, user_agent, referer, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, q, c.LinkID, c.IPAddress, c.UserAgent, c.Referer, c.CreatedAt)
	return err
}
