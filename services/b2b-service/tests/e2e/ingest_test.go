//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngest_NewReportBecomesCaseReadyForReview(t *testing.T) {
	t.Parallel()
	r := newReport()
	want := validAssessment()
	ai.on(r.Conclusion, respond(want))

	status, ref := api.ingest(t, r)
	require.Equal(t, http.StatusCreated, status)
	require.NotEmpty(t, ref.CaseID)
	assert.Equal(t, "received", ref.Status)

	c := api.awaitStatus(t, ref.CaseID, "in_review")
	assert.Equal(t, r.Study.ID, c.StudyID)
	assert.Equal(t, want.Urgency, c.Urgency)
	assert.Equal(t, want.CatalogVersion, c.CatalogVersion)
	assert.Equal(t, want.GuidelinesVersion, c.GuidelinesVersion)
	assert.Equal(t, time.UTC, c.ReceivedAt.Location())
	assert.Equal(t, time.UTC, c.UpdatedAt.Location())

	require.Len(t, c.Recommendations, len(want.Recommendations))
	for i, rec := range c.Recommendations {
		assert.Equal(t, i+1, rec.Position)
		assert.Equal(t, "ai", rec.Source)
		assert.Equal(t, want.Recommendations[i], rec.recommendation)
	}
}

func TestIngest_AIServiceReceivesDepersonalizedRequestWithHistory(t *testing.T) {
	t.Parallel()
	r := newReport()
	h := history{
		Visits: []visit{
			{ServiceCode: "THER-CONSULT", ServiceName: "Приём терапевта", VisitedAt: time.Date(2026, 8, 20, 9, 0, 0, 0, time.UTC)},
		},
		Appointments: []appointment{
			{ServiceCode: "PULM-CONSULT", ServiceName: "Консультация пульмонолога", ScheduledAt: time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)},
		},
	}
	mis.setHistory(r.Patient.ID, h)

	_, ref := api.ingest(t, r)
	api.awaitStatus(t, ref.CaseID, "in_review")

	requests := ai.requestsFor(r.Conclusion)
	require.Len(t, requests, 1)
	raw := string(requests[0])
	for _, personal := range []string{r.Patient.ID, r.Patient.FullName, r.Patient.Phone, r.Patient.Email, r.Patient.BirthDate, r.Study.ID} {
		assert.NotContains(t, raw, personal)
	}
	assert.Contains(t, raw, `"current_recommendations":[]`)

	var got aiRequest
	require.NoError(t, json.Unmarshal(requests[0], &got))
	assert.Equal(t, ref.CaseID, got.CaseID)
	assert.Equal(t, r.Study.Modality, got.Study.Modality)
	assert.Equal(t, r.Study.BodySite, got.Study.BodySite)
	assert.True(t, got.Study.PerformedAt.Equal(time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)))
	assert.Equal(t, r.Conclusion, got.Conclusion)
	assert.Equal(t, ageOn(time.Date(1974, 3, 12, 0, 0, 0, 0, time.UTC), api.clock(t).UTC()), got.Patient.Age)
	assert.Equal(t, r.Patient.Sex, got.Patient.Sex)
	assert.Equal(t, h.Visits, got.Visits)
	assert.Equal(t, h.Appointments, got.Appointments)
}

func TestIngest_RepeatedReportReturnsExistingCase(t *testing.T) {
	t.Parallel()
	r := newReport()

	firstStatus, first := api.ingest(t, r)
	require.Equal(t, http.StatusCreated, firstStatus)

	secondStatus, second := api.ingest(t, r)
	require.Equal(t, http.StatusOK, secondStatus)
	assert.Equal(t, first.CaseID, second.CaseID)
}

