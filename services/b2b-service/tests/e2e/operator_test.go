//go:build e2e

package e2e

import (
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	operatorID  = "operator-e2e"
	clinicPhone = "+7 (495) 123-45-67"
)

func TestOperator_EmergencyAssessmentCreatesUrgentTask(t *testing.T) {
	t.Parallel()
	r, c := emergencyCaseInReview(t)

	item := awaitTask(t, c.ID)

	assert.Equal(t, "emergency", item.Task.Reason)
	assert.Equal(t, "new", item.Task.Status)
	assert.Zero(t, item.Task.Attempts)
	assert.Equal(t, item.Task.CreatedAt.Add(time.Hour), item.Task.DueAt)
	assert.Equal(t, operatorPatientView{ID: r.Patient.ID, FullName: r.Patient.FullName, Phone: r.Patient.Phone}, item.Patient)
	assert.Equal(t, "emergency", item.Case.Urgency)
	assert.Equal(t, "in_review", item.Case.Status)
	assert.Empty(t, item.OfferService)

	sent := notificationsOf(t, c.ID)
	assert.ElementsMatch(t, []string{"patient/push/urgent_contact", "patient/sms/urgent_contact"}, routesOf(sent))
	for _, n := range sent {
		assert.Contains(t, n.Text, clinicPhone)
		assert.NotContains(t, n.Text, r.Conclusion)
	}
	assert.Equal(t, []jobView{{Kind: "escalate_review", State: "pending"}, {Kind: "escalate_contact", State: "pending"}}, protocolJobsOf(t, c.ID))

	created := careEventsOfType(t, c.ID, "operator_task_created")
	require.Len(t, created, 1)
	assert.Equal(t, item.Task.ID, created[0].Payload["task_id"])
}

func TestOperator_LoweringUrgencyCancelsEmergencyTask(t *testing.T) {
	t.Parallel()
	_, c := emergencyCaseInReview(t)
	item := awaitTask(t, c.ID)

	require.Equal(t, http.StatusNoContent, api.changeUrgency(t, doctor, c.ID, urgencyRequest{Urgency: "priority", Reason: "пневмоторакса нет"}))

	cancelled := []jobView{{Kind: "escalate_review", State: "cancelled"}, {Kind: "escalate_contact", State: "cancelled"}}
	require.True(t, poll(func() bool {
		_, found := taskOf(t, c.ID)
		return !found && slices.Equal(cancelled, protocolJobsOf(t, c.ID))
	}), "emergency task and escalations were not cancelled: %+v", protocolJobsOf(t, c.ID))
	_, card := api.operatorCard(t, item.Task.ID)
	assert.Equal(t, "cancelled", card.Task.Status)
	closed := careEventsOfType(t, c.ID, "operator_task_closed")
	require.Len(t, closed, 1)
	assert.Equal(t, doctor, closed[0].Actor)

	require.Equal(t, http.StatusNoContent, api.changeUrgency(t, doctor, c.ID, urgencyRequest{Urgency: "emergency"}))
	again := awaitTask(t, c.ID)
	assert.NotEqual(t, item.Task.ID, again.Task.ID)
}

