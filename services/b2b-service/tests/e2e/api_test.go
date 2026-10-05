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

func TestCases_UnknownCaseIsNotFound(t *testing.T) {
	t.Parallel()
	for _, id := range []string{uuid.NewString(), "not-a-uuid"} {
		status, _ := api.getCase(t, id)
		assert.Equal(t, http.StatusNotFound, status, id)
	}
}

func TestDemoClock_AdvancesBusinessTime(t *testing.T) {
	advance := func(by string) time.Time {
		status, body := api.do(t, http.MethodPost, "/demo/clock/advance", []byte(`{"by":"`+by+`"}`))
		require.Equal(t, http.StatusOK, status)
		return decode[clockView](t, body).Now
	}

	first := advance("1s")
	second := advance("30m")

	elapsed := second.Sub(first)
	assert.GreaterOrEqual(t, elapsed, 30*time.Minute)
	assert.Less(t, elapsed, 31*time.Minute)

	for _, by := range []string{`"-5m"`, `"0s"`, `"tomorrow"`} {
		status, _ := api.do(t, http.MethodPost, "/demo/clock/advance", []byte(`{"by":`+by+`}`))
		assert.Equal(t, http.StatusBadRequest, status, by)
	}
}

func TestCases_HistoryComesFromMIS(t *testing.T) {
	t.Parallel()
	r := newReport()
	h := history{
		Visits:       []visit{{ServiceCode: "THER-CONSULT", ServiceName: "Приём терапевта", VisitedAt: time.Date(2026, 8, 20, 9, 0, 0, 0, time.UTC)}},
		Appointments: []appointment{{ServiceCode: "PULM-CONSULT", ServiceName: "Консультация пульмонолога", ScheduledAt: time.Date(2026, 10, 20, 6, 0, 0, 0, time.UTC)}},
	}
	mis.setHistory(r.Patient.ID, h)
	c := ingestInReview(t, r, validAssessment())

	status, got := api.caseHistory(t, c.ID)
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, h, got)

	mis.failHistory(r.Patient.ID, http.StatusInternalServerError)
	status, _ = api.caseHistory(t, c.ID)
	assert.Equal(t, http.StatusBadGateway, status)

	status, _ = api.caseHistory(t, uuid.NewString())
	assert.Equal(t, http.StatusNotFound, status)
}

func TestClock_ReturnsBusinessTime(t *testing.T) {
	before := api.clock(t)
	status, _ := api.do(t, http.MethodPost, "/demo/clock/advance", []byte(`{"by":"1h"}`))
	require.Equal(t, http.StatusOK, status)

	elapsed := api.clock(t).Sub(before)
	assert.GreaterOrEqual(t, elapsed, time.Hour)
	assert.Less(t, elapsed, time.Hour+time.Minute)
}