func TestIngest_RepeatedReportWithNewContactsUpdatesPatient(t *testing.T) {
	t.Parallel()
	r := newReport()
	_, first := api.ingest(t, r)

	r.Patient.Phone = "+7-updated-" + r.Patient.ID
	status, second := api.ingest(t, r)

	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, first.CaseID, second.CaseID)
	assert.Equal(t, r.Patient.Phone, patientPhone(t, r.Patient.ID))
	assert.True(t, poll(func() bool { return patientContactsSynced(t, r.Patient.ID, r.Patient.Phone) }), "care module did not receive new contacts")
}

func TestIngest_ReportWithDifferentContentIsRejectedAndRolledBack(t *testing.T) {
	t.Parallel()
	r := newReport()
	originalPhone := r.Patient.Phone
	api.ingest(t, r)

	r.Conclusion += " (исправлено)"
	r.Patient.Phone = "+7-changed-" + r.Patient.ID
	status, _ := api.ingest(t, r)

	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, originalPhone, patientPhone(t, r.Patient.ID))
	assert.Equal(t, 1, countByStudy(t, query.CountCasesByStudy, r.Study.ID))
}

func TestIngest_ConcurrentDuplicatesCreateSingleCase(t *testing.T) {
	t.Parallel()
	const senders = 20
	r := newReport()
	body := mustJSON(t, r)

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		errs     []error
		statuses = map[int]int{}
		caseIDs  = map[string]struct{}{}
	)
	for range senders {
		wg.Go(func() {
			status, resp, err := api.send(http.MethodPost, "/api/v1/reports", body, "")
			var ref caseRef
			if err == nil {
				err = json.Unmarshal(resp, &ref)
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			statuses[status]++
			caseIDs[ref.CaseID] = struct{}{}
		})
	}
	wg.Wait()

	require.Empty(t, errs)
	assert.Equal(t, map[int]int{http.StatusCreated: 1, http.StatusOK: senders - 1}, statuses)
	assert.Len(t, caseIDs, 1)
	assert.Equal(t, 1, countByStudy(t, query.CountCasesByStudy, r.Study.ID))
	assert.Equal(t, 1, countByStudy(t, query.CountJobsByStudy, r.Study.ID))
	assert.Equal(t, 1, countByStudy(t, query.CountReportEventsByStudy, r.Study.ID))
}

func TestIngest_InvalidReportIsRejected(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mutate  func(r *report)
		message string
	}{
		{"missing source system", func(r *report) { r.SourceSystem = "" }, "source_system"},
		{"missing study id", func(r *report) { r.Study.ID = "" }, "study.id"},
		{"unknown modality", func(r *report) { r.Study.Modality = "MRI" }, "study.modality"},
		{"missing performed at", func(r *report) { r.Study.PerformedAt = "" }, "study.performed_at is required"},
		{"malformed performed at", func(r *report) { r.Study.PerformedAt = "01.10.2026" }, "study.performed_at must be"},
		{"missing conclusion", func(r *report) { r.Conclusion = "" }, "conclusion"},
		{"missing patient id", func(r *report) { r.Patient.ID = "" }, "patient.id"},
		{"missing patient name", func(r *report) { r.Patient.FullName = "" }, "patient.full_name"},
		{"malformed birth date", func(r *report) { r.Patient.BirthDate = "12.03.1974" }, "patient.birth_date must be"},
		{"unknown sex", func(r *report) { r.Patient.Sex = "unknown" }, "patient.sex"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newReport()
			tt.mutate(&r)

			status, body := api.do(t, http.MethodPost, "/api/v1/reports", mustJSON(t, r))

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Contains(t, string(body), tt.message)
			assert.Zero(t, countByStudy(t, query.CountCasesByStudy, r.Study.ID))
		})
	}

	t.Run("malformed json", func(t *testing.T) {
		t.Parallel()
		status, _ := api.do(t, http.MethodPost, "/api/v1/reports", []byte(`{"source_system":`))
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func ageOn(birth, at time.Time) int {
	age := at.Year() - birth.Year()
	if at.Month() < birth.Month() || (at.Month() == birth.Month() && at.Day() < birth.Day()) {
		age--
	}
	return age
}
