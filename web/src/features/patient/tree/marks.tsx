import clsx from 'clsx'

import type { BranchTone } from './layout'
import { toneStyles, type NodeLook } from './looks'

const lookStyles: Record<NodeLook, string> = {
  done: 'fill-accent-deep stroke-white stroke-2',
  'booked-minor': 'fill-white stroke-mark-minor-dark stroke-3',
  'booked-critical': 'fill-white stroke-mark-critical-dark stroke-3',
  planned: 'fill-white stroke-muted stroke-2',
  overdue: 'fill-overdue-tint stroke-overdue stroke-[2.5]',
  declined: 'fill-white stroke-muted stroke-2 opacity-50',
}

const dashed: NodeLook[] = ['planned', 'declined']

interface NodeMarkProps {
  look: NodeLook
  r: number
  ai?: boolean
}

export function NodeMark({ look, r, ai = false }: NodeMarkProps) {
  return (
    <>
      <circle
        r={r}
        className={clsx('tree-shape', lookStyles[look])}
        strokeDasharray={dashed.includes(look) ? '3 2.6' : undefined}
      />
      {look === 'overdue' && (
        <text
          textAnchor="middle"
          dominantBaseline="central"
          className="fill-overdue text-[11px] font-bold"
          style={{ fontSize: r * 1.35 }}
        >
          !
        </text>
      )}
      {ai && <AiBadge x={r * 0.8} y={-r * 0.8} />}
    </>
  )
}

interface AiBadgeProps {
  x: number
  y: number
}

export function AiBadge({ x, y }: AiBadgeProps) {
  return (
    <g transform={`translate(${x} ${y})`}>
      <circle r={6.5} className="fill-white stroke-accent-deep" strokeWidth={1.5} />
      <path
        d="M0 -4 C0.4 -1.2 1.2 -0.4 4 0 C1.2 0.4 0.4 1.2 0 4 C-0.4 1.2 -1.2 0.4 -4 0 C-1.2 -0.4 -0.4 -1.2 0 -4 Z"
        className="fill-accent-deep"
      />
    </g>
  )
}

interface NodeGlyphProps {
  look: NodeLook
  ai?: boolean
  size?: number
  className?: string
}

export function NodeGlyph({ look, ai = false, size = 24, className }: NodeGlyphProps) {
  const r = size / 2 - 3
  return (
    <svg width={size} height={size} viewBox={`${-size / 2} ${-size / 2} ${size} ${size}`} className={clsx('shrink-0 overflow-visible', className)} aria-hidden>
      <NodeMark look={look} r={r} ai={ai} />
    </svg>
  )
}

interface BranchGlyphProps {
  tone: BranchTone
}

export function BranchGlyph({ tone }: BranchGlyphProps) {
  if (tone === 'dashed' || tone === 'faded') {
    return (
      <svg width={44} height={16} aria-hidden className="shrink-0">
        <path d="M2 13 C16 13 26 4 42 4" fill="none" strokeWidth={2} strokeDasharray="5 4" strokeLinecap="round" className="stroke-muted" />
      </svg>
    )
  }
  const [lit, dark] = toneStyles[tone]
  return (
    <svg width={44} height={16} aria-hidden className="shrink-0">
      <path d="M2 9 C16 9 26 1 42 2 L42 4 C27 4 17 13 2 13 Z" className={lit} />
      <path d="M2 13 C17 13 27 4 42 4 L42 6 C28 7 18 16 2 16 Z" className={dark} />
    </svg>
  )
}
