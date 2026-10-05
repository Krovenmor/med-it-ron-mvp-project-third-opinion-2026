import clsx from 'clsx'
import { Sparkles } from 'lucide-react'
import { useEffect, useMemo, useRef, type CSSProperties } from 'react'

import type { PatientRoute } from '@/api/models'
import { useElementWidth } from '@/shared/lib/useElementWidth'

import { stepCaption, trunkCaption, type Caption } from './captions'
import { buildTreeLayout, treeMetrics, type Side, type TreeLayout } from './layout'
import { stepLook, toneStyles } from './looks'
import { NodeMark } from './marks'

interface HealthTreeProps {
  route: PatientRoute
  selectedId: string | null
  onSelect: (id: string) => void
}

export function HealthTree({ route, selectedId, onSelect }: HealthTreeProps) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLDivElement>(null)
  const focusId = useRef(selectedId)
  const scrolled = useRef(false)
  const width = useElementWidth(canvasRef)
  const layout = useMemo(() => (width > 0 ? buildTreeLayout(route, treeMetrics(width)) : null), [route, width])
  const todayY = layout?.todayY

  useEffect(() => {
    const scroller = scrollRef.current
    if (!scroller || todayY === undefined) {
      return
    }
    const focused =
      !scrolled.current && focusId.current
        ? scroller.querySelector<HTMLElement>(`[data-node="${CSS.escape(focusId.current)}"]`)
        : null
    const target = focused ? Number(focused.dataset.y) : todayY
    scroller.scrollTo({ top: target - scroller.clientHeight * 0.55, behavior: scrolled.current ? 'smooth' : 'auto' })
    scrolled.current = true
  }, [todayY])

  return (
    <div
      ref={scrollRef}
      className="relative max-h-[min(70svh,760px)] overflow-y-auto overscroll-contain lg:max-h-[calc(100svh-15rem)]"
    >
      <div ref={canvasRef} className="relative" style={{ height: layout?.height ?? 320 }}>
        {layout && (
          <TreeCanvas
            layout={layout}
            width={width}
            route={route}
            selectedId={selectedId}
            onSelect={onSelect}
          />
        )}
      </div>
    </div>
  )
}

interface TreeCanvasProps {
  layout: TreeLayout
  width: number
  route: PatientRoute
  selectedId: string | null
  onSelect: (id: string) => void
}

