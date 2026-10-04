package intake

import (
	"bytes"
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Service struct {
	tx       trm.Manager
	clock    Clock
	patients Patients
	cases    Cases
	events   Events
	jobs     Jobs
}

func NewService(tx trm.Manager, clock Clock, patients Patients, cases Cases, events Events, jobs Jobs) *Service {
	return &Service{tx: tx, clock: clock, patients: patients, cases: cases, events: events, jobs: jobs}
}

func (s *Service) Ingest(ctx context.Context, report domain.Report) (IngestResult, error) {
	if err := report.Validate(); err != nil {
		return IngestResult{}, err
	}

	now := s.clock.Now()
	var result IngestResult
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		patientID, err := s.patients.Upsert(ctx, report.Patient, now)
		if err != nil {
			return err
		}

		c := domain.NewCase(patientID, report, now)
		id, created, err := s.cases.CreateIfAbsent(ctx, c)
		if err != nil {
			return err
		}
		if !created {
			existing, err := s.cases.GetByStudy(ctx, report.SourceSystem, report.Study.ID)
			if err != nil {
				return err
			}
			if !bytes.Equal(existing.Fingerprint, c.Fingerprint) {
				return domain.ErrStudyConflict
			}
			result = IngestResult{Case: existing}
			return nil
		}

		c.ID = id
		if err := s.events.Append(ctx, domain.StatusChanged(c.ID, "", c.Status, domain.ActorSystem, now)); err != nil {
			return err
		}
		if err := s.jobs.Enqueue(ctx, domain.NewJob(c.ID, domain.JobKindAssess, now), now); err != nil {
			return err
		}
		result = IngestResult{Case: c, Created: true}
		return nil
	})
	return result, err
}
