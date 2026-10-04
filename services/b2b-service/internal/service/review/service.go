package review

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Service struct {
	tx              trm.Manager
	clock           Clock
	cases           Cases
	recommendations Recommendations
	events          Events
}

func NewService(tx trm.Manager, clock Clock, cases Cases, recommendations Recommendations, events Events) *Service {
	return &Service{tx: tx, clock: clock, cases: cases, recommendations: recommendations, events: events}
}

func (s *Service) Queue(ctx context.Context) ([]domain.ReviewQueueItem, error) {
	return s.cases.ReviewQueue(ctx)
}

func (s *Service) Open(ctx context.Context, caseID uuid.UUID, actor string) error {
	c, err := s.cases.Get(ctx, caseID)
	if err != nil {
		return err
	}
	if err := c.EnsureInReview(); err != nil {
		return err
	}
	return s.events.Append(ctx, domain.CaseOpened(c.ID, actor, s.clock.Now()))
}

func (s *Service) ReviewRecommendation(ctx context.Context, cmd ReviewRecommendation) (domain.Recommendation, error) {
	now := s.clock.Now()
	var review *domain.Review
	if cmd.Mark != "" {
		review = &domain.Review{
			Mark:          cmd.Mark,
			RejectReason:  cmd.RejectReason,
			RejectComment: cmd.RejectComment,
			ReviewedBy:    cmd.Actor,
			ReviewedAt:    now,
		}
	}

	var rec domain.Recommendation
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		if err := s.lockInReview(ctx, cmd.CaseID); err != nil {
			return err
		}
		var err error
		rec, err = s.recommendations.Get(ctx, cmd.CaseID, cmd.RecommendationID)
		if err != nil {
			return err
		}
		if err := rec.Edit(review, cmd.PatientText); err != nil {
			return err
		}
		if review == nil {
			if err := s.recommendations.UpdatePatientText(ctx, rec); err != nil {
				return err
			}
			return s.events.Append(ctx, domain.RecommendationTextEdited(rec, cmd.Actor, now))
		}
		if err := s.recommendations.UpdateReview(ctx, rec); err != nil {
			return err
		}
		return s.events.Append(ctx, domain.RecommendationReviewed(rec, cmd.PatientText != nil))
	})
	return rec, err
}

func (s *Service) AddRecommendation(ctx context.Context, cmd AddRecommendation) (domain.Recommendation, error) {
	rec, err := domain.NewDoctorRecommendation(cmd.CaseID, cmd.Service, cmd.PatientText, cmd.Rationale, domain.Review{
		Mark:       cmd.Mark,
		ReviewedBy: cmd.Actor,
		ReviewedAt: s.clock.Now(),
	})
	if err != nil {
		return domain.Recommendation{}, err
	}

	err = s.tx.Do(ctx, func(ctx context.Context) error {
		if err := s.lockInReview(ctx, cmd.CaseID); err != nil {
			return err
		}
		var err error
		rec, err = s.recommendations.Append(ctx, rec)
		if err != nil {
			return err
		}
		return s.events.Append(ctx, domain.RecommendationAdded(rec))
	})
	return rec, err
}

func (s *Service) ChangeUrgency(ctx context.Context, cmd ChangeUrgency) error {
	return s.tx.Do(ctx, func(ctx context.Context) error {
		c, err := s.cases.GetForUpdate(ctx, cmd.CaseID)
		if err != nil {
			return err
		}
		from := c.Urgency
		now := s.clock.Now()
		if err := c.ChangeUrgency(cmd.Urgency, cmd.Reason, now); err != nil {
			return err
		}
		if err := s.cases.Update(ctx, c); err != nil {
			return err
		}
		return s.events.Append(ctx, domain.UrgencyChanged(c.ID, from, c.Urgency, cmd.Reason, cmd.Actor, now))
	})
}

func (s *Service) Confirm(ctx context.Context, caseID uuid.UUID, actor string) (domain.Case, error) {
	var c domain.Case
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		var err error
		c, err = s.cases.GetForUpdate(ctx, caseID)
		if err != nil {
			return err
		}
		recs, err := s.recommendations.ListByCase(ctx, caseID)
		if err != nil {
			return err
		}
		from := c.Status
		now := s.clock.Now()
		if err := c.Confirm(recs, now); err != nil {
			return err
		}
		if err := s.cases.Update(ctx, c); err != nil {
			return err
		}
		return s.events.Append(ctx, domain.StatusChanged(c.ID, from, c.Status, actor, now))
	})
	return c, err
}

func (s *Service) lockInReview(ctx context.Context, caseID uuid.UUID) error {
	c, err := s.cases.GetForUpdate(ctx, caseID)
	if err != nil {
		return err
	}
	return c.EnsureInReview()
}
