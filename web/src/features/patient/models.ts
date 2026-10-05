import { CalendarDays, History, ListChecks, Settings, TreeDeciduous, type LucideIcon } from 'lucide-react'

import type { Clinic, DeclineReason, PatientRoute, RouteStep, StepStatus, TrunkNode } from '@/api/models'

export interface DemoPatient {
  id: string
  name: string
  story: string
}

export const demoPatients: DemoPatient[] = [
  { id: 'P-ANNA', name: 'Анна', story: 'маммография' },
  { id: 'P-IGOR', name: 'Игорь', story: 'КТ, есть пропущенный шаг' },
  { id: 'P-SERGEY', name: 'Сергей', story: 'неотложный' },
  { id: 'P-OLGA', name: 'Ольга', story: 'без отклонений' },
]

export const defaultPatientId = demoPatients[0].id

export interface PatientTab {
  path: string
  label: string
  short: string
  icon: LucideIcon
}

export const patientTabs: PatientTab[] = [
  { path: '', label: 'Дерево', short: 'Дерево', icon: TreeDeciduous },
  { path: 'steps', label: 'Мои шаги', short: 'Шаги', icon: ListChecks },
  { path: 'appointments', label: 'Записи', short: 'Записи', icon: CalendarDays },
  { path: 'history', label: 'История', short: 'История', icon: History },
  { path: 'settings', label: 'Настройки', short: 'Настройки', icon: Settings },
]

export const stepStatusLabels: Record<StepStatus, string> = {
  recommended: 'Назначено',
  booked: 'Вы записаны',
  done: 'Пройдено',
  overdue: 'Срок прошёл',
  declined: 'Вы отказались',
}

export const stepStatusStyles: Record<StepStatus, string> = {
  recommended: 'bg-white text-subtle ring-1 ring-inset ring-line',
  booked: 'bg-priority-bg text-priority-fg',
  done: 'bg-normal-bg text-normal-fg',
  overdue: 'bg-overdue-tint text-overdue',
  declined: 'bg-card text-subtle',
}

export const patientDeclineReasonLabels: Record<DeclineReason, string> = {
  expensive: 'Дорого',
  far: 'Далеко добираться',
  other_clinic: 'Пройду в другой клинике',
  not_needed: 'Не считаю нужным',
  other: 'Другая причина',
}

export type SheetMode = 'details' | 'book'

export type RouteNode = { type: 'trunk'; node: TrunkNode } | { type: 'step'; step: RouteStep }

export interface PatientContext {
  patientId: string
  route: PatientRoute
  now: Date
  openNode: (id: string, mode?: SheetMode) => void
}

export function findNode(route: PatientRoute, id: string): RouteNode | null {
  const step = route.steps.find((item) => item.id === id)
  if (step) {
    return { type: 'step', step }
  }
  const node = route.trunk.find((item) => item.id === id)
  return node ? { type: 'trunk', node } : null
}

export function isOpen(step: RouteStep): boolean {
  return step.kind === 'recommendation' && (step.status === 'recommended' || step.status === 'overdue')
}

export function recommendations(route: PatientRoute): RouteStep[] {
  return route.steps.filter((step) => step.kind === 'recommendation')
}

export function phoneHref(clinic: Clinic): string {
  return `tel:${clinic.phone.replace(/[^\d+]/g, '')}`
}

export function mapsHref(clinic: Clinic): string {
  return `https://yandex.ru/maps/?text=${encodeURIComponent(clinic.address)}`
}
