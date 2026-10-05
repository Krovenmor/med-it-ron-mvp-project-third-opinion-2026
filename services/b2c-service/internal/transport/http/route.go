package http

import "net/http"

func (h *Handler) patientRoute(w http.ResponseWriter, r *http.Request) {
	route, err := h.routes.Build(r.Context(), r.PathValue("patient_id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, newRouteResponse(route))
}
