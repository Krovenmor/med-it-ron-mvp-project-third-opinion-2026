package scheduling

import "time"

type Slot struct {
	ID          string    `json:"id"`
	ServiceCode string    `json:"service_code"`
	StartsAt    time.Time `json:"starts_at"`
}

type BookRequest struct {
	PatientID   string
	SlotID      string
	ServiceName string
	ReferralID  string
}
