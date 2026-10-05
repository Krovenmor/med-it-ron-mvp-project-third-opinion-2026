import type { DashboardFunnel, Ratio } from '@/api/models'
import type { Option } from '@/shared/ui/Segmented'

export type PeriodValue = '7' | '30' | '90'

export const periodOptions: Option<PeriodValue>[] = [
  { value: '7', label: '7 дней' },
  { value: '30', label: '30 дней' },
  { value: '90', label: '90 дней' },
]

export const funnelSteps: { key: keyof DashboardFunnel; label: string }[] = [
  { key: 'received', label: 'Заключения' },
  { key: 'recommended', label: 'С находкой' },
  { key: 'confirmed', label: 'Подтверждено врачом' },
  { key: 'notified', label: 'Уведомлено' },
  { key: 'booked', label: 'Записались' },
  { key: 'completed', label: 'Пришли' },
]

const integer = new Intl.NumberFormat('ru-RU')
const money = new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 })

export function share(part: number, total: number): number | null {
  return total > 0 ? part / total : null
}

export function slaShare(ratio: Ratio): number | null {
  return share(ratio.on_time, ratio.total)
}

export function conversion(funnel: DashboardFunnel): number | null {
  return share(funnel.booked, funnel.confirmed)
}

export function formatPercent(value: number | null): string {
  return value === null ? '—' : `${Math.round(value * 100)}%`
}

export function formatNumber(value: number): string {
  return integer.format(Math.round(value))
}

export function formatMoney(value: number): string {
  return money.format(Math.round(value))
}

export interface Delta {
  text: string
  direction: 'up' | 'down' | 'flat'
}

export function pointsDelta(current: number | null, previous: number | null): Delta | null {
  if (current === null || previous === null) {
    return null
  }
  const points = Math.round((current - previous) * 100)
  return { text: `${points > 0 ? '+' : ''}${points} п. п.`, direction: directionOf(points) }
}

export function countDelta(current: number, previous: number): Delta | null {
  if (previous === 0) {
    return null
  }
  const percent = Math.round(((current - previous) / previous) * 100)
  return { text: `${percent > 0 ? '+' : ''}${percent}%`, direction: directionOf(percent) }
}

function directionOf(value: number): Delta['direction'] {
  if (value > 0) {
    return 'up'
  }
  return value < 0 ? 'down' : 'flat'
}
