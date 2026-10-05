import clsx from 'clsx'
import { Minus, TrendingDown, TrendingUp } from 'lucide-react'

import type { Dashboard } from '@/api/models'

import {
  conversion,
  countDelta,
  formatNumber,
  formatPercent,
  pointsDelta,
  share,
  slaShare,
  type Delta,
} from './models'

const deltaIcons = { up: TrendingUp, down: TrendingDown, flat: Minus }

interface KpiProps {
  value: string
  label: string
  delta: Delta | null
  neutral?: boolean
}

function Kpi({ value, label, delta, neutral = false }: KpiProps) {
  const Icon = delta && deltaIcons[delta.direction]
  return (
    <div className="flex flex-col rounded-card bg-card px-6 py-5">
      <div className="text-[28px] font-medium leading-tight">{value}</div>
      <div className="mt-1 text-sm text-subtle">{label}</div>
      {delta && Icon && (
        <div
          className={clsx(
            'mt-3 flex items-center gap-1.5 text-sm font-medium',
            (neutral || delta.direction === 'flat') && 'text-subtle',
            !neutral && delta.direction === 'up' && 'text-accent-deep',
            !neutral && delta.direction === 'down' && 'text-overdue',
          )}
        >
          <Icon className="size-4" aria-hidden />
          {delta.text}
        </div>
      )}
    </div>
  )
}

interface KpiRowProps {
  dashboard: Dashboard
}

export function KpiRow({ dashboard }: KpiRowProps) {
  const { current, previous } = dashboard
  const selfShare = (stats: typeof current) => share(stats.bookings.self, stats.funnel.booked)
  return (
    <div className="grid grid-cols-5 gap-5">
      <Kpi
        value={formatNumber(current.funnel.received)}
        label="заключений"
        delta={countDelta(current.funnel.received, previous.funnel.received)}
        neutral
      />
      <Kpi
        value={formatPercent(conversion(current.funnel))}
        label="записались после подтверждения"
        delta={pointsDelta(conversion(current.funnel), conversion(previous.funnel))}
      />
      <Kpi
        value={formatPercent(slaShare(current.doctor_sla))}
        label="SLA врача"
        delta={pointsDelta(slaShare(current.doctor_sla), slaShare(previous.doctor_sla))}
      />
      <Kpi
        value={formatPercent(slaShare(current.operator_sla))}
        label="SLA администратора"
        delta={pointsDelta(slaShare(current.operator_sla), slaShare(previous.operator_sla))}
      />
      <Kpi
        value={formatPercent(selfShare(current))}
        label="записались сами в приложении"
        delta={pointsDelta(selfShare(current), selfShare(previous))}
      />
    </div>
  )
}
