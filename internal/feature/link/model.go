package link

import "time"

type ShortLink struct {
	AccountPublicId string
	Code            string
	LongURL         string
	ExpiresAt       time.Time
}

type LongLink struct {
	LongURL   string
	ExpiresAt *time.Time
}
