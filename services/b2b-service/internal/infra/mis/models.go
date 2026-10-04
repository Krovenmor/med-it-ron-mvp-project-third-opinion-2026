package mis

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type historyResponse struct {
	Visits       []visitDTO       `json:"visits"`
	Appointments []appointmentDTO `json:"appointments"`
}

type visitDTO struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type appointmentDTO struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func (r historyResponse) toDomain() domain.PatientHistory {
	h := domain.PatientHistory{
		Visits:       make([]domain.Visit, 0, len(r.Visits)),
		Appointments: make([]domain.Appointment, 0, len(r.Appointments)),
	}
	for _, v := range r.Visits {
		h.Visits = append(h.Visits, v.toDomain())
	}
	for _, a := range r.Appointments {
		h.Appointments = append(h.Appointments, a.toDomain())
	}
	return h
}

func (v visitDTO) toDomain() domain.Visit {
	return domain.Visit{ServiceCode: v.ServiceCode, ServiceName: v.ServiceName, VisitedAt: v.VisitedAt}
}

func (a appointmentDTO) toDomain() domain.Appointment {
	return domain.Appointment{ServiceCode: a.ServiceCode, ServiceName: a.ServiceName, ScheduledAt: a.ScheduledAt}
}
