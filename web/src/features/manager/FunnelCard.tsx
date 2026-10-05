import { Bar, BarChart, LabelList, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import type { Dashboard } from '@/api/models'
import { Card } from '@/shared/ui/Layout'

import { chartColors, tooltipProps } from './chart'
import { Legend } from './Legend'
import { formatNumber, formatPercent, funnelSteps, share } from './models'

interface StepTickProps {
  x?: number
  y?: number
  payload?: { value: string; index: number }
  rates: (number | null)[]
}

function StepTick({ x = 0, y = 0, payload, rates }: StepTickProps) {
  if (!payload) {
    return null
  }
  const rate = rates[payload.index]
  return (
    <text x={x - 8} y={y} textAnchor="end" fill={chartColors.compare} fontSize={14}>
      <tspan x={x - 8} dy={-2}>
        {payload.value}
      </tspan>
      {rate !== null && (
        <tspan x={x - 8} dy={17} fill={chartColors.axis} fontSize={12}>
          {formatPercent(rate)} от предыдущего шага
        </tspan>
      )}
    </text>
  )
}

interface FunnelCardProps {
  dashboard: Dashboard
  className?: string
}

export function FunnelCard({ dashboard, className }: FunnelCardProps) {
  const data = funnelSteps.map((step) => ({
    label: step.label,
    current: dashboard.current.funnel[step.key],
    previous: dashboard.previous.funnel[step.key],
  }))
  const rates = data.map((step, i) => (i === 0 ? null : share(step.current, data[i - 1].current)))

  return (
    <Card className={className}>
      <div className="mb-4 flex flex-wrap items-baseline justify-between gap-3">
        <h2 className="text-lg font-medium">Воронка</h2>
        <Legend
          items={[
            { color: chartColors.focus, label: 'Текущий период' },
            { color: chartColors.compare, label: 'Прошлый период' },
          ]}
        />
      </div>
      <ResponsiveContainer width="100%" height={380}>
        <BarChart data={data} layout="vertical" margin={{ left: 8, right: 48 }} barGap={3}>
          <XAxis type="number" hide />
          <YAxis
            type="category"
            dataKey="label"
            width={210}
            axisLine={false}
            tickLine={false}
            tick={<StepTick rates={rates} />}
          />
          <Tooltip
            {...tooltipProps}
            formatter={(value, name) => [formatNumber(Number(value)), name === 'current' ? 'Текущий' : 'Прошлый']}
          />
          <Bar dataKey="current" fill={chartColors.focus} radius={[0, 6, 6, 0]} barSize={20}>
            <LabelList dataKey="current" position="right" fill={chartColors.compare} fontSize={14} formatter={(v) => formatNumber(Number(v))} />
          </Bar>
          <Bar dataKey="previous" fill={chartColors.compare} radius={[0, 4, 4, 0]} barSize={6} />
        </BarChart>
      </ResponsiveContainer>
      <p className="mt-2 text-sm text-subtle">
        Кейсы, поступившие за период. Прошлый период – по состоянию на его конец, чтобы сравнение было честным.
      </p>
    </Card>
  )
}
