import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import type { Dashboard } from '@/api/models'
import { Card } from '@/shared/ui/Layout'

import { axisTick, chartColors, tooltipProps } from './chart'
import { Legend } from './Legend'
import { formatNumber } from './models'

const dayFormat = new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: '2-digit' })

interface TrendCardProps {
  dashboard: Dashboard
  className?: string
}

export function TrendCard({ dashboard, className }: TrendCardProps) {
  const data = dashboard.trend.map((point) => ({
    label: dayFormat.format(new Date(point.from)),
    received: point.received,
    booked: point.booked,
  }))
  const names: Record<string, string> = { received: 'Поступило', booked: 'Записались' }

  return (
    <Card className={className}>
      <div className="mb-4 flex flex-wrap items-baseline justify-between gap-3">
        <h2 className="text-lg font-medium">
          Динамика <span className="text-muted">{dashboard.trend_step === 'week' ? 'по неделям' : 'по дням'}</span>
        </h2>
        <Legend
          items={[
            { color: chartColors.focus, label: 'Записались' },
            { color: chartColors.rest, label: 'Поступило заключений' },
          ]}
        />
      </div>
      <ResponsiveContainer width="100%" height={380}>
        <LineChart data={data} margin={{ top: 8, right: 8, left: -16 }}>
          <CartesianGrid stroke={chartColors.grid} vertical={false} />
          <XAxis dataKey="label" tick={axisTick} tickLine={false} axisLine={{ stroke: chartColors.grid }} minTickGap={16} />
          <YAxis tick={axisTick} tickLine={false} axisLine={false} allowDecimals={false} />
          <Tooltip
            {...tooltipProps}
            cursor={{ stroke: chartColors.grid }}
            formatter={(value, name) => [formatNumber(Number(value)), names[String(name)]]}
          />
          <Line type="monotone" dataKey="received" stroke={chartColors.rest} strokeWidth={2} dot={false} />
          <Line type="monotone" dataKey="booked" stroke={chartColors.focus} strokeWidth={3} dot={false} />
        </LineChart>
      </ResponsiveContainer>
      <p className="mt-2 text-sm text-subtle">
        По дате события: поступление заключения и первая запись по нему. Последняя точка – неполная.
      </p>
    </Card>
  )
}
