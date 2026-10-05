package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

type Handler struct {
	httpx.Responder
	clock    Clock
	demo     Demo
	demoMode bool
}

func NewHandler(clock Clock, demo Demo, demoMode bool, log *zap.Logger) *Handler {
	return &Handler{Responder: httpx.NewResponder(log, "system"), clock: clock, demo: demo, demoMode: demoMode}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/clock", h.currentTime)
	if h.demoMode {
		mux.HandleFunc("POST /demo/clock/advance", h.advanceClock)
		mux.HandleFunc("POST /demo/reset", h.resetDemo)
	}
}
