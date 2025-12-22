package link

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

const uniqueViolationErrCode = "23505"

type URLService struct {
	repo Repository
}

func NewURLService(repo Repository) *URLService {
	return &URLService{repo: repo}
}

func (service *URLService) createShortLink(ctx context.Context, longURL string) (shortLink, error) {
	longURL = strings.TrimSpace(longURL)
	if !validateURL(longURL) {
		return shortLink{}, errors.New("invalid url")
	}

	// Generates shortCode with retries in case unique constraint violation
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		code, err := generateShortCode()
		if err != nil {
			return shortLink{}, fmt.Errorf("failed to generate short code: %w", err)
		}

		err = service.repo.Save(ctx, code, longURL)
		if err == nil {
			return shortLink{Code: code, LongURL: longURL}, nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == uniqueViolationErrCode {
			log.Printf("ERROR createShortLink: attempt %d failed to generate unique code for longURL=%s", attempt, longURL)
			continue
		}

		return shortLink{}, err
	}

	log.Printf("ERROR createShortLink: failed to generate unique code after %d attempts, longURL=%s", maxAttempts, longURL)
	return shortLink{}, errors.New("failed to generate unique short code")
}

func (service *URLService) resolveLongLink(ctx context.Context, code string) (string, bool) {
	longURL, err := service.repo.GetLongURL(ctx, code)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", false
		}

		log.Printf("ERROR resolveLongLink failed: code=%s err=%v", code, err)

		return "", false
	}

	return longURL, true
}
