package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

func (h *Handler) sendReports(w http.ResponseWriter, r *http.Request) {
	reports := demo.Reports()
	deliveries := make([]b2b.Delivery, 0, len(reports))
	for _, report := range reports {
		d, err := h.sender.SendReport(r.Context(), report)
		if err != nil {
			h.log.Error("send report to b2b", zap.Error(err))
			h.writeJSON(w, http.StatusBadGateway, deliveriesResponse{Delivered: deliveries, Error: err.Error()})
			return
		}
		deliveries = append(deliveries, d)
	}
	h.writeJSON(w, http.StatusOK, deliveriesResponse{Delivered: deliveries})
}
