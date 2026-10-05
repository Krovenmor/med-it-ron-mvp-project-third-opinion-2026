import type { PatientRoute, RouteStep, TrunkNode } from '@/api/models'
import { formatDate, pluralize } from '@/shared/lib/format'

import { branchCenterline, branchRibbon, trunkShape, type TwoTone } from './geometry'

export type Side = 'left' | 'right'
export type BranchTone = 'green' | 'minor' | 'critical' | 'dashed' | 'faded'

export interface TreeMetrics {
  width: number
  dateColumn: number
  graphWidth: number
  trunkX: number
  eventHeight: number
  emptyDayHeight: number
  gapHeight: number
  todayGap: number
  padding: number
  firstLane: number
  laneGap: number
  maxLanes: number
  trunkBase: number
  trunkTop: number
  flare: number
  branchStart: number
  branchEnd: number
  trunkNodeRadius: number
  stepNodeRadius: number
  labelGap: number
  labelMaxWidth: number
}

export interface Band {
  key: string
  top: number
  height: number
  label: string
  gap: boolean
  today: boolean
}

export interface NodeLayout<T> {
  item: T
  x: number
  y: number
  r: number
  side: Side
  labelWidth: number
}

export interface BranchLayout {
  id: string
  tone: BranchTone
  centerline: string
  ribbon: TwoTone
}

export interface TrunkLayout extends TwoTone {
  x: number
  groundY: number
  groundRadius: number
  breaks: number[]
}

export interface TreeLayout {
  height: number
  bands: Band[]
  lines: { key: string; y: number }[]
  todayY: number
  trunk: TrunkLayout
  trunkNodes: NodeLayout<TrunkNode>[]
  stepNodes: NodeLayout<RouteStep>[]
  branches: BranchLayout[]
}

const dayMs = 86_400_000

export function treeMetrics(width: number): TreeMetrics {
  const compact = width < 560
  const dateColumn = compact ? 70 : 104
  const graphWidth = width - dateColumn
  return {
    width,
    dateColumn,
    graphWidth,
    trunkX: Math.round(graphWidth / 2),
    eventHeight: compact ? 66 : 62,
    emptyDayHeight: compact ? 22 : 26,
    gapHeight: 56,
    todayGap: 18,
    padding: 16,
    firstLane: compact ? 36 : Math.max(64, Math.round(graphWidth * 0.13)),
    laneGap: compact ? 24 : Math.max(48, Math.round(graphWidth * 0.09)),
    maxLanes: 3,
    trunkBase: compact ? 9 : 14,
    trunkTop: compact ? 5 : 7,
    flare: compact ? 11 : 18,
    branchStart: compact ? 5 : 7,
    branchEnd: compact ? 2.2 : 3,
    trunkNodeRadius: compact ? 10 : 12.5,
    stepNodeRadius: compact ? 7.5 : 9,
    labelGap: 8,
    labelMaxWidth: compact ? 140 : 230,
  }
}

interface TreeEvent {
  id: string
  at: number
  trunk: boolean
}

type Span = { gap: false; day: number } | { gap: true; from: number; to: number }

interface Placement {
  side: Side
  lane: number
  from: number
  to: number
  origin: string
}

