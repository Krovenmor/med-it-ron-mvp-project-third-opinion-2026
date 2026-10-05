package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

const (
	JobKindNotifyPatient   = "notify_patient"
	JobKindRemindPatient   = "remind_patient"
	JobKindFollowUp        = "follow_up"
	JobKindEscalateContact = "escalate_contact"

	ContactEscalationDelay = 2 * time.Hour
	MaxFailedCallAttempts  = 3

	day = 24 * time.Hour
)

var (
	ProtocolJobKinds  = []string{JobKindNotifyPatient, JobKindRemindPatient, JobKindFollowUp, JobKindEscalateContact}
	EmergencyJobKinds = []string{JobKindEscalateContact}
)

func EmergencyJobs(caseID uuid.UUID, now time.Time) []jobs.Job {
	return []jobs.Job{jobs.New(JobKindEscalateContact, caseID.String(), now.Add(ContactEscalationDelay))}
}

func JobsAfterConfirm(r Route, now time.Time) []jobs.Job {
	key := r.CaseID.String()
	planned := []jobs.Job{jobs.New(JobKindNotifyPatient, key, now)}
	switch r.Urgency {
	case UrgencyPriority:
		planned = append(planned, jobs.New(JobKindFollowUp, key, now.Add(day)))
	case UrgencyPlanned:
		planned = append(planned,
			jobs.New(JobKindRemindPatient, key, now.Add(3*day)),
			jobs.New(JobKindRemindPatient, key, now.Add(7*day)),
			jobs.New(JobKindFollowUp, key, now.Add(7*day)),
		)
	}
	return planned
}

func taskDueAt(reason TaskReason, urgency Urgency, now time.Time) time.Time {
	switch {
	case reason == TaskReasonEmergency:
		return now.Add(time.Hour)
	case reason == TaskReasonHelpRequest:
		return now.Add(2 * time.Hour)
	case urgency == UrgencyPlanned:
		return addBusinessDays(now, 2)
	default:
		return now.Add(4 * time.Hour)
	}
}

func callRetryDelay(urgency Urgency) time.Duration {
	if urgency == UrgencyEmergency {
		return 15 * time.Minute
	}
	return 2 * time.Hour
}

func bookByOffset(urgency Urgency) time.Duration {
	switch urgency {
	case UrgencyEmergency:
		return 0
	case UrgencyPriority:
		return 7 * day
	case UrgencyPlanned:
		return 30 * day
	default:
		return 365 * day
	}
}

func addBusinessDays(t time.Time, days int) time.Time {
	for days > 0 {
		t = t.Add(day)
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			days--
		}
	}
	return t
}
