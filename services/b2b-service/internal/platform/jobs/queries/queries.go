package queries

import (
	"embed"
	"errors"
	"fmt"
)

//go:embed *.sql
var files embed.FS

type Queries struct {
	Enqueue       string
	Claim         string
	Complete      string
	Reschedule    string
	Fail          string
	CancelPending string
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
		Enqueue:       read("enqueue.sql"),
		Claim:         read("claim.sql"),
		Complete:      read("complete.sql"),
		Reschedule:    read("reschedule.sql"),
		Fail:          read("fail.sql"),
		CancelPending: read("cancel_pending.sql"),
	}
	if err := errors.Join(errs...); err != nil {
		return Queries{}, fmt.Errorf("load job queries: %w", err)
	}
	return q, nil
}
