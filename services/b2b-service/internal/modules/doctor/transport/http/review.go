package http

import "net/http"

func (h *Handler) reviewQueue(w http.ResponseWriter, r *http.Request) {
	items, err := h.review.Queue(r.Context())
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newReviewQueueResponse(items))
}

func (h *Handler) openCase(w http.ResponseWriter, r *http.Request) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return
	}
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.review.Open(r.Context(), caseID, actor); err != nil {
		h.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) reviewRecommendation(w http.ResponseWriter, r *http.Request) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return
	}
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	recID, ok := h.PathID(w, r, "rec_id")
	if !ok {
		return
	}
	var req reviewRecommendationRequest
	if !h.Decode(w, r, &req) {
		return
	}
	rec, err := h.review.ReviewRecommendation(r.Context(), req.toCommand(caseID, recID, actor))
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newRecommendationResponse(rec))
}

func (h *Handler) addRecommendation(w http.ResponseWriter, r *http.Request) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return
	}
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	var req addRecommendationRequest
	if !h.Decode(w, r, &req) {
		return
	}
	rec, err := h.review.AddRecommendation(r.Context(), req.toCommand(caseID, actor))
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusCreated, newRecommendationResponse(rec))
}

func (h *Handler) changeUrgency(w http.ResponseWriter, r *http.Request) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return
	}
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	var req changeUrgencyRequest
	if !h.Decode(w, r, &req) {
		return
	}
	if err := h.review.ChangeUrgency(r.Context(), req.toCommand(caseID, actor)); err != nil {
		h.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) confirmCase(w http.ResponseWriter, r *http.Request) {
	actor, ok := h.Actor(w, r)
	if !ok {
		return
	}
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	c, err := h.review.Confirm(r.Context(), caseID, actor)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, caseRefResponse{CaseID: c.ID, Status: string(c.Status)})
}
