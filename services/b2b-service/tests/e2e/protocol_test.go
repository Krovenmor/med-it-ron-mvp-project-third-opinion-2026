//go:build e2e

package e2e

import (
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtocol_EmergencyEscalatesToDutyStaff(t *testing.T) {
	r, ignored := emergencyCaseInReview(t)
	_, handled := emergencyCaseInReview(t)
	require.Equal(t, http.StatusNoContent, api.openCase(t, doctor, handled.ID))
	handledTask := awaitTask(t, handled.ID)
	status, _ := api.operatorAction(t, operatorID, handledTask.Task.ID, "callback",
		callbackRequest{CallAt: api.clock(t).Add(3 * time.Hour).Format(time.RFC3339)})
	require.Equal(t, http.StatusOK, status)

	api.advanceClock(t, 31*time.Minute)

	review := awaitNotification(t, ignored.ID, "review_escalation")
	assert.Equal(t, "duty_doctor", review.Recipient)
	assert.Equal(t, "staff", review.Channel)
	assert.Contains(t, review.Text, r.Patient.ID)
	awaitJobSettled(t, handled.ID, "escalate_review")
	assert.Nil(t, findNotification(t, handled.ID, "review_escalation"))

	api.advanceClock(t, 90*time.Minute)

	contact := awaitNotification(t, ignored.ID, "contact_escalation")
	assert.Equal(t, "admin_on_duty", contact.Recipient)
	awaitJobSettled(t, handled.ID, "escalate_contact")
	assert.Nil(t, findNotification(t, handled.ID, "contact_escalation"))
}

func TestProtocol_PatientsWhoDoNotBookGetRemindersAndCalls(t *testing.T) {
	waiting := confirmedCase(t, bookableAssessment())
	selfBooked := confirmedCase(t, bookableAssessment())
	status, _ := api.registerBooking(t, selfBooked.ID, validBooking(selfBooked.Recommendations[0].ID, ""))
	require.Equal(t, http.StatusCreated, status)
	plannedAssessment := bookableAssessment()
	plannedAssessment.Urgency = "planned"
	planned := confirmedCase(t, plannedAssessment)

	assert.ElementsMatch(t, []string{"patient/push/plan_ready", "patient/sms/plan_ready"}, routesOf(notificationsOf(t, waiting.ID)))
	assert.Equal(t, []string{"patient/push/plan_ready"}, routesOf(notificationsOf(t, planned.ID)))
	for _, n := range notificationsOf(t, waiting.ID) {
		if n.Channel == "sms" {
			assert.Contains(t, n.Text, clinicPhone)
		}
	}

	api.advanceClock(t, 25*time.Hour)

	task := awaitTask(t, waiting.ID)
	assert.Equal(t, "no_booking", task.Task.Reason)
	assert.Equal(t, task.Task.CreatedAt.Add(4*time.Hour), task.Task.DueAt)
	assert.Equal(t, waiting.Recommendations[0].ServiceName, task.OfferService)
	assert.Equal(t, "notified", task.Case.Status)
	assert.Equal(t, []jobView{{Kind: "notify_patient", State: "done"}, {Kind: "follow_up", State: "cancelled"}}, protocolJobsOf(t, selfBooked.ID))
	_, found := taskOf(t, selfBooked.ID)
	assert.False(t, found)
	_, found = taskOf(t, planned.ID)
	assert.False(t, found)

	api.advanceClock(t, 48*time.Hour)

	awaitNotificationCount(t, planned.ID, "reminder", 1)

	api.advanceClock(t, 96*time.Hour)

	awaitNotificationCount(t, planned.ID, "reminder", 2)
	plannedTask := awaitTask(t, planned.ID)
	assert.Equal(t, "no_booking", plannedTask.Task.Reason)
	assert.GreaterOrEqual(t, plannedTask.Task.DueAt.Sub(plannedTask.Task.CreatedAt), 48*time.Hour)
	assert.LessOrEqual(t, plannedTask.Task.DueAt.Sub(plannedTask.Task.CreatedAt), 96*time.Hour)

	tasks := api.operatorTasks(t, "", false).Tasks
	waitingAt := slices.IndexFunc(tasks, ofCase(waiting.ID))
	plannedAt := slices.IndexFunc(tasks, ofCase(planned.ID))
	assert.Less(t, waitingAt, plannedAt, "overdue task goes first")
}

func TestProtocol_TimersCountFromConfirmationNotFromDelivery(t *testing.T) {
	c := newCaseInReview(t, "priority")
	reviewAll(t, c)
	status, _ := api.confirm(t, doctor, c.ID)
	require.Equal(t, http.StatusOK, status)

	api.advanceClock(t, 25*time.Hour)

	task := awaitTask(t, c.ID)
	assert.Equal(t, "no_booking", task.Task.Reason)
}

func awaitNotification(t *testing.T, caseID, kind string) notificationView {
	t.Helper()
	var found *notificationView
	require.True(t, poll(func() bool {
		found = findNotification(t, caseID, kind)
		return found != nil
	}), "no %s notification for case %s", kind, caseID)
	return *found
}

func awaitNotificationCount(t *testing.T, caseID, kind string, want int) {
	t.Helper()
	var got int
	reached := poll(func() bool {
		got = countNotifications(t, caseID, kind)
		return got == want
	})
	require.True(t, reached, "case %s: %d %s notifications, want %d", caseID, got, kind, want)
}

func awaitJobSettled(t *testing.T, caseID, kind string) {
	t.Helper()
	require.True(t, poll(func() bool {
		return !slices.Contains(jobsOf(t, caseID), jobView{Kind: kind, State: "pending"})
	}), "%s job of case %s is still pending", kind, caseID)
}
