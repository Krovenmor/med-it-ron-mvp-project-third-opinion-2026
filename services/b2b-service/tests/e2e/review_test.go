//go:build e2e

package e2e

import (
	"net/http"
	"slices"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const doctor = "doctor-e2e"

func TestReview_QueueOrdersByUrgencyThenArrival(t *testing.T) {
	t.Parallel()
	normal := newCaseInReview(t, "normal")
	firstPriority := newCaseInReview(t, "priority")
	emergency := newCaseInReview(t, "emergency")
	secondPriority := newCaseInReview(t, "priority")

	queue := api.reviewQueue(t)

	ours := map[string]bool{normal.ID: true, firstPriority.ID: true, emergency.ID: true, secondPriority.ID: true}
	var order []string
	for _, item := range queue.Cases {
		if ours[item.CaseID] {
			order = append(order, item.CaseID)
		}
	}
	assert.Equal(t, []string{emergency.ID, firstPriority.ID, secondPriority.ID, normal.ID}, order)

	item := queue.Cases[slices.IndexFunc(queue.Cases, func(c queueCaseView) bool { return c.CaseID == emergency.ID })]
	assert.Equal(t, "emergency", item.Urgency)
	assert.Equal(t, emergency.Patient, item.Patient)
	assert.Equal(t, 2, item.RecommendationsTotal)
	assert.Zero(t, item.RecommendationsReviewed)
}

func TestReview_DoctorConfirmsCaseAfterReviewingAllRecommendations(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")

	require.Equal(t, http.StatusNoContent, api.openCase(t, doctor, c.ID))

	editedText := "Пожалуйста, запишитесь к пульмонологу в ближайшие дни"
	status, first := api.reviewRecommendation(t, doctor, c.ID, c.Recommendations[0].ID,
		reviewRequest{Mark: "critical", PatientText: &editedText})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, editedText, first.PatientText)
	require.NotNil(t, first.Review)
	assert.Equal(t, "critical", first.Review.Mark)
	assert.Equal(t, doctor, first.Review.ReviewedBy)

	status, second := api.reviewRecommendation(t, doctor, c.ID, c.Recommendations[1].ID,
		reviewRequest{Mark: "rejected", RejectReason: "other", RejectComment: "уже наблюдается у хирурга"})
	require.Equal(t, http.StatusOK, status)
	require.NotNil(t, second.Review)
	assert.Equal(t, reviewView{
		Mark:          "rejected",
		RejectReason:  "other",
		RejectComment: "уже наблюдается у хирурга",
		ReviewedBy:    doctor,
		ReviewedAt:    second.Review.ReviewedAt,
	}, *second.Review)

	status, own := api.addRecommendation(t, doctor, c.ID, addRecommendationRequest{
		ServiceName: "Консультация онколога",
		PatientText: "Рекомендуем консультацию онколога для обсуждения результатов",
		Mark:        "minor",
	})
	require.Equal(t, http.StatusCreated, status)
	assert.Equal(t, 3, own.Position)
	assert.Equal(t, "doctor", own.Source)
	assert.Equal(t, "low", own.Importance)
	require.NotNil(t, own.Review)
	assert.Equal(t, "minor", own.Review.Mark)

	reason := "очаг стабилен по сравнению с прошлым КТ"
	require.Equal(t, http.StatusNoContent, api.changeUrgency(t, doctor, c.ID, urgencyRequest{Urgency: "planned", Reason: reason}))

	status, ref := api.confirm(t, doctor, c.ID)
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "confirmed", ref.Status)

	_, confirmed := api.getCase(t, c.ID)
	assert.Equal(t, "confirmed", confirmed.Status)
	assert.Equal(t, "planned", confirmed.Urgency)
	require.Len(t, confirmed.Recommendations, 3)
	for _, rec := range confirmed.Recommendations {
		assert.NotNil(t, rec.Review, "recommendation %d", rec.Position)
	}
	assert.NotContains(t, queueIDs(api.reviewQueue(t)), c.ID)

	events := eventsOf(t, c.ID)
	types := make([]string, 0, len(events))
	for _, e := range events {
		types = append(types, e.Type)
	}
	assert.Equal(t, []string{
		"status_changed", "status_changed",
		"case_opened",
		"recommendation_reviewed", "recommendation_reviewed",
		"recommendation_added",
		"urgency_changed",
		"status_changed",
	}, types)
	for _, e := range events[2:] {
		assert.Equal(t, doctor, e.Actor, e.Type)
	}
	assert.Equal(t, map[string]string{"from": "priority", "to": "planned", "reason": reason}, events[6].Payload)
	assert.Equal(t, "true", events[3].Payload["patient_text_edited"])
}

