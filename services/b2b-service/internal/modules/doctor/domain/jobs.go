package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

const JobKindEscalateReview = "escalate_review"

func ReviewEscalationJob(caseID uuid.UUID, from time.Time) jobs.Job {
	return jobs.New(JobKindEscalateReview, caseID.String(), from.Add(ReviewEscalationDelay))
}
