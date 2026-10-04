//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssessment_TransientFailuresAreRetried(t *testing.T) {
	t.Parallel()
	r := newReport()
	ai.on(r.Conclusion, failFirst(maxAttempts-1, http.StatusServiceUnavailable, validAssessment()))

	_, ref := api.ingest(t, r)
	c := api.awaitStatus(t, ref.CaseID, "in_review")

	assert.NotEmpty(t, c.Recommendations)
	assert.Len(t, ai.requestsFor(r.Conclusion), maxAttempts)
	assert.Equal(t, jobState{State: "done", Attempts: maxAttempts}, jobOf(t, ref.CaseID))
}

func TestAssessment_GivesUpAfterMaxAttempts(t *testing.T) {
	t.Parallel()
	invalid := func(mutate func(a *assessment)) aiBehavior {
		a := validAssessment()
		mutate(&a)
		return respond(a)
	}
	tests := []struct {
		name     string
		behavior aiBehavior
	}{
		{"ai-service error", failWith(http.StatusInternalServerError)},
		{"unknown urgency", invalid(func(a *assessment) { a.Urgency = "whenever" })},
		{"missing catalog version", invalid(func(a *assessment) { a.CatalogVersion = "" })},
		{"recommendation without patient text", invalid(func(a *assessment) { a.Recommendations[0].PatientText = "" })},
		{"unknown importance", invalid(func(a *assessment) { a.Recommendations[0].Importance = "critical" })},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newReport()
			ai.on(r.Conclusion, tt.behavior)

			_, ref := api.ingest(t, r)

			var job jobState
			failed := poll(func() bool {
				job = jobOf(t, ref.CaseID)
				return job.State == "failed"
			})
			require.True(t, failed, "job state %+v", job)
			assert.Equal(t, maxAttempts, job.Attempts)
			assert.Len(t, ai.requestsFor(r.Conclusion), maxAttempts)

			status, c := api.getCase(t, ref.CaseID)
			require.Equal(t, http.StatusOK, status)
			assert.Equal(t, "draft", c.Status)
			assert.Empty(t, c.Urgency)
			assert.Empty(t, c.Recommendations)
		})
	}
}
