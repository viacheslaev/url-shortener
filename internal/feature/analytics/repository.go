package analytics

import (
	"context"
)

type AnalyticsRepository interface {
	InsertClick(ctx context.Context, c Click) error
}
