import clsx from 'clsx'

import { urgencyLabels, urgencyOrder } from '@/shared/lib/labels'

import { queueStatusLabels, type StatusFilter, type UrgencyFilter } from './models'

interface Option<T extends string> {
  value: T
  label: string
}

interface SegmentedProps<T extends string> {
  label: string
  options: Option<T>[]
  value: T
  onChange: (value: T) => void
}

function Segmented<T extends string>({ label, options, value, onChange }: SegmentedProps<T>) {
  return (
    <div className="flex items-center gap-3">
      <span className="text-sm text-subtle">{label}</span>
      <div className="flex gap-1 rounded-pill bg-card p-1">
        {options.map((option) => (
          <button
            key={option.value}
            type="button"
            onClick={() => onChange(option.value)}
            className={clsx(
              'rounded-pill px-3 py-1 text-sm font-medium transition-colors',
              option.value === value ? 'bg-white text-ink shadow-sm' : 'text-subtle hover:text-ink',
            )}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  )
}

const urgencyOptions: Option<UrgencyFilter>[] = [
  { value: 'all', label: 'Все' },
  ...urgencyOrder.map((urgency) => ({ value: urgency, label: urgencyLabels[urgency] })),
]

const statusOptions: Option<StatusFilter>[] = [
  { value: 'all', label: 'Все' },
  { value: 'new', label: queueStatusLabels.new },
  { value: 'in_work', label: queueStatusLabels.in_work },
]

interface QueueFiltersProps {
  urgency: UrgencyFilter
  status: StatusFilter
  onUrgencyChange: (value: UrgencyFilter) => void
  onStatusChange: (value: StatusFilter) => void
}

export function QueueFilters({ urgency, status, onUrgencyChange, onStatusChange }: QueueFiltersProps) {
  return (
    <div className="flex flex-wrap gap-x-8 gap-y-3">
      <Segmented label="Срочность" options={urgencyOptions} value={urgency} onChange={onUrgencyChange} />
      <Segmented label="Статус" options={statusOptions} value={status} onChange={onStatusChange} />
    </div>
  )
}
