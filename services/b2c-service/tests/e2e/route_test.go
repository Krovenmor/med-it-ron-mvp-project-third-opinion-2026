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

var now = time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)

func day(month time.Month, d, hour int) time.Time {
	return time.Date(2026, month, d, hour, 0, 0, 0, time.UTC)
}

func TestRoute_BuildsTreeFromPlanAndHistory(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")

	done := newRecommendation(uniqueID("PULM"), "Консультация пульмонолога")
	booked := minor(newRecommendation(uniqueID("US"), "УЗИ"))
	open := newRecommendation(uniqueID("CTF"), "Контрольная КТ")
	external := newRecommendation("", "Консультация онколога")
	declined := minor(newRecommendation(uniqueID("SPIRO"), "Спирометрия"))
	declined.Declined, declined.DeclineReason = true, "far"
	current := newConfirmedCase("CT", day(10, 2, 9), done, booked, open, external, declined)

	overdue := newRecommendation(uniqueID("OLD"), "Повторный приём")
	late := minor(newRecommendation(uniqueID("LAB"), "Биохимия крови"))
	previous := newConfirmedCase("DX", day(6, 1, 9), overdue, late)
	pending := planCaseDTO{CaseID: uniqueID("CASE"), Status: "in_review", Study: planStudyDTO{Modality: "MG", PerformedAt: day(10, 3, 9)}}

	b2b.setPlan(patientID, planDTO{
		Now:     now,
		Patient: planPatientDTO{FullName: "Игорь Петров"},
		Cases:   []planCaseDTO{pending, current, previous},
	})
	mis.setHistory(patientID, misHistory{
		Visits: []misVisit{
			{ServiceCode: "THER-CONSULT", ServiceName: "Приём терапевта", VisitedAt: day(9, 30, 10)},
			{ServiceCode: "LAB-CBC", ServiceName: "Общий анализ крови", VisitedAt: day(10, 1, 9)},
			{ServiceCode: done.ServiceCode, ServiceName: done.ServiceName, VisitedAt: day(10, 4, 8)},
			{ServiceCode: late.ServiceCode, ServiceName: late.ServiceName, VisitedAt: day(7, 20, 9)},
		},
		Appointments: []misAppointment{
			{ID: "APT-US", ServiceCode: booked.ServiceCode, ServiceName: booked.ServiceName, ScheduledAt: day(10, 6, 11), ReferralID: booked.ID},
		},
	})

	status, route := api.route(t, patientID)

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, now, route.Now)
	assert.Equal(t, "Игорь Петров", route.Patient.FullName)
	assert.Equal(t, clinicView{Name: clinicName, Phone: clinicPhone, Address: clinicAddress}, route.Clinic)
	assert.True(t, route.PendingReview)
	assert.False(t, route.UrgentContact)
	assert.False(t, route.NoFindings)

	assert.Equal(t, []trunkNodeView{
		{ID: "study:" + previous.CaseID, Kind: "study", Title: "Рентгенография", Date: day(6, 1, 9), AIReport: true},
		{ID: "visit:0", Kind: "visit", Title: "Приём терапевта", Date: day(9, 30, 10)},
		{ID: "study:" + current.CaseID, Kind: "study", Title: "Компьютерная томография", Date: day(10, 2, 9), AIReport: true},
		{ID: "study:" + pending.CaseID, Kind: "study", Title: "Маммография", Date: day(10, 3, 9), AIReport: true, Pending: true},
	}, route.Trunk)

	due := *current.BookBy
	steps := byID(route.Steps)
	assertStep(t, steps[done.ID], "done", "study:"+current.CaseID, day(10, 4, 8), true, false)
	assert.False(t, steps[done.ID].Late)
	assertStep(t, steps[booked.ID], "booked", "study:"+current.CaseID, day(10, 6, 11), true, false)
	require.NotNil(t, steps[booked.ID].Appointment)
	assert.Equal(t, "APT-US", steps[booked.ID].Appointment.ID)
	assertStep(t, steps[open.ID], "recommended", "study:"+current.CaseID, due, true, true)
	assertStep(t, steps[external.ID], "recommended", "study:"+current.CaseID, due, false, false)
	assertStep(t, steps[declined.ID], "declined", "study:"+current.CaseID, due, true, false)
	assert.Equal(t, "far", steps[declined.ID].DeclineReason)
	assertStep(t, steps[overdue.ID], "overdue", "study:"+previous.CaseID, *previous.BookBy, true, true)
	assertStep(t, steps[late.ID], "done", "study:"+previous.CaseID, day(7, 20, 9), true, false)
	assert.True(t, steps[late.ID].Late)
	assert.Equal(t, open.PatientText, steps[open.ID].Text)
	assert.Equal(t, *current.ConfirmedAt, steps[open.ID].AssignedAt)
	assert.Equal(t, due, steps[open.ID].DueAt)

	lab := steps["history:0"]
	assert.Equal(t, stepView{ID: "history:0", Kind: "visit", OriginID: "visit:0", Title: "Общий анализ крови", Status: "done", Date: day(10, 1, 9), VisitedAt: day(10, 1, 9)}, lab)

	require.NotNil(t, route.Progress)
	assert.Equal(t, 44, *route.Progress, "on time 3 + late 0.5 of overdue 3, declined 1, on time 3, late 1")
	assert.Equal(t, overdue.ID, route.NextStepID, "overdue steps come first")
	assert.Equal(t, []appointmentView{{ID: "APT-US", ServiceName: booked.ServiceName, ScheduledAt: day(10, 6, 11)}}, route.Appointments)
}

