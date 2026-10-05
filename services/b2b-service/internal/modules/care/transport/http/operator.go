package http

import "net/http"

func (h *Handler) operatorTasks(w http.ResponseWriter, r *http.Request) {
	assignee := ""
	if r.URL.Query().Get("mine") == "true" {
		actor, ok := h.Actor(w, r)
		if !ok {
			return
		}
		assignee = actor
	}
	items, err := h.operator.List(r.Context(), assignee)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newOperatorTasksResponse(items))
}

func (h *Handler) operatorTaskCard(w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	card, err := h.operator.Card(r.Context(), taskID)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newOperatorCardResponse(card))
}

func (h *Handler) takeTask(w http.ResponseWriter, r *http.Request) {
	actor, taskID, ok := h.taskAction(w, r)
	if !ok {
		return
	}
	task, err := h.operator.Take(r.Context(), taskID, actor)
	h.respondTask(w, r, task, err)
}

func (h *Handler) recordNoAnswer(w http.ResponseWriter, r *http.Request) {
	cmd, ok := h.outcome(w, r)
	if !ok {
		return
	}
	task, err := h.operator.NoAnswer(r.Context(), cmd)
	h.respondTask(w, r, task, err)
}

func (h *Handler) recordContacted(w http.ResponseWriter, r *http.Request) {
	cmd, ok := h.outcome(w, r)
	if !ok {
		return
	}
	task, err := h.operator.Contacted(r.Context(), cmd)
	h.respondTask(w, r, task, err)
}

func (h *Handler) handToDoctor(w http.ResponseWriter, r *http.Request) {
	cmd, ok := h.outcome(w, r)
	if !ok {
		return
	}
	task, err := h.operator.HandToDoctor(r.Context(), cmd)
	h.respondTask(w, r, task, err)
}

func (h *Handler) scheduleCallback(w http.ResponseWriter, r *http.Request) {
	actor, taskID, ok := h.taskAction(w, r)
	if !ok {
		return
	}
	var req callbackRequest
	if !h.Decode(w, r, &req) {
		return
	}
	cmd, err := req.toCommand(taskID, actor)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	task, err := h.operator.Callback(r.Context(), cmd)
	h.respondTask(w, r, task, err)
}

func (h *Handler) recordDecline(w http.ResponseWriter, r *http.Request) {
	actor, taskID, ok := h.taskAction(w, r)
	if !ok {
		return
	}
	var req declineRequest
	if !h.Decode(w, r, &req) {
		return
	}
	task, err := h.operator.Decline(r.Context(), req.toCommand(taskID, actor))
	h.respondTask(w, r, task, err)
}

func (h *Handler) bookByOperator(w http.ResponseWriter, r *http.Request) {
	actor, taskID, ok := h.taskAction(w, r)
	if !ok {
		return
	}
	var req operatorBookingRequest
	if !h.Decode(w, r, &req) {
		return
	}
	cmd, err := req.toCommand(taskID, actor)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}

	result, err := h.operator.Book(r.Context(), cmd)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.WriteJSON(w, status, newBookingResponse(result))
}
