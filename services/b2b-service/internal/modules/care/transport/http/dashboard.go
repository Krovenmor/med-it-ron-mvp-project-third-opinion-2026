package http

import "net/http"

func (h *Handler) getDashboard(w http.ResponseWriter, r *http.Request) {
	days, err := dashboardDays(r.URL.Query().Get("days"))
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	d, err := h.dashboard.Get(r.Context(), days)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newDashboardResponse(d))
}
