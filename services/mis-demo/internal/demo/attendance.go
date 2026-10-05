package demo

import "time"

func Attend(appointments []Appointment, now time.Time) ([]Visit, []Appointment) {
	var visits []Visit
	upcoming := make([]Appointment, 0, len(appointments))
	for _, a := range appointments {
		if a.ScheduledAt.After(now) {
			upcoming = append(upcoming, a)
			continue
		}
		visits = append(visits, Visit{ServiceCode: a.ServiceCode, ServiceName: a.ServiceName, VisitedAt: a.ScheduledAt})
	}
	return visits, upcoming
}