func TestOperator_ThreeMissedCallsHandPatientToDoctor(t *testing.T) {
	t.Parallel()
	r, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)

	status, taken := api.operatorAction(t, operatorID, item.Task.ID, "take", nil)
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "in_progress", taken.Status)
	assert.Equal(t, operatorID, taken.Assignee)
	assert.True(t, slices.ContainsFunc(api.operatorTasks(t, operatorID, true).Tasks, ofCase(c.ID)))
	assert.False(t, slices.ContainsFunc(api.operatorTasks(t, "another-operator", true).Tasks, ofCase(c.ID)))

	var missed operatorTaskView
	for attempt := 1; attempt <= 3; attempt++ {
		now := api.clock(t)
		status, missed = api.operatorAction(t, operatorID, item.Task.ID, "no-answer", nil)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "no_answer", missed.Status)
		assert.Equal(t, attempt, missed.Attempts)
		assert.Equal(t, attempt == 3, missed.NeedsDoctor)
		require.NotNil(t, missed.NextCallAt)
		assert.WithinDuration(t, now.Add(15*time.Minute), *missed.NextCallAt, time.Minute)
	}

	assert.Equal(t, "unreachable", api.routeStatus(t, c))
	assert.Contains(t, routesOf(notificationsOf(t, c.ID)), "patient/sms/unreachable")

	status, handed := api.operatorAction(t, operatorID, item.Task.ID, "hand-to-doctor", outcomeRequest{Comment: "три недозвона"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "handed_to_doctor", handed.Status)
	assert.NotNil(t, handed.ClosedAt)
	_, found := taskOf(t, c.ID)
	assert.False(t, found)

	escalation := findNotification(t, c.ID, "handed_to_doctor")
	require.NotNil(t, escalation)
	assert.Equal(t, "duty_doctor", escalation.Recipient)
	assert.Contains(t, escalation.Text, r.Patient.ID)
	assert.Contains(t, escalation.Text, "три недозвона")

	_, card := api.operatorCard(t, item.Task.ID)
	assert.Equal(t, []string{"no_answer", "no_answer", "no_answer", "handed_to_doctor"}, outcomesOf(card.Attempts))
	status, _ = api.operatorAction(t, operatorID, item.Task.ID, "no-answer", nil)
	assert.Equal(t, http.StatusConflict, status)
}

func TestOperator_BooksSlotThroughMIS(t *testing.T) {
	t.Parallel()
	r, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)
	assert.Equal(t, c.Recommendations[0].ServiceName, item.OfferService)

	status, raw := api.do(t, http.MethodGet, "/api/v1/operator/tasks/"+item.Task.ID, nil)
	require.Equal(t, http.StatusOK, status)
	for _, hidden := range []string{r.Conclusion, c.Recommendations[0].Rationale, c.Recommendations[0].GuidelineRef} {
		assert.NotContains(t, string(raw), hidden)
	}
	card := decode[operatorCardView](t, raw)
	require.Len(t, card.Offers, 2)
	offer := card.Offers[0]
	assert.Equal(t, c.Recommendations[0].ID, offer.RecommendationID)
	assert.Equal(t, c.Recommendations[0].PatientText, offer.PatientText)
	assert.Equal(t, "critical", offer.Mark)
	assert.False(t, offer.Booked)
	require.Len(t, offer.Slots, len(slotStarts))
	slot := offer.Slots[1]
	assert.Equal(t, slotStarts[1], slot.StartsAt)

	status, booked := api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: offer.RecommendationID, SlotID: slot.ID})

	require.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "booked", booked.CaseStatus)
	assert.Equal(t, []misBookingRequest{{
		PatientID:   r.Patient.ID,
		SlotID:      slot.ID,
		ServiceName: offer.ServiceName,
		ReferralID:  offer.RecommendationID,
	}}, mis.bookingsOf(r.Patient.ID))
	_, found := taskOf(t, c.ID)
	assert.False(t, found)

	_, card = api.operatorCard(t, item.Task.ID)
	assert.Equal(t, "booked", card.Task.Status)
	assert.True(t, card.Offers[0].Booked)
	assert.Empty(t, card.Offers[0].Slots)
	assert.Equal(t, []string{"booked"}, outcomesOf(card.Attempts))
	assert.Equal(t, operatorID, card.Attempts[0].Actor)

	bookedEvents := careEventsOfType(t, c.ID, "recommendation_booked")
	require.Len(t, bookedEvents, 1)
	assert.Equal(t, "operator", bookedEvents[0].Payload["channel"])
	assert.Equal(t, slotStarts[1].Format(time.RFC3339), bookedEvents[0].Payload["scheduled_at"])
	for _, job := range protocolJobsOf(t, c.ID) {
		assert.NotEqual(t, "pending", job.State, job.Kind)
	}

	status, _ = api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: card.Offers[1].RecommendationID, SlotID: card.Offers[1].Slots[0].ID})
	assert.Equal(t, http.StatusConflict, status)
}

