package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Bookings struct {
	executor
	q queries.Bookings
}

func NewBookings(pool *pgxpool.Pool, q queries.Queries) *Bookings {
	return &Bookings{executor: executor{pool}, q: q.Bookings}
}

func (r *Bookings) CreateIfAbsent(ctx context.Context, b domain.Booking) (domain.Booking, bool, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.InsertIfAbsent, insertBookingArgs(b))
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("insert booking: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[bookingRow])
	if err == nil {
		return row.toDomain(), true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Booking{}, false, fmt.Errorf("scan booking: %w", err)
	}

	rows, err = r.db(ctx).Query(ctx, r.q.GetByAppointment, pgx.StrictNamedArgs{"appointment_id": b.AppointmentID})
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("get booking: %w", err)
	}
	row, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[bookingRow])
	if err != nil {
		return domain.Booking{}, false, fmt.Errorf("scan existing booking: %w", err)
	}
	return row.toDomain(), false, nil
}
