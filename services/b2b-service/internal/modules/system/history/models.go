package history

import (
	"time"

	"github.com/google/uuid"
)

const SourceSystem = "demo-history"

type Case struct {
	ID                uuid.UUID
	SourceSystem      string
	Patient           Patient
	Study             Study
	Conclusion        string
	Urgency           string
	CatalogVersion    string
	GuidelinesVersion string
	ReceivedAt        time.Time
	AssessedAt        time.Time
	ReviewDueAt       time.Time
	Review            Review
	Recommendations   []Recommendation
	NotifiedAt        time.Time
	Bookings          []Booking
	Tasks             []Task
	CompletedAt       time.Time
	DeclinedAt        time.Time
	DeclinedBy        string
	UnreachableAt     time.Time
	UnreachableBy     string
}

type Patient struct {
	ID         uuid.UUID
	ExternalID string
	FullName   string
	BirthDate  time.Time
	Sex        string
	Phone      string
}

type Study struct {
	ID          string
	Modality    string
	BodySite    string
	PerformedAt time.Time
}

type Review struct {
	Doctor      string
	OpenedAt    time.Time
	ConfirmedAt time.Time
}

type Recommendation struct {
	ID            uuid.UUID
	Position      int
	ServiceCode   string
	ServiceName   string
	Importance    string
	Rationale     string
	GuidelineRef  string
	PatientText   string
	Mark          string
	RejectReason  string
	RejectComment string
}

type Booking struct {
	ID               uuid.UUID
	RecommendationID uuid.UUID
	AppointmentID    string
	ScheduledAt      time.Time
	Channel          string
	CreatedAt        time.Time
}

type Task struct {
	ID         uuid.UUID
	Reason     string
	Status     string
	Assignee   string
	DueAt      time.Time
	CreatedAt  time.Time
	TakenAt    time.Time
	ClosedAt   time.Time
	NextCallAt time.Time
	Attempts   []Attempt
}

type Attempt struct {
	Outcome       string
	DeclineReason string
	Comment       string
	Actor         string
	CreatedAt     time.Time
}

func (c Case) Accepted() []Recommendation {
	accepted := make([]Recommendation, 0, len(c.Recommendations))
	for _, r := range c.Recommendations {
		if r.Mark == "critical" || r.Mark == "minor" {
			accepted = append(accepted, r)
		}
	}
	return accepted
}

type weighted[T any] struct {
	value  T
	weight int
}

type historyService struct {
	Code         string
	Name         string
	Importance   string
	Rationale    string
	GuidelineRef string
	PatientText  string
}

type modalityProfile struct {
	bodySite    string
	urgencies   []weighted[string]
	conclusions map[string]string
	services    map[string][]historyService
}

type taskOutcome string

const (
	outcomeBooked      taskOutcome = "booked"
	outcomeContacted   taskOutcome = "contacted"
	outcomeDeclined    taskOutcome = "declined"
	outcomeUnreachable taskOutcome = "unreachable"
)
