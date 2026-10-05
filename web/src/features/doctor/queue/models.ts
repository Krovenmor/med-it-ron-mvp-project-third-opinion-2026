import type { QueueCase, Urgency } from '@/api/models'

export type QueueStatus = 'new' | 'in_work'
export type UrgencyFilter = Urgency | 'all'
export type StatusFilter = QueueStatus | 'all'

export const queueStatusLabels: Record<QueueStatus, string> = {
  new: 'Новый',
  in_work: 'В работе',
}

export function queueStatusOf(item: QueueCase): QueueStatus {
  return item.opened_at || item.recommendations_reviewed > 0 ? 'in_work' : 'new'
}

export function matchesFilters(item: QueueCase, urgency: UrgencyFilter, status: StatusFilter): boolean {
  return (urgency === 'all' || item.urgency === urgency) && (status === 'all' || queueStatusOf(item) === status)
}
