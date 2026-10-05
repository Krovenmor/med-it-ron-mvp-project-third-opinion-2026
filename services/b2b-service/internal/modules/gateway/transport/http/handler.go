package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

type Handler struct {
	httpx.Responder
	intake Intake
}

func NewHandler(intake Intake, log *zap.Logger) *Handler {
	return &Handler{Responder: httpx.NewResponder(log, "gateway"), intake: intake}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/reports", h.ingestReport)
	mux.HandleFunc("GET /api/v1/reports/{case_id}", h.reportStatus)
}
