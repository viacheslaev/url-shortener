package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
)

type mockAnalyticsRepo struct {
	saveClickFunc func(ctx context.Context, c Click) error
	GetStatsFunc  func(ctx context.Context, linkID int64, since time.Time) (Stats, error)
}

func (m *mockAnalyticsRepo) SaveClick(ctx context.Context, c Click) error {
	if m.saveClickFunc == nil {
		return errors.New("SaveClick not configured")
	}
	return m.saveClickFunc(ctx, c)
}

func (m *mockAnalyticsRepo) GetStats(ctx context.Context, linkID int64, since time.Time) (Stats, error) {
	if m.GetStatsFunc == nil {
		return Stats{}, errors.New("GetStats not configured")
	}
	return m.GetStatsFunc(ctx, linkID, since)
}

type mockLinksRepo struct {
	getLinkIdFunc func(ctx context.Context, code string, acc string) (int64, error)
}

func (m *mockLinksRepo) GetLinkByCodeAndAccountPublicId(ctx context.Context, code string, accountPublicId string) (int64, error) {
	if m.getLinkIdFunc == nil {
		return 0, errors.New("GetLinkByCodeAndAccountPublicId not configured")
	}
	return m.getLinkIdFunc(ctx, code, accountPublicId)
}

func TestAnalyticsService_TrackClick_EnqueuesWhenNotFull(t *testing.T) {
	svc := NewAnalyticsService(&mockAnalyticsRepo{}, &mockLinksRepo{})

	svc.TrackClick(link.ClickEvent{LinkID: 1, IP: "1.2.3.4"})

	if got := len(svc.clickEventChan); got != 1 {
		t.Fatalf("expected len=1, got len=%d", got)
	}
}

func TestAnalyticsService_GetLinkAnalytics_OK(t *testing.T) {
	linksRepo := &mockLinksRepo{getLinkIdFunc: func(ctx context.Context, code, accPublicId string) (int64, error) {
		if code != "abc" {
			t.Fatalf("unexpected code: %q", code)
		}
		if accPublicId != "488e1984-99f7-4369-b6b1-facd467870cc" {
			t.Fatalf("unexpected accPublicId: %q", accPublicId)
		}
		return 777, nil
	}}

	analyticsRepo := &mockAnalyticsRepo{GetStatsFunc: func(ctx context.Context, linkID int64, since time.Time) (Stats, error) {
		if linkID != 777 {
			t.Fatalf("unexpected linkID: %d", linkID)
		}
		return Stats{TotalClicks: 10, UniqueClicks: 3}, nil
	}}

	svc := NewAnalyticsService(analyticsRepo, linksRepo)

	stats, err := svc.GetLinkAnalytics(context.Background(), "488e1984-99f7-4369-b6b1-facd467870cc", "abc", 7)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalClicks != 10 || stats.UniqueClicks != 3 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

}

func TestAnalyticsService_handleClickEvent_SavesClick(t *testing.T) {
	var savedClick Click
	analyticsRepo := &mockAnalyticsRepo{saveClickFunc: func(ctx context.Context, c Click) error {
		savedClick = c
		return nil
	}}
	linksRepo := &mockLinksRepo{getLinkIdFunc: func(ctx context.Context, code, acc string) (int64, error) { return 0, nil }}

	svc := NewAnalyticsService(analyticsRepo, linksRepo)

	svc.handleClickEvent(link.ClickEvent{LinkID: 123, IP: "1.2.3.4", UserAgent: "ua", Referer: "ref"})

	if savedClick.LinkID != 123 {
		t.Fatalf("expected LinkID=123, savedClick %d", savedClick.LinkID)
	}
	if savedClick.IPAddress != "1.2.3.4" {
		t.Fatalf("expected IPAddress=1.2.3.4, savedClick %q", savedClick.IPAddress)
	}
	if savedClick.UserAgent != "ua" {
		t.Fatalf("expected UserAgent=ua, savedClick %q", savedClick.UserAgent)
	}
	if savedClick.Referer != "ref" {
		t.Fatalf("expected Referer=ref, savedClick %q", savedClick.Referer)
	}
}
