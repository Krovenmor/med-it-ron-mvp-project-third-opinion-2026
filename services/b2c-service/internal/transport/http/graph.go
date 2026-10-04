package http

import "net/http"

func (h *Handler) patientGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.graphs.Build(r.Context(), r.PathValue("patient_id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, newGraphResponse(graph))
}
