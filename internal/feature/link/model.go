package link

import "time"

type ShortLink struct {
	AccountPublicId string
	Code            string
	LongURL         string
	ExpiresAt       time.Time
}

type LongLink struct {
	Id        int64
	LongURL   string
	ExpiresAt *time.Time
}

type ClientContext struct {
	IP        string
	UserAgent string
	Referer   string
}
