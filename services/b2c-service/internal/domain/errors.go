package domain

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrNotFound        = errors.New("not found")
	ErrInvalidState    = errors.New("invalid state")
	ErrSlotUnavailable = errors.New("slot is unavailable")
	ErrUpstream        = errors.New("upstream service failed")
)
