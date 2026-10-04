package http

import (
	"net/http"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

func (h *Handler) patientHistory(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, demo.HistoryOf(r.PathValue("id")))
}
