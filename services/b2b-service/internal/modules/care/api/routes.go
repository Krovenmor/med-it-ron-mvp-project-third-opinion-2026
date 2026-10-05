package api

import (
	"context"

	"github.com/google/uuid"
)

type Routes interface {
	ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]Service, error)
}

type Service struct {
	Code string
	Name string
}
