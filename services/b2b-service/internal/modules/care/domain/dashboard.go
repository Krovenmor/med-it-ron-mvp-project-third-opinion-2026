package domain

import (
	"fmt"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

const (
	MaxDashboardDays = 365
	weeklyTrendAfter = 31
)

type Period struct {
	From time.Time
	To   time.Time
	Days int
}

func NewPeriod(now time.Time, days int) (Period, error) {
	if days < 1 || days > MaxDashboardDays {
		return Period{}, apperr.Invalid(fmt.Sprintf("days must be between 1 and %d", MaxDashboardDays))
	}
	from := now.UTC().Truncate(day).Add(-time.Duration(days-1) * day)
	return Period{From: from, To: now, Days: days}, nil
}

func (p Period) Previous() Period {
	return Period{From: p.From.Add(-time.Duration(p.Days) * day), To: p.From, Days: p.Days}
}

func (p Period) Contains(t time.Time) bool {
	return !t.Before(p.From) && t.Before(p.To)
}

func (p Period) TrendStep() time.Duration {
	if p.Days > weeklyTrendAfter {
		return 7 * day
	}
	return day
}

type CaseFact struct {
	Urgency        Urgency
	Modality       Modality
	ReceivedAt     time.Time
	AssessedAt     time.Time
	ReviewDueAt    time.Time
	ConfirmedAt    time.Time
	NotifiedAt     time.Time
	BookedAt       time.Time
	CompletedAt    time.Time
	Accepted       bool
	BookingChannel BookingChannel
}

func (f CaseFact) recommended(asOf time.Time) bool {
	return happenedBy(f.AssessedAt, asOf) && f.Urgency.Valid() && f.Urgency != UrgencyNormal
}

func (f CaseFact) confirmed(asOf time.Time) bool {
	return f.recommended(asOf) && f.Accepted && happenedBy(f.ConfirmedAt, asOf)
}

func (f CaseFact) booked(asOf time.Time) bool {
	return f.confirmed(asOf) && happenedBy(f.BookedAt, asOf)
}

func (f CaseFact) reviewSLA(asOf time.Time) (Ratio, bool) {
	if !happenedBy(f.AssessedAt, asOf) {
		return Ratio{}, false
	}
	return ratioOf(f.ReviewDueAt, doneBy(f.ConfirmedAt, asOf), asOf)
}

type TaskFact struct {
	Urgency   Urgency
	Status    TaskStatus
	DueAt     time.Time
	CreatedAt time.Time
	ClosedAt  time.Time
}

func (f TaskFact) contactSLA(asOf time.Time) (Ratio, bool) {
	if f.Status == TaskStatusCancelled && happenedBy(f.ClosedAt, asOf) {
		return Ratio{}, false
	}
	return ratioOf(f.DueAt, doneBy(f.ClosedAt, asOf), asOf)
}

func happenedBy(t, asOf time.Time) bool {
	return !t.IsZero() && !t.After(asOf)
}

func doneBy(t, asOf time.Time) time.Time {
	if happenedBy(t, asOf) {
		return t
	}
	return time.Time{}
}

type DeclineCount struct {
	Reason DeclineReason
	Count  int
}

type Ratio struct {
	Hit   int
	Total int
}

func (r *Ratio) add(other Ratio) {
	r.Hit += other.Hit
	r.Total += other.Total
}

func ratioOf(deadline, doneAt, now time.Time) (Ratio, bool) {
	switch {
	case !doneAt.IsZero() && !doneAt.After(deadline):
		return Ratio{Hit: 1, Total: 1}, true
	case !doneAt.IsZero() || now.After(deadline):
		return Ratio{Total: 1}, true
	}
	return Ratio{}, false
}

type Funnel struct {
	Received    int
	Recommended int
	Confirmed   int
	Notified    int
	Booked      int
	Completed   int
}

type PeriodStats struct {
	Funnel           Funnel
	DoctorSLA        Ratio
	OperatorSLA      Ratio
	SelfBookings     int
	OperatorBookings int
}

type UrgencySlice struct {
	Urgency     Urgency
	Confirmed   int
	Booked      int
	DoctorSLA   Ratio
	OperatorSLA Ratio
}

type ModalitySlice struct {
	Modality  Modality
	Confirmed int
	Booked    int
}

type TrendPoint struct {
	From     time.Time
	Received int
	Booked   int
}

type Dashboard struct {
	Period         Period
	PreviousPeriod Period
	Current        PeriodStats
	Previous       PeriodStats
	ByUrgency      []UrgencySlice
	ByModality     []ModalitySlice
	DeclineReasons []DeclineCount
	Trend          []TrendPoint
}

var (
	dashboardUrgencies  = []Urgency{UrgencyEmergency, UrgencyPriority, UrgencyPlanned, UrgencyNormal}
	dashboardModalities = []Modality{ModalityCT, ModalityDX, ModalityMG}
)

func BuildDashboard(period Period, now time.Time, cases []CaseFact, tasks []TaskFact, declines []DeclineCount) Dashboard {
	previous := period.Previous()
	d := Dashboard{
		Period:         period,
		PreviousPeriod: previous,
		DeclineReasons: declines,
		Trend:          newTrend(period),
	}
	byUrgency := map[Urgency]*UrgencySlice{}
	for _, u := range dashboardUrgencies {
		d.ByUrgency = append(d.ByUrgency, UrgencySlice{Urgency: u})
	}
	for i := range d.ByUrgency {
		byUrgency[d.ByUrgency[i].Urgency] = &d.ByUrgency[i]
	}
	byModality := map[Modality]*ModalitySlice{}
	for _, m := range dashboardModalities {
		d.ByModality = append(d.ByModality, ModalitySlice{Modality: m})
	}
	for i := range d.ByModality {
		byModality[d.ByModality[i].Modality] = &d.ByModality[i]
	}

	for _, f := range cases {
		d.addTrend(f)
		switch {
		case period.Contains(f.ReceivedAt):
			d.Current.addCase(f, now)
			if slice, ok := byUrgency[f.Urgency]; ok {
				slice.addCase(f, now)
			}
			if slice, ok := byModality[f.Modality]; ok && f.confirmed(now) {
				slice.Confirmed++
				if f.booked(now) {
					slice.Booked++
				}
			}
		case previous.Contains(f.ReceivedAt):
			d.Previous.addCase(f, previous.To)
		}
	}
	for _, f := range tasks {
		switch {
		case period.Contains(f.CreatedAt):
			if sla, counted := f.contactSLA(now); counted {
				d.Current.OperatorSLA.add(sla)
				if slice, ok := byUrgency[f.Urgency]; ok {
					slice.OperatorSLA.add(sla)
				}
			}
		case previous.Contains(f.CreatedAt):
			if sla, counted := f.contactSLA(previous.To); counted {
				d.Previous.OperatorSLA.add(sla)
			}
		}
	}
	return d
}

func (s *PeriodStats) addCase(f CaseFact, asOf time.Time) {
	s.Funnel.Received++
	if sla, counted := f.reviewSLA(asOf); counted {
		s.DoctorSLA.add(sla)
	}
	if !f.recommended(asOf) {
		return
	}
	s.Funnel.Recommended++
	if !f.confirmed(asOf) {
		return
	}
	s.Funnel.Confirmed++
	if happenedBy(f.NotifiedAt, asOf) {
		s.Funnel.Notified++
	}
	if !f.booked(asOf) {
		return
	}
	s.Funnel.Booked++
	if f.BookingChannel == BookingChannelSelf {
		s.SelfBookings++
	} else {
		s.OperatorBookings++
	}
	if happenedBy(f.CompletedAt, asOf) {
		s.Funnel.Completed++
	}
}

func (s *UrgencySlice) addCase(f CaseFact, asOf time.Time) {
	if sla, counted := f.reviewSLA(asOf); counted {
		s.DoctorSLA.add(sla)
	}
	if f.confirmed(asOf) {
		s.Confirmed++
		if f.booked(asOf) {
			s.Booked++
		}
	}
}

func newTrend(period Period) []TrendPoint {
	step := period.TrendStep()
	var points []TrendPoint
	for from := period.From; from.Before(period.To); from = from.Add(step) {
		points = append(points, TrendPoint{From: from})
	}
	return points
}

func (d *Dashboard) addTrend(f CaseFact) {
	if i, ok := d.trendIndex(f.ReceivedAt); ok {
		d.Trend[i].Received++
	}
	if i, ok := d.trendIndex(f.BookedAt); ok && f.BookingChannel != "" {
		d.Trend[i].Booked++
	}
}

func (d *Dashboard) trendIndex(t time.Time) (int, bool) {
	if !d.Period.Contains(t) {
		return 0, false
	}
	return int(t.Sub(d.Period.From) / d.Period.TrendStep()), true
}
