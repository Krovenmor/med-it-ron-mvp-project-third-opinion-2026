package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Bookings struct {
	postgres.Executor
	q queries.Bookings
}

func NewBookings(db *postgres.Database, q queries.Queries) *Bookings {
	return &Bookings{Executor: db.Executor(), q: q.Bookings}
}

func (r *Bookings) CreateIfAbsent(ctx context.Context, b domain.Booking) (domain.Booking, bool, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.InsertIfAbsent, insertBookingArgs(b))
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("insert booking: %w", err)
	}
	saved, created, err := postgres.CollectOne(rows, bookingRow.toDomain)
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("scan booking: %w", err)
	}
	if created {
		return saved, true, nil
	}

	rows, err = r.DB(ctx).Query(ctx, r.q.GetByAppointment, pgx.StrictNamedArgs{"appointment_id": b.AppointmentID})
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("get booking: %w", err)
	}
	existing, found, err := postgres.CollectOne(rows, bookingRow.toDomain)
	switch {
	case err != nil:
		return domain.Booking{}, false, fmt.Errorf("scan existing booking: %w", err)
	case !found:
		return domain.Booking{}, false, fmt.Errorf("booking %s: %w", b.AppointmentID, apperr.ErrNotFound)
	}
	return existing, false, nil
}

func (r *Bookings) BookedRecommendationIDs(ctx context.Context, caseID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListRecommendationIDsByCase, pgx.StrictNamedArgs{"case_id": caseID})
	if err != nil {
		return nil, fmt.Errorf("list booked recommendations: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, fmt.Errorf("scan booked recommendations: %w", err)
	}
	return ids, nil
}
