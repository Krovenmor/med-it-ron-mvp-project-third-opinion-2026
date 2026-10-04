package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecommendationSource string

const (
	RecommendationSourceAI     RecommendationSource = "ai"
	RecommendationSourceDoctor RecommendationSource = "doctor"
)

type Importance string

const (
	ImportanceHigh Importance = "high"
	ImportanceLow  Importance = "low"
)

func (i Importance) Valid() bool {
	return i == ImportanceHigh || i == ImportanceLow
}

type Recommendation struct {
	ID            uuid.UUID
	CaseID        uuid.UUID
	Position      int
	Source        RecommendationSource
	ServiceCode   string
	ServiceName   string
	Importance    Importance
	Rationale     string
	GuidelineRef  string
	PatientText   string
	AlreadyBooked bool
	CreatedAt     time.Time
}

type Service struct {
	Code string
	Name string
}
