package link

import "context"

type LinkRepository interface {
	Save(ctx context.Context, link ShortLink) error
	GetLongLink(ctx context.Context, code string) (LongLink, error)
	DeleteExpiredLinks(ctx context.Context) (int64, error)
}
