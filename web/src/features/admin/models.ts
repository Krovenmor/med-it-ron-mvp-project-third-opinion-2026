import type {
  ActiveTaskStatus,
  CallOutcome,
  NotificationChannel,
  NotificationKind,
  NotificationRecipient,
  OperatorTaskItem,
  TaskReason,
  TaskStatus,
  Urgency,
} from '@/api/models'

export type UrgencyFilter = Urgency | 'all'
export type StatusFilter = ActiveTaskStatus | 'all'
export type OwnerFilter = 'all' | 'mine'
export type RecipientFilter = 'all' | 'patient' | 'staff'

export const adminId = import.meta.env.VITE_ADMIN_ID
export const clinicName = import.meta.env.VITE_CLINIC_NAME

export const activeTaskStatuses: ActiveTaskStatus[] = ['new', 'in_progress', 'no_answer', 'callback']

export const taskReasonLabels: Record<TaskReason, string> = {
  emergency: 'Неотложный случай',
  no_booking: 'Не записался вовремя',
  help_request: 'Просит перезвонить',
}

export const taskStatusLabels: Record<TaskStatus, string> = {
  new: 'Новый',
  in_progress: 'В работе',
  no_answer: 'Не дозвонились',
  callback: 'Перезвонить',
  contacted: 'Дозвонились',
  booked: 'Записан',
  declined: 'Отказ',
  handed_to_doctor: 'Передано врачу',
  cancelled: 'Снята',
}

export const taskStatusStyles: Record<TaskStatus, string> = {
  new: 'bg-accent-soft text-accent-deep',
  in_progress: 'bg-white text-ink ring-1 ring-inset ring-line',
  no_answer: 'bg-priority-bg text-priority-fg',
  callback: 'bg-planned-bg text-planned-fg ring-1 ring-inset ring-line',
  contacted: 'bg-normal-bg text-normal-fg',
  booked: 'bg-normal-bg text-normal-fg',
  declined: 'bg-card text-subtle',
  handed_to_doctor: 'bg-emergency-bg text-emergency-fg',
  cancelled: 'bg-card text-subtle',
}

export const outcomeLabels: Record<CallOutcome, string> = {
  no_answer: 'Не дозвонился',
  callback: 'Перезвонить позже',
  contacted: 'Дозвонился',
  declined: 'Отказ',
  booked: 'Записал',
  handed_to_doctor: 'Передал врачу',
}

export const notificationKindLabels: Record<NotificationKind, string> = {
  plan_ready: 'План готов',
  no_findings: 'Без отклонений',
  reminder: 'Напоминание',
  urgent_contact: 'Срочно связаться',
  unreachable: 'Не дозвонились',
  review_escalation: 'Эскалация: кейс не открыт',
  contact_escalation: 'Эскалация: нет контакта',
  handed_to_doctor: 'Передано врачу',
}

export const escalationKinds: NotificationKind[] = ['review_escalation', 'contact_escalation', 'handed_to_doctor']

export const recipientLabels: Record<NotificationRecipient, string> = {
  patient: 'Пациенту',
  duty_doctor: 'Дежурному врачу',
  admin_on_duty: 'Дежурному администратору',
}

export const channelLabels: Record<NotificationChannel, string> = {
  push: 'push',
  sms: 'SMS',
  staff: 'внутреннее',
}

export function matchesFilters(item: OperatorTaskItem, urgency: UrgencyFilter, status: StatusFilter): boolean {
  return (urgency === 'all' || item.case.urgency === urgency) && (status === 'all' || item.task.status === status)
}

export function isOverdue(dueAt: string, now: Date): boolean {
  return Date.parse(dueAt) < now.getTime()
}
