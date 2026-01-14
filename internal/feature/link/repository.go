package link

import "context"

type LinkRepository interface {
	CreateShortLink(ctx context.Context, link ShortLink) error
	GetLongLink(ctx context.Context, code string) (LongLink, error)
}

type ExpiredLinksRepository interface {
	DeleteExpiredLinks(ctx context.Context) (int64, error)
}
