//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"sync"
)

type fakeB2B struct {
	mu            sync.Mutex
	plans         map[string]planDTO
	planStatus    map[string]int
	bookingStatus map[string]int
	bookings      map[string][]b2bBookingCall
	declines      map[string][]b2bDeclineCall
	declineStatus int
	helpRequests  map[string]int
	sourceSystems map[string]string
	mux           *http.ServeMux
}

func newFakeB2B() *fakeB2B {
	f := &fakeB2B{
		plans:         map[string]planDTO{},
		planStatus:    map[string]int{},
		bookingStatus: map[string]int{},
		bookings:      map[string][]b2bBookingCall{},
		declines:      map[string][]b2bDeclineCall{},
		helpRequests:  map[string]int{},
		sourceSystems: map[string]string{},
		mux:           http.NewServeMux(),
	}
	f.mux.HandleFunc("GET /api/v1/patients/{id}/plan", f.plan)
	f.mux.HandleFunc("POST /api/v1/cases/{id}/bookings", f.book)
	f.mux.HandleFunc("POST /api/v1/cases/{id}/recommendations/{rec_id}/decline", f.decline)
	f.mux.HandleFunc("POST /api/v1/cases/{id}/help-requests", f.help)
	return f
}

func (f *fakeB2B) declinesOf(caseID string) []b2bDeclineCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]b2bDeclineCall(nil), f.declines[caseID]...)
}

func (f *fakeB2B) helpRequestsOf(caseID string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.helpRequests[caseID]
}

func (f *fakeB2B) rejectDeclines(status int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.declineStatus = status
}

func (f *fakeB2B) decline(w http.ResponseWriter, r *http.Request) {
	var call b2bDeclineCall
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	call.RecommendationID = r.PathValue("rec_id")
	f.mu.Lock()
	status := f.declineStatus
	if status == 0 {
		f.declines[r.PathValue("id")] = append(f.declines[r.PathValue("id")], call)
	}
	f.mu.Unlock()
	if status != 0 {
		w.WriteHeader(status)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"recommendation_id": call.RecommendationID, "case_status": "notified"})
}

func (f *fakeB2B) help(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.helpRequests[r.PathValue("id")]++
	f.mu.Unlock()
	writeJSON(w, http.StatusCreated, map[string]any{"task_id": "task", "created": true})
}

func (f *fakeB2B) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mux.ServeHTTP(w, r)
}

func (f *fakeB2B) setPlan(patientID string, p planDTO) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.plans[patientID] = p
}

func (f *fakeB2B) failPlan(patientID string, status int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.planStatus[patientID] = status
}

func (f *fakeB2B) failBookings(caseID string, status int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if status == 0 {
		delete(f.bookingStatus, caseID)
		return
	}
	f.bookingStatus[caseID] = status
}

func (f *fakeB2B) bookingsOf(caseID string) []b2bBookingCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]b2bBookingCall(nil), f.bookings[caseID]...)
}

func (f *fakeB2B) sourceSystemOf(patientID string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.sourceSystems[patientID]
}

func (f *fakeB2B) plan(w http.ResponseWriter, r *http.Request) {
	patientID := r.PathValue("id")
	f.mu.Lock()
	f.sourceSystems[patientID] = r.URL.Query().Get("source_system")
	status, failing := f.planStatus[patientID]
	plan, found := f.plans[patientID]
	f.mu.Unlock()

	switch {
	case failing:
		w.WriteHeader(status)
	case !found:
		w.WriteHeader(http.StatusNotFound)
	default:
		writeJSON(w, http.StatusOK, plan)
	}
}

func (f *fakeB2B) book(w http.ResponseWriter, r *http.Request) {
	var call b2bBookingCall
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	caseID := r.PathValue("id")
	f.mu.Lock()
	status, failing := f.bookingStatus[caseID]
	if !failing {
		f.bookings[caseID] = append(f.bookings[caseID], call)
	}
	f.mu.Unlock()

	if failing {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type fakeMIS struct {
	mu         sync.Mutex
	histories  map[string]misHistory
	slots      map[string][]misSlot
	bookStatus map[string]int
	bookings   []misBookingCall
	mux        *http.ServeMux
}

func newFakeMIS() *fakeMIS {
	f := &fakeMIS{
		histories:  map[string]misHistory{},
		slots:      map[string][]misSlot{},
		bookStatus: map[string]int{},
		mux:        http.NewServeMux(),
	}
	f.mux.HandleFunc("GET /api/v1/patients/{id}/history", f.history)
	f.mux.HandleFunc("GET /api/v1/slots", f.listSlots)
	f.mux.HandleFunc("POST /api/v1/appointments", f.book)
	return f
}

func (f *fakeMIS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mux.ServeHTTP(w, r)
}

func (f *fakeMIS) setHistory(patientID string, h misHistory) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.histories[patientID] = h
}

func (f *fakeMIS) setSlots(serviceCode string, slots ...misSlot) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.slots[serviceCode] = slots
}

func (f *fakeMIS) failBooking(slotID string, status int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bookStatus[slotID] = status
}

func (f *fakeMIS) bookingsOf(patientID string) []misBookingCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	var calls []misBookingCall
	for _, c := range f.bookings {
		if c.PatientID == patientID {
			calls = append(calls, c)
		}
	}
	return calls
}

func (f *fakeMIS) history(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	h := f.histories[r.PathValue("id")]
	f.mu.Unlock()
	if h.Visits == nil {
		h.Visits = []misVisit{}
	}
	if h.Appointments == nil {
		h.Appointments = []misAppointment{}
	}
	writeJSON(w, http.StatusOK, h)
}

func (f *fakeMIS) listSlots(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	slots := f.slots[r.URL.Query().Get("service_code")]
	f.mu.Unlock()
	writeJSON(w, http.StatusOK, misSlots{Slots: append([]misSlot{}, slots...)})
}

func (f *fakeMIS) book(w http.ResponseWriter, r *http.Request) {
	var call misBookingCall
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.bookings = append(f.bookings, call)
	if status, failing := f.bookStatus[call.SlotID]; failing {
		w.WriteHeader(status)
		return
	}
	slot, found := f.findSlot(call.SlotID)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appointment := misAppointment{
		ID:          "APT-" + call.SlotID,
		PatientID:   call.PatientID,
		ServiceCode: slot.ServiceCode,
		ServiceName: call.ServiceName,
		ScheduledAt: slot.StartsAt,
		ReferralID:  call.ReferralID,
	}
	h := f.histories[call.PatientID]
	h.Appointments = append(h.Appointments, appointment)
	f.histories[call.PatientID] = h
	writeJSON(w, http.StatusCreated, appointment)
}

func (f *fakeMIS) findSlot(id string) (misSlot, bool) {
	for _, slots := range f.slots {
		for _, s := range slots {
			if s.ID == id {
				return s, true
			}
		}
	}
	return misSlot{}, false
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
