//go:build e2e

package e2e

import (
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboard_CountsWhatHappensToLiveCases(t *testing.T) {
	_, before := api.dashboard(t, "1")

	normalAssessment := validAssessment()
	normalAssessment.Urgency = "normal"
	confirmedCase(t, normalAssessment)

	selfBooked := confirmedCase(t, bookableAssessment())
	status, _ := api.registerBooking(t, selfBooked.ID, validBooking(selfBooked.Recommendations[0].ID, ""))
	require.Equal(t, http.StatusCreated, status)
	status, _ = api.declineRecommendation(t, selfBooked.ID, selfBooked.Recommendations[1].ID, patientDeclineRequest{Reason: "not_needed"})
	require.Equal(t, http.StatusOK, status)

	_, operatorBooked := confirmedEmergencyCase(t)
	task := awaitTask(t, operatorBooked.ID)
	_, card := api.operatorCard(t, task.Task.ID)
	status, _ = api.operatorBook(t, operatorID, task.Task.ID, operatorBookingRequest{
		RecommendationID: card.Offers[0].RecommendationID,
		SlotID:           card.Offers[0].Slots[0].ID,
	})
	require.Equal(t, http.StatusCreated, status)

	_, declined := confirmedEmergencyCase(t)
	declinedTask := awaitTask(t, declined.ID)
	status, _ = api.operatorAction(t, operatorID, declinedTask.Task.ID, "decline", declineRequest{Reason: "far"})
	require.Equal(t, http.StatusOK, status)

	status, after := api.dashboard(t, "1")

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, 1, after.Period.Days)
	assert.Equal(t, "day", after.TrendStep)
	assert.Equal(t, funnelView{Received: 4, Recommended: 3, Confirmed: 3, Notified: 3, Booked: 2}, funnelDelta(before.Current.Funnel, after.Current.Funnel))
	assert.Equal(t, channelsView{Self: 1, Operator: 1}, channelsView{
		Self:     after.Current.Bookings.Self - before.Current.Bookings.Self,
		Operator: after.Current.Bookings.Operator - before.Current.Bookings.Operator,
	})
	assert.Equal(t, ratioView{OnTime: 4, Total: 4}, ratioDelta(before.Current.DoctorSLA, after.Current.DoctorSLA))
	assert.Equal(t, ratioView{OnTime: 2, Total: 2}, ratioDelta(before.Current.OperatorSLA, after.Current.OperatorSLA))
	assert.Equal(t, 1, declinesOf(after, "far")-declinesOf(before, "far"))
	assert.Equal(t, 1, declinesOf(after, "not_needed")-declinesOf(before, "not_needed"), "patient declines count too")

	emergencyBefore, emergencyAfter := urgencyOf(before, "emergency"), urgencyOf(after, "emergency")
	assert.Equal(t, 2, emergencyAfter.Confirmed-emergencyBefore.Confirmed)
	assert.Equal(t, 1, emergencyAfter.Booked-emergencyBefore.Booked)
	priorityBefore, priorityAfter := urgencyOf(before, "priority"), urgencyOf(after, "priority")
	assert.Equal(t, 1, priorityAfter.Booked-priorityBefore.Booked)

	require.NotEmpty(t, after.Trend)
	today, todayBefore := after.Trend[len(after.Trend)-1], before.Trend[len(before.Trend)-1]
	assert.Equal(t, 4, today.Received-todayBefore.Received)
	assert.Equal(t, 2, today.Booked-todayBefore.Booked)
}

func TestDashboard_PeriodMustBeValid(t *testing.T) {
	t.Parallel()
	for _, days := range []string{"", "0", "366", "week"} {
		status, _ := api.dashboard(t, days)
		assert.Equal(t, http.StatusBadRequest, status, days)
	}

	status, d := api.dashboard(t, "90")
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "week", d.TrendStep)
	assert.Len(t, d.Trend, 13)
	assert.Equal(t, d.Period.From, d.PreviousPeriod.To)
	assert.Equal(t, []string{"emergency", "priority", "planned", "normal"}, urgencyKeys(d))
}

func funnelDelta(before, after funnelView) funnelView {
	return funnelView{
		Received:    after.Received - before.Received,
		Recommended: after.Recommended - before.Recommended,
		Confirmed:   after.Confirmed - before.Confirmed,
		Notified:    after.Notified - before.Notified,
		Booked:      after.Booked - before.Booked,
		Completed:   after.Completed - before.Completed,
	}
}

func ratioDelta(before, after ratioView) ratioView {
	return ratioView{OnTime: after.OnTime - before.OnTime, Total: after.Total - before.Total}
}

func declinesOf(d dashboardView, reason string) int {
	for _, c := range d.DeclineReasons {
		if c.Reason == reason {
			return c.Count
		}
	}
	return 0
}

func urgencyOf(d dashboardView, urgency string) urgencySliceView {
	i := slices.IndexFunc(d.ByUrgency, func(s urgencySliceView) bool { return s.Urgency == urgency })
	if i < 0 {
		return urgencySliceView{}
	}
	return d.ByUrgency[i]
}

func urgencyKeys(d dashboardView) []string {
	keys := make([]string, 0, len(d.ByUrgency))
	for _, s := range d.ByUrgency {
		keys = append(keys, s.Urgency)
	}
	return keys
}