func TestOperator_TakenSlotKeepsTaskOpen(t *testing.T) {
	t.Parallel()
	_, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)
	_, card := api.operatorCard(t, item.Task.ID)
	offer := card.Offers[0]
	mis.takeSlot(offer.Slots[0].ID)

	status, _ := api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: offer.RecommendationID, SlotID: offer.Slots[0].ID})

	assert.Equal(t, http.StatusConflict, status)
	_, found := taskOf(t, c.ID)
	assert.True(t, found)
	assert.Equal(t, "notified", api.routeStatus(t, c))
}

func TestOperator_DeclineClosesTaskAndStopsProtocol(t *testing.T) {
	t.Parallel()
	_, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)

	for _, req := range []declineRequest{{Reason: "other"}, {Reason: "too_busy"}, {}} {
		status, _ := api.operatorAction(t, operatorID, item.Task.ID, "decline", req)
		assert.Equal(t, http.StatusBadRequest, status, req.Reason)
	}

	status, declined := api.operatorAction(t, operatorID, item.Task.ID, "decline", declineRequest{Reason: "far", Comment: "живёт в другом городе"})

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "declined", declined.Status)
	assert.Equal(t, "declined", api.routeStatus(t, c))
	for _, job := range protocolJobsOf(t, c.ID) {
		assert.NotEqual(t, "pending", job.State, job.Kind)
	}
	_, card := api.operatorCard(t, item.Task.ID)
	require.Len(t, card.Attempts, 1)
	assert.Equal(t, callAttemptView{Outcome: "declined", DeclineReason: "far", Comment: "живёт в другом городе", Actor: operatorID}, card.Attempts[0])
	attempts := careEventsOfType(t, c.ID, "call_attempt")
	require.Len(t, attempts, 1)
	assert.Equal(t, "far", attempts[0].Payload["decline_reason"])
}

func TestOperator_CallbackThenContacted(t *testing.T) {
	t.Parallel()
	_, c := emergencyCaseInReview(t)
	item := awaitTask(t, c.ID)
	now := api.clock(t)

	for _, req := range []callbackRequest{{}, {CallAt: "tomorrow"}, {CallAt: now.Add(-time.Minute).Format(time.RFC3339)}} {
		status, _ := api.operatorAction(t, operatorID, item.Task.ID, "callback", req)
		assert.Equal(t, http.StatusBadRequest, status, req.CallAt)
	}

	callAt := now.Add(3 * time.Hour).Truncate(time.Second)
	status, scheduled := api.operatorAction(t, operatorID, item.Task.ID, "callback", callbackRequest{CallAt: callAt.Format(time.RFC3339), Comment: "просил перезвонить вечером"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "callback", scheduled.Status)
	require.NotNil(t, scheduled.NextCallAt)
	assert.True(t, callAt.Equal(*scheduled.NextCallAt))

	status, contacted := api.operatorAction(t, operatorID, item.Task.ID, "contacted", outcomeRequest{Comment: "едет в клинику"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "contacted", contacted.Status)
	_, found := taskOf(t, c.ID)
	assert.False(t, found)

	_, card := api.operatorCard(t, item.Task.ID)
	assert.Equal(t, []string{"callback", "contacted"}, outcomesOf(card.Attempts))
	require.NotNil(t, card.Attempts[0].CallbackAt)
	assert.True(t, callAt.Equal(*card.Attempts[0].CallbackAt))
}

func TestOperator_SelfBookingClosesTask(t *testing.T) {
	t.Parallel()
	_, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)

	status, _ := api.registerBooking(t, c.ID, validBooking(c.Recommendations[0].ID, ""))
	require.Equal(t, http.StatusCreated, status)

	_, found := taskOf(t, c.ID)
	assert.False(t, found)
	_, card := api.operatorCard(t, item.Task.ID)
	assert.Equal(t, "booked", card.Task.Status)
	closed := careEventsOfType(t, c.ID, "operator_task_closed")
	require.Len(t, closed, 1)
	assert.Equal(t, "patient", closed[0].Actor)
}

func TestOperator_RequestErrors(t *testing.T) {
	t.Parallel()
	_, c := emergencyCaseInReview(t)
	item := awaitTask(t, c.ID)

	status, _ := api.operatorAction(t, "", item.Task.ID, "take", nil)
	assert.Equal(t, http.StatusUnauthorized, status)
	status, _ = api.do(t, http.MethodGet, "/api/v1/operator/tasks?mine=true", nil)
	assert.Equal(t, http.StatusUnauthorized, status)

	for _, id := range []string{uuid.NewString(), "not-a-uuid"} {
		status, _ = api.operatorCard(t, id)
		assert.Equal(t, http.StatusNotFound, status, id)
		status, _ = api.operatorAction(t, operatorID, id, "no-answer", nil)
		assert.Equal(t, http.StatusNotFound, status, id)
	}

	status, _ = api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: "nope", SlotID: "slot"})
	assert.Equal(t, http.StatusBadRequest, status)
	status, _ = api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: c.Recommendations[0].ID})
	assert.Equal(t, http.StatusBadRequest, status)
	status, _ = api.operatorBook(t, operatorID, item.Task.ID, operatorBookingRequest{RecommendationID: c.Recommendations[0].ID, SlotID: "PULM-CONSULT@" + slotStarts[0].Format(time.RFC3339)})
	assert.Equal(t, http.StatusConflict, status, "recommendation is not confirmed yet")
}

