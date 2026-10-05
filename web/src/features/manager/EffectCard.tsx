import { useState } from 'react'

import type { Dashboard } from '@/api/models'
import { Card } from '@/shared/ui/Layout'

import { conversion, formatMoney, formatNumber, formatPercent, share } from './models'

interface NumberFieldProps {
  label: string
  suffix: string
  value: number
  onChange: (value: number) => void
}

function NumberField({ label, suffix, value, onChange }: NumberFieldProps) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="text-sm text-subtle">{label}</span>
      <span className="flex items-center gap-2 rounded-control border border-line bg-white px-3 focus-within:border-accent">
        <input
          type="number"
          min={0}
          value={value}
          onChange={(event) => onChange(Math.max(0, Number(event.target.value)))}
          className="h-10 w-full bg-transparent outline-none"
        />
        <span className="text-sm text-subtle">{suffix}</span>
      </span>
    </label>
  )
}

interface MeasuredProps {
  label: string
  value: string
}

function Measured({ label, value }: MeasuredProps) {
  return (
    <div className="flex flex-col gap-1.5">
      <span className="text-sm text-subtle">{label}</span>
      <span className="flex h-10 items-center rounded-control bg-white/60 px-3 font-medium">{value}</span>
    </div>
  )
}

interface EffectCardProps {
  dashboard: Dashboard
  className?: string
}

export function EffectCard({ dashboard, className }: EffectCardProps) {
  const [studies, setStudies] = useState(3000)
  const [baseline, setBaseline] = useState(35)
  const [check, setCheck] = useState(12000)

  const { funnel } = dashboard.current
  const findings = share(funnel.recommended, funnel.received) ?? 0
  const withSystem = conversion(funnel) ?? 0
  const uplift = Math.max(0, withSystem - baseline / 100)
  const extraBookings = studies * findings * uplift
  const extraRevenue = extraBookings * check

  return (
    <Card className={className}>
      <h2 className="mb-1 text-lg font-medium">Расчёт эффекта</h2>
      <p className="mb-4 text-sm text-subtle">
        Исследований в месяц × доля с находкой × прирост конверсии × средний чек цепочки услуг
      </p>
      <div className="grid grid-cols-3 gap-4">
        <NumberField label="Исследований в месяц" suffix="шт." value={studies} onChange={setStudies} />
        <NumberField label="Конверсия без системы" suffix="%" value={baseline} onChange={setBaseline} />
        <NumberField label="Средний чек цепочки" suffix="₽" value={check} onChange={setCheck} />
        <Measured label="Доля с находкой" value={formatPercent(findings)} />
        <Measured label="Конверсия с системой" value={formatPercent(withSystem)} />
        <Measured label="Прирост конверсии" value={`${Math.round(uplift * 100)} п. п.`} />
      </div>
      <p className="mt-3 text-sm text-subtle">
        Доля с находкой и конверсия с системой – из дашборда за выбранный период, остальное – допущения клиники.
      </p>
      <div className="mt-5 grid grid-cols-3 gap-4 border-t border-line pt-5">
        <div>
          <div className="text-[28px] font-medium leading-tight">+{formatNumber(extraBookings)}</div>
          <div className="text-sm text-subtle">записей в месяц</div>
        </div>
        <div>
          <div className="text-[28px] font-medium leading-tight text-accent-deep">{formatMoney(extraRevenue)}</div>
          <div className="text-sm text-subtle">дополнительной выручки в месяц</div>
        </div>
        <div>
          <div className="text-[28px] font-medium leading-tight">{formatMoney(extraRevenue * 12)}</div>
          <div className="text-sm text-subtle">в год</div>
        </div>
      </div>
    </Card>
  )
}