func TestRoute_EarliestVisitBecomesTheRoot(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	c := newConfirmedCase("MG", day(10, 2, 9), newRecommendation(uniqueID("US"), "УЗИ"))
	b2b.setPlan(patientID, planDTO{Now: now, Cases: []planCaseDTO{c}})
	mis.setHistory(patientID, misHistory{Visits: []misVisit{
		{ServiceCode: "MG-SCREENING", ServiceName: "Маммография (скрининг)", VisitedAt: day(1, 15, 9)},
		{ServiceCode: "LAB-CBC", ServiceName: "Общий анализ крови", VisitedAt: day(2, 1, 9)},
	}})

	_, route := api.route(t, patientID)

	require.Len(t, route.Trunk, 2)
	assert.Equal(t, trunkNodeView{ID: "root:0", Kind: "visit", Title: "Маммография (скрининг)", Date: day(1, 15, 9)}, route.Trunk[0])
	assert.Equal(t, "root:0", byID(route.Steps)["history:1"].OriginID)
	assert.Nil(t, route.Progress, "nothing is due or done yet")
	assert.Equal(t, c.Recommendations[0].ID, route.NextStepID)
}

func TestRoute_UrgentCaseOffersACallInsteadOfBooking(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	rec := newRecommendation(uniqueID("SURG"), "Консультация хирурга")
	urgent := newConfirmedCase("DX", day(10, 3, 9), rec)
	urgent.UrgentContact = true
	urgent.BookBy = urgent.ConfirmedAt
	b2b.setPlan(patientID, planDTO{Now: now, Cases: []planCaseDTO{urgent}})

	_, route := api.route(t, patientID)

	assert.True(t, route.UrgentContact)
	assert.False(t, route.PendingReview)
	require.Len(t, route.Steps, 1)
	step := route.Steps[0]
	assert.False(t, step.Bookable)
	assert.Equal(t, "recommended", step.Status, "an urgent step waits for the call, it is never overdue")
	assert.True(t, step.DueAt.IsZero())
	assert.Equal(t, *urgent.ConfirmedAt, step.Date)
	assert.Nil(t, route.Progress)
}

func TestRoute_NormalStudyIsMarkedAsNoFindings(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	screening := minor(newRecommendation("SCREENING-ANNUAL", "Плановый профилактический осмотр"))
	normal := newConfirmedCase("DX", day(10, 3, 9), screening)
	normal.NoFindings = true
	nextYear := day(10, 3, 15).AddDate(1, 0, 0)
	normal.BookBy = &nextYear
	b2b.setPlan(patientID, planDTO{Now: now, Cases: []planCaseDTO{normal}})

	_, route := api.route(t, patientID)

	assert.True(t, route.NoFindings)
	assert.Equal(t, nextYear, byID(route.Steps)[screening.ID].Date)
}

func TestRoute_Errors(t *testing.T) {
	t.Parallel()
	status, _ := api.route(t, uniqueID("P"))
	assert.Equal(t, http.StatusNotFound, status)

	failing := uniqueID("P")
	b2b.failPlan(failing, http.StatusInternalServerError)
	status, _ = api.route(t, failing)
	assert.Equal(t, http.StatusBadGateway, status)
}

func TestSteps_DeclineAndHelpAreForwardedToB2B(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	rec := newRecommendation(uniqueID("PULM"), "Консультация пульмонолога")
	c := newConfirmedCase("CT", day(10, 2, 9), rec)
	b2b.setPlan(patientID, planDTO{Now: now, Cases: []planCaseDTO{c}})

	assert.Equal(t, http.StatusNoContent, api.decline(t, patientID, rec.ID, declineRequest{Reason: "other", Comment: "уже записан в другой клинике"}))
	assert.Equal(t, []b2bDeclineCall{{RecommendationID: rec.ID, Reason: "other", Comment: "уже записан в другой клинике"}}, b2b.declinesOf(c.CaseID))

	assert.Equal(t, http.StatusAccepted, api.requestHelp(t, patientID, rec.ID))
	assert.Equal(t, 1, b2b.helpRequestsOf(c.CaseID))

	assert.Equal(t, http.StatusNotFound, api.decline(t, patientID, uniqueID("REC"), declineRequest{Reason: "far"}))
	assert.Equal(t, http.StatusNotFound, api.requestHelp(t, uniqueID("P"), rec.ID))
}

func TestSteps_DeclineRejectedByB2BKeepsItsMeaning(t *testing.T) {
	patientID := uniqueID("P")
	rec := newRecommendation(uniqueID("PULM"), "Консультация пульмонолога")
	b2b.setPlan(patientID, planDTO{Now: now, Cases: []planCaseDTO{newConfirmedCase("CT", day(10, 2, 9), rec)}})
	t.Cleanup(func() { b2b.rejectDeclines(0) })

	b2b.rejectDeclines(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, api.decline(t, patientID, rec.ID, declineRequest{Reason: "other"}))
	b2b.rejectDeclines(http.StatusConflict)
	assert.Equal(t, http.StatusConflict, api.decline(t, patientID, rec.ID, declineRequest{Reason: "far"}))
}

func minor(r planRecommendationDTO) planRecommendationDTO {
	r.Mark = "minor"
	return r
}

func byID(steps []stepView) map[string]stepView {
	out := make(map[string]stepView, len(steps))
	for _, s := range steps {
		out[s.ID] = s
	}
	return out
}

func assertStep(t *testing.T, s stepView, status, origin string, date time.Time, inClinic, bookable bool) {
	t.Helper()
	assert.Equal(t, status, s.Status, s.Title)
	assert.Equal(t, origin, s.OriginID, s.Title)
	assert.Equal(t, date, s.Date, s.Title)
	assert.Equal(t, inClinic, s.InClinic, s.Title)
	assert.Equal(t, bookable, s.Bookable, s.Title)
	assert.True(t, slices.Contains([]string{"critical", "minor"}, s.Mark), s.Title)
}
