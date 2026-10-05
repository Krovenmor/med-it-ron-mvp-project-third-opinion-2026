package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID           uuid.UUID
	SourceSystem string
	ExternalID   string
	FullName     string
	BirthDate    time.Time
	Sex          string
	Phone        string
	Email        string
}
