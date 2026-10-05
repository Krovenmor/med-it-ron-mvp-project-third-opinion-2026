package http

import (
	"net/http"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

func (h *Handler) patientPlan(w http.ResponseWriter, r *http.Request) {
	sourceSystem := r.URL.Query().Get("source_system")
	if sourceSystem == "" {
		h.WriteJSON(w, http.StatusBadRequest, httpx.ErrorResponse{Error: "source_system query parameter is required"})
		return
	}
	plan, err := h.plan.Get(r.Context(), sourceSystem, r.PathValue("external_id"))
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newPatientPlanResponse(plan))
}
