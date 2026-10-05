package http

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/booking"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/operator"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/patient"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

type bookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	ScheduledAt      string `json:"scheduled_at"`
	Channel          string `json:"channel"`
}

type bookingResponse struct {
	BookingID        uuid.UUID `json:"booking_id"`
	CaseID           uuid.UUID `json:"case_id"`
	RecommendationID uuid.UUID `json:"recommendation_id"`
	AppointmentID    string    `json:"appointment_id"`
	CaseStatus       string    `json:"case_status"`
}

type patientPlanResponse struct {
	Now     time.Time           `json:"now"`
	Patient planPatientResponse `json:"patient"`
	Cases   []planCaseResponse  `json:"cases"`
}

type planPatientResponse struct {
	FullName string `json:"full_name"`
}

type planCaseResponse struct {
	CaseID          uuid.UUID                    `json:"case_id"`
	Status          string                       `json:"status"`
	UrgentContact   bool                         `json:"urgent_contact"`
	NoFindings      bool                         `json:"no_findings"`
	ConfirmedAt     *time.Time                   `json:"confirmed_at"`
	BookBy          *time.Time                   `json:"book_by"`
	Study           planStudyResponse            `json:"study"`
	Recommendations []planRecommendationResponse `json:"recommendations"`
}

type planStudyResponse struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendationResponse struct {
	ID            uuid.UUID `json:"id"`
	ServiceCode   string    `json:"service_code"`
	ServiceName   string    `json:"service_name"`
	PatientText   string    `json:"patient_text"`
	Mark          string    `json:"mark"`
	Declined      bool      `json:"declined"`
	DeclineReason string    `json:"decline_reason,omitempty"`
}

type patientDeclineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

type patientDeclineResponse struct {
	RecommendationID uuid.UUID `json:"recommendation_id"`
	CaseStatus       string    `json:"case_status"`
}

type helpRequestResponse struct {
	TaskID  uuid.UUID `json:"task_id"`
	Created bool      `json:"created"`
}

func (r bookingRequest) toCommand(caseID uuid.UUID) (booking.Book, error) {
	recommendationID, err := uuid.Parse(r.RecommendationID)
	if err != nil {
		return booking.Book{}, fmt.Errorf("%w: recommendation_id must be a UUID", apperr.ErrInvalidInput)
	}
	scheduledAt, err := httpx.ParseOptionalTime("scheduled_at", r.ScheduledAt, time.RFC3339, "an RFC 3339 date-time")
	if err != nil {
		return booking.Book{}, err
	}
	return booking.Book{
		CaseID:           caseID,
		RecommendationID: recommendationID,
		AppointmentID:    r.AppointmentID,
		ScheduledAt:      scheduledAt,
		Channel:          domain.BookingChannel(r.Channel),
	}, nil
}

func newBookingResponse(r booking.Result) bookingResponse {
	return bookingResponse{
		BookingID:        r.Booking.ID,
		CaseID:           r.Booking.CaseID,
		RecommendationID: r.Booking.RecommendationID,
		AppointmentID:    r.Booking.AppointmentID,
		CaseStatus:       string(r.CaseStatus),
	}
}

func newPatientPlanResponse(p domain.PatientPlan) patientPlanResponse {
	resp := patientPlanResponse{
		Now:     p.Now,
		Patient: planPatientResponse{FullName: p.Patient.FullName},
		Cases:   make([]planCaseResponse, 0, len(p.Cases)),
	}
	for _, c := range p.Cases {
		resp.Cases = append(resp.Cases, newPlanCaseResponse(c))
	}
	return resp
}

func newPlanCaseResponse(c domain.PlanCase) planCaseResponse {
	recs := make([]planRecommendationResponse, 0, len(c.Items))
	for _, item := range c.Items {
		recs = append(recs, planRecommendationResponse{
			ID:            item.ID,
			ServiceCode:   item.ServiceCode,
			ServiceName:   item.ServiceName,
			PatientText:   item.PatientText,
			Mark:          string(item.Mark),
			Declined:      item.Declined(),
			DeclineReason: string(item.DeclineReason),
		})
	}
	resp := planCaseResponse{
		CaseID:        c.Route.CaseID,
		Status:        string(c.Route.Status),
		UrgentContact: c.UrgentContact,
		NoFindings:    c.Route.NoFindings(),
		ConfirmedAt:   httpx.OptionalTime(c.Route.ConfirmedAt),
		Study: planStudyResponse{
			Modality:    string(c.Route.Modality),
			PerformedAt: c.Route.PerformedAt,
		},
		Recommendations: recs,
	}
	if c.Route.Confirmed() {
		bookBy := c.Route.BookBy()
		resp.BookBy = &bookBy
	}
	return resp
}

