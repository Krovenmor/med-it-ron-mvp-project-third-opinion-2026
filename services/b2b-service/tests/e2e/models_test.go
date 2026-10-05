//go:build e2e

package e2e

import (
	"time"

	"github.com/google/uuid"
)

type report struct {
	SourceSystem string        `json:"source_system"`
	Study        reportStudy   `json:"study"`
	Conclusion   string        `json:"conclusion"`
	Patient      reportPatient `json:"patient"`
}

type reportStudy struct {
	ID          string `json:"id"`
	Modality    string `json:"modality"`
	BodySite    string `json:"body_site"`
	PerformedAt string `json:"performed_at"`
}

type reportPatient struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func newReport() report {
	id := uuid.NewString()
	return report{
		SourceSystem: "e2e",
		Study: reportStudy{
			ID:          "STUDY-" + id,
			Modality:    "CT",
			BodySite:    "chest",
			PerformedAt: "2026-10-01T09:00:00Z",
		},
		Conclusion: "Солидный очаг 9 мм, исследование " + id,
		Patient: reportPatient{
			ID:        "PATIENT-" + id,
			FullName:  "Пациент " + id,
			BirthDate: "1974-03-12",
			Sex:       "female",
			Phone:     "+7-" + id,
			Email:     id + "@example.com",
		},
	}
}

type caseRef struct {
	CaseID string `json:"case_id"`
	Status string `json:"status"`
}

