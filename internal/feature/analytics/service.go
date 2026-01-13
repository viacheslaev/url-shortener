package analytics

import (
	"context"
	"log"
	"time"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

type AnalyticsService struct {
	repo           AnalyticsRepository
	clickEventChan chan link.ClickEvent
}

func NewAnalyticsService(repo AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repo:           repo,
		clickEventChan: make(chan link.ClickEvent, 10_000),
	}
}

// TrackClick publishes click event to background analytics worker.
// The call is non-blocking and will drop the event if the internal buffer is full,
// so redirect path is never slowed down by analytics processing.
func (service *AnalyticsService) TrackClick(ev link.ClickEvent) {
	select {
	case service.clickEventChan <- ev:
	default:
		// drop if buffer is full
	}
}

func (service *AnalyticsService) GetLinkAnalytics(ctx context.Context, linkID int64, days int) (Stats, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().UTC().AddDate(0, 0, -days)
	return service.repo.GetStats(ctx, linkID, since)
}

func (service *AnalyticsService) handleClickEvent(ev link.ClickEvent) {
	log.Printf("[analytics] click link_id=%d ip=%service", ev.LinkID, ev.IP)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := service.repo.InsertClick(ctx, Click{
		LinkID:    ev.LinkID,
		IPAddress: ev.IP,
		UserAgent: ev.UserAgent,
		Referer:   ev.Referer,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		log.Printf("[analytics] failed to save click link_id=%d err=%v", ev.LinkID, err)
	}
}
