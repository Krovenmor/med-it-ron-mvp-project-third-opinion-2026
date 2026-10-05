import clsx from 'clsx'
import type { ReactNode } from 'react'

const noticeStyles = {
  overdue: 'bg-overdue-tint text-ink [&_svg]:text-overdue',
  urgent: 'bg-emergency-bg text-emergency-fg',
  neutral: 'bg-card text-ink',
  accent: 'bg-accent-tint text-ink [&_svg]:text-accent-deep',
}

interface NoticeProps {
  tone: keyof typeof noticeStyles
  icon?: ReactNode
  children: ReactNode
  className?: string
}

export function Notice({ tone, icon, children, className }: NoticeProps) {
  return (
    <div className={clsx('flex gap-3 rounded-card p-4', noticeStyles[tone], className)}>
      {icon && <span className="mt-0.5 shrink-0">{icon}</span>}
      <div className="min-w-0 flex-1">{children}</div>
    </div>
  )
}
