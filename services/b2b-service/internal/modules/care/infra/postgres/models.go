package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type patientRow struct {
	ID           uuid.UUID `db:"id"`
	SourceSystem string    `db:"source_system"`
	ExternalID   string    `db:"external_id"`
	FullName     string    `db:"full_name"`
	BirthDate    time.Time `db:"birth_date"`
	Sex          string    `db:"sex"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
}

type routeRow struct {
	CaseID      uuid.UUID  `db:"case_id"`
	PatientID   uuid.UUID  `db:"patient_id"`
	Status      string     `db:"status"`
	Urgency     *string    `db:"urgency"`
	Modality    string     `db:"modality"`
	PerformedAt time.Time  `db:"performed_at"`
	ReceivedAt  time.Time  `db:"received_at"`
	AssessedAt  *time.Time `db:"assessed_at"`
	ReviewDueAt *time.Time `db:"review_due_at"`
	ConfirmedAt *time.Time `db:"confirmed_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type planItemRow struct {
	ID             uuid.UUID  `db:"id"`
	CaseID         uuid.UUID  `db:"case_id"`
	Position       int        `db:"position"`
	ServiceCode    string     `db:"service_code"`
	ServiceName    string     `db:"service_name"`
	PatientText    string     `db:"patient_text"`
	Mark           string     `db:"mark"`
	DeclinedAt     *time.Time `db:"declined_at"`
	DeclineReason  *string    `db:"decline_reason"`
	DeclineComment string     `db:"decline_comment"`
}

type serviceRow struct {
	Code string `db:"service_code"`
	Name string `db:"service_name"`
}

type bookingRow struct {
	ID               uuid.UUID `db:"id"`
	CaseID           uuid.UUID `db:"case_id"`
	RecommendationID uuid.UUID `db:"recommendation_id"`
	AppointmentID    string    `db:"appointment_id"`
	ScheduledAt      time.Time `db:"scheduled_at"`
	Channel          string    `db:"channel"`
	CreatedAt        time.Time `db:"created_at"`
}