function TreeCanvas({ layout, width, route, selectedId, onSelect }: TreeCanvasProps) {
  const metrics = treeMetrics(width)
  const { trunk } = layout

  return (
    <>
      <svg width={width} height={layout.height} className="absolute inset-0 overflow-visible" aria-hidden>
        {layout.bands
          .filter((band) => band.gap)
          .map((band) => (
            <rect
              key={band.key}
              width={width}
              height={band.height}
              className="tree-move fill-card-nested"
              style={move(0, band.top)}
            />
          ))}
        {layout.lines.map((line) => (
          <line key={line.key} x2={width} className="tree-move stroke-muted/25" style={move(0, line.y)} />
        ))}

        <ellipse
          cx={trunk.x}
          cy={trunk.groundY}
          rx={trunk.groundRadius}
          ry={4}
          className="tree-shape fill-accent-soft"
        />
        {layout.branches.map((branch) =>
          branch.tone === 'dashed' || branch.tone === 'faded' ? (
            <path
              key={`${branch.id}:line`}
              d={branch.centerline}
              fill="none"
              strokeWidth={2}
              strokeDasharray="5 5"
              strokeLinecap="round"
              className={clsx('tree-shape stroke-muted', branch.tone === 'faded' && 'opacity-40')}
              style={shape(branch.centerline)}
            />
          ) : (
            <g key={`${branch.id}:ribbon`}>
              <path d={branch.ribbon.light} className={clsx('tree-shape', toneStyles[branch.tone][0])} style={shape(branch.ribbon.light)} />
              <path d={branch.ribbon.shade} className={clsx('tree-shape', toneStyles[branch.tone][1])} style={shape(branch.ribbon.shade)} />
            </g>
          ),
        )}
        <path d={trunk.light} className="tree-shape fill-accent" style={shape(trunk.light)} />
        <path d={trunk.shade} className="tree-shape fill-accent-dark" style={shape(trunk.shade)} />
        {trunk.breaks.map((y) => (
          <g key={y} className="tree-move stroke-white" strokeWidth={2} style={move(trunk.x, y)}>
            <line x1={-metrics.trunkBase - 4} y1={-1} x2={metrics.trunkBase + 4} y2={-6} />
            <line x1={-metrics.trunkBase - 4} y1={5} x2={metrics.trunkBase + 4} y2={0} />
          </g>
        ))}

        <g className="tree-move" style={move(0, layout.todayY)}>
          <line x2={width} strokeWidth={1.5} className="stroke-today" />
          <circle cx={4} r={3.5} className="fill-today" />
        </g>

        {layout.trunkNodes.map(({ item, x, y, r }) => (
          <g key={item.id} className="tree-move" style={move(x, y)}>
            {item.id === selectedId && <circle r={r + 5} fill="none" strokeWidth={1.5} className="stroke-ink" />}
            <NodeMark look="done" r={r} ai={item.ai_report} />
          </g>
        ))}
        {layout.stepNodes.map(({ item, x, y, r }) => (
          <g key={item.id} className="tree-move" style={move(x, y)}>
            {item.id === route.next_step_id && (
              <circle
                r={r + 6}
                className={clsx('tree-halo animate-halo', item.status === 'overdue' ? 'fill-overdue/25' : 'fill-accent/40')}
              />
            )}
            {item.id === selectedId && <circle r={r + 5} fill="none" strokeWidth={1.5} className="stroke-ink" />}
            <NodeMark look={stepLook(item)} r={r} />
          </g>
        ))}
      </svg>

      {layout.bands
        .filter((band) => !band.today)
        .map((band) => (
          <div
            key={band.key}
            className="tree-move pointer-events-none absolute right-0 top-0 flex flex-col items-end pr-2"
            style={{ width: metrics.dateColumn, transform: `translateY(${band.top + band.height / 2}px) translateY(-50%)` }}
          >
            {band.gap ? (
              <>
                <span className="text-lg leading-none text-muted">⋯</span>
                <span className="text-[10.5px] text-subtle">{band.label}</span>
              </>
            ) : (
              <span className="text-[11px] tabular-nums text-subtle">{band.label}</span>
            )}
          </div>
        ))}
      <div
        className="tree-move pointer-events-none absolute right-1 top-0"
        style={{ transform: `translateY(${layout.todayY}px) translateY(-50%)` }}
      >
        <span className="rounded-pill bg-today px-2 py-0.5 text-[10.5px] font-semibold text-white">Сегодня</span>
      </div>

      {layout.trunkNodes.map(({ item, x, y, r, side, labelWidth }) => (
        <NodeLabel
          key={item.id}
          id={item.id}
          title={item.title}
          caption={trunkCaption(item)}
          ai={item.ai_report && !item.pending}
          x={x}
          y={y}
          r={r}
          gap={metrics.labelGap}
          side={side}
          width={labelWidth}
          onSelect={onSelect}
        />
      ))}
      {layout.stepNodes.map(({ item, x, y, r, side, labelWidth }) => (
        <NodeLabel
          key={item.id}
          id={item.id}
          title={item.title}
          caption={stepCaption(item)}
          ai={false}
          x={x}
          y={y}
          r={r}
          gap={metrics.labelGap}
          side={side}
          width={labelWidth}
          onSelect={onSelect}
        />
      ))}
    </>
  )
}

interface NodeLabelProps {
  id: string
  title: string
  caption: Caption | null
  ai: boolean
  x: number
  y: number
  r: number
  gap: number
  side: Side
  width: number
  onSelect: (id: string) => void
}

function NodeLabel({ id, title, caption, ai, x, y, r, gap, side, width, onSelect }: NodeLabelProps) {
  const anchor = side === 'right' ? x + r + gap : x - r - gap
  const hit = Math.max(36, r * 2 + 12)
  return (
    <>
      <button
        type="button"
        tabIndex={-1}
        aria-hidden
        onClick={() => onSelect(id)}
        className="tree-move absolute left-0 top-0 rounded-full"
        style={{ width: hit, height: hit, transform: `translate(${x - hit / 2}px, ${y - hit / 2}px)` }}
      />
      <button
        type="button"
        data-node={id}
        data-y={y}
        onClick={() => onSelect(id)}
        className={clsx(
          'tree-move absolute left-0 top-0 rounded-md bg-white/85 px-0.5 py-0.5 hover:bg-white sm:px-1',
          'focus-visible:outline-2 focus-visible:outline-accent-deep',
          side === 'left' ? 'text-right' : 'text-left',
        )}
        style={{
          maxWidth: width,
          transform: `translate(${anchor}px, ${y}px) translate(${side === 'right' ? '0' : '-100%'}, -50%)`,
        }}
      >
        <span className="line-clamp-3 hyphens-auto break-words text-[12px] font-medium leading-tight sm:text-[13px]">{title}</span>
        {caption && (
          <span
            className={clsx(
              'mt-0.5 flex items-center gap-1 text-[11px] leading-tight',
              side === 'left' && 'justify-end',
              caption.alert ? 'font-semibold text-overdue' : ai ? 'text-accent-deep' : 'text-subtle',
            )}
          >
            {ai && <Sparkles className="size-3 shrink-0" aria-hidden />}
            {caption.text}
          </span>
        )}
      </button>
    </>
  )
}

function move(x: number, y: number): CSSProperties {
  return { transform: `translate(${x}px, ${y}px)` }
}

function shape(d: string): CSSProperties {
  return { ['d' as string]: `path("${d}")` } as CSSProperties
}
