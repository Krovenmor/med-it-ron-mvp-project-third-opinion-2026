package history

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

const (
	days              = 180
	seed              = 20261005
	quietTime         = 12 * time.Hour
	cutoffMargin      = 10 * time.Minute
	catalogVersion    = "catalog-2026.09"
	guidelinesVersion = "guidelines-2026.09"
	day               = 24 * time.Hour
)

type generator struct {
	rng    *rand.Rand
	cutoff time.Time
	seq    int
}

func Generate(now time.Time) []Case {
	g := &generator{rng: rand.New(rand.NewPCG(seed, seed)), cutoff: now.Add(-cutoffMargin)}
	latest := now.Add(-quietTime)
	start := now.UTC().Truncate(day).Add(-days * day)

	var cases []Case
	for d := range days + 1 {
		date := start.Add(time.Duration(d) * day)
		progress := float64(d) / days
		for range g.casesOn(date) {
			receivedAt := date.Add(5*time.Hour + g.between(0, 12*time.Hour))
			if receivedAt.After(latest) {
				continue
			}
			cases = append(cases, g.caseAt(receivedAt, progress))
		}
	}
	return cases
}

func (g *generator) casesOn(date time.Time) int {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return 2 + g.rng.IntN(4)
	}
	return 9 + g.rng.IntN(7)
}

func (g *generator) caseAt(receivedAt time.Time, progress float64) Case {
	g.seq++
	modality := pick(g.rng, modalities)
	profile := profiles[modality]
	urgency := pick(g.rng, profile.urgencies)
	c := Case{
		ID:           uuid.New(),
		SourceSystem: SourceSystem,
		Patient:      g.patient(modality, receivedAt),
		Study: Study{
			ID:          fmt.Sprintf("HIST-%05d", g.seq),
			Modality:    modality,
			BodySite:    profile.bodySite,
			PerformedAt: receivedAt.Add(-g.between(time.Hour, 6*time.Hour)),
		},
		Conclusion:        profile.conclusions[urgency],
		Urgency:           urgency,
		CatalogVersion:    catalogVersion,
		GuidelinesVersion: guidelinesVersion,
		ReceivedAt:        receivedAt,
	}
	t := &timeline{g: g, progress: progress, c: &c}
	t.run(profile.services[urgency])
	return c
}

func (g *generator) patient(modality string, receivedAt time.Time) Patient {
	sex := "female"
	if modality != "MG" && g.chance(0.5) {
		sex = "male"
	}
	first, last := pickOne(g.rng, maleNames), pickOne(g.rng, surnames)
	if sex == "female" {
		first, last = pickOne(g.rng, femaleNames), last+"а"
	}
	age := 30 + g.rng.IntN(46)
	return Patient{
		ID:         uuid.New(),
		ExternalID: fmt.Sprintf("H-%05d", g.seq),
		FullName:   first + " " + last,
		BirthDate:  receivedAt.UTC().Truncate(day).AddDate(-age, -g.rng.IntN(12), -g.rng.IntN(28)),
		Sex:        sex,
		Phone:      fmt.Sprintf("+7 9%02d %03d-%02d-%02d", g.rng.IntN(100), g.rng.IntN(1000), g.rng.IntN(100), g.rng.IntN(100)),
	}
}

func (g *generator) mark(position int, urgency string) (string, string, string) {
	switch {
	case urgency == "normal":
		return "minor", "", ""
	case position == 0:
		return g.markWith([]weighted[string]{{"critical", 88}, {"minor", 7}, {"rejected", 5}})
	}
	return g.markWith([]weighted[string]{{"minor", 60}, {"critical", 20}, {"rejected", 20}})
}

func (g *generator) markWith(options []weighted[string]) (string, string, string) {
	mark := pick(g.rng, options)
	if mark != "rejected" {
		return mark, "", ""
	}
	if g.chance(0.5) {
		return mark, "other", pickOne(g.rng, rejectComments)
	}
	return mark, "contraindicated", ""
}

func (g *generator) happened(at time.Time) bool {
	return !at.After(g.cutoff)
}

func (g *generator) chance(p float64) bool {
	return g.rng.Float64() < p
}

func (g *generator) between(from, to time.Duration) time.Duration {
	return from + time.Duration(g.rng.Int64N(int64(to-from)+1))
}

func (g *generator) share(d time.Duration, from, to float64) time.Duration {
	return time.Duration(float64(d) * (from + g.rng.Float64()*(to-from)))
}