func (r patientDeclineRequest) toCommand(caseID, recommendationID uuid.UUID) patient.Decline {
	return patient.Decline{
		CaseID:           caseID,
		RecommendationID: recommendationID,
		Reason:           domain.DeclineReason(r.Reason),
		Comment:          r.Comment,
	}
}

func newPatientDeclineResponse(r patient.DeclineResult) patientDeclineResponse {
	return patientDeclineResponse{RecommendationID: r.Item.ID, CaseStatus: string(r.CaseStatus)}
}

type operatorTasksResponse struct {
	Tasks []operatorTaskItemResponse `json:"tasks"`
}

type operatorTaskItemResponse struct {
	Task         operatorTaskResponse    `json:"task"`
	Case         operatorCaseResponse    `json:"case"`
	Patient      operatorPatientResponse `json:"patient"`
	OfferService string                  `json:"offer_service"`
}

type operatorTaskResponse struct {
	ID          uuid.UUID  `json:"id"`
	CaseID      uuid.UUID  `json:"case_id"`
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

type operatorCaseResponse struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	Urgency  string    `json:"urgency"`
	Modality string    `json:"modality"`
	BookBy   time.Time `json:"book_by"`
}

type operatorPatientResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type operatorCardResponse struct {
	Task     operatorTaskResponse    `json:"task"`
	Case     operatorCaseResponse    `json:"case"`
	Patient  operatorPatientResponse `json:"patient"`
	Attempts []callAttemptResponse   `json:"attempts"`
	Offers   []operatorOfferResponse `json:"offers"`
}

type callAttemptResponse struct {
	ID            uuid.UUID  `json:"id"`
	Outcome       string     `json:"outcome"`
	DeclineReason string     `json:"decline_reason,omitempty"`
	Comment       string     `json:"comment"`
	CallbackAt    *time.Time `json:"callback_at"`
	Actor         string     `json:"actor"`
	CreatedAt     time.Time  `json:"created_at"`
}

type operatorOfferResponse struct {
	RecommendationID uuid.UUID      `json:"recommendation_id"`
	ServiceCode      string         `json:"service_code"`
	ServiceName      string         `json:"service_name"`
	PatientText      string         `json:"patient_text"`
	Mark             string         `json:"mark"`
	Booked           bool           `json:"booked"`
	Slots            []slotResponse `json:"slots"`
	SlotsUnavailable bool           `json:"slots_unavailable"`
}

type slotResponse struct {
	ID       string    `json:"id"`
	StartsAt time.Time `json:"starts_at"`
}

type outcomeRequest struct {
	Comment string `json:"comment"`
}

type callbackRequest struct {
	CallAt  string `json:"call_at"`
	Comment string `json:"comment"`
}

type declineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

type operatorBookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	SlotID           string `json:"slot_id"`
	Comment          string `json:"comment"`
}

type notificationsResponse struct {
	Notifications []notificationResponse `json:"notifications"`
}

type notificationResponse struct {
	ID        uuid.UUID                   `json:"id"`
	CaseID    uuid.UUID                   `json:"case_id"`
	Recipient string                      `json:"recipient"`
	Channel   string                      `json:"channel"`
	Kind      string                      `json:"kind"`
	Text      string                      `json:"text"`
	Urgency   string                      `json:"urgency"`
	Patient   notificationPatientResponse `json:"patient"`
	CreatedAt time.Time                   `json:"created_at"`
}

type notificationPatientResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

func (r outcomeRequest) toCommand(taskID uuid.UUID, actor string) operator.Outcome {
	return operator.Outcome{TaskID: taskID, Actor: actor, Comment: r.Comment}
}

func (r callbackRequest) toCommand(taskID uuid.UUID, actor string) (operator.Callback, error) {
	if r.CallAt == "" {
		return operator.Callback{}, fmt.Errorf("%w: call_at is required", apperr.ErrInvalidInput)
	}
	at, err := httpx.ParseOptionalTime("call_at", r.CallAt, time.RFC3339, "an RFC 3339 date-time")
	if err != nil {
		return operator.Callback{}, err
	}
	return operator.Callback{
		Outcome: operator.Outcome{TaskID: taskID, Actor: actor, Comment: r.Comment},
		At:      at,
	}, nil
}

func (r declineRequest) toCommand(taskID uuid.UUID, actor string) operator.Decline {
	return operator.Decline{
		Outcome: operator.Outcome{TaskID: taskID, Actor: actor, Comment: r.Comment},
		Reason:  domain.DeclineReason(r.Reason),
	}
}

