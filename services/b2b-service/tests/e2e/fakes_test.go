//go:build e2e

package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
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

type fakeMIS struct {
	mu        sync.Mutex
	histories map[string]history
}

func newFakeMIS() *fakeMIS {
	return &fakeMIS{histories: map[string]history{}}
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
		f.mu.Unlock()
		if !ok {
			h = history{Visits: []visit{}, Appointments: []appointment{}}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(h)
	})
	mux.ServeHTTP(w, r)
}
