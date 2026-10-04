//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph_BuildsPatientJourney(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	studyDate := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	booked := newRecommendation(uniqueID("PULM"), "Консультация пульмонолога")
	done := newRecommendation(uniqueID("CT"), "Контрольная КТ")
	open := newRecommendation(uniqueID("US"), "УЗИ молочных желез")
	external := newRecommendation("", "Консультация онколога")
	mammography := newConfirmedCase("MG", studyDate, booked, done, open, external)
	urgent := planCaseDTO{
		CaseID:        uniqueID("CASE"),
		Status:        "in_review",
		UrgentContact: true,
		Study:         planStudyDTO{Modality: "DX", PerformedAt: studyDate},
	}
	b2b.setPlan(patientID, planDTO{Patient: planPatientDTO{FullName: "Анна Смирнова"}, Cases: []planCaseDTO{mammography, urgent}})

	appointmentDate := time.Date(2026, 10, 10, 9, 0, 0, 0, time.UTC)
	visitDate := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	mis.setHistory(patientID, misHistory{
		Visits: []misVisit{
			{ServiceCode: open.ServiceCode, ServiceName: "УЗИ до исследования", VisitedAt: studyDate.AddDate(-1, 0, 0)},
			{ServiceCode: done.ServiceCode, ServiceName: "КТ органов грудной клетки", VisitedAt: visitDate},
		},
		Appointments: []misAppointment{
			{ID: "APT-1", PatientID: patientID, ServiceCode: booked.ServiceCode, ServiceName: booked.ServiceName, ScheduledAt: appointmentDate, ReferralID: booked.ID},
		},
	})

	status, graph := api.graph(t, patientID)

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, sourceSystem, b2b.sourceSystemOf(patientID))
	assert.Equal(t, "Анна Смирнова", graph.Patient.FullName)
	assert.Equal(t, clinicPhone, graph.ClinicPhone)

	study := "study:" + mammography.CaseID
	assert.Equal(t, []nodeView{
		{ID: study, Type: "study", Title: "Маммография", Date: studyDate},
		recommendationNode(booked, "booked", true, false),
		{ID: "appointment:APT-1", Type: "appointment", Title: booked.ServiceName, Date: appointmentDate},
		recommendationNode(done, "done", true, false),
		{ID: "visit:" + done.ID, Type: "visit", Title: "КТ органов грудной клетки", Date: visitDate},
		recommendationNode(open, "recommended", true, true),
		recommendationNode(external, "recommended", false, false),
		{ID: "study:" + urgent.CaseID, Type: "study", Title: "Рентгенография", Date: studyDate},
		{
			ID:    "urgent:" + urgent.CaseID,
			Type:  "urgent_contact",
			Title: "Врач просит срочно связаться с клиникой",
			Text:  "Позвоните в клинику: " + clinicPhone,
		},
	}, graph.Nodes)

	assert.Equal(t, []edgeView{
		{From: study, To: "recommendation:" + booked.ID},
		{From: "recommendation:" + booked.ID, To: "appointment:APT-1"},
		{From: study, To: "recommendation:" + done.ID},
		{From: "recommendation:" + done.ID, To: "visit:" + done.ID},
		{From: study, To: "recommendation:" + open.ID},
		{From: study, To: "recommendation:" + external.ID},
		{From: "study:" + urgent.CaseID, To: "urgent:" + urgent.CaseID},
	}, graph.Edges)
}

func TestGraph_AppointmentBookedByPhoneIsMatchedByService(t *testing.T) {
	t.Parallel()
	patientID := uniqueID("P")
	studyDate := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	rec := newRecommendation(uniqueID("PULM"), "Консультация пульмонолога")
	c := newConfirmedCase("CT", studyDate, rec)
	b2b.setPlan(patientID, planDTO{Cases: []planCaseDTO{c}})
	mis.setHistory(patientID, misHistory{Appointments: []misAppointment{
		{ID: "APT-before", ServiceCode: rec.ServiceCode, ServiceName: rec.ServiceName, ScheduledAt: studyDate.AddDate(0, 0, -3)},
		{ID: "APT-phone", ServiceCode: rec.ServiceCode, ServiceName: rec.ServiceName, ScheduledAt: studyDate.AddDate(0, 0, 5)},
	}})

	status, graph := api.graph(t, patientID)

	require.Equal(t, http.StatusOK, status)
	require.Len(t, graph.Nodes, 3)
	assert.Equal(t, "booked", graph.Nodes[1].Status)
	assert.Equal(t, "appointment:APT-phone", graph.Nodes[2].ID)
}

func TestGraph_Errors(t *testing.T) {
	t.Parallel()
	status, _ := api.graph(t, uniqueID("P"))
	assert.Equal(t, http.StatusNotFound, status)

	failing := uniqueID("P")
	b2b.failPlan(failing, http.StatusInternalServerError)
	status, _ = api.graph(t, failing)
	assert.Equal(t, http.StatusBadGateway, status)
}

func recommendationNode(rec planRecommendationDTO, status string, inClinic, bookable bool) nodeView {
	return nodeView{
		ID:       "recommendation:" + rec.ID,
		Type:     "recommendation",
		Title:    rec.ServiceName,
		Text:     rec.PatientText,
		Status:   status,
		Mark:     rec.Mark,
		InClinic: &inClinic,
		Bookable: &bookable,
	}
}
