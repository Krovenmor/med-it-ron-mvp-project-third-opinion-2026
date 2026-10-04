package domain

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotFound      = errors.New("not found")
	ErrStudyConflict = errors.New("study already received with different content")
	ErrInvalidState  = errors.New("invalid state")
	ErrLeaseLost     = errors.New("job lease lost")
)
