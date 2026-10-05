//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatient_NormalCaseIsMarkedAsNoFindings(t *testing.T) {
	t.Parallel()
	a := validAssessment()
	a.Urgency = "normal"
	c := confirmedCase(t, a)

	status, plan, _ := api.plan(t, c.SourceSystem, c.Patient.ID)

	require.Equal(t, http.StatusOK, status)
	require.Len(t, plan.Cases, 1)
	assert.True(t, plan.Cases[0].NoFindings)
	require.NotNil(t, plan.Cases[0].BookBy)
	assert.Equal(t, plan.Cases[0].ConfirmedAt.AddDate(1, 0, 0), *plan.Cases[0].BookBy)
}

func TestPatient_DecliningEveryStepClosesTheCase(t *testing.T) {
	t.Parallel()
	_, c := confirmedEmergencyCase(t)
	item := awaitTask(t, c.ID)

	for _, req := range []patientDeclineRequest{{Reason: "other"}, {Reason: "too_busy"}} {
		status, _ := api.declineRecommendation(t, c.ID, c.Recommendations[0].ID, req)
		assert.Equal(t, http.StatusBadRequest, status, req.Reason)
	}

	status, first := api.declineRecommendation(t, c.ID, c.Recommendations[0].ID, patientDeclineRequest{Reason: "expensive"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "notified", first.CaseStatus)
	_, found := taskOf(t, c.ID)
	assert.True(t, found, "one step is still open")

	status, _ = api.declineRecommendation(t, c.ID, c.Recommendations[0].ID, patientDeclineRequest{Reason: "far"})
	assert.Equal(t, http.StatusConflict, status, "already declined")

	status, last := api.declineRecommendation(t, c.ID, c.Recommendations[1].ID, patientDeclineRequest{Reason: "other", Comment: "лечусь в другой сети"})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "declined", last.CaseStatus)
	assert.Equal(t, "declined", api.routeStatus(t, c))
	_, found = taskOf(t, c.ID)
	assert.False(t, found)
	_, card := api.operatorCard(t, item.Task.ID)
	assert.Equal(t, "declined", card.Task.Status)
	for _, job := range protocolJobsOf(t, c.ID) {
		assert.NotEqual(t, "pending", job.State, job.Kind)
	}

	_, plan, _ := api.plan(t, c.SourceSystem, c.Patient.ID)
	require.Len(t, plan.Cases, 1)
	for _, rec := range plan.Cases[0].Recommendations {
		assert.True(t, rec.Declined, rec.ServiceName)
	}
	declined := careEventsOfType(t, c.ID, "recommendation_declined")
	require.Len(t, declined, 2)
	assert.Equal(t, "patient", declined[0].Actor)
	assert.Equal(t, "expensive", declined[0].Payload["reason"])
}

func TestPatient_BookedStepCannotBeDeclined(t *testing.T) {
	t.Parallel()
	c := confirmedCase(t, bookableAssessment())
	status, _ := api.registerBooking(t, c.ID, validBooking(c.Recommendations[0].ID, ""))
	require.Equal(t, http.StatusCreated, status)

	status, _ = api.declineRecommendation(t, c.ID, c.Recommendations[0].ID, patientDeclineRequest{Reason: "far"})
	assert.Equal(t, http.StatusConflict, status)

	inReview := newCaseInReview(t, "priority")
	status, _ = api.declineRecommendation(t, inReview.ID, inReview.Recommendations[0].ID, patientDeclineRequest{Reason: "far"})
	assert.Equal(t, http.StatusConflict, status)
}

func TestPatient_HelpRequestCreatesOneCallbackTask(t *testing.T) {
	t.Parallel()
	c := confirmedCase(t, bookableAssessment())

	status, first := api.requestHelp(t, c.ID)
	require.Equal(t, http.StatusCreated, status)
	assert.True(t, first.Created)

	task := awaitTask(t, c.ID)
	assert.Equal(t, "help_request", task.Task.Reason)
	assert.Equal(t, first.TaskID, task.Task.ID)
	assert.Equal(t, task.Task.CreatedAt.Add(2*time.Hour), task.Task.DueAt)

	status, again := api.requestHelp(t, c.ID)
	require.Equal(t, http.StatusOK, status)
	assert.False(t, again.Created)
	assert.Equal(t, first.TaskID, again.TaskID)

	status, _ = api.requestHelp(t, "00000000-0000-0000-0000-000000000000")
	assert.Equal(t, http.StatusNotFound, status)
}
