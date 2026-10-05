import clsx from 'clsx'
import { CalendarClock, CircleArrowUp, CircleCheck, TriangleAlert, type LucideIcon } from 'lucide-react'
import type { ReactNode } from 'react'

import type { Urgency } from '@/api/models'
import { urgencyLabels } from '@/shared/lib/labels'

interface PillProps {
  children: ReactNode
  className?: string
}

export function Pill({ children, className }: PillProps) {
  return (
    <span
      className={clsx(
        'inline-flex items-center gap-1.5 whitespace-nowrap rounded-pill px-3 py-1 text-sm font-medium',
        className,
      )}
    >
      {children}
    </span>
  )
}

const urgencyStyles: Record<Urgency, string> = {
  emergency: 'bg-emergency-bg text-emergency-fg',
  priority: 'bg-priority-bg text-priority-fg',
  planned: 'bg-planned-bg text-planned-fg ring-1 ring-inset ring-line',
  normal: 'bg-normal-bg text-normal-fg',
}

const urgencyIcons: Record<Urgency, LucideIcon> = {
  emergency: TriangleAlert,
  priority: CircleArrowUp,
  planned: CalendarClock,
  normal: CircleCheck,
}

interface UrgencyBadgeProps {
  urgency: Urgency
}

export function UrgencyBadge({ urgency }: UrgencyBadgeProps) {
  const Icon = urgencyIcons[urgency]
  return (
    <Pill className={urgencyStyles[urgency]}>
      <Icon className="size-4" aria-hidden />
      {urgencyLabels[urgency]}
    </Pill>
  )
}
