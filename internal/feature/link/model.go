package link

import "time"

type ShortLink struct {
	Code      string
	LongURL   string
	ExpiresAt time.Time
}