func (g *generator) appointmentAfter(at time.Time) time.Time {
	date := at.UTC().Truncate(day).Add(time.Duration(1+g.rng.IntN(12)) * day)
	return date.Add(6*time.Hour + time.Duration(g.rng.IntN(9))*time.Hour)
}

func pick[T any](rng *rand.Rand, options []weighted[T]) T {
	total := 0
	for _, o := range options {
		total += o.weight
	}
	n := rng.IntN(total)
	for _, o := range options {
		if n < o.weight {
			return o.value
		}
		n -= o.weight
	}
	return options[len(options)-1].value
}

func pickOne[T any](rng *rand.Rand, options []T) T {
	return options[rng.IntN(len(options))]
}

func reviewSLA(urgency string) time.Duration {
	switch urgency {
	case "emergency":
		return 30 * time.Minute
	case "priority":
		return 4 * time.Hour
	}
	return day
}

func followUpDelay(urgency string) time.Duration {
	if urgency == "planned" {
		return 7 * day
	}
	return day
}

func taskDueAt(reason, urgency string, created time.Time) time.Time {
	switch {
	case reason == "emergency":
		return created.Add(time.Hour)
	case urgency == "planned":
		return addBusinessDays(created, 2)
	}
	return created.Add(4 * time.Hour)
}

func addBusinessDays(t time.Time, n int) time.Time {
	for n > 0 {
		t = t.Add(day)
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			n--
		}
	}
	return t
}

type timeline struct {
	g        *generator
	progress float64
	c        *Case
}

func (t *timeline) run(services []historyService) {
	c := t.c
	c.AssessedAt = c.ReceivedAt.Add(t.g.between(time.Minute, 4*time.Minute))
	c.ReviewDueAt = c.AssessedAt.Add(reviewSLA(c.Urgency))
	for i, s := range services {
		mark, reason, comment := t.g.mark(i, c.Urgency)
		c.Recommendations = append(c.Recommendations, Recommendation{
			ID:            uuid.New(),
			Position:      i + 1,
			ServiceCode:   s.Code,
			ServiceName:   s.Name,
			Importance:    s.Importance,
			Rationale:     s.Rationale,
			GuidelineRef:  s.GuidelineRef,
			PatientText:   s.PatientText,
			Mark:          mark,
			RejectReason:  reason,
			RejectComment: comment,
		})
	}
	t.review()

	notifiedAt := c.Review.ConfirmedAt.Add(t.g.between(2*time.Second, 20*time.Second))
	if !t.g.happened(notifiedAt) {
		return
	}
	c.NotifiedAt = notifiedAt

	switch {
	case len(c.Accepted()) == 0 || c.Urgency == "normal":
	case c.Urgency == "emergency":
		t.resolveEmergency()
	default:
		t.followUp()
	}
}

func (t *timeline) review() {
	c := t.c
	sla := reviewSLA(c.Urgency)
	took := t.g.share(sla, 1.05, 2.2)
	if t.g.chance(reviewOnTime[c.Urgency] + 0.06*t.progress) {
		took = t.g.share(sla, 0.15, 0.95)
	}
	confirmedAt := c.AssessedAt.Add(took)
	if confirmedAt.After(t.g.cutoff) {
		confirmedAt = t.g.cutoff
	}
	c.Review = Review{
		Doctor:      pickOne(t.g.rng, doctors),
		OpenedAt:    c.AssessedAt.Add(confirmedAt.Sub(c.AssessedAt) * 3 / 10),
		ConfirmedAt: confirmedAt,
	}
}

func (t *timeline) resolveEmergency() {
	c := t.c
	task := Task{ID: uuid.New(), Reason: "emergency", CreatedAt: c.AssessedAt, DueAt: taskDueAt("emergency", c.Urgency, c.AssessedAt)}
	window := task.DueAt.Sub(task.CreatedAt)
	contactAt := task.CreatedAt.Add(t.g.share(window, 1.1, 3))
	if t.g.chance(0.82 + 0.08*t.progress) {
		contactAt = task.CreatedAt.Add(t.g.share(window, 0.15, 0.9))
	}
	if earliest := c.NotifiedAt.Add(time.Minute); contactAt.Before(earliest) {
		contactAt = earliest
	}

	outcome := pick(t.g.rng, emergencyOutcomes)
	if !t.workTask(task, outcome, contactAt) || outcome != outcomeContacted {
		return
	}
	if t.g.chance(0.7) {
		t.book("self", contactAt.Add(t.g.between(time.Hour, 8*time.Hour)))
	}
}

