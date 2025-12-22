package link

import "context"

var ErrNotFound = errorString("link not found")

type Repository interface {
	Save(ctx context.Context, code string, longURL string) error
	GetLongURL(ctx context.Context, code string) (string, error)
}

type errorString string

func (e errorString) Error() string { return string(e) }
