package analytics

import (
	"net/http"
	"strconv"

	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/feature/link"

	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type AnalyticsHandler struct {
	linkRepo         link.LinkRepository
	analyticsService *AnalyticsService
}

func NewAnalyticsHandler(linkRepo link.LinkRepository, service *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		linkRepo:         linkRepo,
		analyticsService: service}
}

// GetStats returns aggregated analytics for a short link.
// Route: GET /api/v1/links/{code}/stats?days=30
// Access: owner only (by JWT subject == accounts.public_id).
func (handler *AnalyticsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	accPublicId, ok := auth.AccountPublicIDFromContext(r.Context())
	if !ok {
		httpx.WriteErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := r.PathValue("code")
	if code == "" {
		httpx.WriteErr(w, http.StatusBadRequest, "missing code")
		return
	}

	linkID, err := handler.linkRepo.GetLinkByCodeAndAccountPublicId(r.Context(), code, accPublicId)
	if err != nil {
		if err == link.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		httpx.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	days := 30
	if v := r.URL.Query().Get("days"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			days = parsed
		}
	}

	stats, err := handler.analyticsService.GetLinkAnalytics(r.Context(), linkID, days)
	if err != nil {
		httpx.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := StatsResponse{
		TotalClicks:  stats.TotalClicks,
		UniqueClicks: stats.UniqueClicks,
		ByDay:        make([]dayCount, 0, len(stats.ByDay)),
	}
	for _, d := range stats.ByDay {
		resp.ByDay = append(resp.ByDay, dayCount{Date: d.Date.UTC().Format("2006-01-02"), Count: d.Count})
	}

	httpx.WriteResponse(w, http.StatusOK, resp)
}