func emergencyCaseInReview(t *testing.T) (report, caseView) {
	t.Helper()
	r := newReport()
	a := bookableAssessment()
	a.Urgency = "emergency"
	return r, ingestInReview(t, r, a)
}

func confirmedEmergencyCase(t *testing.T) (report, caseView) {
	t.Helper()
	r, c := emergencyCaseInReview(t)
	reviewAll(t, c)
	confirmCase(t, c)
	return r, c
}

func taskOf(t *testing.T, caseID string) (operatorTaskItemView, bool) {
	t.Helper()
	tasks := api.operatorTasks(t, "", false).Tasks
	i := slices.IndexFunc(tasks, ofCase(caseID))
	if i < 0 {
		return operatorTaskItemView{}, false
	}
	return tasks[i], true
}

func awaitTask(t *testing.T, caseID string) operatorTaskItemView {
	t.Helper()
	var (
		item  operatorTaskItemView
		found bool
	)
	require.True(t, poll(func() bool {
		item, found = taskOf(t, caseID)
		return found
	}), "no active operator task for case %s", caseID)
	return item
}

func ofCase(caseID string) func(operatorTaskItemView) bool {
	return func(item operatorTaskItemView) bool { return item.Task.CaseID == caseID }
}

func notificationsOf(t *testing.T, caseID string) []notificationView {
	t.Helper()
	var found []notificationView
	for _, n := range api.notifications(t) {
		if n.CaseID == caseID {
			found = append(found, n)
		}
	}
	return found
}

func findNotification(t *testing.T, caseID, kind string) *notificationView {
	t.Helper()
	for _, n := range notificationsOf(t, caseID) {
		if n.Kind == kind {
			return &n
		}
	}
	return nil
}

func countNotifications(t *testing.T, caseID, kind string) int {
	t.Helper()
	n := 0
	for _, sent := range notificationsOf(t, caseID) {
		if sent.Kind == kind {
			n++
		}
	}
	return n
}

func routesOf(notifications []notificationView) []string {
	routes := make([]string, 0, len(notifications))
	for _, n := range notifications {
		routes = append(routes, n.Recipient+"/"+n.Channel+"/"+n.Kind)
	}
	return routes
}

func outcomesOf(attempts []callAttemptView) []string {
	outcomes := make([]string, 0, len(attempts))
	for _, a := range attempts {
		outcomes = append(outcomes, a.Outcome)
	}
	return outcomes
}

func protocolJobsOf(t *testing.T, caseID string) []jobView {
	t.Helper()
	var jobs []jobView
	for _, job := range jobsOf(t, caseID) {
		if job.Kind != "assess" {
			jobs = append(jobs, job)
		}
	}
	return jobs
}
