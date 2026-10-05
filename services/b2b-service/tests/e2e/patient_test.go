//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlan_ShowsOnlyWhatPatientMaySee(t *testing.T) {
	t.Parallel()
	confirmedReport := newReport()
	confirmed := ingestInReview(t, confirmedReport, validAssessment())
	reviewCase(t, confirmed, "critical", "rejected")
	confirmCase(t, confirmed)

	pendingReport := newReport()
	pendingReport.Patient = confirmedReport.Patient
	pending := ingestInReview(t, pendingReport, validAssessment())

	urgentReport := newReport()
	urgentReport.Patient = confirmedReport.Patient
	urgentAssessment := validAssessment()
	urgentAssessment.Urgency = "emergency"
	urgent := ingestInReview(t, urgentReport, urgentAssessment)

	var (
		status int
		plan   planView
		raw    []byte
		byID   map[string]planCaseView
	)
	require.True(t, poll(func() bool {
		status, plan, raw = api.plan(t, confirmedReport.SourceSystem, confirmedReport.Patient.ID)
		byID = map[string]planCaseView{}
		for _, c := range plan.Cases {
			byID[c.CaseID] = c
		}
		_, ok := byID[urgent.ID]
		return ok
	}), "urgent case did not reach the patient plan")

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, confirmedReport.Patient.FullName, plan.Patient.FullName)
	assert.WithinDuration(t, api.clock(t), plan.Now, time.Minute)

	require.Contains(t, byID, pending.ID)
	assert.Equal(t, "in_review", byID[pending.ID].Status)
	assert.False(t, byID[pending.ID].UrgentContact)
	assert.Empty(t, byID[pending.ID].Recommendations)
	assert.Nil(t, byID[pending.ID].BookBy)

	require.Contains(t, byID, confirmed.ID)
	approved := byID[confirmed.ID]
	assert.Equal(t, "notified", approved.Status)
	assert.False(t, approved.UrgentContact)
	assert.False(t, approved.NoFindings)
	require.NotNil(t, approved.ConfirmedAt)
	require.NotNil(t, approved.BookBy)
	assert.Equal(t, approved.ConfirmedAt.Add(7*24*time.Hour), *approved.BookBy)
	assert.Equal(t, confirmedReport.Study.Modality, approved.Study.Modality)
	assert.Equal(t, []planRecommendation{{
		ID:          confirmed.Recommendations[0].ID,
		ServiceCode: confirmed.Recommendations[0].ServiceCode,
		ServiceName: confirmed.Recommendations[0].ServiceName,
		PatientText: confirmed.Recommendations[0].PatientText,
		Mark:        "critical",
	}}, approved.Recommendations)

	require.Contains(t, byID, urgent.ID)
	assert.True(t, byID[urgent.ID].UrgentContact)
	assert.Empty(t, byID[urgent.ID].Recommendations)

	for _, hidden := range []string{
		confirmedReport.Conclusion,
		confirmed.Recommendations[0].Rationale,
		confirmed.Recommendations[0].GuidelineRef,
		confirmedReport.Patient.Phone,
		"priority",
	} {
		assert.NotContains(t, string(raw), hidden)
	}
}

func TestPlan_RequestErrors(t *testing.T) {
	t.Parallel()
	status, _, _ := api.plan(t, "e2e", "PATIENT-"+uuid.NewString())
	assert.Equal(t, http.StatusNotFound, status)

	status, _, _ = api.plan(t, "", "PATIENT-"+uuid.NewString())
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestBooking_SelfBookingMovesCaseToBooked(t *testing.T) {
	t.Parallel()
	c := confirmedCase(t, bookableAssessment())
	req := bookingRequest{
		RecommendationID: c.Recommendations[0].ID,
		AppointmentID:    "APT-" + uuid.NewString(),
		ScheduledAt:      "2026-10-20T09:00:00+03:00",
		Channel:          "self",
	}

	status, first := api.registerBooking(t, c.ID, req)
	require.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "booked", first.CaseStatus)
	assert.Equal(t, req.AppointmentID, first.AppointmentID)

	status, again := api.registerBooking(t, c.ID, req)
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, first.BookingID, again.BookingID)

	assert.Equal(t, "booked", api.routeStatus(t, c))

	booked := careEventsOfType(t, c.ID, "recommendation_booked")
	require.Len(t, booked, 1)
	assert.Equal(t, "patient", booked[0].Actor)
	assert.Equal(t, map[string]string{
		"recommendation_id": req.RecommendationID,
		"appointment_id":    req.AppointmentID,
		"scheduled_at":      "2026-10-20T06:00:00Z",
		"channel":           "self",
	}, booked[0].Payload)
	transitions := careEventsOfType(t, c.ID, "status_changed")
	assert.Equal(t, "patient", transitions[len(transitions)-1].Actor)
}

