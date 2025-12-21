package link

import (
	"encoding/json"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type URLHandler struct {
	config  *config.Config
	service *URLService
}

func NewURLHandler(cfg *config.Config, svc *URLService) *URLHandler {
	return &URLHandler{
		config:  cfg,
		service: svc}
}

func (handler *URLHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req createShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	link, err := handler.service.createShortLink(req.LongURL)
	if err != nil {
		httpx.WriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := createShortLinkResponse(handler.config.BaseURL, link)
	httpx.WriteResponse(w, http.StatusCreated, resp)
}

func (handler *URLHandler) ResolveShortLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	longLink, ok := handler.service.resolveLongLink(code)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longLink, http.StatusFound)
}
