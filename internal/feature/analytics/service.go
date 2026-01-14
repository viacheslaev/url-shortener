package analytics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

type AnalyticsService struct {
	analyticsRepo  AnalyticsRepository
	linksRepo      LinkRepository
	clickEventChan chan link.ClickEvent
}

func NewAnalyticsService(analyticRepo AnalyticsRepository, linkRepo LinkRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:  analyticRepo,
		linksRepo:      linkRepo,
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

func (service *AnalyticsService) GetLinkAnalytics(ctx context.Context, accPublicId string, shortCode string, days int) (Stats, error) {
	linkID, err := service.linksRepo.GetLinkByCodeAndAccountPublicId(ctx, shortCode, accPublicId)
	if err != nil {
		if errors.Is(err, link.ErrNotFound) {
			return Stats{}, ErrAnalyticsNotFound
		}

		return Stats{}, fmt.Errorf("get analytics failed: %w", err)
	}

	since := time.Now().UTC().AddDate(0, 0, -days)

	return service.analyticsRepo.GetStats(ctx, linkID, since)
}

func (service *AnalyticsService) handleClickEvent(ev link.ClickEvent) {
	log.Printf("[analytics] click link_id=%d ip=%s", ev.LinkID, ev.IP)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := service.analyticsRepo.SaveClick(ctx, Click{
		LinkID:    ev.LinkID,
		IPAddress: ev.IP,
		UserAgent: ev.UserAgent,
		Referer:   ev.Referer,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		log.Printf("[analytics] failed to save click link_id=%d err=%v", ev.LinkID, err)
	}
}
