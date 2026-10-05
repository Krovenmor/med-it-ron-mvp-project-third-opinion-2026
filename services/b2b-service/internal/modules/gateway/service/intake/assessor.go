package intake

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Assessor struct {
	tx        trm.Manager
	clock     Clock
	patients  Patients
	intakes   Intakes
	publisher Publisher
	ai        AIService
	mis       MIS
	routes    Routes
}

func NewAssessor(
	tx trm.Manager,
	clock Clock,
	patients Patients,
	intakes Intakes,
	publisher Publisher,
	ai AIService,
	mis MIS,
	routes Routes,
) *Assessor {
	return &Assessor{
		tx:        tx,
		clock:     clock,
		patients:  patients,
		intakes:   intakes,
		publisher: publisher,
		ai:        ai,
		mis:       mis,
		routes:    routes,
	}
}

func (a *Assessor) Kind() string {
	return KindAssess
}

func (a *Assessor) Handle(ctx context.Context, job jobs.Job) error {
	caseID, err := uuid.Parse(job.Key)
	if err != nil {
		return fmt.Errorf("assess job key: %w", err)
	}
	intake, err := a.intakes.Get(ctx, caseID)
	if err != nil {
		return err
	}
	if intake.Status != domain.IntakeStatusReceived {
		return nil
	}
	patient, err := a.patients.Get(ctx, intake.PatientID)
	if err != nil {
		return err
	}

	req, err := a.buildRequest(ctx, intake, patient)
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
		intake, err := a.intakes.GetForUpdate(ctx, caseID)
		if err != nil {
			return err
		}
		if !intake.MarkAssessed(a.clock.Now()) {
			return nil
		}
		if err := a.intakes.Update(ctx, intake); err != nil {
			return err
		}
		return a.publisher.Publish(ctx, api.TopicCaseAssessed, caseID.String(), caseAssessed(intake, patient, assessment))
	})
}

func (a *Assessor) buildRequest(ctx context.Context, intake domain.Intake, patient domain.Patient) (domain.AssessmentRequest, error) {
	history, err := a.mis.PatientHistory(ctx, patient.ExternalID)
	if err != nil {
		return domain.AssessmentRequest{}, fmt.Errorf("fetch patient history: %w", err)
	}
	active, err := a.routes.ActiveServices(ctx, patient.ID, intake.CaseID)
	if err != nil {
		return domain.AssessmentRequest{}, err
	}
	current := make([]domain.Service, 0, len(active))
	for _, s := range active {
		current = append(current, domain.Service{Code: s.Code, Name: s.Name})
	}

	return domain.AssessmentRequest{
		CaseID:                 intake.CaseID,
		Modality:               intake.Study.Modality,
		BodySite:               intake.Study.BodySite,
		PerformedAt:            intake.Study.PerformedAt,
		Conclusion:             intake.Conclusion,
		PatientAge:             patient.AgeAt(a.clock.Now()),
		PatientSex:             patient.Sex,
		CurrentRecommendations: current,
		History:                history,
	}, nil
}
