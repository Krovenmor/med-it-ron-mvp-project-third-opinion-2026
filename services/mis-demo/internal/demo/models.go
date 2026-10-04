package demo

import "time"

type Visit struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type Appointment struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patient_id"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ReferralID  string    `json:"referral_id,omitempty"`
}

type History struct {
	Visits       []Visit       `json:"visits"`
	Appointments []Appointment `json:"appointments"`
}

type Patient struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Study struct {
	ID          string    `json:"id"`
	Modality    string    `json:"modality"`
	BodySite    string    `json:"body_site"`
	PerformedAt time.Time `json:"performed_at"`
}

type Report struct {
	SourceSystem string  `json:"source_system"`
	Study        Study   `json:"study"`
	Conclusion   string  `json:"conclusion"`
	Patient      Patient `json:"patient"`
}

type Service struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