func TestReview_ConfirmRequiresAllRecommendationsReviewed(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")
	status, _ := api.reviewRecommendation(t, doctor, c.ID, c.Recommendations[0].ID, reviewRequest{Mark: "critical"})
	require.Equal(t, http.StatusOK, status)

	status, _ = api.confirm(t, doctor, c.ID)

	assert.Equal(t, http.StatusConflict, status)
	_, current := api.getCase(t, c.ID)
	assert.Equal(t, "in_review", current.Status)
}

func TestReview_UrgencyChangeRules(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")

	steps := []struct {
		name   string
		req    urgencyRequest
		status int
	}{
		{"lowering without reason", urgencyRequest{Urgency: "planned"}, http.StatusBadRequest},
		{"same urgency", urgencyRequest{Urgency: "priority"}, http.StatusBadRequest},
		{"unknown urgency", urgencyRequest{Urgency: "asap"}, http.StatusBadRequest},
		{"raising without reason", urgencyRequest{Urgency: "emergency"}, http.StatusNoContent},
		{"lowering with reason", urgencyRequest{Urgency: "normal", Reason: "ложноположительная находка"}, http.StatusNoContent},
	}
	for _, step := range steps {
		assert.Equal(t, step.status, api.changeUrgency(t, doctor, c.ID, step.req), step.name)
	}

	_, current := api.getCase(t, c.ID)
	assert.Equal(t, "normal", current.Urgency)
}

func TestReview_InvalidReviewIsRejected(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")
	recID := c.Recommendations[0].ID
	empty := ""

	reviews := []struct {
		name string
		req  reviewRequest
	}{
		{"unknown mark", reviewRequest{Mark: "maybe"}},
		{"rejected without reason", reviewRequest{Mark: "rejected"}},
		{"rejected as other without comment", reviewRequest{Mark: "rejected", RejectReason: "other"}},
		{"reason on accepted recommendation", reviewRequest{Mark: "critical", RejectReason: "contraindicated"}},
		{"empty patient text", reviewRequest{Mark: "critical", PatientText: &empty}},
		{"nothing to update", reviewRequest{}},
	}
	for _, tt := range reviews {
		status, _ := api.reviewRecommendation(t, doctor, c.ID, recID, tt.req)
		assert.Equal(t, http.StatusBadRequest, status, tt.name)
	}

	additions := []struct {
		name string
		req  addRecommendationRequest
	}{
		{"missing service name", addRecommendationRequest{PatientText: "текст", Mark: "minor"}},
		{"missing patient text", addRecommendationRequest{ServiceName: "Консультация", Mark: "minor"}},
		{"rejected own recommendation", addRecommendationRequest{ServiceName: "Консультация", PatientText: "текст", Mark: "rejected"}},
	}
	for _, tt := range additions {
		status, _ := api.addRecommendation(t, doctor, c.ID, tt.req)
		assert.Equal(t, http.StatusBadRequest, status, tt.name)
	}

	_, current := api.getCase(t, c.ID)
	assert.Len(t, current.Recommendations, 2)
	assert.Nil(t, current.Recommendations[0].Review)
}

func TestReview_ActionsRequireDoctorIdentity(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")

	assert.Equal(t, http.StatusUnauthorized, api.openCase(t, "", c.ID))
	status, _ := api.reviewRecommendation(t, "", c.ID, c.Recommendations[0].ID, reviewRequest{Mark: "critical"})
	assert.Equal(t, http.StatusUnauthorized, status)
	status, _ = api.addRecommendation(t, "", c.ID, addRecommendationRequest{ServiceName: "Консультация", PatientText: "текст", Mark: "minor"})
	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, http.StatusUnauthorized, api.changeUrgency(t, "", c.ID, urgencyRequest{Urgency: "emergency"}))
	status, _ = api.confirm(t, "", c.ID)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestReview_ConfirmedCaseIsReadOnly(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")
	reviewAll(t, c)
	status, _ := api.confirm(t, doctor, c.ID)
	require.Equal(t, http.StatusOK, status)

	assert.Equal(t, http.StatusConflict, api.openCase(t, doctor, c.ID))
	status, _ = api.reviewRecommendation(t, doctor, c.ID, c.Recommendations[0].ID, reviewRequest{Mark: "minor"})
	assert.Equal(t, http.StatusConflict, status)
	status, _ = api.addRecommendation(t, doctor, c.ID, addRecommendationRequest{ServiceName: "Консультация", PatientText: "текст", Mark: "minor"})
	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, http.StatusConflict, api.changeUrgency(t, doctor, c.ID, urgencyRequest{Urgency: "emergency"}))
	status, _ = api.confirm(t, doctor, c.ID)
	assert.Equal(t, http.StatusConflict, status)
}

