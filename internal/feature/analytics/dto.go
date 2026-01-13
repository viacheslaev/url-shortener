package analytics

import "time"

type dayCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type StatsResponse struct {
	TotalClicks  int64      `json:"total_clicks"`
	UniqueClicks int64      `json:"unique_clicks"`
	ByDay        []dayCount `json:"by_day"`
}

type Stats struct {
	TotalClicks  int64
	UniqueClicks int64
	ByDay        []DayCount
}

type DayCount struct {
	Date  time.Time
	Count int64
}