func TestBooking_SecondRecommendationKeepsCaseBooked(t *testing.T) {
	t.Parallel()
	c := confirmedCase(t, bookableAssessment())

	for _, rec := range c.Recommendations {
		status, resp := api.registerBooking(t, c.ID, bookingRequest{
			RecommendationID: rec.ID,
			AppointmentID:    "APT-" + uuid.NewString(),
			ScheduledAt:      "2026-10-20T09:00:00Z",
			Channel:          "self",
		})
		require.Equal(t, http.StatusCreated, status)
		assert.Equal(t, "booked", resp.CaseStatus)
	}

	assert.Len(t, careEventsOfType(t, c.ID, "recommendation_booked"), 2)
	assert.Len(t, careEventsOfType(t, c.ID, "status_changed"), 5)
}

func TestBooking_InvalidBookingsAreRejected(t *testing.T) {
	t.Parallel()
	confirmed := ingestInReview(t, newReport(), validAssessment())
	reviewCase(t, confirmed, "critical", "rejected")
	confirmCase(t, confirmed)
	notInClinic := confirmedCase(t, validAssessment())
	inReview := newCaseInReview(t, "priority")
	taken := "APT-" + uuid.NewString()
	status, _ := api.registerBooking(t, confirmed.ID, validBooking(confirmed.Recommendations[0].ID, taken))
	require.Equal(t, http.StatusCreated, status)

	tests := []struct {
		name   string
		caseID string
		req    bookingRequest
		status int
	}{
		{"rejected recommendation is not in the plan", confirmed.ID, validBooking(confirmed.Recommendations[1].ID, ""), http.StatusNotFound},
		{"service not in clinic", notInClinic.ID, validBooking(notInClinic.Recommendations[1].ID, ""), http.StatusConflict},
		{"case still in review", inReview.ID, validBooking(inReview.Recommendations[0].ID, ""), http.StatusConflict},
		{"appointment linked to another case", notInClinic.ID, validBooking(notInClinic.Recommendations[0].ID, taken), http.StatusConflict},
		{"recommendation of another case", confirmed.ID, validBooking(notInClinic.Recommendations[0].ID, ""), http.StatusNotFound},
		{"unknown case", uuid.NewString(), validBooking(confirmed.Recommendations[0].ID, ""), http.StatusNotFound},
		{"unknown channel", confirmed.ID, withChannel(validBooking(confirmed.Recommendations[0].ID, ""), "fax"), http.StatusBadRequest},
		{"missing appointment", confirmed.ID, withAppointment(validBooking(confirmed.Recommendations[0].ID, ""), ""), http.StatusBadRequest},
		{"malformed recommendation id", confirmed.ID, validBooking("not-a-uuid", ""), http.StatusBadRequest},
	}
	for _, tt := range tests {
		status, _ := api.registerBooking(t, tt.caseID, tt.req)
		assert.Equal(t, tt.status, status, tt.name)
	}
}

func ingestInReview(t *testing.T, r report, a assessment) caseView {
	t.Helper()
	ai.on(r.Conclusion, respond(a))
	_, ref := api.ingest(t, r)
	return api.awaitStatus(t, ref.CaseID, "in_review")
}

func confirmedCase(t *testing.T, a assessment) caseView {
	t.Helper()
	c := ingestInReview(t, newReport(), a)
	reviewAll(t, c)
	confirmCase(t, c)
	return c
}

func reviewCase(t *testing.T, c caseView, marks ...string) {
	t.Helper()
	for i, mark := range marks {
		req := reviewRequest{Mark: mark}
		if mark == "rejected" {
			req.RejectReason = "contraindicated"
		}
		status, _ := api.reviewRecommendation(t, doctor, c.ID, c.Recommendations[i].ID, req)
		require.Equal(t, http.StatusOK, status)
	}
}

func confirmCase(t *testing.T, c caseView) {
	t.Helper()
	status, _ := api.confirm(t, doctor, c.ID)
	require.Equal(t, http.StatusOK, status)
	api.awaitRouteStatus(t, c, "notified")
}

func bookableAssessment() assessment {
	a := validAssessment()
	a.Recommendations[1].ServiceCode = "CT-CHEST-FOLLOWUP"
	return a
}

func validBooking(recommendationID, appointmentID string) bookingRequest {
	if appointmentID == "" {
		appointmentID = "APT-" + uuid.NewString()
	}
	return bookingRequest{
		RecommendationID: recommendationID,
		AppointmentID:    appointmentID,
		ScheduledAt:      "2026-10-20T09:00:00Z",
		Channel:          "self",
	}
}

func withChannel(r bookingRequest, channel string) bookingRequest {
	r.Channel = channel
	return r
}

func withAppointment(r bookingRequest, appointmentID string) bookingRequest {
	r.AppointmentID = appointmentID
	return r
}

func careEventsOfType(t *testing.T, caseID, eventType string) []caseEventView {
	t.Helper()
	return ofType(careEventsOf(t, caseID), eventType)
}

func doctorEventsOfType(t *testing.T, caseID, eventType string) []caseEventView {
	t.Helper()
	return ofType(doctorEventsOf(t, caseID), eventType)
}

func ofType(events []caseEventView, eventType string) []caseEventView {
	var matched []caseEventView
	for _, e := range events {
		if e.Type == eventType {
			matched = append(matched, e)
		}
	}
	return matched
}
