package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobKind string

const JobKindAssess JobKind = "assess"

type Job struct {
	ID       int64
	CaseID   uuid.UUID
	Kind     JobKind
	RunAt    time.Time
	Attempts int
}

func NewJob(caseID uuid.UUID, kind JobKind, runAt time.Time) Job {
	return Job{CaseID: caseID, Kind: kind, RunAt: runAt}
}
