import { urgencyLabels, urgencyOrder } from '@/shared/lib/labels'
import { Segmented, type Option } from '@/shared/ui/Segmented'

import { queueStatusLabels, type StatusFilter, type UrgencyFilter } from './models'

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
