package analytics

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type AnalyticsHandler struct {
	analyticsService *AnalyticsService
}

func NewAnalyticsHandler(service *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
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

	shortCode := r.PathValue("code")
	if shortCode == "" {
		httpx.WriteErr(w, http.StatusBadRequest, "missing shortCode")
		return
	}

	days, err := parseDays(r.URL.Query().Get("days"))
	if err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, "invalid days parameter")
		return
	}

	stats, err := handler.analyticsService.GetLinkAnalytics(r.Context(), accPublicId, shortCode, days)
	if err != nil {
		switch {
		case errors.Is(err, ErrAnalyticsNotFound):
			httpx.WriteErr(w, http.StatusNotFound, "analytics not found")
			return
		default:
			log.Printf("GetStats failed: %v", err)
			httpx.WriteErr(w, http.StatusInternalServerError, "failed to get analytics")
		}
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

func parseDays(daysParam string) (int, error) {
	if daysParam == "" {
		return 0, fmt.Errorf("days parameter is required")
	}

	parsed, err := strconv.Atoi(daysParam)
	if err != nil || parsed <= 0 || parsed > 365 {
		return 0, fmt.Errorf("invalid days parameter")
	}

	return parsed, nil
}
