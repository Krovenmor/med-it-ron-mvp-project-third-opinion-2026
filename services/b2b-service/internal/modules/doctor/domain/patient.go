package domain

import (
	"time"

	"github.com/google/uuid"
)

type Sex string

type Patient struct {
	ID           uuid.UUID
	SourceSystem string
	ExternalID   string
	BirthDate    time.Time
	Sex          Sex
}
