import clsx from 'clsx'
import { useState } from 'react'

import type { Slot } from '@/api/models'
import { useSlots } from '@/api/queries'
import { formatDate, formatTime } from '@/shared/lib/format'
import { EmptyState, ErrorState, LoadingState } from '@/shared/ui/States'

const weekdayFormat = new Intl.DateTimeFormat('ru-RU', { weekday: 'short' })
const dayFormat = new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'short' })

interface SlotDay {
  key: string
  date: Date
  slots: Slot[]
}

interface SlotPickerProps {
  patientId: string
  recommendationId: string
  selected: Slot | null
  onSelect: (slot: Slot) => void
}

export function SlotPicker({ patientId, recommendationId, selected, onSelect }: SlotPickerProps) {
  const query = useSlots(patientId, recommendationId)
  const [pickedDay, setPickedDay] = useState<string | null>(selected ? formatDate(selected.starts_at) : null)

  if (query.isPending) {
    return <LoadingState text="Ищем свободное время" />
  }
  if (query.isError) {
    return <ErrorState title="Не удалось загрузить свободное время" onRetry={() => void query.refetch()} />
  }

  const days = groupByDay(query.data.slots)
  if (days.length === 0) {
    return (
      <EmptyState title="Свободного времени пока нет" text="Позвоните в клинику: администратор подберёт удобное время." />
    )
  }
  const day = days.find((item) => item.key === pickedDay) ?? days[0]

  return (
    <div className="flex flex-col gap-4">
      <p className="text-subtle">Выберите удобные день и время</p>
      <div className="-mx-5 flex gap-2 overflow-x-auto px-5 pb-1 lg:-mx-7 lg:px-7" role="tablist" aria-label="День">
        {days.map((item) => (
          <button
            key={item.key}
            type="button"
            role="tab"
            aria-selected={item.key === day.key}
            onClick={() => setPickedDay(item.key)}
            className={clsx(
              'flex min-w-16 shrink-0 flex-col items-center rounded-control border px-3 py-2 transition-colors',
              item.key === day.key ? 'border-accent-deep bg-accent-tint' : 'border-line hover:bg-card',
            )}
          >
            <span className="text-xs uppercase text-subtle">{weekdayFormat.format(item.date)}</span>
            <span className="font-medium">{dayFormat.format(item.date)}</span>
          </button>
        ))}
      </div>
      <div className="grid grid-cols-3 gap-2 sm:grid-cols-4">
        {day.slots.map((slot) => (
          <button
            key={slot.id}
            type="button"
            onClick={() => onSelect(slot)}
            aria-pressed={slot.id === selected?.id}
            className={clsx(
              'h-11 rounded-control border font-medium tabular-nums transition-colors',
              slot.id === selected?.id ? 'border-accent-deep bg-accent text-ink' : 'border-line hover:bg-card',
            )}
          >
            {formatTime(slot.starts_at)}
          </button>
        ))}
      </div>
    </div>
  )
}

function groupByDay(slots: Slot[]): SlotDay[] {
  const days: SlotDay[] = []
  for (const slot of slots) {
    const key = formatDate(slot.starts_at)
    const day = days.find((item) => item.key === key)
    if (day) {
      day.slots.push(slot)
    } else {
      days.push({ key, date: new Date(slot.starts_at), slots: [slot] })
    }
  }
  return days
}
