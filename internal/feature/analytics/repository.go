package analytics

import (
	"context"
	"time"
)

type AnalyticsRepository interface {
	InsertClick(ctx context.Context, c Click) error
	GetStats(ctx context.Context, linkID int64, since time.Time) (Stats, error)
}

type LinkRepository interface {
	GetLinkByCodeAndAccountPublicId(ctx context.Context, code string, accountPublicId string) (int64, error)
}
