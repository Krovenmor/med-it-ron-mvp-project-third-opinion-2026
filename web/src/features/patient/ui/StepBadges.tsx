import clsx from 'clsx'

import type { AcceptedMark, StepStatus } from '@/api/models'
import { markLabels } from '@/shared/lib/labels'
import { Pill } from '@/shared/ui/Badge'

import { stepStatusLabels, stepStatusStyles } from '../models'

interface StatusPillProps {
  status: StepStatus
}

export function StatusPill({ status }: StatusPillProps) {
  return <Pill className={clsx('text-xs', stepStatusStyles[status])}>{stepStatusLabels[status]}</Pill>
}

const markDots: Record<AcceptedMark, string> = {
  critical: 'bg-mark-critical',
  minor: 'bg-mark-minor',
}

interface MarkBadgeProps {
  mark: AcceptedMark
}

export function MarkBadge({ mark }: MarkBadgeProps) {
  return (
    <Pill className="bg-card text-xs text-ink">
      <span className={clsx('size-2.5 rounded-full', markDots[mark])} aria-hidden />
      {markLabels[mark]}
    </Pill>
  )
}
