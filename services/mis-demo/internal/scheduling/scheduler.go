package scheduling

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

const slotIDTimestamp = "200601021504"

var (
	clinicZone = time.FixedZone("MSK", 3*60*60)
	slotHours  = []int{9, 11, 14, 16}
)

type Scheduler struct {
	mu           sync.Mutex
	seq          int
	slotsCount   int
	appointments map[string]demo.Appointment
}

func New(slotsCount int) *Scheduler {
	return &Scheduler{slotsCount: slotsCount, appointments: map[string]demo.Appointment{}}
}

func (s *Scheduler) Slots(serviceCode string, now time.Time) []Slot {
	s.mu.Lock()
	defer s.mu.Unlock()

	slots := make([]Slot, 0, s.slotsCount)
	for startsAt := nextWorkingSlot(now); len(slots) < s.slotsCount; startsAt = nextWorkingSlot(startsAt) {
		id := slotID(serviceCode, startsAt)
		if _, taken := s.appointments[id]; !taken {
			slots = append(slots, Slot{ID: id, ServiceCode: serviceCode, StartsAt: startsAt})
		}
	}
	return slots
}

func (s *Scheduler) Book(req BookRequest, now time.Time) (demo.Appointment, bool, error) {
	serviceCode, startsAt, err := parseSlotID(req.SlotID)
	if err != nil {
		return demo.Appointment{}, false, err
	}
	switch {
	case req.PatientID == "":
		return demo.Appointment{}, false, fmt.Errorf("%w: patient_id is required", ErrInvalidRequest)
	case req.ServiceName == "":
		return demo.Appointment{}, false, fmt.Errorf("%w: service_name is required", ErrInvalidRequest)
	case !isWorkingSlot(startsAt) || !startsAt.After(now):
		return demo.Appointment{}, false, fmt.Errorf("%w: slot %s is not available", ErrInvalidRequest, req.SlotID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, taken := s.appointments[req.SlotID]; taken {
		if existing.PatientID == req.PatientID && existing.ReferralID == req.ReferralID {
			return existing, false, nil
		}
		return demo.Appointment{}, false, ErrSlotTaken
	}

	s.seq++
	appointment := demo.Appointment{
		ID:          fmt.Sprintf("APT-%04d", s.seq),
		PatientID:   req.PatientID,
		ServiceCode: serviceCode,
		ServiceName: req.ServiceName,
		ScheduledAt: startsAt,
		ReferralID:  req.ReferralID,
	}
	s.appointments[req.SlotID] = appointment
	return appointment, true, nil
}

func (s *Scheduler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq = 0
	s.appointments = map[string]demo.Appointment{}
}

func (s *Scheduler) AppointmentsOf(patientID string) []demo.Appointment {
	s.mu.Lock()
	defer s.mu.Unlock()

	var found []demo.Appointment
	for _, a := range s.appointments {
		if a.PatientID == patientID {
			found = append(found, a)
		}
	}
	return found
}

func slotID(serviceCode string, startsAt time.Time) string {
	return serviceCode + "." + startsAt.In(clinicZone).Format(slotIDTimestamp)
}

func parseSlotID(id string) (string, time.Time, error) {
	sep := strings.LastIndex(id, ".")
	if sep <= 0 {
		return "", time.Time{}, fmt.Errorf("%w: malformed slot_id %q", ErrInvalidRequest, id)
	}
	startsAt, err := time.ParseInLocation(slotIDTimestamp, id[sep+1:], clinicZone)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: malformed slot_id %q", ErrInvalidRequest, id)
	}
	return id[:sep], startsAt, nil
}

func nextWorkingSlot(after time.Time) time.Time {
	start := after.In(clinicZone)
	for day := 0; ; day++ {
		date := start.AddDate(0, 0, day)
		for _, hour := range slotHours {
			candidate := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, clinicZone)
			if candidate.After(after) && isWorkingSlot(candidate) {
				return candidate
			}
		}
	}
}

func isWorkingSlot(t time.Time) bool {
	t = t.In(clinicZone)
	weekend := t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
	return !weekend && t.Minute() == 0 && slices.Contains(slotHours, t.Hour())
}