func (r operatorBookingRequest) toCommand(taskID uuid.UUID, actor string) (operator.Book, error) {
	recommendationID, err := uuid.Parse(r.RecommendationID)
	if err != nil {
		return operator.Book{}, fmt.Errorf("%w: recommendation_id must be a UUID", apperr.ErrInvalidInput)
	}
	if r.SlotID == "" {
		return operator.Book{}, fmt.Errorf("%w: slot_id is required", apperr.ErrInvalidInput)
	}
	return operator.Book{
		Outcome:          operator.Outcome{TaskID: taskID, Actor: actor, Comment: r.Comment},
		RecommendationID: recommendationID,
		SlotID:           r.SlotID,
	}, nil
}

func newOperatorTasksResponse(items []domain.TaskListItem) operatorTasksResponse {
	resp := operatorTasksResponse{Tasks: make([]operatorTaskItemResponse, 0, len(items))}
	for _, item := range items {
		resp.Tasks = append(resp.Tasks, operatorTaskItemResponse{
			Task:         newOperatorTaskResponse(item.Task),
			Case:         newOperatorCaseResponse(item.Route),
			Patient:      newOperatorPatientResponse(item.Patient),
			OfferService: item.OfferService,
		})
	}
	return resp
}

func newOperatorTaskResponse(t domain.OperatorTask) operatorTaskResponse {
	return operatorTaskResponse{
		ID:          t.ID,
		CaseID:      t.CaseID,
		Reason:      string(t.Reason),
		Status:      string(t.Status),
		Assignee:    t.Assignee,
		DueAt:       t.DueAt,
		NextCallAt:  httpx.OptionalTime(t.NextCallAt),
		Attempts:    t.Attempts,
		NeedsDoctor: t.NeedsDoctor(),
		CreatedAt:   t.CreatedAt,
		ClosedAt:    httpx.OptionalTime(t.ClosedAt),
	}
}

func newOperatorCaseResponse(r domain.Route) operatorCaseResponse {
	return operatorCaseResponse{
		ID:       r.CaseID,
		Status:   string(r.Status),
		Urgency:  string(r.Urgency),
		Modality: string(r.Modality),
		BookBy:   r.BookBy(),
	}
}

func newOperatorPatientResponse(p domain.Patient) operatorPatientResponse {
	return operatorPatientResponse{ID: p.ExternalID, FullName: p.FullName, Phone: p.Phone}
}

func newOperatorCardResponse(c operator.Card) operatorCardResponse {
	resp := operatorCardResponse{
		Task:     newOperatorTaskResponse(c.Task),
		Case:     newOperatorCaseResponse(c.Route),
		Patient:  newOperatorPatientResponse(c.Patient),
		Attempts: make([]callAttemptResponse, 0, len(c.Attempts)),
		Offers:   make([]operatorOfferResponse, 0, len(c.Offers)),
	}
	for _, a := range c.Attempts {
		resp.Attempts = append(resp.Attempts, newCallAttemptResponse(a))
	}
	for _, o := range c.Offers {
		resp.Offers = append(resp.Offers, newOperatorOfferResponse(o))
	}
	return resp
}

func newCallAttemptResponse(a domain.CallAttempt) callAttemptResponse {
	return callAttemptResponse{
		ID:            a.ID,
		Outcome:       string(a.Outcome),
		DeclineReason: string(a.DeclineReason),
		Comment:       a.Comment,
		CallbackAt:    httpx.OptionalTime(a.CallbackAt),
		Actor:         a.Actor,
		CreatedAt:     a.CreatedAt,
	}
}

func newOperatorOfferResponse(o operator.Offer) operatorOfferResponse {
	resp := operatorOfferResponse{
		RecommendationID: o.Item.ID,
		ServiceCode:      o.Item.ServiceCode,
		ServiceName:      o.Item.ServiceName,
		PatientText:      o.Item.PatientText,
		Mark:             string(o.Item.Mark),
		Booked:           o.Booked,
		Slots:            make([]slotResponse, 0, len(o.Slots)),
		SlotsUnavailable: o.SlotsUnavailable,
	}
	for _, s := range o.Slots {
		resp.Slots = append(resp.Slots, slotResponse{ID: s.ID, StartsAt: s.StartsAt})
	}
	return resp
}

func newNotificationsResponse(entries []domain.NotificationEntry) notificationsResponse {
	resp := notificationsResponse{Notifications: make([]notificationResponse, 0, len(entries))}
	for _, e := range entries {
		n := e.Notification
		resp.Notifications = append(resp.Notifications, notificationResponse{
			ID:        n.ID,
			CaseID:    n.CaseID,
			Recipient: string(n.Recipient),
			Channel:   string(n.Channel),
			Kind:      string(n.Kind),
			Text:      n.Text,
			Urgency:   string(e.Urgency),
			Patient:   notificationPatientResponse{ID: e.Patient.ExternalID, FullName: e.Patient.FullName},
			CreatedAt: n.CreatedAt,
		})
	}
	return resp
}

