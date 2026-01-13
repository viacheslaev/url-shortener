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

type ClickEvent struct {
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
}
