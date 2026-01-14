package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/viacheslaev/url-shortener/internal/feature/analytics"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) SaveClick(ctx context.Context, c analytics.Click) error {
	const q = `
		INSERT INTO link_clicks (link_id, ip_address, user_agent, referer, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, q, c.LinkID, c.IPAddress, c.UserAgent, c.Referer, c.CreatedAt)
	return err
}

func (r *AnalyticsRepository) GetStats(ctx context.Context, linkID int64, since time.Time) (analytics.Stats, error) {
	const totalCountQuery = `
		SELECT
			COUNT(*) AS total,
			COUNT(DISTINCT ip_address) AS unique
		FROM link_clicks
		WHERE link_id = $1
		  AND created_at >= $2
	`
	var total, unique int64
	if err := r.db.QueryRowContext(ctx, totalCountQuery, linkID, since).Scan(&total, &unique); err != nil {
		return analytics.Stats{}, err
	}

	const countByDayQuery = `
		SELECT DATE(created_at) AS d, COUNT(*)
		FROM link_clicks
		WHERE link_id = $1
		  AND created_at >= $2
		GROUP BY d
		ORDER BY d
	`
	rows, err := r.db.QueryContext(ctx, countByDayQuery, linkID, since)
	if err != nil {
		return analytics.Stats{}, err
	}
	defer rows.Close()

	byDay := make([]analytics.DayCount, 0)
	for rows.Next() {
		var d time.Time
		var c int64
		if err := rows.Scan(&d, &c); err != nil {
			return analytics.Stats{}, err
		}
		byDay = append(byDay, analytics.DayCount{Date: d, Count: c})
	}
	if err := rows.Err(); err != nil {
		return analytics.Stats{}, err
	}

	return analytics.Stats{TotalClicks: total, UniqueClicks: unique, ByDay: byDay}, nil
}
