import { formatDateTime, formatDuration } from '@/shared/lib/format'

interface DeadlineProps {
  dueAt: string
  now: Date
}

const minute = 60_000

export function Deadline({ dueAt, now }: DeadlineProps) {
  const left = Date.parse(dueAt) - now.getTime()
  if (left <= 0) {
    return (
      <div className="font-medium text-overdue">
        {-left < minute ? 'Просрочено' : `Просрочено на ${formatDuration(-left)}`}
      </div>
    )
  }
  return (
    <div>
      <div>{left < minute ? 'меньше минуты' : `через ${formatDuration(left)}`}</div>
      <div className="text-sm text-subtle">{formatDateTime(dueAt)}</div>
    </div>
  )
}
