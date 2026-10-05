//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bookingScenario struct {
	patientID string
	planCase  planCaseDTO
	rec       planRecommendationDTO
	external  planRecommendationDTO
	slots     []misSlot
}

func newBookingScenario() bookingScenario {
	rec := newRecommendation(uniqueID("US"), "УЗИ молочных желез")
	external := newRecommendation("", "Консультация онколога")
	s := bookingScenario{
		patientID: uniqueID("P"),
		planCase:  newConfirmedCase("MG", time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC), rec, external),
		rec:       rec,
		external:  external,
		slots: []misSlot{
			{ID: rec.ServiceCode + ".1", ServiceCode: rec.ServiceCode, StartsAt: time.Date(2026, 10, 12, 6, 0, 0, 0, time.UTC)},
			{ID: rec.ServiceCode + ".2", ServiceCode: rec.ServiceCode, StartsAt: time.Date(2026, 10, 12, 8, 0, 0, 0, time.UTC)},
		},
	}
	b2b.setPlan(s.patientID, planDTO{Cases: []planCaseDTO{s.planCase}})
	mis.setSlots(rec.ServiceCode, s.slots...)
	return s
}

func TestSlots_ListsSlotsForRecommendedService(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()

	status, slots := api.slots(t, s.patientID, s.rec.ID)

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, []slotView{
		{ID: s.slots[0].ID, StartsAt: s.slots[0].StartsAt},
		{ID: s.slots[1].ID, StartsAt: s.slots[1].StartsAt},
	}, slots.Slots)
}

func TestSlots_Errors(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()

	status, _ := api.slots(t, s.patientID, s.external.ID)
	assert.Equal(t, http.StatusConflict, status, "service not in clinic")
	status, _ = api.slots(t, s.patientID, uniqueID("REC"))
	assert.Equal(t, http.StatusNotFound, status, "unknown recommendation")
	status, _ = api.slots(t, uniqueID("P"), s.rec.ID)
	assert.Equal(t, http.StatusNotFound, status, "unknown patient")
}

func TestBooking_BooksSlotAndRegistersItInB2B(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()
	slot := s.slots[1]

	status, appointment := api.book(t, s.patientID, s.rec.ID, bookRequest{SlotID: slot.ID})

	require.Equal(t, http.StatusCreated, status)
	assert.Equal(t, appointmentView{ID: "APT-" + slot.ID, ServiceName: s.rec.ServiceName, ScheduledAt: slot.StartsAt}, appointment)
	assert.Equal(t, []misBookingCall{{
		PatientID:   s.patientID,
		SlotID:      slot.ID,
		ServiceName: s.rec.ServiceName,
		ReferralID:  s.rec.ID,
	}}, mis.bookingsOf(s.patientID))
	assert.Equal(t, []b2bBookingCall{{
		RecommendationID: s.rec.ID,
		AppointmentID:    appointment.ID,
		ScheduledAt:      slot.StartsAt,
		Channel:          "self",
	}}, b2b.bookingsOf(s.planCase.CaseID))

	_, route := api.route(t, s.patientID)
	step := byID(route.Steps)[s.rec.ID]
	assert.Equal(t, "booked", step.Status)
	assert.False(t, step.Bookable)
	require.NotNil(t, step.Appointment)
	assert.Equal(t, appointment.ID, step.Appointment.ID)
}

func TestBooking_RetryAfterB2BFailureRepairsRegistration(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()
	b2b.failBookings(s.planCase.CaseID, http.StatusServiceUnavailable)

	status, _ := api.book(t, s.patientID, s.rec.ID, bookRequest{SlotID: s.slots[0].ID})
	require.Equal(t, http.StatusBadGateway, status)
	assert.Empty(t, b2b.bookingsOf(s.planCase.CaseID))

	b2b.failBookings(s.planCase.CaseID, 0)
	status, _ = api.book(t, s.patientID, s.rec.ID, bookRequest{SlotID: s.slots[0].ID})

	assert.Equal(t, http.StatusConflict, status)
	assert.Len(t, mis.bookingsOf(s.patientID), 1)
	assert.Equal(t, []b2bBookingCall{{
		RecommendationID: s.rec.ID,
		AppointmentID:    "APT-" + s.slots[0].ID,
		ScheduledAt:      s.slots[0].StartsAt,
		Channel:          "self",
	}}, b2b.bookingsOf(s.planCase.CaseID))
}

func TestBooking_TakenSlotIsConflict(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()
	mis.failBooking(s.slots[0].ID, http.StatusConflict)

	status, _ := api.book(t, s.patientID, s.rec.ID, bookRequest{SlotID: s.slots[0].ID})

	assert.Equal(t, http.StatusConflict, status)
	assert.Empty(t, b2b.bookingsOf(s.planCase.CaseID))
}

func TestBooking_RequestErrors(t *testing.T) {
	t.Parallel()
	s := newBookingScenario()
	path := "/api/v1/patients/" + s.patientID + "/recommendations/" + s.rec.ID + "/appointments"

	status, _ := api.book(t, s.patientID, s.rec.ID, bookRequest{})
	assert.Equal(t, http.StatusBadRequest, status, "missing slot")
	status, _ = api.do(t, http.MethodPost, path, []byte(`{"slot_id":`))
	assert.Equal(t, http.StatusBadRequest, status, "malformed body")
	status, _ = api.book(t, s.patientID, s.external.ID, bookRequest{SlotID: s.slots[0].ID})
	assert.Equal(t, http.StatusConflict, status, "service not in clinic")
	status, _ = api.book(t, s.patientID, s.rec.ID, bookRequest{SlotID: "unknown-slot"})
	assert.Equal(t, http.StatusBadRequest, status, "unknown slot")

	assert.Empty(t, b2b.bookingsOf(s.planCase.CaseID))
}
