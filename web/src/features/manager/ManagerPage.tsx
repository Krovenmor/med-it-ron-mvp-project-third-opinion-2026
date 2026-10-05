import { useState } from 'react'

import { useDashboard } from '@/api/queries'
import { formatDate } from '@/shared/lib/format'
import { Pill } from '@/shared/ui/Badge'
import { TwoToneHeading } from '@/shared/ui/Layout'
import { Segmented } from '@/shared/ui/Segmented'
import { ErrorState, LoadingState } from '@/shared/ui/States'

import { DeclinesCard } from './DeclinesCard'
import { EffectCard } from './EffectCard'
import { FunnelCard } from './FunnelCard'
import { KpiRow } from './KpiRow'
import { periodOptions, type PeriodValue } from './models'
import { SlaCard } from './SlaCard'
import { SlicesCard } from './SlicesCard'
import { TrendCard } from './TrendCard'

const lastMoment = (end: string) => new Date(Date.parse(end) - 1)

export function ManagerPage() {
  const [days, setDays] = useState<PeriodValue>('30')
  const { data, isPending, isError, error, refetch } = useDashboard(Number(days))

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-3">
            <TwoToneHeading main="Дашборд" rest="воронки" />
            <Pill className="bg-priority-bg text-priority-fg">демо-данные</Pill>
          </div>
          {data && (
            <p className="text-subtle">
              {formatDate(data.period.from)} – {formatDate(data.period.to)} · сравнение с{' '}
              {formatDate(data.previous_period.from)} – {formatDate(lastMoment(data.previous_period.to))}
            </p>
          )}
        </div>
        <Segmented label="Период" options={periodOptions} value={days} onChange={setDays} />
      </div>

      {isPending && <LoadingState />}
      {isError && <ErrorState message={error.message} onRetry={() => void refetch()} />}
      {data && (
        <>
          <KpiRow dashboard={data} />
          <div className="grid grid-cols-12 gap-5">
            <FunnelCard dashboard={data} className="col-span-7" />
            <TrendCard dashboard={data} className="col-span-5" />
            <SlicesCard dashboard={data} className="col-span-6" />
            <SlaCard dashboard={data} className="col-span-6" />
            <DeclinesCard dashboard={data} className="col-span-5" />
            <EffectCard dashboard={data} className="col-span-7" />
          </div>
        </>
      )}
    </div>
  )
}
