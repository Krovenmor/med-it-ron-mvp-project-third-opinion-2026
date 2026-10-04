package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReviewQueueItem struct {
	CaseID                  uuid.UUID
	Urgency                 Urgency
	Modality                Modality
	PerformedAt             time.Time
	ReceivedAt              time.Time
	Patient                 Patient
	RecommendationsTotal    int
	RecommendationsReviewed int
	OpenedAt                time.Time
}
