package http

import (
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/scheduling"
)

type deliveriesResponse struct {
	Delivered []b2b.Delivery `json:"delivered"`
	Error     string         `json:"error,omitempty"`
}

type slotsResponse struct {
	Slots []scheduling.Slot `json:"slots"`
}

type bookAppointmentRequest struct {
	PatientID   string `json:"patient_id"`
	SlotID      string `json:"slot_id"`
	ServiceName string `json:"service_name"`
	ReferralID  string `json:"referral_id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (r bookAppointmentRequest) toRequest() scheduling.BookRequest {
	return scheduling.BookRequest{
		PatientID:   r.PatientID,
		SlotID:      r.SlotID,
		ServiceName: r.ServiceName,
		ReferralID:  r.ReferralID,
	}
}
