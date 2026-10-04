package http

import (
	"context"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/scheduling"
)

type Scheduler interface {
	Slots(serviceCode string, now time.Time) []scheduling.Slot
	Book(req scheduling.BookRequest, now time.Time) (demo.Appointment, bool, error)
	AppointmentsOf(patientID string) []demo.Appointment
}

type ReportSender interface {
	SendReport(ctx context.Context, r demo.Report) (b2b.Delivery, error)
}
