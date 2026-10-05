package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed */*.sql
var files embed.FS

type Queries struct {
	Schema          Schema
	Cases           Cases
	Recommendations Recommendations
	CaseEvents      CaseEvents
	Demo            Demo
	History         History
}

type Schema struct {
	Create string
}

type Cases struct {
	InsertIfAbsent  string
	Get             string
	GetForUpdate    string
	Update          string
	ListReviewQueue string
}

type Recommendations struct {
	Insert            string
	Append            string
	Get               string
	UpdateReview      string
	UpdatePatientText string
	ListByCase        string
}

type CaseEvents struct {
	Insert string
	Exists string
}

type Demo struct {
	Reset string
}

type History struct {
	InsertCase           string
	InsertRecommendation string
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
		Schema: Schema{Create: read("schema/create.sql")},
		Cases: Cases{
			InsertIfAbsent:  read("cases/insert_if_absent.sql"),
			Get:             read("cases/get.sql"),
			GetForUpdate:    read("cases/get_for_update.sql"),
			Update:          read("cases/update.sql"),
			ListReviewQueue: read("cases/list_review_queue.sql"),
		},
		Recommendations: Recommendations{
			Insert:            read("recommendations/insert.sql"),
			Append:            read("recommendations/append.sql"),
			Get:               read("recommendations/get.sql"),
			UpdateReview:      read("recommendations/update_review.sql"),
			UpdatePatientText: read("recommendations/update_patient_text.sql"),
			ListByCase:        read("recommendations/list_by_case.sql"),
		},
		CaseEvents: CaseEvents{
			Insert: read("case_events/insert.sql"),
			Exists: read("case_events/exists.sql"),
		},
		Demo: Demo{Reset: read("demo/reset.sql")},
		History: History{
			InsertCase:           read("history/insert_case.sql"),
			InsertRecommendation: read("history/insert_recommendation.sql"),
		},
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load doctor queries: %w", err)
	}
	return q, nil
}
