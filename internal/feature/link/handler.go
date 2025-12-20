package link

import (
	"encoding/json"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type URLHandler struct {
	svc *URLService
}

func NewURLHandler(svc *URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req createShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	link, err := h.svc.createShortLink(req.LongURL)
	if err != nil {
		httpx.WriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := createShortLinkResponse(link)
	httpx.WriteResponse(w, http.StatusCreated, resp)
}

func (h *URLHandler) ResolveShortLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	longLink, ok := h.svc.resolveLongLink(code)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longLink, http.StatusFound)
}
