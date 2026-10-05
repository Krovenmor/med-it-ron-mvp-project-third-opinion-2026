import type { CallAttempt } from '@/api/models'
import { formatDateTime } from '@/shared/lib/format'
import { declineReasonLabels } from '@/shared/lib/labels'
import { Card } from '@/shared/ui/Layout'

import { outcomeLabels } from '../models'

interface AttemptsHistoryProps {
  attempts: CallAttempt[]
}

export function AttemptsHistory({ attempts }: AttemptsHistoryProps) {
  return (
    <Card className="flex flex-col gap-4">
      <h2 className="text-lg font-medium">История попыток</h2>
      {attempts.length === 0 && <p className="text-subtle">Звонков ещё не было</p>}
      <ol className="flex flex-col gap-3">
        {[...attempts].reverse().map((attempt) => (
          <li key={attempt.id} className="flex flex-col gap-0.5 border-l-2 border-line pl-3">
            <span className="font-medium">
              {outcomeLabels[attempt.outcome]}
              {attempt.decline_reason && `: ${declineReasonLabels[attempt.decline_reason].toLowerCase()}`}
            </span>
            {attempt.callback_at && (
              <span className="text-sm">перезвонить {formatDateTime(attempt.callback_at)}</span>
            )}
            {attempt.comment && <span className="text-sm">«{attempt.comment}»</span>}
            <span className="text-sm text-subtle">
              {formatDateTime(attempt.created_at)} · {attempt.actor}
            </span>
          </li>
        ))}
      </ol>
    </Card>
  )
}
