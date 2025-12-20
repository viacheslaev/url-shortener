package link

const baseURL = "http://localhost:8080/"

type createShortLinkRequest struct {
	LongURL string `json:"long_url"`
}

type shortLinkResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

func createShortLinkResponse(link shortLink) shortLinkResponse {
	return shortLinkResponse{
		ShortCode: link.Code,
		ShortURL:  baseURL + link.Code,
		LongURL:   link.LongURL,
	}
}
