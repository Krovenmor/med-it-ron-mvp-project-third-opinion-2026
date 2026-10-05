package httpx

import (
	"fmt"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

func ParseOptionalTime(field, value, layout, format string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s must be %s", apperr.ErrInvalidInput, field, format)
	}
	return t, nil
}

func OptionalTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
