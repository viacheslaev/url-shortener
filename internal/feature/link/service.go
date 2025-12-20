package link

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type URLService struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewURLService() *URLService {
	return &URLService{data: make(map[string]string)}
}

func (s *URLService) createShortLink(longURL string) (shortLink, error) {
	longURL = strings.TrimSpace(longURL)
	if !validateURL(longURL) {
		return shortLink{}, errors.New("invalid url")
	}

	code, err := generateShortCode()
	if err != nil {
		return shortLink{}, fmt.Errorf("failed to generate short code: %w", err)
	}

	// todo: temp impl, replace with DB storage
	s.mu.Lock()
	s.data[code] = longURL
	s.mu.Unlock()

	return shortLink{
		Code:    code,
		LongURL: longURL,
	}, nil
}

func (s *URLService) resolveLongLink(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	long, ok := s.data[code]
	return long, ok
}
