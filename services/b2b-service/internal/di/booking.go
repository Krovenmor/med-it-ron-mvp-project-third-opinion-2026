package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/clock"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/booking"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/transport/http"
)

var bookingModule = fx.Module("booking",
	fx.Provide(
		func(c *clock.Clock) booking.Clock { return c },
		func(r *postgres.Cases) booking.Cases { return r },
		func(r *postgres.Recommendations) booking.Recommendations { return r },
		func(r *postgres.Bookings) booking.Bookings { return r },
		func(r *postgres.Events) booking.Events { return r },
		fx.Annotate(booking.NewService, fx.As(new(httptransport.Bookings))),
	),
)
