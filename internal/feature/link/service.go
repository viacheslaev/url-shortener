package link

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/viacheslaev/url-shortener/internal/config"
)

const uniqueViolationErrCode = "23505"

type LinkService struct {
	clickTracker ClickTracker
	repo         LinkRepository
	cfg          *config.Config
}

func NewLinkService(
	clickTracker ClickTracker,
	repo LinkRepository,
	cfg *config.Config,
) *LinkService {
	return &LinkService{
		clickTracker: clickTracker,
		repo:         repo,
		cfg:          cfg,
	}
}

func (service *LinkService) createShortLink(ctx context.Context, longURL string, accountId string) (ShortLink, error) {
	longURL = strings.TrimSpace(longURL)
	if !validateURL(longURL) {
		return ShortLink{}, errors.New("invalid url")
	}

	// Generates shortCode with retries in case unique constraint violation
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		code, err := generateShortCode()
		if err != nil {
			return ShortLink{}, fmt.Errorf("failed to generate short code: %w", err)
		}

		var shortLink = ShortLink{
			AccountPublicId: accountId,
			Code:            code,
			LongURL:         longURL,
			ExpiresAt:       service.calculateExpireTime(),
		}
		err = service.repo.Save(ctx, shortLink)
		if err == nil {
			return shortLink, nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == uniqueViolationErrCode {
			log.Printf("ERROR createShortLink: attempt %d failed to generate unique code for longURL=%s", attempt, longURL)
			continue
		}

		return ShortLink{}, err
	}

	log.Printf("ERROR createShortLink: failed to generate unique code after %d attempts, longURL=%s", maxAttempts, longURL)
	return ShortLink{}, errors.New("failed to generate unique short code")
}

func (service *LinkService) resolveShortLink(ctx context.Context, code string, clientContext ClientContext) (string, error) {
	longLink, err := service.repo.GetLongLink(ctx, code)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrNotFound
		}

		log.Printf("ERROR resolveShortLink failed: code=%s err=%v", code, err)

		return "", err
	}

	// Track users click for analytics
	service.clickTracker.TrackClick(ClickEvent{
		LinkID:    longLink.Id,
		IP:        clientContext.IP,
		UserAgent: clientContext.UserAgent,
		Referer:   clientContext.Referer,
	})

	// If ExpiresAt == nil this is permanent link
	if longLink.ExpiresAt != nil {
		expiresAt := longLink.ExpiresAt.UTC()
		now := time.Now().UTC()

		if !expiresAt.After(time.Now()) {
			log.Printf(
				"Expired link: code=%s expires_at=%s now=%s",
				code,
				expiresAt.Format(time.RFC3339),
				now.Format(time.RFC3339),
			)
			return "", ErrLinkExpired
		}
	}

	return longLink.LongURL, nil
}

// calculateExpireTime returns the absolute expiration timestamp in UTC.
func (service *LinkService) calculateExpireTime() time.Time {
	return time.Now().UTC().Add(time.Duration(service.cfg.LinkTTLHours) * time.Hour)
}
