package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

type Handler struct {
	httpx.Responder
	cases  Cases
	review Review
}

func NewHandler(cases Cases, review Review, log *zap.Logger) *Handler {
	return &Handler{Responder: httpx.NewResponder(log, "doctor"), cases: cases, review: review}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/review/queue", h.reviewQueue)
	mux.HandleFunc("GET /api/v1/cases/{id}", h.getCase)
	mux.HandleFunc("GET /api/v1/cases/{id}/history", h.caseHistory)
	mux.HandleFunc("POST /api/v1/cases/{id}/open", h.openCase)
	mux.HandleFunc("PATCH /api/v1/cases/{id}/recommendations/{rec_id}", h.reviewRecommendation)
	mux.HandleFunc("POST /api/v1/cases/{id}/recommendations", h.addRecommendation)
	mux.HandleFunc("PUT /api/v1/cases/{id}/urgency", h.changeUrgency)
	mux.HandleFunc("POST /api/v1/cases/{id}/confirm", h.confirmCase)
}
