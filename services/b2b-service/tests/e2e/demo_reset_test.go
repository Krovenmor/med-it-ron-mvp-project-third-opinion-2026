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
	assert.WithinDuration(t, time.Now(), svc.api.clock(t), time.Minute)

	status, again := svc.api.ingest(t, r)
	require.Equal(t, http.StatusCreated, status)
	svc.api.awaitStatus(t, again.CaseID, "in_review")
}