export function buildTreeLayout(route: PatientRoute, m: TreeMetrics): TreeLayout | null {
  if (route.trunk.length === 0) {
    return null
  }
  const now = Date.parse(route.now)
  const today = dayIndex(now)
  const trunkIds = new Set(route.trunk.map((node) => node.id))
  const steps = route.steps.filter((step) => trunkIds.has(step.origin_id))

  const byDay = new Map<number, TreeEvent[]>()
  const events: TreeEvent[] = [
    ...route.trunk.map((node) => ({ id: node.id, at: Date.parse(node.date), trunk: true })),
    ...steps.map((step) => ({ id: step.id, at: Date.parse(step.date), trunk: false })),
  ]
  events.sort((a, b) => a.at - b.at || Number(b.trunk) - Number(a.trunk))
  for (const event of events) {
    const day = dayIndex(event.at)
    byDay.set(day, [...(byDay.get(day) ?? []), event])
  }

  const levels = new Map<string, number>()
  const placedBands: { span: Span; bottom: number; height: number }[] = []
  let level = m.padding
  let todayLevel = level
  for (const span of spansFor(new Set([...byDay.keys(), today]))) {
    if (span.gap) {
      placedBands.push({ span, bottom: level, height: m.gapHeight })
      level += m.gapHeight
      continue
    }
    const dayEvents = byDay.get(span.day) ?? []
    const isToday = span.day === today
    const content = dayEvents.length * m.eventHeight + (isToday ? m.todayGap : 0)
    const height = Math.max(m.emptyDayHeight, content)
    let cursor = level + (height - content) / 2
    let todayPending = isToday
    for (const event of dayEvents) {
      if (todayPending && event.at > now) {
        todayLevel = cursor + m.todayGap / 2
        cursor += m.todayGap
        todayPending = false
      }
      levels.set(event.id, cursor + m.eventHeight / 2)
      cursor += m.eventHeight
    }
    if (todayPending) {
      todayLevel = cursor + m.todayGap / 2
    }
    placedBands.push({ span, bottom: level, height })
    level += height
  }

  const height = level + m.padding
  const toY = (value: number) => height - value
  const levelOf = (id: string) => levels.get(id) ?? 0

  const bands: Band[] = placedBands.map(({ span, bottom, height: bandHeight }) => ({
    key: span.gap ? `gap:${span.from}` : `day:${span.day}`,
    top: toY(bottom + bandHeight),
    height: bandHeight,
    label: span.gap ? periodLabel(span.to - span.from + 1) : formatDate(new Date(dayStart(span.day))),
    gap: span.gap,
    today: !span.gap && span.day === today,
  }))
  const lines = placedBands.map(({ span, bottom }) => ({
    key: span.gap ? `gap:${span.from}` : `day:${span.day}`,
    y: toY(bottom),
  }))
  const last = placedBands[placedBands.length - 1]
  lines.push({ key: 'top', y: toY(last.bottom + last.height) })

  const placements = placeBranches(steps, levelOf, m)
  const laneX = (placement: Placement) =>
    m.trunkX + (placement.side === 'right' ? 1 : -1) * (m.firstLane + placement.lane * m.laneGap)

  const stepNodes = steps.flatMap((step) => {
    const placement = placements.get(step.id)
    if (!placement) {
      return []
    }
    const x = laneX(placement)
    return [nodeLayout(step, x, toY(levelOf(step.id)), m.stepNodeRadius, placement.side, m)]
  })

  const branches = stepNodes.map(({ item: step, x, y }) => {
    const from: [number, number] = [m.trunkX, toY(levelOf(step.origin_id))]
    const to: [number, number] = [x, y]
    const bend = m.eventHeight * 1.15
    return {
      id: step.id,
      tone: branchTone(step),
      centerline: branchCenterline(from, to, bend),
      ribbon: branchRibbon(from, to, bend, m.branchStart, m.branchEnd),
    }
  })

  const trunkLevels = route.trunk.map((node) => levelOf(node.id))
  const rootY = toY(Math.min(...trunkLevels))
  const tipY = toY(Math.max(...trunkLevels)) - m.eventHeight * 0.55
  const groundY = Math.min(rootY + m.eventHeight * 0.45, height - 4)
  const trunk: TrunkLayout = {
    ...trunkShape({
      x: m.trunkX,
      tipY,
      rootY,
      groundY,
      topWidth: m.trunkTop,
      baseWidth: m.trunkBase,
      flare: m.flare,
    }),
    x: m.trunkX,
    groundY,
    groundRadius: m.trunkBase + m.flare + 16,
    breaks: bands.filter((band) => band.gap).map((band) => band.top + band.height / 2).filter((y) => y > tipY && y < rootY),
  }

  const occupied = (side: Side, at: number) =>
    [...placements.values()].some((p) => p.side === side && p.lane === 0 && p.from < at && at < p.to)
  const trunkNodes = route.trunk.map((node) => {
    const at = levelOf(node.id)
    const side: Side = occupied('right', at) && !occupied('left', at) ? 'left' : 'right'
    return nodeLayout(node, m.trunkX, toY(at), m.trunkNodeRadius, side, m)
  })

  return { height, bands, lines, todayY: toY(todayLevel), trunk, trunkNodes, stepNodes, branches }
}

