package demo

import (
	"sync"
	"time"
)

const studyHour = 9

type Scenario struct {
	mu     sync.RWMutex
	anchor time.Time
}

func NewScenario(now time.Time) *Scenario {
	s := &Scenario{}
	s.Reset(now)
	return s
}

func (s *Scenario) Reset(now time.Time) {
	local := now.In(clinicZone)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.anchor = time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, clinicZone)
}

func (s *Scenario) History(patientID string) History {
	h := History{Visits: []Visit{}, Appointments: []Appointment{}}
	for _, v := range histories[patientID] {
		h.Visits = append(h.Visits, Visit{ServiceCode: v.serviceCode, ServiceName: v.serviceName, VisitedAt: s.at(v.day, v.hour)})
	}
	return h
}

func (s *Scenario) Reports() []Report {
	out := make([]Report, 0, len(reports))
	for _, r := range reports {
		out = append(out, Report{
			SourceSystem: SourceSystem,
			Study: Study{
				ID:          r.Study.ID,
				Modality:    r.Study.Modality,
				BodySite:    r.Study.BodySite,
				PerformedAt: s.at(r.Study.Day, studyHour),
			},
			Conclusion: r.Conclusion,
			Patient:    r.Patient,
		})
	}
	return out
}

func (s *Scenario) at(day, hour int) time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.anchor.AddDate(0, 0, day).Add(time.Duration(hour) * time.Hour).UTC()
}
