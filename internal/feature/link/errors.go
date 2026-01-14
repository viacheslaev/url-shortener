package link

import "errors"

var (
	ErrNotFound                  = errors.New("link not found")
	ErrLinkExpired               = errors.New("link expired")
	ErrShortcodeAlreadyExists    = errors.New("short code already exists")
	ErrFailedToGenerateShortCode = errors.New("failed to generate short code")
)
