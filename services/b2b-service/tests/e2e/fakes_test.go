//go:build e2e

package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type aiBehavior func(call int) (status int, body any)

func respond(body any) aiBehavior {
	return func(int) (int, any) { return http.StatusOK, body }
}

func failWith(status int) aiBehavior {
	return func(int) (int, any) { return status, map[string]string{"error": "ai-service is down"} }
}

func failFirst(failures, status int, then any) aiBehavior {
	return func(call int) (int, any) {
		if call <= failures {
			return status, map[string]string{"error": "temporary failure"}
		}
		return http.StatusOK, then
	}
}

type fakeAI struct {
	mu        sync.Mutex
	behaviors map[string]aiBehavior
	requests  map[string][][]byte
}

func newFakeAI() *fakeAI {
	return &fakeAI{behaviors: map[string]aiBehavior{}, requests: map[string][][]byte{}}
}

func (f *fakeAI) on(conclusion string, b aiBehavior) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.behaviors[conclusion] = b
}

func (f *fakeAI) requestsFor(conclusion string) [][]byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([][]byte(nil), f.requests[conclusion]...)
}

func (f *fakeAI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/v1/assessments" {
		http.NotFound(w, r)
		return
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req aiRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f.mu.Lock()
	f.requests[req.Conclusion] = append(f.requests[req.Conclusion], raw)
	call := len(f.requests[req.Conclusion])
	behave, ok := f.behaviors[req.Conclusion]
	f.mu.Unlock()
	if !ok {
		behave = respond(validAssessment())
	}

	status, body := behave(call)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

var slotStarts = []time.Time{
	time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC),
	time.Date(2026, 10, 12, 10, 0, 0, 0, time.UTC),
	time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC),
}

type fakeMIS struct {
	mu        sync.Mutex
	histories map[string]history
	failures  map[string]int
	taken     map[string]bool
	bookings  []misBookingRequest
}

func newFakeMIS() *fakeMIS {
	return &fakeMIS{histories: map[string]history{}, failures: map[string]int{}, taken: map[string]bool{}}
}

func (f *fakeMIS) takeSlot(slotID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.taken[slotID] = true
}

func (f *fakeMIS) bookingsOf(patientID string) []misBookingRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	var found []misBookingRequest
	for _, b := range f.bookings {
		if b.PatientID == patientID {
			found = append(found, b)
		}
	}
	return found
}

func slotsFor(serviceCode string) []misSlot {
	slots := make([]misSlot, 0, len(slotStarts))
	for _, start := range slotStarts {
		slots = append(slots, misSlot{ID: serviceCode + "@" + start.Format(time.RFC3339), ServiceCode: serviceCode, StartsAt: start})
	}
	return slots
}

func (f *fakeMIS) failHistory(patientID string, status int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failures[patientID] = status
}

func (f *fakeMIS) setHistory(patientID string, h history) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.histories[patientID] = h
}

func (f *fakeMIS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/patients/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		h, ok := f.histories[r.PathValue("id")]
		status, failing := f.failures[r.PathValue("id")]
		f.mu.Unlock()
		if failing {
			w.WriteHeader(status)
			return
		}
		if !ok {
			h = history{Visits: []visit{}, Appointments: []appointment{}}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(h)
	})
	mux.HandleFunc("GET /api/v1/slots", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(misSlotsResponse{Slots: slotsFor(r.URL.Query().Get("service_code"))})
	})
	mux.HandleFunc("POST /api/v1/appointments", func(w http.ResponseWriter, r *http.Request) {
		var req misBookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		serviceCode, start, ok := parseSlot(req.SlotID)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		taken := f.taken[req.SlotID]
		if !taken {
			f.bookings = append(f.bookings, req)
		}
		f.mu.Unlock()
		if taken {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(misAppointment{
			ID:          "APT-" + req.ReferralID,
			PatientID:   req.PatientID,
			ServiceCode: serviceCode,
			ServiceName: req.ServiceName,
			ScheduledAt: start,
			ReferralID:  req.ReferralID,
		})
	})
	mux.ServeHTTP(w, r)
}

func parseSlot(id string) (string, time.Time, bool) {
	serviceCode, rawStart, ok := strings.Cut(id, "@")
	if !ok {
		return "", time.Time{}, false
	}
	start, err := time.Parse(time.RFC3339, rawStart)
	return serviceCode, start, err == nil
}
