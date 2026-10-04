package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	intake Intake
	cases  Cases
	review Review
	clock  DemoClock
	log    *zap.Logger
}

func NewHandler(intake Intake, cases Cases, review Review, clock DemoClock, log *zap.Logger) *Handler {
	return &Handler{intake: intake, cases: cases, review: review, clock: clock, log: log.Named("http")}
}

func (h *Handler) Routes(demoMode bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/reports", h.ingestReport)
	mux.HandleFunc("GET /api/v1/cases/{id}", h.getCase)
	mux.HandleFunc("GET /api/v1/review/queue", h.reviewQueue)
	mux.HandleFunc("POST /api/v1/cases/{id}/open", h.openCase)
	mux.HandleFunc("PATCH /api/v1/cases/{id}/recommendations/{rec_id}", h.reviewRecommendation)
	mux.HandleFunc("POST /api/v1/cases/{id}/recommendations", h.addRecommendation)
	mux.HandleFunc("PUT /api/v1/cases/{id}/urgency", h.changeUrgency)
	mux.HandleFunc("POST /api/v1/cases/{id}/confirm", h.confirmCase)
	if demoMode {
		mux.HandleFunc("POST /demo/clock/advance", h.advanceClock)
	}
	return h.recoverPanics(h.logRequests(mux))
}
