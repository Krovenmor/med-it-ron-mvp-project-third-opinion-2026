package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationRecipient string

const (
	RecipientPatient     NotificationRecipient = "patient"
	RecipientDutyDoctor  NotificationRecipient = "duty_doctor"
	RecipientAdminOnDuty NotificationRecipient = "admin_on_duty"
)

type NotificationChannel string

const (
	ChannelPush  NotificationChannel = "push"
	ChannelSMS   NotificationChannel = "sms"
	ChannelStaff NotificationChannel = "staff"
)

type NotificationKind string

const (
	NotificationPlanReady         NotificationKind = "plan_ready"
	NotificationNoFindings        NotificationKind = "no_findings"
	NotificationReminder          NotificationKind = "reminder"
	NotificationUrgentContact     NotificationKind = "urgent_contact"
	NotificationUnreachable       NotificationKind = "unreachable"
	NotificationReviewEscalation  NotificationKind = "review_escalation"
	NotificationContactEscalation NotificationKind = "contact_escalation"
	NotificationHandedToDoctor    NotificationKind = "handed_to_doctor"
)

type Notification struct {
	ID        uuid.UUID
	CaseID    uuid.UUID
	Recipient NotificationRecipient
	Channel   NotificationChannel
	Kind      NotificationKind
	Text      string
	CreatedAt time.Time
}

type NotificationEntry struct {
	Notification Notification
	Patient      Patient
	Urgency      Urgency
}

func PlanReadyNotifications(r Route, clinicPhone string, now time.Time) []Notification {
	switch r.Urgency {
	case UrgencyNormal:
		return []Notification{patientNotification(r.CaseID, ChannelPush, NotificationNoFindings,
			"Результаты обследования готовы: отклонений не найдено. Мы напомним о следующем плановом обследовании.", now)}
	case UrgencyPriority:
		return []Notification{
			patientNotification(r.CaseID, ChannelPush, NotificationPlanReady,
				"Врач подготовил для вас план дальнейших шагов. Откройте приложение, чтобы посмотреть рекомендации и записаться.", now),
			patientNotification(r.CaseID, ChannelSMS, NotificationPlanReady,
				"Клиника: врач подготовил для вас рекомендации. Запишитесь в приложении или по телефону "+clinicPhone+".", now),
		}
	default:
		return []Notification{patientNotification(r.CaseID, ChannelPush, NotificationPlanReady,
			"Врач подготовил для вас план дальнейших шагов. Откройте приложение, чтобы посмотреть рекомендации и записаться.", now)}
	}
}

func ReminderNotification(caseID uuid.UUID, now time.Time) Notification {
	return patientNotification(caseID, ChannelPush, NotificationReminder,
		"Напоминаем: врач рекомендовал вам следующий шаг. Запишитесь в приложении – это займёт минуту.", now)
}

func UrgentContactNotifications(caseID uuid.UUID, clinicPhone string, now time.Time) []Notification {
	text := "Врач просит вас срочно связаться с клиникой: " + clinicPhone + "."
	return []Notification{
		patientNotification(caseID, ChannelPush, NotificationUrgentContact, text, now),
		patientNotification(caseID, ChannelSMS, NotificationUrgentContact, text, now),
	}
}

func UnreachableNotification(caseID uuid.UUID, clinicPhone string, now time.Time) Notification {
	return patientNotification(caseID, ChannelSMS, NotificationUnreachable,
		"Клиника не смогла до вас дозвониться. Пожалуйста, свяжитесь с нами: "+clinicPhone+".", now)
}

func ReviewEscalationNotification(caseID uuid.UUID, patientID string, now time.Time) Notification {
	return staffNotification(caseID, RecipientDutyDoctor, NotificationReviewEscalation,
		"Неотложный случай пациента "+patientID+" не открыт врачом 30 минут. Проверьте кейс.", now)
}

func ContactEscalationNotification(caseID uuid.UUID, patientID string, now time.Time) Notification {
	return staffNotification(caseID, RecipientAdminOnDuty, NotificationContactEscalation,
		"Нет контакта с пациентом "+patientID+" по неотложному случаю 2 часа.", now)
}

func HandedToDoctorNotification(caseID uuid.UUID, patientID, comment string, now time.Time) Notification {
	text := "Администратор передал пациента " + patientID + " врачу."
	if comment != "" {
		text += " Комментарий: " + comment
	}
	return staffNotification(caseID, RecipientDutyDoctor, NotificationHandedToDoctor, text, now)
}

func staffNotification(caseID uuid.UUID, recipient NotificationRecipient, kind NotificationKind, text string, now time.Time) Notification {
	return Notification{
		CaseID:    caseID,
		Recipient: recipient,
		Channel:   ChannelStaff,
		Kind:      kind,
		Text:      text,
		CreatedAt: now,
	}
}

func patientNotification(caseID uuid.UUID, channel NotificationChannel, kind NotificationKind, text string, now time.Time) Notification {
	return Notification{
		CaseID:    caseID,
		Recipient: RecipientPatient,
		Channel:   channel,
		Kind:      kind,
		Text:      text,
		CreatedAt: now,
	}
}
