package link

import "context"

type LinkRepository interface {
	CreateShortLink(ctx context.Context, link ShortLink) error
	GetLongLink(ctx context.Context, code string) (LongLink, error)
	GetLinkByCodeAndAccountPublicId(ctx context.Context, code string, accountPublicId string) (int64, error)
	DeleteExpiredLinks(ctx context.Context) (int64, error)
}
