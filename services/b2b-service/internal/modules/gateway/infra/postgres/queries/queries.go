package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed */*.sql
var files embed.FS

type Queries struct {
	Schema   Schema
	Patients Patients
	Intakes  Intakes
	Demo     Demo
	History  History
}

type Schema struct {
	Create string
}

type Patients struct {
	Upsert string
	Get    string
}

type Intakes struct {
	InsertIfAbsent string
	Get            string
	GetForUpdate   string
	GetByStudy     string
	Update         string
}

type Demo struct {
	Reset string
}

type History struct {
	InsertPatient string
	InsertIntake  string
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
			Upsert: read("patients/upsert.sql"),
			Get:    read("patients/get.sql"),
		},
		Intakes: Intakes{
			InsertIfAbsent: read("intakes/insert_if_absent.sql"),
			Get:            read("intakes/get.sql"),
			GetForUpdate:   read("intakes/get_for_update.sql"),
			GetByStudy:     read("intakes/get_by_study.sql"),
			Update:         read("intakes/update.sql"),
		},
		Demo: Demo{Reset: read("demo/reset.sql")},
		History: History{
			InsertPatient: read("history/insert_patient.sql"),
			InsertIntake:  read("history/insert_intake.sql"),
		},
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load gateway queries: %w", err)
	}
	return q, nil
}