type taskRow struct {
	ID         uuid.UUID  `db:"id"`
	CaseID     uuid.UUID  `db:"case_id"`
	Reason     string     `db:"reason"`
	Status     string     `db:"status"`
	Assignee   *string    `db:"assignee"`
	DueAt      time.Time  `db:"due_at"`
	NextCallAt *time.Time `db:"next_call_at"`
	Attempts   int        `db:"attempts"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	ClosedAt   *time.Time `db:"closed_at"`
}

type taskListRow struct {
	taskRow
	PatientID         uuid.UUID  `db:"patient_id"`
	RouteStatus       string     `db:"route_status"`
	Urgency           *string    `db:"urgency"`
	Modality          string     `db:"modality"`
	ReceivedAt        time.Time  `db:"received_at"`
	ConfirmedAt       *time.Time `db:"confirmed_at"`
	PatientExternalID string     `db:"patient_external_id"`
	PatientFullName   string     `db:"patient_full_name"`
	PatientPhone      string     `db:"patient_phone"`
	OfferService      *string    `db:"offer_service"`
}

type callAttemptRow struct {
	ID            uuid.UUID  `db:"id"`
	TaskID        uuid.UUID  `db:"task_id"`
	Outcome       string     `db:"outcome"`
	DeclineReason *string    `db:"decline_reason"`
	Comment       string     `db:"comment"`
	CallbackAt    *time.Time `db:"callback_at"`
	Actor         string     `db:"actor"`
	CreatedAt     time.Time  `db:"created_at"`
}

type notificationEntryRow struct {
	ID                uuid.UUID `db:"id"`
	CaseID            uuid.UUID `db:"case_id"`
	Recipient         string    `db:"recipient"`
	Channel           string    `db:"channel"`
	Kind              string    `db:"kind"`
	Text              string    `db:"text"`
	CreatedAt         time.Time `db:"created_at"`
	Urgency           *string   `db:"urgency"`
	PatientExternalID string    `db:"patient_external_id"`
	PatientFullName   string    `db:"patient_full_name"`
}

type caseFactRow struct {
	Urgency        *string    `db:"urgency"`
	Modality       string     `db:"modality"`
	ReceivedAt     time.Time  `db:"received_at"`
	AssessedAt     *time.Time `db:"assessed_at"`
	ReviewDueAt    *time.Time `db:"review_due_at"`
	ConfirmedAt    *time.Time `db:"confirmed_at"`
	NotifiedAt     *time.Time `db:"notified_at"`
	BookedAt       *time.Time `db:"booked_at"`
	CompletedAt    *time.Time `db:"completed_at"`
	Accepted       bool       `db:"accepted"`
	BookingChannel *string    `db:"booking_channel"`
}

type taskFactRow struct {
	Urgency   *string    `db:"urgency"`
	Status    string     `db:"status"`
	DueAt     time.Time  `db:"due_at"`
	CreatedAt time.Time  `db:"created_at"`
	ClosedAt  *time.Time `db:"closed_at"`
}

type declineCountRow struct {
	Reason string `db:"decline_reason"`
	Count  int    `db:"count"`
}

func (r patientRow) toDomain() domain.Patient {
	return domain.Patient{
		ID:           r.ID,
		SourceSystem: r.SourceSystem,
		ExternalID:   r.ExternalID,
		FullName:     r.FullName,
		BirthDate:    r.BirthDate,
		Sex:          r.Sex,
		Phone:        r.Phone,
		Email:        r.Email,
	}
}

func (r routeRow) toDomain() domain.Route {
	return domain.Route{
		CaseID:      r.CaseID,
		PatientID:   r.PatientID,
		Status:      domain.RouteStatus(r.Status),
		Urgency:     domain.Urgency(postgres.ValueOf(r.Urgency)),
		Modality:    domain.Modality(r.Modality),
		PerformedAt: r.PerformedAt,
		ReceivedAt:  r.ReceivedAt,
		AssessedAt:  postgres.ValueOf(r.AssessedAt),
		ReviewDueAt: postgres.ValueOf(r.ReviewDueAt),
		ConfirmedAt: postgres.ValueOf(r.ConfirmedAt),
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r planItemRow) toDomain() domain.PlanItem {
	return domain.PlanItem{
		ID:             r.ID,
		CaseID:         r.CaseID,
		Position:       r.Position,
		ServiceCode:    r.ServiceCode,
		ServiceName:    r.ServiceName,
		PatientText:    r.PatientText,
		Mark:           domain.Mark(r.Mark),
		DeclinedAt:     postgres.ValueOf(r.DeclinedAt),
		DeclineReason:  domain.DeclineReason(postgres.ValueOf(r.DeclineReason)),
		DeclineComment: r.DeclineComment,
	}
}

func (r serviceRow) toDomain() domain.Service {
	return domain.Service{Code: r.Code, Name: r.Name}
}

func (r bookingRow) toDomain() domain.Booking {
	return domain.Booking{
		ID:               r.ID,
		CaseID:           r.CaseID,
		RecommendationID: r.RecommendationID,
		AppointmentID:    r.AppointmentID,
		ScheduledAt:      r.ScheduledAt,
		Channel:          domain.BookingChannel(r.Channel),
		CreatedAt:        r.CreatedAt,
	}
}

func (r taskRow) toDomain() domain.OperatorTask {
	return domain.OperatorTask{
		ID:         r.ID,
		CaseID:     r.CaseID,
		Reason:     domain.TaskReason(r.Reason),
		Status:     domain.TaskStatus(r.Status),
		Assignee:   postgres.ValueOf(r.Assignee),
		DueAt:      r.DueAt,
		NextCallAt: postgres.ValueOf(r.NextCallAt),
		Attempts:   r.Attempts,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		ClosedAt:   postgres.ValueOf(r.ClosedAt),
	}
}

func (r taskListRow) toDomain() domain.TaskListItem {
	return domain.TaskListItem{
		Task: r.taskRow.toDomain(),
		Route: domain.Route{
			CaseID:      r.CaseID,
			PatientID:   r.PatientID,
			Status:      domain.RouteStatus(r.RouteStatus),
			Urgency:     domain.Urgency(postgres.ValueOf(r.Urgency)),
			Modality:    domain.Modality(r.Modality),
			ReceivedAt:  r.ReceivedAt,
			ConfirmedAt: postgres.ValueOf(r.ConfirmedAt),
		},
		Patient: domain.Patient{
			ID:         r.PatientID,
			ExternalID: r.PatientExternalID,
			FullName:   r.PatientFullName,
			Phone:      r.PatientPhone,
		},
		OfferService: postgres.ValueOf(r.OfferService),
	}
}

func (r callAttemptRow) toDomain() domain.CallAttempt {
	return domain.CallAttempt{
		ID:            r.ID,
		TaskID:        r.TaskID,
		Outcome:       domain.CallOutcome(r.Outcome),
		DeclineReason: domain.DeclineReason(postgres.ValueOf(r.DeclineReason)),
		Comment:       r.Comment,
		CallbackAt:    postgres.ValueOf(r.CallbackAt),
		Actor:         r.Actor,
		CreatedAt:     r.CreatedAt,
	}
}

func (r notificationEntryRow) toDomain() domain.NotificationEntry {
	return domain.NotificationEntry{
		Notification: domain.Notification{
			ID:        r.ID,
			CaseID:    r.CaseID,
			Recipient: domain.NotificationRecipient(r.Recipient),
			Channel:   domain.NotificationChannel(r.Channel),
			Kind:      domain.NotificationKind(r.Kind),
			Text:      r.Text,
			CreatedAt: r.CreatedAt,
		},
		Patient: domain.Patient{ExternalID: r.PatientExternalID, FullName: r.PatientFullName},
		Urgency: domain.Urgency(postgres.ValueOf(r.Urgency)),
	}
}

func (r caseFactRow) toDomain() domain.CaseFact {
	return domain.CaseFact{
		Urgency:        domain.Urgency(postgres.ValueOf(r.Urgency)),
		Modality:       domain.Modality(r.Modality),
		ReceivedAt:     r.ReceivedAt,
		AssessedAt:     postgres.ValueOf(r.AssessedAt),
		ReviewDueAt:    postgres.ValueOf(r.ReviewDueAt),
		ConfirmedAt:    postgres.ValueOf(r.ConfirmedAt),
		NotifiedAt:     postgres.ValueOf(r.NotifiedAt),
		BookedAt:       postgres.ValueOf(r.BookedAt),
		CompletedAt:    postgres.ValueOf(r.CompletedAt),
		Accepted:       r.Accepted,
		BookingChannel: domain.BookingChannel(postgres.ValueOf(r.BookingChannel)),
	}
}

func (r taskFactRow) toDomain() domain.TaskFact {
	return domain.TaskFact{
		Urgency:   domain.Urgency(postgres.ValueOf(r.Urgency)),
		Status:    domain.TaskStatus(r.Status),
		DueAt:     r.DueAt,
		CreatedAt: r.CreatedAt,
		ClosedAt:  postgres.ValueOf(r.ClosedAt),
	}
}

func (r declineCountRow) toDomain() domain.DeclineCount {
	return domain.DeclineCount{Reason: domain.DeclineReason(r.Reason), Count: r.Count}
}

func upsertPatientArgs(p domain.Patient, now time.Time) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":            p.ID,
		"source_system": p.SourceSystem,
		"external_id":   p.ExternalID,
		"full_name":     p.FullName,
		"birth_date":    p.BirthDate,
		"sex":           p.Sex,
		"phone":         p.Phone,
		"email":         p.Email,
		"updated_at":    now,
	}
}

func routeArgs(r domain.Route) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":       r.CaseID,
		"patient_id":    r.PatientID,
		"status":        string(r.Status),
		"urgency":       string(r.Urgency),
		"modality":      string(r.Modality),
		"performed_at":  r.PerformedAt,
		"received_at":   r.ReceivedAt,
		"assessed_at":   postgres.NullableTime(r.AssessedAt),
		"review_due_at": postgres.NullableTime(r.ReviewDueAt),
		"confirmed_at":  postgres.NullableTime(r.ConfirmedAt),
		"updated_at":    r.UpdatedAt,
	}
}

func planItemArgs(i domain.PlanItem) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":           i.ID,
		"case_id":      i.CaseID,
		"position":     i.Position,
		"service_code": i.ServiceCode,
		"service_name": i.ServiceName,
		"patient_text": i.PatientText,
		"mark":         string(i.Mark),
	}
}

func declineArgs(i domain.PlanItem) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":              i.ID,
		"declined_at":     postgres.NullableTime(i.DeclinedAt),
		"decline_reason":  string(i.DeclineReason),
		"decline_comment": i.DeclineComment,
	}
}

func insertBookingArgs(b domain.Booking) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":           b.CaseID,
		"recommendation_id": b.RecommendationID,
		"appointment_id":    b.AppointmentID,
		"scheduled_at":      b.ScheduledAt,
		"channel":           string(b.Channel),
		"created_at":        b.CreatedAt,
	}
}

func insertTaskArgs(t domain.OperatorTask) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":    t.CaseID,
		"reason":     string(t.Reason),
		"status":     string(t.Status),
		"due_at":     t.DueAt,
		"created_at": t.CreatedAt,
	}
}

func updateTaskArgs(t domain.OperatorTask) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":           t.ID,
		"status":       string(t.Status),
		"assignee":     t.Assignee,
		"next_call_at": postgres.NullableTime(t.NextCallAt),
		"attempts":     t.Attempts,
		"updated_at":   t.UpdatedAt,
		"closed_at":    postgres.NullableTime(t.ClosedAt),
	}
}

func insertCallAttemptArgs(a domain.CallAttempt) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"task_id":        a.TaskID,
		"outcome":        string(a.Outcome),
		"decline_reason": string(a.DeclineReason),
		"comment":        a.Comment,
		"callback_at":    postgres.NullableTime(a.CallbackAt),
		"actor":          a.Actor,
		"created_at":     a.CreatedAt,
	}
}

func insertNotificationArgs(n domain.Notification) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"case_id":    n.CaseID,
		"recipient":  string(n.Recipient),
		"channel":    string(n.Channel),
		"kind":       string(n.Kind),
		"text":       n.Text,
		"created_at": n.CreatedAt,
	}
}

func insertCaseEventArgs(e domain.CaseEvent) pgx.StrictNamedArgs {
	payload := e.Payload
	if payload == nil {
		payload = map[string]string{}
	}
	return pgx.StrictNamedArgs{
		"payload":     payload,
		"case_id":     e.CaseID,
		"type":        string(e.Type),
		"from_status": string(e.FromStatus),
		"to_status":   string(e.ToStatus),
		"actor":       e.Actor,
		"occurred_at": e.OccurredAt,
	}
}

func historyTaskArgs(t domain.OperatorTask) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{
		"id":           t.ID,
		"case_id":      t.CaseID,
		"reason":       string(t.Reason),
		"status":       string(t.Status),
		"assignee":     t.Assignee,
		"due_at":       t.DueAt,
		"next_call_at": postgres.NullableTime(t.NextCallAt),
		"attempts":     t.Attempts,
		"created_at":   t.CreatedAt,
		"updated_at":   t.UpdatedAt,
		"closed_at":    postgres.NullableTime(t.ClosedAt),
	}
}

func historyBookingArgs(b domain.Booking) pgx.StrictNamedArgs {
	args := insertBookingArgs(b)
	args["id"] = b.ID
	return args
}

func historyPatient(c history.Case) domain.Patient {
	return domain.Patient{
		ID:           c.Patient.ID,
		SourceSystem: c.SourceSystem,
		ExternalID:   c.Patient.ExternalID,
		FullName:     c.Patient.FullName,
		BirthDate:    c.Patient.BirthDate,
		Sex:          c.Patient.Sex,
		Phone:        c.Patient.Phone,
	}
}

func historyRoute(c history.Case) domain.Route {
	return domain.Route{
		CaseID:      c.ID,
		PatientID:   c.Patient.ID,
		Status:      historyStatus(c),
		Urgency:     domain.Urgency(c.Urgency),
		Modality:    domain.Modality(c.Study.Modality),
		PerformedAt: c.Study.PerformedAt,
		ReceivedAt:  c.ReceivedAt,
		AssessedAt:  c.AssessedAt,
		ReviewDueAt: c.ReviewDueAt,
		ConfirmedAt: c.Review.ConfirmedAt,
		UpdatedAt:   c.ReceivedAt,
	}
}

func historyStatus(c history.Case) domain.RouteStatus {
	switch {
	case !c.CompletedAt.IsZero():
		return domain.RouteStatusCompleted
	case len(c.Bookings) > 0:
		return domain.RouteStatusBooked
	case !c.DeclinedAt.IsZero():
		return domain.RouteStatusDeclined
	case !c.UnreachableAt.IsZero():
		return domain.RouteStatusUnreachable
	case !c.NotifiedAt.IsZero():
		return domain.RouteStatusNotified
	case !c.Review.ConfirmedAt.IsZero():
		return domain.RouteStatusConfirmed
	}
	return domain.RouteStatusInReview
}

func historyPlanItems(c history.Case) []domain.PlanItem {
	accepted := c.Accepted()
	items := make([]domain.PlanItem, 0, len(accepted))
	for _, r := range accepted {
		items = append(items, domain.PlanItem{
			ID:          r.ID,
			CaseID:      c.ID,
			Position:    r.Position,
			ServiceCode: r.ServiceCode,
			ServiceName: r.ServiceName,
			PatientText: r.PatientText,
			Mark:        domain.Mark(r.Mark),
		})
	}
	return items
}

func historyBooking(c history.Case, b history.Booking) domain.Booking {
	return domain.Booking{
		ID:               b.ID,
		CaseID:           c.ID,
		RecommendationID: b.RecommendationID,
		AppointmentID:    b.AppointmentID,
		ScheduledAt:      b.ScheduledAt,
		Channel:          domain.BookingChannel(b.Channel),
		CreatedAt:        b.CreatedAt,
	}
}

func historyTask(c history.Case, t history.Task) domain.OperatorTask {
	return domain.OperatorTask{
		ID:         t.ID,
		CaseID:     c.ID,
		Reason:     domain.TaskReason(t.Reason),
		Status:     domain.TaskStatus(t.Status),
		Assignee:   t.Assignee,
		DueAt:      t.DueAt,
		NextCallAt: t.NextCallAt,
		Attempts:   noAnswers(t),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.ClosedAt,
		ClosedAt:   t.ClosedAt,
	}
}

func historyAttempt(t history.Task, a history.Attempt) domain.CallAttempt {
	return domain.CallAttempt{
		TaskID:        t.ID,
		Outcome:       domain.CallOutcome(a.Outcome),
		DeclineReason: domain.DeclineReason(a.DeclineReason),
		Comment:       a.Comment,
		Actor:         a.Actor,
		CreatedAt:     a.CreatedAt,
	}
}

func noAnswers(t history.Task) int {
	n := 0
	for _, a := range t.Attempts {
		if a.Outcome == string(domain.CallOutcomeNoAnswer) {
			n++
		}
	}
	return n
}

func historyEvents(c history.Case) []domain.CaseEvent {
	events := []domain.CaseEvent{
		domain.StatusChanged(c.ID, "", domain.RouteStatusReceived, domain.ActorSystem, c.ReceivedAt),
		domain.StatusChanged(c.ID, domain.RouteStatusReceived, domain.RouteStatusInReview, domain.ActorSystem, c.AssessedAt),
		domain.StatusChanged(c.ID, domain.RouteStatusInReview, domain.RouteStatusConfirmed, c.Review.Doctor, c.Review.ConfirmedAt),
	}
	status := domain.RouteStatusConfirmed
	move := func(to domain.RouteStatus, actor string, at time.Time) {
		if at.IsZero() {
			return
		}
		events = append(events, domain.StatusChanged(c.ID, status, to, actor, at))
		status = to
	}
	move(domain.RouteStatusNotified, domain.ActorSystem, c.NotifiedAt)
	move(domain.RouteStatusUnreachable, c.UnreachableBy, c.UnreachableAt)
	move(domain.RouteStatusDeclined, c.DeclinedBy, c.DeclinedAt)
	for i, b := range c.Bookings {
		booking := historyBooking(c, b)
		if i == 0 {
			move(domain.RouteStatusBooked, booking.Channel.Actor(), booking.CreatedAt)
		}
		events = append(events, domain.RecommendationBooked(booking))
	}
	move(domain.RouteStatusCompleted, "mis", c.CompletedAt)
	for _, t := range c.Tasks {
		task := historyTask(c, t)
		events = append(events, domain.OperatorTaskCreated(task))
		for _, a := range t.Attempts {
			events = append(events, domain.CallAttempted(c.ID, historyAttempt(t, a)))
		}
		events = append(events, domain.OperatorTaskClosed(task, t.Assignee))
	}
	return events
}
