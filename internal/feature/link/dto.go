package link

type createShortLinkRequest struct {
	LongURL string `json:"long_url"`
}

type shortLinkResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

func createShortLinkResponse(baseURL string, link ShortLink) shortLinkResponse {
	return shortLinkResponse{
		ShortCode: link.Code,
		ShortURL:  baseURL + "/" + link.Code,
		LongURL:   link.LongURL,
	}
}
