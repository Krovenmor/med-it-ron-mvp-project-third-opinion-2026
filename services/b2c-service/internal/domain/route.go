package domain

import (
	"cmp"
	"math"
	"slices"
	"strconv"
	"time"
)

type StepKind string

const (
	StepKindRecommendation StepKind = "recommendation"
	StepKindVisit          StepKind = "visit"
)

type StepStatus string

const (
	StepStatusRecommended StepStatus = "recommended"
	StepStatusBooked      StepStatus = "booked"
	StepStatusDone        StepStatus = "done"
	StepStatusOverdue     StepStatus = "overdue"
	StepStatusDeclined    StepStatus = "declined"
)

type TrunkKind string

const (
	TrunkKindVisit TrunkKind = "visit"
	TrunkKindStudy TrunkKind = "study"
)

var (
	modalityTitles = map[string]string{
		"CT": "Компьютерная томография",
		"DX": "Рентгенография",
		"MG": "Маммография",
	}
	primaryCareServices = []string{"THER-CONSULT"}
	markWeights         = map[string]float64{"critical": 3, "minor": 1}
)

type Clinic struct {
	Name    string
	Phone   string
	Address string
}

type Route struct {
	Now           time.Time
	PatientName   string
	Clinic        Clinic
	UrgentContact bool
	PendingReview bool
	NoFindings    bool
	Progress      *int
	NextStepID    string
	Trunk         []TrunkNode
	Steps         []Step
	Appointments  []Appointment
}

type TrunkNode struct {
	ID       string
	Kind     TrunkKind
	Title    string
	Date     time.Time
	AIReport bool
	Pending  bool
}

type Step struct {
	ID            string
	Kind          StepKind
	OriginID      string
	CaseID        string
	Title         string
	Text          string
	Mark          string
	Status        StepStatus
	Late          bool
	AssignedAt    time.Time
	DueAt         time.Time
	Date          time.Time
	InClinic      bool
	Bookable      bool
	Appointment   *Appointment
	VisitedAt     time.Time
	DeclineReason string
}

func (s Step) open() bool {
	return s.Kind == StepKindRecommendation && (s.Status == StepStatusRecommended || s.Status == StepStatusOverdue)
}

func BuildRoute(plan Plan, history History, clinic Clinic) Route {
	b := routeBuilder{
		route:   Route{Now: plan.Now, PatientName: plan.PatientName, Clinic: clinic},
		history: history,
		used:    map[string]bool{},
	}
	for _, c := range plan.Cases {
		b.addCase(c)
	}
	b.addHistory()
	b.summarize(plan)
	return b.route
}

type routeBuilder struct {
	route   Route
	history History
	used    map[string]bool
}

func (b *routeBuilder) addCase(c PlanCase) {
	b.route.Trunk = append(b.route.Trunk, TrunkNode{
		ID:       "study:" + c.ID,
		Kind:     TrunkKindStudy,
		Title:    modalityTitles[c.Modality],
		Date:     c.PerformedAt,
		AIReport: true,
		Pending:  c.InReview(),
	})
	b.route.UrgentContact = b.route.UrgentContact || c.UrgentContact
	b.route.PendingReview = b.route.PendingReview || (c.InReview() && !c.UrgentContact)
	for _, rec := range c.Recommendations {
		b.route.Steps = append(b.route.Steps, b.recommendationStep(c, rec))
	}
}

func (b *routeBuilder) recommendationStep(c PlanCase, rec Recommendation) Step {
	step := Step{
		ID:            rec.ID,
		Kind:          StepKindRecommendation,
		OriginID:      "study:" + c.ID,
		CaseID:        c.ID,
		Title:         rec.ServiceName,
		Text:          rec.PatientText,
		Mark:          rec.Mark,
		Status:        StepStatusRecommended,
		AssignedAt:    c.AssignedAt(),
		InClinic:      rec.InClinic(),
		DeclineReason: rec.DeclineReason,
	}
	if !c.UrgentContact {
		step.DueAt = c.BookBy
	}
	switch visit, visited := b.visitFor(rec, c.PerformedAt); {
	case visited:
		step.Status = StepStatusDone
		step.VisitedAt = visit.VisitedAt
		step.Late = !step.DueAt.IsZero() && visit.VisitedAt.After(step.DueAt)
	default:
		if appointment, ok := b.appointmentFor(rec, c.PerformedAt); ok {
			step.Status = StepStatusBooked
			step.Appointment = &appointment
		}
	}
	if step.Status == StepStatusRecommended {
		switch {
		case rec.Declined:
			step.Status = StepStatusDeclined
		case !step.DueAt.IsZero() && step.DueAt.Before(b.route.Now):
			step.Status = StepStatusOverdue
		}
	}
	step.Bookable = step.InClinic && step.open() && !c.UrgentContact
	step.Date = step.position()
	return step
}

func (s Step) position() time.Time {
	switch {
	case s.Status == StepStatusDone:
		return s.VisitedAt
	case s.Appointment != nil:
		return s.Appointment.ScheduledAt
	case !s.DueAt.IsZero():
		return s.DueAt
	}
	return s.AssignedAt
}

