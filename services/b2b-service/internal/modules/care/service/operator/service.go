package operator

import (
	"context"
	"slices"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

const journalLimit = 200

type Service struct {
	tx            trm.Manager
	clock         Clock
	tasks         Tasks
	attempts      Attempts
	routes        Routes
	patients      Patients
	planItems     PlanItems
	bookings      Bookings
	notifications Notifications
	events        Events
	protocol      Protocol
	mis           MIS
	booker        Booker
	clinicPhone   string
}

func NewService(
	tx trm.Manager,
	clock Clock,
	tasks Tasks,
	attempts Attempts,
	routes Routes,
	patients Patients,
	planItems PlanItems,
	bookings Bookings,
	notifications Notifications,
	events Events,
	protocol Protocol,
	mis MIS,
	booker Booker,
	clinic config.Clinic,
) *Service {
	return &Service{
		tx:            tx,
		clock:         clock,
		tasks:         tasks,
		attempts:      attempts,
		routes:        routes,
		patients:      patients,
		planItems:     planItems,
		bookings:      bookings,
		notifications: notifications,
		events:        events,
		protocol:      protocol,
		mis:           mis,
		booker:        booker,
		clinicPhone:   clinic.Phone,
	}
}

func (s *Service) List(ctx context.Context, assignee string) ([]domain.TaskListItem, error) {
	return s.tasks.ListActive(ctx, s.clock.Now(), assignee)
}

func (s *Service) Journal(ctx context.Context) ([]domain.NotificationEntry, error) {
	return s.notifications.List(ctx, journalLimit)
}

func (s *Service) Card(ctx context.Context, taskID uuid.UUID) (Card, error) {
	task, err := s.tasks.Get(ctx, taskID)
	if err != nil {
		return Card{}, err
	}
	route, err := s.routes.Get(ctx, task.CaseID)
	if err != nil {
		return Card{}, err
	}
	patient, err := s.patients.Get(ctx, route.PatientID)
	if err != nil {
		return Card{}, err
	}
	attempts, err := s.attempts.ListByTask(ctx, task.ID)
	if err != nil {
		return Card{}, err
	}
	offers, err := s.offers(ctx, route.CaseID)
	if err != nil {
		return Card{}, err
	}
	return Card{Task: task, Route: route, Patient: patient, Attempts: attempts, Offers: offers}, nil
}

func (s *Service) offers(ctx context.Context, caseID uuid.UUID) ([]Offer, error) {
	items, err := s.planItems.ListByCase(ctx, caseID)
	if err != nil {
		return nil, err
	}
	booked, err := s.bookings.BookedRecommendationIDs(ctx, caseID)
	if err != nil {
		return nil, err
	}

	offers := make([]Offer, 0, len(items))
	for _, item := range items {
		if item.Declined() {
			continue
		}
		offer := Offer{Item: item, Booked: slices.Contains(booked, item.ID)}
		if !offer.Booked && item.ServiceCode != "" {
			offer.Slots, err = s.mis.Slots(ctx, item.ServiceCode, s.clock.Now())
			offer.SlotsUnavailable = err != nil
		}
		offers = append(offers, offer)
	}
	slices.SortStableFunc(offers, func(a, b Offer) int {
		return criticalFirst(a.Item) - criticalFirst(b.Item)
	})
	return offers, nil
}

func criticalFirst(i domain.PlanItem) int {
	if i.Mark == domain.MarkCritical {
		return 0
	}
	return 1
}
