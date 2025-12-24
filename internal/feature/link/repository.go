package link

import "context"

type Repository interface {
	Save(ctx context.Context, link ShortLink) error
	GetLongLink(ctx context.Context, code string) (LongLink, error)
}
