package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed */*.sql
var files embed.FS

type Queries struct {
	Patients        Patients
	Cases           Cases
	Recommendations Recommendations
	CaseEvents      CaseEvents
	Jobs            Jobs
}

type Patients struct {
	Upsert string
	Get    string
}

type Cases struct {
	InsertIfAbsent string
	Get            string
	GetForUpdate   string
	GetByStudy     string
	Update         string
}

type Recommendations struct {
	Insert             string
	ListByCase         string
	ListActiveServices string
}

type CaseEvents struct {
	Insert string
}

type Jobs struct {
	Enqueue    string
	Claim      string
	Complete   string
	Reschedule string
	Fail       string
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
		Patients: Patients{
			Upsert: read("patients/upsert.sql"),
			Get:    read("patients/get.sql"),
		},
		Cases: Cases{
			InsertIfAbsent: read("cases/insert_if_absent.sql"),
			Get:            read("cases/get.sql"),
			GetForUpdate:   read("cases/get_for_update.sql"),
			GetByStudy:     read("cases/get_by_study.sql"),
			Update:         read("cases/update.sql"),
		},
		Recommendations: Recommendations{
			Insert:             read("recommendations/insert.sql"),
			ListByCase:         read("recommendations/list_by_case.sql"),
			ListActiveServices: read("recommendations/list_active_services.sql"),
		},
		CaseEvents: CaseEvents{
			Insert: read("case_events/insert.sql"),
		},
		Jobs: Jobs{
			Enqueue:    read("jobs/enqueue.sql"),
			Claim:      read("jobs/claim.sql"),
			Complete:   read("jobs/complete.sql"),
			Reschedule: read("jobs/reschedule.sql"),
			Fail:       read("jobs/fail.sql"),
		},
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load queries: %w", err)
	}
	return q, nil
}
