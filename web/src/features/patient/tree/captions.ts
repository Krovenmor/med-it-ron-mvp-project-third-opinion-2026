import type { RouteStep, TrunkNode } from '@/api/models'
import { formatDate, formatShortDateTime } from '@/shared/lib/format'

export interface Caption {
  text: string
  alert: boolean
}

export function stepCaption(step: RouteStep): Caption {
  switch (step.status) {
    case 'done':
      return { text: step.late ? 'пройдено позже срока' : 'пройдено', alert: false }
    case 'booked':
      return { text: step.appointment ? `запись ${formatShortDateTime(step.appointment.scheduled_at)}` : 'вы записаны', alert: false }
    case 'overdue':
      return { text: 'срок прошёл', alert: true }
    case 'declined':
      return { text: 'вы отказались', alert: false }
    default:
      return { text: step.due_at ? `до ${formatDate(step.due_at)}` : 'назначено', alert: false }
  }
}

export function trunkCaption(node: TrunkNode): Caption | null {
  if (node.pending) {
    return { text: 'врач готовит план', alert: false }
  }
  if (node.ai_report) {
    return { text: 'заключение ИИ', alert: false }
  }
  return null
}