type caseView struct {
	ID                string               `json:"id"`
	SourceSystem      string               `json:"source_system"`
	StudyID           string               `json:"study_id"`
	Status            string               `json:"status"`
	Urgency           string               `json:"urgency"`
	CatalogVersion    string               `json:"catalog_version"`
	GuidelinesVersion string               `json:"guidelines_version"`
	ReceivedAt        time.Time            `json:"received_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	Patient           patientView          `json:"patient"`
	Recommendations   []caseRecommendation `json:"recommendations"`
}

type patientView struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
}

type caseRecommendation struct {
	recommendation
	ID       string      `json:"id"`
	Position int         `json:"position"`
	Source   string      `json:"source"`
	Review   *reviewView `json:"review"`
}

type reviewView struct {
	Mark          string    `json:"mark"`
	RejectReason  string    `json:"reject_reason"`
	RejectComment string    `json:"reject_comment"`
	ReviewedBy    string    `json:"reviewed_by"`
	ReviewedAt    time.Time `json:"reviewed_at"`
}

type reviewQueueView struct {
	Cases []queueCaseView `json:"cases"`
}

type queueCaseView struct {
	CaseID                  string      `json:"case_id"`
	Urgency                 string      `json:"urgency"`
	Modality                string      `json:"modality"`
	Patient                 patientView `json:"patient"`
	RecommendationsTotal    int         `json:"recommendations_total"`
	RecommendationsReviewed int         `json:"recommendations_reviewed"`
	OpenedAt                time.Time   `json:"opened_at"`
}

type reviewRequest struct {
	Mark          string  `json:"mark"`
	RejectReason  string  `json:"reject_reason,omitempty"`
	RejectComment string  `json:"reject_comment,omitempty"`
	PatientText   *string `json:"patient_text,omitempty"`
}

type addRecommendationRequest struct {
	ServiceCode string `json:"service_code,omitempty"`
	ServiceName string `json:"service_name"`
	Rationale   string `json:"rationale,omitempty"`
	PatientText string `json:"patient_text"`
	Mark        string `json:"mark"`
}

type urgencyRequest struct {
	Urgency string `json:"urgency"`
	Reason  string `json:"reason,omitempty"`
}

type planView struct {
	Now     time.Time       `json:"now"`
	Patient planPatientView `json:"patient"`
	Cases   []planCaseView  `json:"cases"`
}

type planPatientView struct {
	FullName string `json:"full_name"`
}

type planCaseView struct {
	CaseID          string               `json:"case_id"`
	Status          string               `json:"status"`
	UrgentContact   bool                 `json:"urgent_contact"`
	NoFindings      bool                 `json:"no_findings"`
	ConfirmedAt     *time.Time           `json:"confirmed_at"`
	BookBy          *time.Time           `json:"book_by"`
	Study           planStudyView        `json:"study"`
	Recommendations []planRecommendation `json:"recommendations"`
}

type planStudyView struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendation struct {
	ID            string `json:"id"`
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	PatientText   string `json:"patient_text"`
	Mark          string `json:"mark"`
	Declined      bool   `json:"declined"`
	DeclineReason string `json:"decline_reason,omitempty"`
}

type bookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	ScheduledAt      string `json:"scheduled_at"`
	Channel          string `json:"channel"`
}

type bookingView struct {
	BookingID        string `json:"booking_id"`
	CaseID           string `json:"case_id"`
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	CaseStatus       string `json:"case_status"`
}

type caseEventView struct {
	Type    string            `db:"type"`
	Actor   string            `db:"actor"`
	Payload map[string]string `db:"payload"`
}

type assessment struct {
	Urgency           string           `json:"urgency"`
	CatalogVersion    string           `json:"catalog_version"`
	GuidelinesVersion string           `json:"guidelines_version"`
	Recommendations   []recommendation `json:"recommendations"`
}

type recommendation struct {
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	Importance    string `json:"importance"`
	Rationale     string `json:"rationale"`
	GuidelineRef  string `json:"guideline_ref"`
	PatientText   string `json:"patient_text"`
	AlreadyBooked bool   `json:"already_booked"`
}

func validAssessment() assessment {
	return assessment{
		Urgency:           "priority",
		CatalogVersion:    "catalog-e2e",
		GuidelinesVersion: "guidelines-e2e",
		Recommendations: []recommendation{
			{
				ServiceCode:  "PULM-CONSULT",
				ServiceName:  "Консультация пульмонолога",
				Importance:   "high",
				Rationale:    "Солидный очаг более 8 мм",
				GuidelineRef: "КР МЗ РФ",
				PatientText:  "Рекомендуем консультацию пульмонолога",
			},
			{
				ServiceName:   "Консультация торакального хирурга",
				Importance:    "low",
				PatientText:   "Рекомендуем консультацию хирурга в другой клинике",
				AlreadyBooked: true,
			},
		},
	}
}

type history struct {
	Visits       []visit       `json:"visits"`
	Appointments []appointment `json:"appointments"`
}

type visit struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type appointment struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type aiRequest struct {
	CaseID       string        `json:"case_id"`
	Study        aiStudy       `json:"study"`
	Conclusion   string        `json:"conclusion"`
	Patient      aiPatient     `json:"patient"`
	Visits       []visit       `json:"visits"`
	Appointments []appointment `json:"appointments"`
}

type aiStudy struct {
	Modality    string    `json:"modality"`
	BodySite    string    `json:"body_site"`
	PerformedAt time.Time `json:"performed_at"`
}

type aiPatient struct {
	Age int    `json:"age"`
	Sex string `json:"sex"`
}

type clockView struct {
	Now time.Time `json:"now"`
}

type jobState struct {
	State    string `db:"state"`
	Attempts int    `db:"attempts"`
}

type advanceClockRequest struct {
	By string `json:"by"`
}

type operatorTasksView struct {
	Tasks []operatorTaskItemView `json:"tasks"`
}

type operatorTaskItemView struct {
	Task         operatorTaskView    `json:"task"`
	Case         operatorCaseView    `json:"case"`
	Patient      operatorPatientView `json:"patient"`
	OfferService string              `json:"offer_service"`
}

type operatorTaskView struct {
	ID          string     `json:"id"`
	CaseID      string     `json:"case_id"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	Assignee    string     `json:"assignee"`
	DueAt       time.Time  `json:"due_at"`
	NextCallAt  *time.Time `json:"next_call_at"`
	Attempts    int        `json:"attempts"`
	NeedsDoctor bool       `json:"needs_doctor"`
	CreatedAt   time.Time  `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}

type operatorCaseView struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	Urgency  string    `json:"urgency"`
	Modality string    `json:"modality"`
	BookBy   time.Time `json:"book_by"`
}

type operatorPatientView struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type operatorCardView struct {
	Task     operatorTaskView    `json:"task"`
	Case     operatorCaseView    `json:"case"`
	Patient  operatorPatientView `json:"patient"`
	Attempts []callAttemptView   `json:"attempts"`
	Offers   []offerView         `json:"offers"`
}

type callAttemptView struct {
	Outcome       string     `json:"outcome"`
	DeclineReason string     `json:"decline_reason"`
	Comment       string     `json:"comment"`
	CallbackAt    *time.Time `json:"callback_at"`
	Actor         string     `json:"actor"`
}

type offerView struct {
	RecommendationID string     `json:"recommendation_id"`
	ServiceCode      string     `json:"service_code"`
	ServiceName      string     `json:"service_name"`
	PatientText      string     `json:"patient_text"`
	Mark             string     `json:"mark"`
	Booked           bool       `json:"booked"`
	Slots            []slotView `json:"slots"`
	SlotsUnavailable bool       `json:"slots_unavailable"`
}

type slotView struct {
	ID       string    `json:"id"`
	StartsAt time.Time `json:"starts_at"`
}

type outcomeRequest struct {
	Comment string `json:"comment,omitempty"`
}

type callbackRequest struct {
	CallAt  string `json:"call_at"`
	Comment string `json:"comment,omitempty"`
}

type declineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment,omitempty"`
}

type operatorBookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	SlotID           string `json:"slot_id"`
	Comment          string `json:"comment,omitempty"`
}

type notificationsView struct {
	Notifications []notificationView `json:"notifications"`
}

type notificationView struct {
	CaseID    string                  `json:"case_id"`
	Recipient string                  `json:"recipient"`
	Channel   string                  `json:"channel"`
	Kind      string                  `json:"kind"`
	Text      string                  `json:"text"`
	Urgency   string                  `json:"urgency"`
	Patient   notificationPatientView `json:"patient"`
}

type notificationPatientView struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type misSlotsResponse struct {
	Slots []misSlot `json:"slots"`
}

type misSlot struct {
	ID          string    `json:"id"`
	ServiceCode string    `json:"service_code"`
	StartsAt    time.Time `json:"starts_at"`
}

type misBookingRequest struct {
	PatientID   string `json:"patient_id"`
	SlotID      string `json:"slot_id"`
	ServiceName string `json:"service_name"`
	ReferralID  string `json:"referral_id"`
}

type misAppointment struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patient_id"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ReferralID  string    `json:"referral_id"`
}

type jobView struct {
	Kind  string `db:"kind"`
	State string `db:"state"`
}

type dashboardView struct {
	Period         periodView          `json:"period"`
	PreviousPeriod periodView          `json:"previous_period"`
	Current        periodStatsView     `json:"current"`
	Previous       periodStatsView     `json:"previous"`
	ByUrgency      []urgencySliceView  `json:"by_urgency"`
	ByModality     []modalitySliceView `json:"by_modality"`
	DeclineReasons []declineCountView  `json:"decline_reasons"`
	TrendStep      string              `json:"trend_step"`
	Trend          []trendPointView    `json:"trend"`
}

type periodView struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
	Days int       `json:"days"`
}

type periodStatsView struct {
	Funnel      funnelView   `json:"funnel"`
	DoctorSLA   ratioView    `json:"doctor_sla"`
	OperatorSLA ratioView    `json:"operator_sla"`
	Bookings    channelsView `json:"bookings"`
}

type funnelView struct {
	Received    int `json:"received"`
	Recommended int `json:"recommended"`
	Confirmed   int `json:"confirmed"`
	Notified    int `json:"notified"`
	Booked      int `json:"booked"`
	Completed   int `json:"completed"`
}

type ratioView struct {
	OnTime int `json:"on_time"`
	Total  int `json:"total"`
}

type channelsView struct {
	Self     int `json:"self"`
	Operator int `json:"operator"`
}

type urgencySliceView struct {
	Urgency     string    `json:"urgency"`
	Confirmed   int       `json:"confirmed"`
	Booked      int       `json:"booked"`
	DoctorSLA   ratioView `json:"doctor_sla"`
	OperatorSLA ratioView `json:"operator_sla"`
}

type modalitySliceView struct {
	Modality  string `json:"modality"`
	Confirmed int    `json:"confirmed"`
	Booked    int    `json:"booked"`
}

type declineCountView struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

type trendPointView struct {
	From     time.Time `json:"from"`
	Received int       `json:"received"`
	Booked   int       `json:"booked"`
}

type reportStatusView struct {
	CaseID     string     `json:"case_id"`
	StudyID    string     `json:"study_id"`
	Status     string     `json:"status"`
	ReceivedAt time.Time  `json:"received_at"`
	AssessedAt *time.Time `json:"assessed_at"`
}

type patientDeclineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment,omitempty"`
}

type patientDeclineView struct {
	RecommendationID string `json:"recommendation_id"`
	CaseStatus       string `json:"case_status"`
}

type helpRequestView struct {
	TaskID  string `json:"task_id"`
	Created bool   `json:"created"`
}