func (t *timeline) followUp() {
	c := t.c
	booked := t.g.chance(bookingRate[c.Urgency] + modalityBookingShift[c.Study.Modality] + 0.14*t.progress)
	createdAt := c.NotifiedAt.Add(followUpDelay(c.Urgency))
	if booked && t.g.chance(0.58+0.1*t.progress) {
		t.book("self", c.NotifiedAt.Add(t.g.share(createdAt.Sub(c.NotifiedAt), 0.05, 0.9)))
		return
	}

	task := Task{ID: uuid.New(), Reason: "no_booking", CreatedAt: createdAt, DueAt: taskDueAt("no_booking", c.Urgency, createdAt)}
	window := task.DueAt.Sub(task.CreatedAt)
	closeAt := task.CreatedAt.Add(t.g.share(window, 1.1, 2.5))
	if t.g.chance(0.74 + 0.12*t.progress) {
		closeAt = task.CreatedAt.Add(t.g.share(window, 0.1, 0.9))
	}
	outcome := outcomeBooked
	if !booked {
		outcome = pick(t.g.rng, missedBookingOutcomes)
	}
	t.workTask(task, outcome, closeAt)
}

func (t *timeline) workTask(task Task, outcome taskOutcome, closeAt time.Time) bool {
	retry := 2 * time.Hour
	if task.Reason == "emergency" {
		retry = 15 * time.Minute
	}
	finishedAt := closeAt
	if outcome == outcomeUnreachable {
		finishedAt = closeAt.Add(2*retry + 10*time.Minute)
	}
	if !t.g.happened(finishedAt) {
		return false
	}

	operator := pickOne(t.g.rng, operators)
	task.Assignee = operator
	task.TakenAt = task.CreatedAt.Add(closeAt.Sub(task.CreatedAt) / 4)
	attempt := func(outcome, reason, comment string, at time.Time) {
		task.Attempts = append(task.Attempts, Attempt{Outcome: outcome, DeclineReason: reason, Comment: comment, Actor: operator, CreatedAt: at})
	}

	switch outcome {
	case outcomeBooked:
		if t.g.chance(0.3) {
			missedAt := task.CreatedAt.Add(closeAt.Sub(task.CreatedAt) / 2)
			attempt("no_answer", "", "", missedAt)
		}
		t.book("operator", closeAt)
		attempt("booked", "", "", closeAt)
		task.Status, task.ClosedAt = "booked", closeAt
	case outcomeContacted:
		attempt("contacted", "", "", closeAt)
		task.Status, task.ClosedAt = "contacted", closeAt
	case outcomeDeclined:
		reason := pick(t.g.rng, declineReasons)
		comment := ""
		if reason == "other" {
			comment = pickOne(t.g.rng, declineComments)
		}
		attempt("declined", reason, comment, closeAt)
		task.Status, task.ClosedAt = "declined", closeAt
		t.c.DeclinedAt, t.c.DeclinedBy = closeAt, operator
	default:
		last := closeAt
		for i := range 3 {
			last = closeAt.Add(time.Duration(i) * retry)
			attempt("no_answer", "", "", last)
		}
		handedAt := last.Add(10 * time.Minute)
		attempt("handed_to_doctor", "", "", handedAt)
		task.NextCallAt = last.Add(retry)
		task.Status, task.ClosedAt = "handed_to_doctor", handedAt
		t.c.UnreachableAt, t.c.UnreachableBy = last, operator
	}
	t.c.Tasks = append(t.c.Tasks, task)
	return true
}

func (t *timeline) book(channel string, at time.Time) {
	if !t.g.happened(at) {
		return
	}
	accepted := t.c.Accepted()
	recs := accepted[:1]
	if len(accepted) > 1 && t.g.chance(0.35) {
		recs = accepted[:2]
	}
	scheduledAt := t.g.appointmentAfter(at)
	for _, rec := range recs {
		t.c.Bookings = append(t.c.Bookings, Booking{
			ID:               uuid.New(),
			RecommendationID: rec.ID,
			AppointmentID:    "HIST-APT-" + rec.ID.String(),
			ScheduledAt:      scheduledAt,
			Channel:          channel,
			CreatedAt:        at,
		})
	}

	visitedAt := scheduledAt.Add(time.Hour)
	if t.g.happened(visitedAt) && t.g.chance(0.88) {
		t.c.CompletedAt = visitedAt
	}
}
