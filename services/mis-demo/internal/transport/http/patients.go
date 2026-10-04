package http

import (
	"net/http"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

func (h *Handler) patientHistory(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("id")
	history := demo.HistoryOf(patientID)
	history.Appointments = append(history.Appointments, h.scheduler.AppointmentsOf(patientID)...)
	h.writeJSON(w, http.StatusOK, history)
}
