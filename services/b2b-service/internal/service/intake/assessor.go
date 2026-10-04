package intake

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Assessor struct {
	tx              trm.Manager
	clock           Clock
	patients        Patients
	cases           Cases
	recommendations Recommendations
	events          Events
	ai              AIService
	mis             MIS
}

func NewAssessor(
	tx trm.Manager,
	clock Clock,
	patients Patients,
	cases Cases,
	recommendations Recommendations,
	events Events,
	ai AIService,
	mis MIS,
) *Assessor {
	return &Assessor{
		tx:              tx,
		clock:           clock,
		patients:        patients,
		cases:           cases,
		recommendations: recommendations,
		events:          events,
		ai:              ai,
		mis:             mis,
	}
}

func (a *Assessor) Kind() domain.JobKind {
	return domain.JobKindAssess
}

func (a *Assessor) Handle(ctx context.Context, job domain.Job) error {
	c, err := a.cases.Get(ctx, job.CaseID)
	if err != nil {
		return err
	}
	if c.Status != domain.CaseStatusDraft {
		return nil
	}

	req, err := a.buildRequest(ctx, c)
	if err != nil {
		return err
	}
	assessment, err := a.ai.Assess(ctx, req)
	if err != nil {
		return fmt.Errorf("assess case: %w", err)
	}
	if err := assessment.Validate(); err != nil {
		return fmt.Errorf("invalid assessment: %w", err)
	}

	return a.tx.Do(ctx, func(ctx context.Context) error {
		c, err := a.cases.GetForUpdate(ctx, job.CaseID)
		if err != nil {
			return err
		}
		if c.Status != domain.CaseStatusDraft {
			return nil
		}

		now := a.clock.Now()
		if err := a.recommendations.InsertMany(ctx, assessment.RecommendationsFor(c.ID, now)); err != nil {
			return err
		}
		from := c.Status
		c.ApplyAssessment(assessment, now)
		if err := a.cases.Update(ctx, c); err != nil {
			return err
		}
		return a.events.Append(ctx, domain.StatusChanged(c.ID, from, c.Status, domain.ActorSystem, now))
	})
}

func (a *Assessor) buildRequest(ctx context.Context, c domain.Case) (domain.AssessmentRequest, error) {
	patient, err := a.patients.Get(ctx, c.PatientID)
	if err != nil {
		return domain.AssessmentRequest{}, err
	}
	history, err := a.mis.PatientHistory(ctx, patient)
	if err != nil {
		return domain.AssessmentRequest{}, fmt.Errorf("fetch patient history: %w", err)
	}
	current, err := a.recommendations.ActiveServices(ctx, c.PatientID, c.ID)
	if err != nil {
		return domain.AssessmentRequest{}, err
	}

	return domain.AssessmentRequest{
		CaseID:                 c.ID,
		Modality:               c.Study.Modality,
		BodySite:               c.Study.BodySite,
		PerformedAt:            c.Study.PerformedAt,
		Conclusion:             c.Conclusion,
		PatientAge:             patient.AgeAt(a.clock.Now()),
		PatientSex:             patient.Sex,
		CurrentRecommendations: current,
		History:                history,
	}, nil
}
