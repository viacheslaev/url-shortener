package link

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type LinkHandler struct {
	config  *config.Config
	service *LinkService
}

func NewLinkHandler(cfg *config.Config, svc *LinkService) *LinkHandler {
	return &LinkHandler{
		config:  cfg,
		service: svc}
}

func (handler *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	accountPublicId, ok := auth.AccountPublicIDFromContext(r.Context())
	if !ok {
		httpx.WriteErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	link, err := handler.service.createShortLink(r.Context(), req.LongURL, accountPublicId)
	if err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := createShortLinkResponse(handler.config.BaseURL, link)
	httpx.WriteResponse(w, http.StatusCreated, resp)
}

func (handler *LinkHandler) ResolveShortLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	longLink, err := handler.service.resolveShortLink(r.Context(), code, ClientContext{
		IP:        httpx.ClientIP(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
	})

	if err == nil {
		http.Redirect(w, r, longLink, http.StatusFound)
		return
	}

	switch {
	case errors.Is(err, ErrLinkExpired):
		httpx.WriteErr(w, http.StatusGone, ErrLinkExpired.Error())
		return

	case errors.Is(err, ErrNotFound):
		http.NotFound(w, r)
		return

	default:
		log.Printf("Resolve failed: code=%s err=%v", code, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
