import { urgencyLabels, urgencyOrder } from '@/shared/lib/labels'
import { Segmented, type Option } from '@/shared/ui/Segmented'

import { activeTaskStatuses, taskStatusLabels, type OwnerFilter, type StatusFilter, type UrgencyFilter } from '../models'

const urgencyOptions: Option<UrgencyFilter>[] = [
  { value: 'all', label: 'Все' },
  ...urgencyOrder.map((urgency) => ({ value: urgency, label: urgencyLabels[urgency] })),
]

const statusOptions: Option<StatusFilter>[] = [
  { value: 'all', label: 'Все' },
  ...activeTaskStatuses.map((status) => ({ value: status, label: taskStatusLabels[status] })),
]

const ownerOptions: Option<OwnerFilter>[] = [
  { value: 'all', label: 'Все' },
  { value: 'mine', label: 'Мои' },
]

interface TaskFiltersProps {
  urgency: UrgencyFilter
  status: StatusFilter
  owner: OwnerFilter
  onUrgencyChange: (value: UrgencyFilter) => void
  onStatusChange: (value: StatusFilter) => void
  onOwnerChange: (value: OwnerFilter) => void
}

export function TaskFilters(props: TaskFiltersProps) {
  return (
    <div className="flex flex-wrap gap-x-8 gap-y-3">
      <Segmented label="Задачи" options={ownerOptions} value={props.owner} onChange={props.onOwnerChange} />
      <Segmented label="Уровень" options={urgencyOptions} value={props.urgency} onChange={props.onUrgencyChange} />
      <Segmented label="Статус" options={statusOptions} value={props.status} onChange={props.onStatusChange} />
    </div>
  )
}
