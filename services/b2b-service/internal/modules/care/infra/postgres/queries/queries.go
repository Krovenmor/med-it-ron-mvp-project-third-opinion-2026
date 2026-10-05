package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed */*.sql
var files embed.FS

type Queries struct {
	Schema        Schema
	Patients      Patients
	Routes        Routes
	PlanItems     PlanItems
	Bookings      Bookings
	OperatorTasks OperatorTasks
	CallAttempts  CallAttempts
	Notifications Notifications
	CaseEvents    CaseEvents
	Dashboard     Dashboard
	Demo          Demo
	History       History
}

type Schema struct {
	Create string
}

type Patients struct {
	Upsert          string
	Get             string
	GetByExternalID string
}

type Routes struct {
	InsertIfAbsent string
	Get            string
	GetForUpdate   string
	Update         string
	ListByPatient  string
}

type PlanItems struct {
	Insert             string
	Get                string
	GetForUpdate       string
	UpdateDecline      string
	ListByCase         string
	ListActiveServices string
}

type Bookings struct {
	InsertIfAbsent              string
	GetByAppointment            string
	ListRecommendationIDsByCase string
}

type OperatorTasks struct {
	InsertIfNoActive         string
	Get                      string
	GetForUpdate             string
	GetActiveByCaseForUpdate string
	Update                   string
	ListActive               string
}

type CallAttempts struct {
	Insert     string
	ListByTask string
}

type Notifications struct {
	Insert string
	List   string
}

type CaseEvents struct {
	Insert string
}

type Dashboard struct {
	CaseFacts      string
	TaskFacts      string
	DeclineReasons string
}

type Demo struct {
	Reset string
}

type History struct {
	InsertTask    string
	InsertBooking string
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
		Patients: Patients{
			Upsert:          read("patients/upsert.sql"),
			Get:             read("patients/get.sql"),
			GetByExternalID: read("patients/get_by_external_id.sql"),
		},
		Routes: Routes{
			InsertIfAbsent: read("routes/insert_if_absent.sql"),
			Get:            read("routes/get.sql"),
			GetForUpdate:   read("routes/get_for_update.sql"),
			Update:         read("routes/update.sql"),
			ListByPatient:  read("routes/list_by_patient.sql"),
		},
		PlanItems: PlanItems{
			Insert:             read("plan_items/insert.sql"),
			Get:                read("plan_items/get.sql"),
			GetForUpdate:       read("plan_items/get_for_update.sql"),
			UpdateDecline:      read("plan_items/update_decline.sql"),
			ListByCase:         read("plan_items/list_by_case.sql"),
			ListActiveServices: read("plan_items/list_active_services.sql"),
		},
		Bookings: Bookings{
			InsertIfAbsent:              read("bookings/insert_if_absent.sql"),
			GetByAppointment:            read("bookings/get_by_appointment.sql"),
			ListRecommendationIDsByCase: read("bookings/list_recommendation_ids_by_case.sql"),
		},
		OperatorTasks: OperatorTasks{
			InsertIfNoActive:         read("operator_tasks/insert_if_no_active.sql"),
			Get:                      read("operator_tasks/get.sql"),
			GetForUpdate:             read("operator_tasks/get_for_update.sql"),
			GetActiveByCaseForUpdate: read("operator_tasks/get_active_by_case_for_update.sql"),
			Update:                   read("operator_tasks/update.sql"),
			ListActive:               read("operator_tasks/list_active.sql"),
		},
		CallAttempts: CallAttempts{
			Insert:     read("operator_task_attempts/insert.sql"),
			ListByTask: read("operator_task_attempts/list_by_task.sql"),
		},
		Notifications: Notifications{
			Insert: read("notifications/insert.sql"),
			List:   read("notifications/list.sql"),
		},
		CaseEvents: CaseEvents{Insert: read("case_events/insert.sql")},
		Dashboard: Dashboard{
			CaseFacts:      read("dashboard/case_facts.sql"),
			TaskFacts:      read("dashboard/task_facts.sql"),
			DeclineReasons: read("dashboard/decline_reasons.sql"),
		},
		Demo: Demo{Reset: read("demo/reset.sql")},
		History: History{
			InsertTask:    read("history/insert_task.sql"),
			InsertBooking: read("history/insert_booking.sql"),
		},
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load care queries: %w", err)
	}
	return q, nil
}
