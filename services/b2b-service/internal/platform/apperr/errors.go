package apperr

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrInvalidState = errors.New("invalid state")
	ErrConflict     = errors.New("conflict")
	ErrUpstream     = errors.New("upstream service failed")
)

func Invalid(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, reason)
}
