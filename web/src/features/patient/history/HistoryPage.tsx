import clsx from 'clsx'
import { History, Sparkles } from 'lucide-react'

import { formatDate } from '@/shared/lib/format'
import { EmptyState } from '@/shared/ui/States'

import { PageTitle } from '../ui/PageTitle'
import { usePatient } from '../usePatient'

interface HistoryItem {
  id: string
  title: string
  date: string
  tag: string
  ai: boolean
}

export function HistoryPage() {
  const { route, now, openNode } = usePatient()

  const items: HistoryItem[] = [
    ...route.trunk
      .filter((node) => Date.parse(node.date) <= now.getTime())
      .map((node) => ({
        id: node.id,
        title: node.title,
        date: node.date,
        tag: node.kind === 'study' ? 'Исследование' : 'Приём',
        ai: node.ai_report,
      })),
    ...route.steps
      .filter((step) => step.status === 'done')
      .map((step) => ({
        id: step.id,
        title: step.title,
        date: step.date,
        tag: step.kind === 'recommendation' ? 'По плану врача' : 'Приём',
        ai: false,
      })),
  ].sort((a, b) => Date.parse(b.date) - Date.parse(a.date))

  return (
    <div className="flex max-w-3xl flex-col gap-5">
      <PageTitle main="История" rest="пройденные приёмы и обследования" />
      {items.length === 0 ? (
        <EmptyState icon={History} title="История пока пуста" />
      ) : (
        <ul className="divide-y divide-line border-y border-line">
          {items.map((item) => (
            <li key={item.id}>
              <button
                type="button"
                onClick={() => openNode(item.id)}
                className="grid w-full grid-cols-[92px_minmax(0,1fr)] gap-x-4 gap-y-1 py-4 text-left hover:bg-card-nested sm:grid-cols-[120px_minmax(0,1fr)_auto] sm:items-center"
              >
                <span className="text-sm tabular-nums text-subtle sm:text-base">{formatDate(item.date)}</span>
                <span className="font-medium">{item.title}</span>
                <span
                  className={clsx(
                    'col-start-2 inline-flex items-center gap-1 text-sm sm:col-start-3',
                    item.ai ? 'text-accent-deep' : 'text-subtle',
                  )}
                >
                  {item.ai && <Sparkles className="size-3.5" aria-hidden />}
                  {item.tag}
                  {item.ai && ' · заключение ИИ'}
                </span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
