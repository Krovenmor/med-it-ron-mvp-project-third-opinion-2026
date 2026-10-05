package scheduling

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid booking request")
	ErrSlotTaken      = errors.New("slot is already taken")
)
