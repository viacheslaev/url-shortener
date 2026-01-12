package analytics

import "time"

type Click struct {
	LinkID    int64
	IPAddress string
	UserAgent string
	Referer   string
	CreatedAt time.Time
}

type ClickEvent struct {
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
}