func TestReview_UnknownCaseOrRecommendationIsNotFound(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")
	other := newCaseInReview(t, "priority")

	status, _ := api.reviewRecommendation(t, doctor, c.ID, other.Recommendations[0].ID, reviewRequest{Mark: "critical"})
	assert.Equal(t, http.StatusNotFound, status)
	status, _ = api.confirm(t, doctor, uuid.NewString())
	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, http.StatusNotFound, api.openCase(t, doctor, "not-a-uuid"))
}

func TestReview_ConcurrentConfirmationsConfirmOnce(t *testing.T) {
	t.Parallel()
	const doctors = 10
	c := newCaseInReview(t, "priority")
	reviewAll(t, c)

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errs     []error
		statuses = map[int]int{}
	)
	for range doctors {
		wg.Go(func() {
			status, _, err := api.send(http.MethodPost, "/api/v1/cases/"+c.ID+"/confirm", nil, doctor)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			statuses[status]++
		})
	}
	wg.Wait()

	require.Empty(t, errs)
	assert.Equal(t, map[int]int{http.StatusOK: 1, http.StatusConflict: doctors - 1}, statuses)
	confirmations := 0
	for _, e := range eventsOf(t, c.ID) {
		if e.Type == "status_changed" && e.Actor == doctor {
			confirmations++
		}
	}
	assert.Equal(t, 1, confirmations)
}

func newCaseInReview(t *testing.T, urgency string) caseView {
	t.Helper()
	r := newReport()
	a := validAssessment()
	a.Urgency = urgency
	ai.on(r.Conclusion, respond(a))
	_, ref := api.ingest(t, r)
	return api.awaitStatus(t, ref.CaseID, "in_review")
}

func reviewAll(t *testing.T, c caseView) {
	t.Helper()
	for _, rec := range c.Recommendations {
		status, _ := api.reviewRecommendation(t, doctor, c.ID, rec.ID, reviewRequest{Mark: "critical"})
		require.Equal(t, http.StatusOK, status)
	}
}

func queueIDs(q reviewQueueView) []string {
	ids := make([]string, 0, len(q.Cases))
	for _, c := range q.Cases {
		ids = append(ids, c.CaseID)
	}
	return ids
}

func TestReview_QueueShowsPatientIdentifierAndWhenCaseWasOpened(t *testing.T) {
	t.Parallel()
	r := newReport()
	c := ingestInReview(t, r, validAssessment())

	item := queueItem(t, c.ID)
	assert.Equal(t, r.Patient.ID, item.Patient.ID)
	assert.True(t, item.OpenedAt.IsZero())

	require.Equal(t, http.StatusNoContent, api.openCase(t, doctor, c.ID))
	assert.False(t, queueItem(t, c.ID).OpenedAt.IsZero())
}

func queueItem(t *testing.T, caseID string) queueCaseView {
	t.Helper()
	queue := api.reviewQueue(t)
	i := slices.IndexFunc(queue.Cases, func(c queueCaseView) bool { return c.CaseID == caseID })
	require.GreaterOrEqual(t, i, 0, "case %s is not in the queue", caseID)
	return queue.Cases[i]
}

func TestReview_PatientTextCanBeEditedBeforeAndAfterMarking(t *testing.T) {
	t.Parallel()
	c := newCaseInReview(t, "priority")
	rec := c.Recommendations[0]

	draft := "Рекомендуем записаться к пульмонологу в ближайшие две недели"
	status, edited := api.reviewRecommendation(t, doctor, c.ID, rec.ID, reviewRequest{PatientText: &draft})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, draft, edited.PatientText)
	assert.Nil(t, edited.Review)

	status, marked := api.reviewRecommendation(t, doctor, c.ID, rec.ID, reviewRequest{Mark: "minor"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, draft, marked.PatientText)
	require.NotNil(t, marked.Review)
	assert.Equal(t, "minor", marked.Review.Mark)

	final := "Запишитесь к пульмонологу в течение месяца"
	status, reworded := api.reviewRecommendation(t, doctor, c.ID, rec.ID, reviewRequest{PatientText: &final})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, final, reworded.PatientText)
	require.NotNil(t, reworded.Review)
	assert.Equal(t, "minor", reworded.Review.Mark)

	edits := eventsOfType(t, c.ID, "recommendation_text_edited")
	require.Len(t, edits, 2)
	assert.Equal(t, doctor, edits[0].Actor)
	assert.Equal(t, map[string]string{"recommendation_id": rec.ID}, edits[0].Payload)
}
