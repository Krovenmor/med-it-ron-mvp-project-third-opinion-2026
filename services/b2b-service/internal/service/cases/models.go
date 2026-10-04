package cases

import "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"

type Details struct {
	Case            domain.Case
	Patient         domain.Patient
	Recommendations []domain.Recommendation
}
