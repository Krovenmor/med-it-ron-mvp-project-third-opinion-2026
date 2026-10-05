package domain

import (
	"fmt"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

var ErrStudyConflict = fmt.Errorf("%w: study already received with different content", apperr.ErrConflict)
