//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func countByStudy(t *testing.T, q, studyID string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRow(context.Background(), q, pgx.StrictNamedArgs{"study_id": studyID}).Scan(&n))
	return n
}

func patientPhone(t *testing.T, externalID string) string {
	t.Helper()
	var phone string
	err := db.QueryRow(context.Background(), query.GetPatientPhone, pgx.StrictNamedArgs{"external_id": externalID}).Scan(&phone)
	require.NoError(t, err)
	return phone
}

func jobOf(t *testing.T, caseID string) jobState {
	t.Helper()
	rows, err := db.Query(context.Background(), query.GetJobByCase, pgx.StrictNamedArgs{"case_id": caseID})
	require.NoError(t, err)
	job, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[jobState])
	require.NoError(t, err)
	return job
}
