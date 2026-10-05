package di

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/mis"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/service/booking"
	httptransport "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/transport/http"
)

var bookingModule = fx.Module("booking",
	fx.Provide(
		func(c *b2b.Client) booking.Plans { return c },
		func(c *mis.Client) booking.MIS { return c },
		fx.Annotate(booking.NewService, fx.As(new(httptransport.Bookings))),
	),
)
