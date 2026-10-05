import { useState } from 'react'
import { Bar, BarChart, LabelList, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import type { Dashboard } from '@/api/models'
import { modalityLabels, urgencyLabels } from '@/shared/lib/labels'
import { Card } from '@/shared/ui/Layout'
import { Segmented, type Option } from '@/shared/ui/Segmented'

import { axisTick, chartColors, tooltipProps } from './chart'
import { Legend } from './Legend'
import { formatNumber, formatPercent, share } from './models'

type Dimension = 'urgency' | 'modality' | 'channel'

const dimensions: Option<Dimension>[] = [
  { value: 'urgency', label: 'Срочность' },
  { value: 'modality', label: 'Модальность' },
  { value: 'channel', label: 'Канал записи' },
]

interface ConversionRow {
  label: string
  rate: number
  booked: number
  confirmed: number
}

function conversionRows(dashboard: Dashboard, dimension: Exclude<Dimension, 'channel'>): ConversionRow[] {
  const slices =
    dimension === 'urgency'
      ? dashboard.by_urgency
          .filter((slice) => slice.urgency !== 'normal')
          .map((slice) => ({ label: urgencyLabels[slice.urgency], ...slice }))
      : dashboard.by_modality.map((slice) => ({ label: modalityLabels[slice.modality], ...slice }))
  return slices.map((slice) => ({
    label: slice.label,
    rate: Math.round((share(slice.booked, slice.confirmed) ?? 0) * 100),
    booked: slice.booked,
    confirmed: slice.confirmed,
  }))
}

interface SlicesCardProps {
  dashboard: Dashboard
  className?: string
}

export function SlicesCard({ dashboard, className }: SlicesCardProps) {
  const [dimension, setDimension] = useState<Dimension>('urgency')

  return (
    <Card className={className}>
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <h2 className="text-lg font-medium">Разрезы</h2>
        <Segmented options={dimensions} value={dimension} onChange={setDimension} />
      </div>
      {dimension === 'channel' ? <ChannelChart dashboard={dashboard} /> : <ConversionChart rows={conversionRows(dashboard, dimension)} />}
    </Card>
  )
}

function ConversionChart({ rows }: { rows: ConversionRow[] }) {
  return (
    <>
      <p className="mb-2 text-sm text-subtle">Доля записавшихся среди подтверждённых врачом кейсов</p>
      <ResponsiveContainer width="100%" height={260}>
        <BarChart data={rows} layout="vertical" margin={{ right: 120 }}>
          <XAxis type="number" domain={[0, 100]} hide />
          <YAxis type="category" dataKey="label" width={130} tick={axisTick} axisLine={false} tickLine={false} />
          <Tooltip
            {...tooltipProps}
            formatter={(_, __, item) => {
              const row = item.payload as ConversionRow
              return [`${row.rate}% · ${row.booked} из ${row.confirmed}`, 'Записались']
            }}
          />
          <Bar dataKey="rate" fill={chartColors.focus} radius={[0, 6, 6, 0]} barSize={24} background={{ fill: '#f2f4f4', radius: 6 }}>
            <LabelList
              position="right"
              content={({ x = 0, y = 0, width = 0, height = 0, index = 0 }) => {
                const row = rows[index]
                return (
                  <text x={Number(x) + Number(width) + 8} y={Number(y) + Number(height) / 2 + 5} fill={chartColors.compare} fontSize={14}>
                    {row.rate}%{' '}
                    <tspan fill={chartColors.axis} fontSize={12}>
                      {row.booked} из {row.confirmed}
                    </tspan>
                  </text>
                )
              }}
            />
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </>
  )
}

function ChannelChart({ dashboard }: { dashboard: Dashboard }) {
  const { current, previous } = dashboard
  const rows = [
    { label: 'Сами в приложении', current: current.bookings.self, previous: previous.bookings.self },
    { label: 'Через администратора', current: current.bookings.operator, previous: previous.bookings.operator },
  ]
  return (
    <>
      <div className="mb-2 flex flex-wrap items-baseline justify-between gap-3">
        <p className="text-sm text-subtle">
          Записи за период · сами: {formatPercent(share(current.bookings.self, current.funnel.booked))}
        </p>
        <Legend
          items={[
            { color: chartColors.focus, label: 'Текущий период' },
            { color: chartColors.compare, label: 'Прошлый период' },
          ]}
        />
      </div>
      <ResponsiveContainer width="100%" height={260}>
        <BarChart data={rows} layout="vertical" margin={{ right: 56 }} barGap={3}>
          <XAxis type="number" hide />
          <YAxis type="category" dataKey="label" width={190} tick={axisTick} axisLine={false} tickLine={false} />
          <Tooltip
            {...tooltipProps}
            formatter={(value, name) => [formatNumber(Number(value)), name === 'current' ? 'Текущий' : 'Прошлый']}
          />
          <Bar dataKey="current" fill={chartColors.focus} radius={[0, 6, 6, 0]} barSize={24}>
            <LabelList dataKey="current" position="right" fill={chartColors.compare} fontSize={14} />
          </Bar>
          <Bar dataKey="previous" fill={chartColors.compare} radius={[0, 4, 4, 0]} barSize={6} />
        </BarChart>
      </ResponsiveContainer>
    </>
  )
}
