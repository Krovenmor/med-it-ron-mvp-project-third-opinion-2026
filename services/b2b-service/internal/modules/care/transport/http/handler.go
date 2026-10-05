package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

type Handler struct {
	httpx.Responder
	plan      Plan
	bookings  Bookings
	patient   Patient
	operator  Operator
	dashboard Dashboard
}

func NewHandler(plan Plan, bookings Bookings, patient Patient, operator Operator, dashboard Dashboard, log *zap.Logger) *Handler {
	return &Handler{
		Responder: httpx.NewResponder(log, "care"),
		plan:      plan,
		bookings:  bookings,
		patient:   patient,
		operator:  operator,
		dashboard: dashboard,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/patients/{external_id}/plan", h.patientPlan)
	mux.HandleFunc("POST /api/v1/cases/{id}/bookings", h.registerBooking)
	mux.HandleFunc("POST /api/v1/cases/{id}/recommendations/{rec_id}/decline", h.declineRecommendation)
	mux.HandleFunc("POST /api/v1/cases/{id}/help-requests", h.requestHelp)
	mux.HandleFunc("GET /api/v1/operator/tasks", h.operatorTasks)
	mux.HandleFunc("GET /api/v1/operator/tasks/{id}", h.operatorTaskCard)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/take", h.takeTask)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/no-answer", h.recordNoAnswer)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/callback", h.scheduleCallback)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/contacted", h.recordContacted)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/decline", h.recordDecline)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/hand-to-doctor", h.handToDoctor)
	mux.HandleFunc("POST /api/v1/operator/tasks/{id}/bookings", h.bookByOperator)
	mux.HandleFunc("GET /api/v1/notifications", h.notificationJournal)
	mux.HandleFunc("GET /api/v1/dashboard", h.getDashboard)
}
