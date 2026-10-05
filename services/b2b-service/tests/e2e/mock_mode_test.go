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

func TestBuiltInAIMockService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	_, err := db.Exec(ctx, query.CreateMockModeDatabase)
	require.NoError(t, err)

	svc, err := env.startService(ctx, map[string]string{
		"POSTGRES_DSN":    env.dsn("mock_mode"),
		"AI_SERVICE_MOCK": "true",
		"DEMO_MODE":       "false",
	})
	if svc != nil && svc.container != nil {
		t.Cleanup(func() { assert.NoError(t, testcontainers.TerminateContainer(svc.container)) })
	}
	require.NoError(t, err)

	t.Run("assesses report with built-in mock", func(t *testing.T) {
		r := newReport()
		r.Study.Modality = "DX"
		r.Conclusion = "Справа пневмоторакс, " + r.Study.ID

		status, ref := svc.api.ingest(t, r)
		require.Equal(t, http.StatusCreated, status)

		c := svc.api.awaitStatus(t, ref.CaseID, "in_review")
		assert.Equal(t, "emergency", c.Urgency)
		assert.NotEmpty(t, c.Recommendations)
		assert.Empty(t, ai.requestsFor(r.Conclusion))
	})

	t.Run("demo routes are disabled", func(t *testing.T) {
		status, _ := svc.api.do(t, http.MethodPost, "/demo/clock/advance", []byte(`{"by":"30m"}`))
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("stops gracefully on SIGTERM", func(t *testing.T) {
		timeout := 10 * time.Second
		require.NoError(t, svc.container.Stop(ctx, &timeout))

		state, err := svc.container.State(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, state.ExitCode)
	})
}
