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
	t.Parallel()
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
