package analytics

import "time"

type Click struct {
	LinkID    int64
	IPAddress string
	UserAgent string
	Referer   string
	CreatedAt time.Time
}
