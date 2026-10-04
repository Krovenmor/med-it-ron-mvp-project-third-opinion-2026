//go:build e2e

package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed *.sql
var files embed.FS

type Queries struct {
	CountCasesByStudy       string
	CountJobsByStudy        string
	CountDraftEventsByStudy string
	GetPatientPhone         string
	GetJobByCase            string
	CreateMockModeDatabase  string
	CreateDemoResetDatabase string
	ListEventsByCase        string
}

func Load() (Queries, error) {
	var errs []error
	read := func(name string) string {
		data, err := files.ReadFile(name)
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		return string(data)
	}

	q := Queries{
		CountCasesByStudy:       read("count_cases_by_study.sql"),
		CountJobsByStudy:        read("count_jobs_by_study.sql"),
		CountDraftEventsByStudy: read("count_draft_events_by_study.sql"),
		GetPatientPhone:         read("get_patient_phone.sql"),
		GetJobByCase:            read("get_job_by_case.sql"),
		CreateMockModeDatabase:  read("create_mock_mode_database.sql"),
		CreateDemoResetDatabase: read("create_demo_reset_database.sql"),
		ListEventsByCase:        read("list_events_by_case.sql"),
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load test queries: %w", err)
	}
	return q, nil
}
