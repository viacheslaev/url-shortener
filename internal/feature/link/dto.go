package link

import "time"

type createShortLinkRequest struct {
	LongURL string `json:"long_url"`
}

type shortLinkResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
	ExpiresAt string `json:"expires_at"`
}

func createShortLinkResponse(baseURL string, link ShortLink) shortLinkResponse {
	return shortLinkResponse{
		ShortCode: link.Code,
		ShortURL:  baseURL + "/" + link.Code,
		LongURL:   link.LongURL,
		ExpiresAt: link.ExpiresAt.UTC().Format(time.RFC3339),
	}
}
