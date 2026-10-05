import { formatCountdown, formatDuration } from '@/shared/lib/format'
import { emergencyReviewSlaMs } from '@/shared/lib/time'

interface SlaTimerProps {
  receivedAt: string
  now: Date
}

export function SlaTimer({ receivedAt, now }: SlaTimerProps) {
  const left = Date.parse(receivedAt) + emergencyReviewSlaMs - now.getTime()
  return left > 0 ? (
    <div className="text-sm font-medium text-emergency-fg">SLA: осталось {formatCountdown(left)}</div>
  ) : (
    <div className="text-sm font-medium text-overdue">
      {-left < 60_000 ? 'SLA просрочен' : `SLA просрочен на ${formatDuration(-left)}`}
    </div>
  )
}
