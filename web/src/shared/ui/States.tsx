import { Inbox, LoaderCircle, type LucideIcon } from 'lucide-react'

import { Button } from './Button'

interface LoadingStateProps {
  text?: string
}

export function LoadingState({ text = 'Загружаем данные' }: LoadingStateProps) {
  return (
    <div className="flex items-center justify-center gap-3 py-16 text-subtle">
      <LoaderCircle className="size-5 animate-spin" />
      {text}
    </div>
  )
}

interface ErrorStateProps {
  title?: string
  message?: string
  onRetry?: () => void
}

export function ErrorState({ title = 'Не удалось загрузить данные', message, onRetry }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center gap-3 py-16 text-center">
      <div className="text-lg font-medium">{title}</div>
      {message && <div className="max-w-md text-sm text-subtle">{message}</div>}
      {onRetry && (
        <Button variant="secondary" size="sm" onClick={onRetry}>
          Повторить
        </Button>
      )}
    </div>
  )
}

interface EmptyStateProps {
  title: string
  text?: string
  icon?: LucideIcon
}

export function EmptyState({ title, text, icon: Icon = Inbox }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center gap-2 py-16 text-center">
      <Icon className="size-8 text-muted" aria-hidden />
      <div className="text-lg font-medium">{title}</div>
      {text && <div className="max-w-md text-sm text-subtle">{text}</div>}
    </div>
  )
}
