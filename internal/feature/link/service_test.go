package link

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/viacheslaev/url-shortener/internal/config"
)

type mockLinkRepo struct {
	createShortLinkFunc func(ctx context.Context, l ShortLink) error
	getLongLinkFunc     func(ctx context.Context, code string) (LongLink, error)
}

func (m *mockLinkRepo) CreateShortLink(ctx context.Context, l ShortLink) error {
	if m.createShortLinkFunc == nil {
		return errors.New("CreateShortLink not configured")
	}
	return m.createShortLinkFunc(ctx, l)
}

func (m *mockLinkRepo) GetLongLink(ctx context.Context, code string) (LongLink, error) {
	if m.getLongLinkFunc == nil {
		return LongLink{}, errors.New("GetLongLink not configured")
	}
	return m.getLongLinkFunc(ctx, code)
}

type mockClickTracker struct {
	trackFn func(ev ClickEvent)
}

func (m *mockClickTracker) TrackClick(ev ClickEvent) {
	if m.trackFn != nil {
		m.trackFn(ev)
	}
}

func testCfg() *config.Config {
	return &config.Config{LinkTTLHours: 24}
}

func TestLinkService_createShortLink_OK(t *testing.T) {
	repo := &mockLinkRepo{
		createShortLinkFunc: func(ctx context.Context, l ShortLink) error {
			if l.AccountPublicId != "acc-1" {
				t.Fatalf("unexpected account id: %q", l.AccountPublicId)
			}
			if l.LongURL != "https://example.com" {
				t.Fatalf("expected trimmed url https://example.com, got %q", l.LongURL)
			}
			if len(l.Code) != 6 {
				t.Fatalf("expected 6-char code, got %q", l.Code)
			}
			return nil
		},
	}
	svc := NewLinkService(&mockClickTracker{}, repo, testCfg())

	got, err := svc.createShortLink(context.Background(), "https://example.com", "acc-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Code == "" {
		t.Fatal("expected non-empty shortcode")
	}
}

func TestLinkService_createShortLink_RetriesOnDuplicate(t *testing.T) {
	attempts := 0

	repo := &mockLinkRepo{
		createShortLinkFunc: func(ctx context.Context, l ShortLink) error {
			attempts++
			if attempts < 3 {
				return ErrShortcodeAlreadyExists
			}
			return nil
		},
	}

	svc := NewLinkService(&mockClickTracker{}, repo, testCfg())

	if _, err := svc.createShortLink(context.Background(), "https://example.com", "acc-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestLinkService_createShortLink_MaxAttemptsExceeded(t *testing.T) {
	repo := &mockLinkRepo{
		createShortLinkFunc: func(ctx context.Context, l ShortLink) error {
			return ErrShortcodeAlreadyExists
		},
	}
	svc := NewLinkService(&mockClickTracker{}, repo, testCfg())

	_, err := svc.createShortLink(context.Background(), "https://example.com", "acc-1")

	if !errors.Is(err, ErrFailedToGenerateShortCode) {
		t.Fatalf("expected ErrFailedToGenerateShortCode, got %v", err)
	}
}

func TestLinkService_resolveShortLink_NotFound(t *testing.T) {
	repo := &mockLinkRepo{
		getLongLinkFunc: func(ctx context.Context, code string) (LongLink, error) {
			return LongLink{}, ErrNotFound
		},
	}
	svc := NewLinkService(&mockClickTracker{}, repo, testCfg())

	_, err := svc.resolveShortLink(context.Background(), "missing", ClientContext{})

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkService_resolveShortLink_Expired(t *testing.T) {
	exp := time.Now().UTC().Add(-1 * time.Hour)

	repo := &mockLinkRepo{
		getLongLinkFunc: func(ctx context.Context, code string) (LongLink, error) {
			return LongLink{Id: 1, LongURL: "https://example.com", ExpiresAt: &exp}, nil
		},
	}
	svc := NewLinkService(&mockClickTracker{}, repo, testCfg())

	_, err := svc.resolveShortLink(context.Background(), "abc", ClientContext{})

	if !errors.Is(err, ErrLinkExpired) {
		t.Fatalf("expected ErrLinkExpired, got %v", err)
	}
}
