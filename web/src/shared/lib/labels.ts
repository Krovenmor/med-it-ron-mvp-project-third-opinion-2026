import type { AcceptedMark, DeclineReason, Mark, Modality, RejectReason, Sex, Urgency } from '@/api/models'

export const urgencyOrder: Urgency[] = ['emergency', 'priority', 'planned', 'normal']

export const urgencyLabels: Record<Urgency, string> = {
  emergency: 'Неотложный',
  priority: 'Приоритетный',
  planned: 'Плановый',
  normal: 'Норма',
}

export const modalityLabels: Record<Modality, string> = {
  CT: 'КТ',
  DX: 'Рентгенография',
  MG: 'Маммография',
}

export const sexLabels: Record<Sex, string> = {
  male: 'мужской',
  female: 'женский',
}

export const markLabels: Record<Mark, string> = {
  critical: 'Приоритетно',
  minor: 'Вторично',
  rejected: 'Не нужна',
}

export const acceptedMarks: AcceptedMark[] = ['critical', 'minor']

export const rejectReasonLabels: Record<RejectReason, string> = {
  contraindicated: 'Противопоказано',
  other: 'Другое',
}

export function urgencyRank(urgency: Urgency): number {
  return urgencyOrder.length - urgencyOrder.indexOf(urgency)
}

export const declineReasons: DeclineReason[] = ['expensive', 'far', 'other_clinic', 'not_needed', 'other']

export const declineReasonLabels: Record<DeclineReason, string> = {
  expensive: 'Дорого',
  far: 'Далеко',
  other_clinic: 'Лечится в другой клинике',
  not_needed: 'Не считает нужным',
  other: 'Другое',
}
