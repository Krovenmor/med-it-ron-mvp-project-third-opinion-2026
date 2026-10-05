//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestDemoReset_ClearsDataAndClock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	_, err := db.Exec(ctx, query.CreateDemoResetDatabase)
	require.NoError(t, err)

	svc, err := env.startService(ctx, map[string]string{
		"POSTGRES_DSN":    env.dsn("demo_reset"),
		"AI_SERVICE_MOCK": "true",
		"DEMO_MODE":       "true",
	})
	if svc != nil && svc.container != nil {
		t.Cleanup(func() { assert.NoError(t, testcontainers.TerminateContainer(svc.container)) })
	}
	require.NoError(t, err)

	r := newReport()
	_, ref := svc.api.ingest(t, r)
	svc.api.awaitStatus(t, ref.CaseID, "in_review")
	status, _ := svc.api.do(t, http.MethodPost, "/demo/clock/advance", []byte(`{"by":"24h"}`))
	require.Equal(t, http.StatusOK, status)

	status, _ = svc.api.do(t, http.MethodPost, "/demo/reset", nil)
	require.Equal(t, http.StatusNoContent, status)

	status, _ = svc.api.getCase(t, ref.CaseID)
	assert.Equal(t, http.StatusNotFound, status)
	assert.Empty(t, svc.api.reviewQueue(t).Cases)
	assert.Empty(t, svc.api.operatorTasks(t, "", false).Tasks)
	assert.Empty(t, svc.api.notifications(t))
	assert.WithinDuration(t, time.Now(), svc.api.clock(t), time.Minute)

	_, history := svc.api.dashboard(t, "90")
	for _, stats := range []periodStatsView{history.Current, history.Previous} {
		f := stats.Funnel
		assert.Greater(t, f.Received, 500)
		assert.True(t, f.Received > f.Recommended && f.Recommended >= f.Confirmed && f.Confirmed >= f.Booked && f.Booked > f.Completed && f.Completed > 0, "%+v", f)
		assert.Equal(t, f.Booked, stats.Bookings.Self+stats.Bookings.Operator)
		assert.Greater(t, stats.DoctorSLA.Total, 0)
		assert.Greater(t, stats.OperatorSLA.Total, 0)
	}
	assert.NotEmpty(t, history.DeclineReasons)

	status, igor, _ := svc.api.plan(t, "mis-demo", "P-IGOR")
	require.Equal(t, http.StatusOK, status, "demo patient Igor has a previous case")
	require.Len(t, igor.Cases, 1)
	previous := igor.Cases[0]
	assert.Equal(t, "booked", previous.Status)
	require.NotNil(t, previous.BookBy)
	assert.True(t, previous.BookBy.Before(svc.api.clock(t)), "the follow-up deadline has passed")
	assert.Len(t, previous.Recommendations, 2)

	status, _ = svc.api.do(t, http.MethodPost, "/demo/reset", nil)
	require.Equal(t, http.StatusNoContent, status)
	_, reseeded := svc.api.dashboard(t, "90")
	assert.Equal(t, history.Current, reseeded.Current, "history is generated from a fixed seed")

	status, again := svc.api.ingest(t, r)
	require.Equal(t, http.StatusCreated, status)
	svc.api.awaitStatus(t, again.CaseID, "in_review")
}
