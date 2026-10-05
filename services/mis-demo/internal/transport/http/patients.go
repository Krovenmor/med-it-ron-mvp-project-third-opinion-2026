package http

import (
	"net/http"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

func (h *Handler) patientHistory(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("id")
	history := h.scenario.History(patientID)
	visits, upcoming := demo.Attend(h.scheduler.AppointmentsOf(patientID), h.clock.Now())
	history.Visits = append(history.Visits, visits...)
	history.Appointments = append(history.Appointments, upcoming...)
	h.writeJSON(w, http.StatusOK, history)
}
