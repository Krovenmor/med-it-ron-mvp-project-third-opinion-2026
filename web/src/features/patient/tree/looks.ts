import type { RouteStep } from '@/api/models'

import type { BranchTone } from './layout'

export type NodeLook = 'done' | 'booked-minor' | 'booked-critical' | 'planned' | 'overdue' | 'declined'

export function stepLook(step: RouteStep): NodeLook {
  switch (step.status) {
    case 'done':
      return 'done'
    case 'booked':
      return step.mark === 'critical' ? 'booked-critical' : 'booked-minor'
    case 'overdue':
      return 'overdue'
    case 'declined':
      return 'declined'
    default:
      return 'planned'
  }
}

export const toneStyles: Record<Exclude<BranchTone, 'dashed' | 'faded'>, [string, string]> = {
  green: ['fill-accent', 'fill-accent-dark'],
  minor: ['fill-mark-minor', 'fill-mark-minor-dark'],
  critical: ['fill-mark-critical', 'fill-mark-critical-dark'],
}
