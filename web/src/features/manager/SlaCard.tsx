import type { Dashboard, Ratio } from '@/api/models'
import { urgencyLabels } from '@/shared/lib/labels'
import { Card } from '@/shared/ui/Layout'

import { formatPercent, slaShare } from './models'

function SlaBar({ ratio }: { ratio: Ratio }) {
  const value = slaShare(ratio)
  if (value === null) {
    return <span className="text-sm text-muted">нет задач</span>
  }
  return (
    <div className="flex items-center gap-3">
      <div className="h-2.5 flex-1 overflow-hidden rounded-pill bg-white">
        <div className="h-full rounded-pill bg-accent" style={{ width: `${value * 100}%` }} />
      </div>
      <span className="w-11 text-right font-medium">{formatPercent(value)}</span>
    </div>
  )
}

interface SlaCardProps {
  dashboard: Dashboard
  className?: string
}

export function SlaCard({ dashboard, className }: SlaCardProps) {
  const { current } = dashboard
  return (
    <Card className={className}>
      <h2 className="mb-4 text-lg font-medium">Выполнено вовремя</h2>
      <div className="grid grid-cols-[140px_1fr_1fr] items-center gap-x-6 gap-y-4">
        <span />
        <div>
          <div className="text-[28px] font-medium leading-tight">{formatPercent(slaShare(current.doctor_sla))}</div>
          <div className="text-sm text-subtle">врач: от оценки ИИ до подтверждения</div>
        </div>
        <div>
          <div className="text-[28px] font-medium leading-tight">{formatPercent(slaShare(current.operator_sla))}</div>
          <div className="text-sm text-subtle">администратор: задача закрыта в срок</div>
        </div>
        {dashboard.by_urgency.map((slice) => (
          <div key={slice.urgency} className="contents">
            <span className="text-sm text-subtle">{urgencyLabels[slice.urgency]}</span>
            <SlaBar ratio={slice.doctor_sla} />
            <SlaBar ratio={slice.operator_sla} />
          </div>
        ))}
      </div>
      <p className="mt-4 text-sm text-subtle">
        Срок проверки врачом: неотложный – 30 минут, приоритетный – 4 часа, остальные – 24 часа. Срок задачи
        администратора: неотложный – 1 час, приоритетный – 4 часа, плановый – 2 рабочих дня.
      </p>
    </Card>
  )
}