type dashboardResponse struct {
	Period         periodResponse          `json:"period"`
	PreviousPeriod periodResponse          `json:"previous_period"`
	Current        periodStatsResponse     `json:"current"`
	Previous       periodStatsResponse     `json:"previous"`
	ByUrgency      []urgencySliceResponse  `json:"by_urgency"`
	ByModality     []modalitySliceResponse `json:"by_modality"`
	DeclineReasons []declineCountResponse  `json:"decline_reasons"`
	TrendStep      string                  `json:"trend_step"`
	Trend          []trendPointResponse    `json:"trend"`
}

type periodResponse struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
	Days int       `json:"days"`
}

type periodStatsResponse struct {
	Funnel      funnelResponse   `json:"funnel"`
	DoctorSLA   ratioResponse    `json:"doctor_sla"`
	OperatorSLA ratioResponse    `json:"operator_sla"`
	Bookings    channelsResponse `json:"bookings"`
}

type funnelResponse struct {
	Received    int `json:"received"`
	Recommended int `json:"recommended"`
	Confirmed   int `json:"confirmed"`
	Notified    int `json:"notified"`
	Booked      int `json:"booked"`
	Completed   int `json:"completed"`
}

type ratioResponse struct {
	OnTime int `json:"on_time"`
	Total  int `json:"total"`
}

type channelsResponse struct {
	Self     int `json:"self"`
	Operator int `json:"operator"`
}

type urgencySliceResponse struct {
	Urgency     string        `json:"urgency"`
	Confirmed   int           `json:"confirmed"`
	Booked      int           `json:"booked"`
	DoctorSLA   ratioResponse `json:"doctor_sla"`
	OperatorSLA ratioResponse `json:"operator_sla"`
}

type modalitySliceResponse struct {
	Modality  string `json:"modality"`
	Confirmed int    `json:"confirmed"`
	Booked    int    `json:"booked"`
}

type declineCountResponse struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

type trendPointResponse struct {
	From     time.Time `json:"from"`
	Received int       `json:"received"`
	Booked   int       `json:"booked"`
}

func newDashboardResponse(d domain.Dashboard) dashboardResponse {
	resp := dashboardResponse{
		Period:         newPeriodResponse(d.Period),
		PreviousPeriod: newPeriodResponse(d.PreviousPeriod),
		Current:        newPeriodStatsResponse(d.Current),
		Previous:       newPeriodStatsResponse(d.Previous),
		ByUrgency:      make([]urgencySliceResponse, 0, len(d.ByUrgency)),
		ByModality:     make([]modalitySliceResponse, 0, len(d.ByModality)),
		DeclineReasons: make([]declineCountResponse, 0, len(d.DeclineReasons)),
		TrendStep:      trendStepName(d.Period),
		Trend:          make([]trendPointResponse, 0, len(d.Trend)),
	}
	for _, s := range d.ByUrgency {
		resp.ByUrgency = append(resp.ByUrgency, urgencySliceResponse{
			Urgency:     string(s.Urgency),
			Confirmed:   s.Confirmed,
			Booked:      s.Booked,
			DoctorSLA:   newRatioResponse(s.DoctorSLA),
			OperatorSLA: newRatioResponse(s.OperatorSLA),
		})
	}
	for _, s := range d.ByModality {
		resp.ByModality = append(resp.ByModality, modalitySliceResponse{Modality: string(s.Modality), Confirmed: s.Confirmed, Booked: s.Booked})
	}
	for _, c := range d.DeclineReasons {
		resp.DeclineReasons = append(resp.DeclineReasons, declineCountResponse{Reason: string(c.Reason), Count: c.Count})
	}
	for _, p := range d.Trend {
		resp.Trend = append(resp.Trend, trendPointResponse{From: p.From, Received: p.Received, Booked: p.Booked})
	}
	return resp
}

func newPeriodResponse(p domain.Period) periodResponse {
	return periodResponse{From: p.From, To: p.To, Days: p.Days}
}

func newPeriodStatsResponse(s domain.PeriodStats) periodStatsResponse {
	return periodStatsResponse{
		Funnel: funnelResponse{
			Received:    s.Funnel.Received,
			Recommended: s.Funnel.Recommended,
			Confirmed:   s.Funnel.Confirmed,
			Notified:    s.Funnel.Notified,
			Booked:      s.Funnel.Booked,
			Completed:   s.Funnel.Completed,
		},
		DoctorSLA:   newRatioResponse(s.DoctorSLA),
		OperatorSLA: newRatioResponse(s.OperatorSLA),
		Bookings:    channelsResponse{Self: s.SelfBookings, Operator: s.OperatorBookings},
	}
}

func newRatioResponse(r domain.Ratio) ratioResponse {
	return ratioResponse{OnTime: r.Hit, Total: r.Total}
}

func trendStepName(p domain.Period) string {
	if p.TrendStep() > 24*time.Hour {
		return "week"
	}
	return "day"
}
