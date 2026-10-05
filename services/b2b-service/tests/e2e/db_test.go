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

func doctorEventsOf(t *testing.T, caseID string) []caseEventView {
	t.Helper()
	return eventsFrom(t, query.ListDoctorEventsByCase, caseID)
}

func careEventsOf(t *testing.T, caseID string) []caseEventView {
	t.Helper()
	return eventsFrom(t, query.ListCareEventsByCase, caseID)
}

func eventsFrom(t *testing.T, q, caseID string) []caseEventView {
	t.Helper()
	rows, err := db.Query(context.Background(), q, pgx.StrictNamedArgs{"case_id": caseID})
	require.NoError(t, err)
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[caseEventView])
	require.NoError(t, err)
	return events
}

func patientContactsSynced(t *testing.T, externalID, phone string) bool {
	t.Helper()
	var n int
	err := db.QueryRow(context.Background(), query.CountPatientContacts, pgx.StrictNamedArgs{"external_id": externalID, "phone": phone}).Scan(&n)
	require.NoError(t, err)
	return n > 0
}

func jobOf(t *testing.T, caseID string) jobState {
	t.Helper()
	rows, err := db.Query(context.Background(), query.GetJobByCase, pgx.StrictNamedArgs{"case_id": caseID})
	require.NoError(t, err)
	job, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[jobState])
	require.NoError(t, err)
	return job
}

func jobsOf(t *testing.T, caseID string) []jobView {
	t.Helper()
	rows, err := db.Query(context.Background(), query.ListJobsByCase, pgx.StrictNamedArgs{"case_id": caseID})
	require.NoError(t, err)
	jobs, err := pgx.CollectRows(rows, pgx.RowToStructByName[jobView])
	require.NoError(t, err)
	return jobs
}
