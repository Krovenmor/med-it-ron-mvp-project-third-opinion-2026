import clsx from 'clsx'
import { ChevronDown } from 'lucide-react'
import { useState, type ReactNode } from 'react'

interface TwoToneHeadingProps {
  main: string
  rest?: string
  className?: string
}

export function TwoToneHeading({ main, rest, className }: TwoToneHeadingProps) {
  return (
    <h1 className={clsx('text-3xl font-medium tracking-tight', className)}>
      {main}
      {rest && <span className="text-muted"> {rest}</span>}
    </h1>
  )
}

interface CardProps {
  children: ReactNode
  className?: string
}

export function Card({ children, className }: CardProps) {
  return <section className={clsx('rounded-card bg-card p-6', className)}>{children}</section>
}

interface MetricProps {
  value: ReactNode
  label: string
}

export function Metric({ value, label }: MetricProps) {
  return (
    <div className="rounded-card bg-card px-6 py-5">
      <div className="text-[28px] font-medium leading-tight">{value}</div>
      <div className="mt-1 text-sm text-subtle">{label}</div>
    </div>
  )
}

interface CollapsibleProps {
  title: string
  defaultOpen?: boolean
  aside?: ReactNode
  children: ReactNode
}

export function Collapsible({ title, defaultOpen = false, aside, children }: CollapsibleProps) {
  const [open, setOpen] = useState(defaultOpen)
  return (
    <section className="rounded-card bg-card">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        className="flex w-full items-center justify-between gap-4 px-6 py-5 text-left"
        aria-expanded={open}
      >
        <span className="text-lg font-medium">{title}</span>
        <span className="flex items-center gap-3 text-sm text-subtle">
          {aside}
          <ChevronDown className={clsx('size-5 transition-transform', open && 'rotate-180')} />
        </span>
      </button>
      {open && <div className="px-6 pb-6">{children}</div>}
    </section>
  )
}