func (b *routeBuilder) addHistory() {
	var others []Visit
	for i, v := range b.history.Visits {
		switch {
		case b.used["visit:"+strconv.Itoa(i)]:
		case slices.Contains(primaryCareServices, v.ServiceCode):
			b.route.Trunk = append(b.route.Trunk, TrunkNode{ID: "visit:" + strconv.Itoa(i), Kind: TrunkKindVisit, Title: v.ServiceName, Date: v.VisitedAt})
		default:
			others = append(others, v)
		}
	}
	slices.SortStableFunc(others, func(a, b Visit) int { return a.VisitedAt.Compare(b.VisitedAt) })
	slices.SortStableFunc(b.route.Trunk, func(a, b TrunkNode) int { return a.Date.Compare(b.Date) })

	for i, v := range others {
		origin, ok := b.trunkBefore(v.VisitedAt)
		if !ok {
			b.route.Trunk = slices.Insert(b.route.Trunk, 0, TrunkNode{ID: "root:" + strconv.Itoa(i), Kind: TrunkKindVisit, Title: v.ServiceName, Date: v.VisitedAt})
			continue
		}
		b.route.Steps = append(b.route.Steps, Step{
			ID:        "history:" + strconv.Itoa(i),
			Kind:      StepKindVisit,
			OriginID:  origin,
			Title:     v.ServiceName,
			Status:    StepStatusDone,
			Date:      v.VisitedAt,
			VisitedAt: v.VisitedAt,
		})
	}

	for _, a := range b.history.Appointments {
		if a.ScheduledAt.After(b.route.Now) {
			b.route.Appointments = append(b.route.Appointments, a)
		}
	}
	slices.SortStableFunc(b.route.Appointments, func(a, b Appointment) int { return a.ScheduledAt.Compare(b.ScheduledAt) })
}

func (b *routeBuilder) trunkBefore(at time.Time) (string, bool) {
	origin := ""
	for _, n := range b.route.Trunk {
		if n.Date.After(at) {
			break
		}
		origin = n.ID
	}
	return origin, origin != ""
}

func (b *routeBuilder) summarize(plan Plan) {
	if latest, ok := latestConfirmed(plan); ok {
		b.route.NoFindings = latest.NoFindings
	}
	b.route.Progress = progress(b.route.Steps, b.route.Now)
	b.route.NextStepID = nextStep(b.route.Steps)
}

func latestConfirmed(plan Plan) (PlanCase, bool) {
	var latest PlanCase
	found := false
	for _, c := range plan.Cases {
		if c.InReview() || c.UrgentContact {
			continue
		}
		if !found || c.PerformedAt.After(latest.PerformedAt) {
			latest, found = c, true
		}
	}
	return latest, found
}

func progress(steps []Step, now time.Time) *int {
	var earned, possible float64
	for _, s := range steps {
		if s.Kind != StepKindRecommendation {
			continue
		}
		weight := markWeights[s.Mark]
		due := !s.DueAt.IsZero() && s.DueAt.Before(now)
		switch {
		case s.Status == StepStatusDone && s.Late:
			earned += 0.5 * weight
		case s.Status == StepStatusDone:
			earned += weight
		case s.Status == StepStatusDeclined, due:
		default:
			continue
		}
		possible += weight
	}
	if possible == 0 {
		return nil
	}
	percent := int(math.Round(100 * earned / possible))
	return &percent
}

func nextStep(steps []Step) string {
	var open, booked []Step
	for _, s := range steps {
		switch {
		case s.open():
			open = append(open, s)
		case s.Kind == StepKindRecommendation && s.Status == StepStatusBooked:
			booked = append(booked, s)
		}
	}
	slices.SortStableFunc(open, func(a, b Step) int {
		if (a.Status == StepStatusOverdue) != (b.Status == StepStatusOverdue) {
			if a.Status == StepStatusOverdue {
				return -1
			}
			return 1
		}
		return cmp.Compare(a.DueAt.UnixNano(), b.DueAt.UnixNano())
	})
	if len(open) > 0 {
		return open[0].ID
	}
	slices.SortStableFunc(booked, func(a, b Step) int { return a.Date.Compare(b.Date) })
	if len(booked) > 0 {
		return booked[0].ID
	}
	return ""
}

func (b *routeBuilder) visitFor(rec Recommendation, studyDate time.Time) (Visit, bool) {
	if !rec.InClinic() {
		return Visit{}, false
	}
	for i, v := range b.history.Visits {
		key := "visit:" + strconv.Itoa(i)
		if !b.used[key] && v.ServiceCode == rec.ServiceCode && v.VisitedAt.After(studyDate) {
			b.used[key] = true
			return v, true
		}
	}
	return Visit{}, false
}

func (b *routeBuilder) appointmentFor(rec Recommendation, studyDate time.Time) (Appointment, bool) {
	if a, ok := b.history.AppointmentByReferral(rec.ID); ok && !b.used["appointment:"+a.ID] {
		b.used["appointment:"+a.ID] = true
		return a, true
	}
	if !rec.InClinic() {
		return Appointment{}, false
	}
	for _, a := range b.history.Appointments {
		key := "appointment:" + a.ID
		if !b.used[key] && a.ReferralID == "" && a.ServiceCode == rec.ServiceCode && a.ScheduledAt.After(studyDate) {
			b.used[key] = true
			return a, true
		}
	}
	return Appointment{}, false
}
