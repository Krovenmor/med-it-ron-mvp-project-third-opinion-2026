package intake

import "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"

type IngestResult struct {
	Case    domain.Case
	Created bool
}
