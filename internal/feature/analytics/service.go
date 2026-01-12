package analytics

import (
	"context"
	"log"
	"time"
)

type AnalyticsService struct {
	repo AnalyticsRepository
	ch   chan ClickEvent
}

func NewAnalyticsService(repo AnalyticsRepository) *AnalyticsService {
	service := &AnalyticsService{
		ch: make(chan ClickEvent, 10_000), // big buffer
	}

	go service.clickEventConsumer(repo)

	return service
}

func (service *AnalyticsService) TrackClick(ev ClickEvent) {
	// non-blocking publish
	select {
	case service.ch <- ev:
	default:
		// drop event if full – never block redirect
	}
}

func (service *AnalyticsService) clickEventConsumer(repo AnalyticsRepository) {
	for ev := range service.ch {
		log.Printf("[analytics] click link_id=%d ip=%s", ev.LinkID, ev.IP)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = repo.InsertClick(ctx, Click{
			LinkID:    ev.LinkID,
			IPAddress: ev.IP,
			UserAgent: ev.UserAgent,
			Referer:   ev.Referer,
			CreatedAt: time.Now().UTC(),
		})
		cancel()
	}
}
