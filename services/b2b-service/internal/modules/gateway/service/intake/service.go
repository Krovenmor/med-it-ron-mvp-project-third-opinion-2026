package intake

import (
	"bytes"
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Service struct {
	tx        trm.Manager
	clock     Clock
	patients  Patients
	intakes   Intakes
	jobs      Jobs
	publisher Publisher
}

func NewService(tx trm.Manager, clock Clock, patients Patients, intakes Intakes, jobs Jobs, publisher Publisher) *Service {
	return &Service{tx: tx, clock: clock, patients: patients, intakes: intakes, jobs: jobs, publisher: publisher}
}

func (s *Service) Ingest(ctx context.Context, report domain.Report) (IngestResult, error) {
	if err := report.Validate(); err != nil {
		return IngestResult{}, err
	}

	now := s.clock.Now()
	var result IngestResult
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		patient := report.Patient
		patient.SourceSystem = report.SourceSystem
		id, err := s.patients.Upsert(ctx, patient, now)
		if err != nil {
			return err
		}
		patient.ID = id
		if err := s.publisher.Publish(ctx, api.TopicPatientRegistered, id.String(), patientRegistered(patient)); err != nil {
			return err
		}

		intake := domain.NewIntake(patient.ID, report, now)
		caseID, created, err := s.intakes.CreateIfAbsent(ctx, intake)
		if err != nil {
			return err
		}
		if !created {
			existing, err := s.intakes.GetByStudy(ctx, report.SourceSystem, report.Study.ID)
			if err != nil {
				return err
			}
			if !bytes.Equal(existing.Fingerprint, intake.Fingerprint) {
				return domain.ErrStudyConflict
			}
			result = IngestResult{Intake: existing}
			return nil
		}

		intake.CaseID = caseID
		if err := s.jobs.Enqueue(ctx, jobs.New(KindAssess, caseID.String(), now), now); err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, api.TopicReportReceived, caseID.String(), reportReceived(intake)); err != nil {
			return err
		}
		result = IngestResult{Intake: intake, Created: true}
		return nil
	})
	return result, err
}

func (s *Service) Status(ctx context.Context, caseID uuid.UUID) (domain.Intake, error) {
	return s.intakes.Get(ctx, caseID)
}
