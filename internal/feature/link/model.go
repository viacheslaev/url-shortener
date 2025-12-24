package link

import "time"

type ShortLink struct {
	Code      string
	LongURL   string
	ExpiresAt time.Time
}

type LongLink struct {
	LongURL   string
	ExpiresAt *time.Time
}
