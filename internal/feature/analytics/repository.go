package analytics

import (
	"context"
	"time"
)

type AnalyticsRepository interface {
	InsertClick(ctx context.Context, c Click) error
	GetStats(ctx context.Context, linkID int64, since time.Time) (Stats, error)
}