function placeBranches(steps: RouteStep[], levelOf: (id: string) => number, m: TreeMetrics): Map<string, Placement> {
  const placed = new Map<string, Placement>()
  const ordered = [...steps].sort(
    (a, b) => levelOf(b.origin_id) - levelOf(a.origin_id) || levelOf(b.id) - levelOf(a.id),
  )

  const fits = (side: Side, lane: number, from: number, to: number) =>
    [...placed.values()].every((p) => {
      if (p.side !== side) {
        return true
      }
      if (p.lane === lane) {
        return !(p.from < to && from < p.to)
      }
      if (p.lane > lane) {
        return !(p.from > from && p.from < to)
      }
      return !(p.from < from && from < p.to)
    })
  const count = (side: Side, origin?: string) =>
    [...placed.values()].filter((p) => p.side === side && (origin === undefined || p.origin === origin)).length
  const preferred = (origin: string): Side[] => {
    const byOrigin = count('right', origin) - count('left', origin)
    const total = count('right') - count('left')
    return byOrigin > 0 || (byOrigin === 0 && total > 0) ? ['left', 'right'] : ['right', 'left']
  }

  for (const step of ordered) {
    const from = levelOf(step.origin_id)
    const to = levelOf(step.id)
    const sides = preferred(step.origin_id)
    let choice: Pick<Placement, 'side' | 'lane'> | null = null
    for (let lane = 0; lane < m.maxLanes && !choice; lane++) {
      const side = sides.find((candidate) => fits(candidate, lane, from, to))
      if (side) {
        choice = { side, lane }
      }
    }
    choice ??= { side: sides[0], lane: m.maxLanes - 1 }
    placed.set(step.id, { ...choice, from, to, origin: step.origin_id })
  }
  return placed
}

function nodeLayout<T>(item: T, x: number, y: number, r: number, side: Side, m: TreeMetrics): NodeLayout<T> {
  const available = side === 'right' ? m.graphWidth - (x + r + m.labelGap) - 4 : x - r - m.labelGap - 4
  return { item, x, y, r, side, labelWidth: Math.max(0, Math.min(available, m.labelMaxWidth)) }
}

function branchTone(step: RouteStep): BranchTone {
  switch (step.status) {
    case 'done':
      return 'green'
    case 'booked':
      return step.mark === 'critical' ? 'critical' : 'minor'
    case 'declined':
      return 'faded'
    default:
      return 'dashed'
  }
}

function spansFor(anchors: Set<number>): Span[] {
  const sorted = [...anchors].sort((a, b) => a - b)
  const first = sorted[0] - 1
  const last = sorted[sorted.length - 1] + 1
  const near = (day: number) => anchors.has(day - 1) || anchors.has(day) || anchors.has(day + 1)

  const spans: Span[] = []
  let day = first
  while (day <= last) {
    if (near(day)) {
      spans.push({ gap: false, day })
      day++
      continue
    }
    let end = day
    while (end + 1 <= last && !near(end + 1)) {
      end++
    }
    if (end - day + 1 <= 2) {
      for (let empty = day; empty <= end; empty++) {
        spans.push({ gap: false, day: empty })
      }
    } else {
      spans.push({ gap: true, from: day, to: end })
    }
    day = end + 1
  }
  return spans
}

function dayIndex(at: number): number {
  const date = new Date(at)
  return Math.round(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()) / dayMs)
}

function dayStart(index: number): number {
  const utc = new Date(index * dayMs)
  return new Date(utc.getUTCFullYear(), utc.getUTCMonth(), utc.getUTCDate()).getTime()
}

function periodLabel(days: number): string {
  if (days < 14) {
    return `${days} ${pluralize(days, 'день', 'дня', 'дней')}`
  }
  if (days < 60) {
    const weeks = Math.round(days / 7)
    return `${weeks} ${pluralize(weeks, 'неделя', 'недели', 'недель')}`
  }
  if (days < 365) {
    const months = Math.round(days / 30)
    return `${months} ${pluralize(months, 'месяц', 'месяца', 'месяцев')}`
  }
  const years = Math.round(days / 365)
  return `${years} ${pluralize(years, 'год', 'года', 'лет')}`
}
