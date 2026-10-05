package mis

import (
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

func historyOf(h domain.PatientHistory) api.History {
	out := api.History{
		Visits:       make([]api.Visit, 0, len(h.Visits)),
		Appointments: make([]api.Appointment, 0, len(h.Appointments)),
	}
	for _, v := range h.Visits {
		out.Visits = append(out.Visits, api.Visit{ServiceCode: v.ServiceCode, ServiceName: v.ServiceName, VisitedAt: v.VisitedAt})
	}
	for _, a := range h.Appointments {
		out.Appointments = append(out.Appointments, appointmentOf(a))
	}
	return out
}

func appointmentOf(a domain.Appointment) api.Appointment {
	return api.Appointment{
		ID:          a.ID,
		ServiceCode: a.ServiceCode,
		ServiceName: a.ServiceName,
		ScheduledAt: a.ScheduledAt,
		ReferralID:  a.ReferralID,
	}
}

func slotsOf(slots []domain.Slot) []api.Slot {
	out := make([]api.Slot, 0, len(slots))
	for _, s := range slots {
		out = append(out, api.Slot{ID: s.ID, ServiceCode: s.ServiceCode, StartsAt: s.StartsAt})
	}
	return out
}
