import { Bar, BarChart, Cell, LabelList, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import type { Dashboard } from '@/api/models'
import { declineReasonLabels } from '@/shared/lib/labels'
import { Card } from '@/shared/ui/Layout'
import { EmptyState } from '@/shared/ui/States'

import { axisTick, chartColors, tooltipProps } from './chart'
import { formatPercent, share } from './models'

interface DeclinesCardProps {
  dashboard: Dashboard
  className?: string
}

export function DeclinesCard({ dashboard, className }: DeclinesCardProps) {
  const total = dashboard.decline_reasons.reduce((sum, item) => sum + item.count, 0)
  const rows = dashboard.decline_reasons.map((item) => ({
    label: declineReasonLabels[item.reason],
    count: item.count,
    share: formatPercent(share(item.count, total)),
  }))

  return (
    <Card className={className}>
      <div className="mb-4 flex items-baseline justify-between gap-3">
        <h2 className="text-lg font-medium">Причины отказов</h2>
        <span className="text-sm text-subtle">всего {total}</span>
      </div>
      {rows.length === 0 ? (
        <EmptyState title="Отказов за период нет" />
      ) : (
        <ResponsiveContainer width="100%" height={Math.max(160, rows.length * 48)}>
          <BarChart data={rows} layout="vertical" margin={{ right: 72 }}>
            <XAxis type="number" hide />
            <YAxis type="category" dataKey="label" width={210} tick={axisTick} axisLine={false} tickLine={false} />
            <Tooltip {...tooltipProps} formatter={(value) => [value, 'Отказов']} />
            <Bar dataKey="count" radius={[0, 6, 6, 0]} barSize={22}>
              {rows.map((row, i) => (
                <Cell key={row.label} fill={i === 0 ? chartColors.focus : chartColors.rest} />
              ))}
              <LabelList
                position="right"
                content={({ x = 0, y = 0, width = 0, height = 0, index = 0 }) => (
                  <text x={Number(x) + Number(width) + 8} y={Number(y) + Number(height) / 2 + 5} fill={chartColors.compare} fontSize={14}>
                    {rows[index].count}{' '}
                    <tspan fill={chartColors.axis} fontSize={12}>
                      {rows[index].share}
                    </tspan>
                  </text>
                )}
              />
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      )}
    </Card>
  )
}
