import { X } from 'lucide-react'
import { useEffect, useId, type ReactNode } from 'react'
import { createPortal } from 'react-dom'

interface SheetProps {
  title: string
  onClose: () => void
  children: ReactNode
  footer?: ReactNode
}

export function Sheet({ title, onClose, children, footer }: SheetProps) {
  const titleId = useId()

  useEffect(() => {
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose()
      }
    }
    const overflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', closeOnEscape)
    return () => {
      document.body.style.overflow = overflow
      window.removeEventListener('keydown', closeOnEscape)
    }
  }, [onClose])

  return createPortal(
    <div className="fixed inset-0 z-40 animate-fade-in bg-ink/30" onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onMouseDown={(event) => event.stopPropagation()}
        className={
          'fixed inset-x-0 bottom-0 flex max-h-[88svh] animate-sheet-up flex-col rounded-t-[20px] bg-white shadow-popup ' +
          'lg:inset-y-0 lg:left-auto lg:right-0 lg:max-h-none lg:w-[460px] lg:animate-sheet-left lg:rounded-none lg:rounded-l-card'
        }
      >
        <div className="mx-auto mt-2 h-1 w-10 shrink-0 rounded-pill bg-line lg:hidden" aria-hidden />
        <header className="flex items-start justify-between gap-4 px-5 pb-2 pt-3 lg:px-7 lg:pt-7">
          <h2 id={titleId} className="text-sm font-medium uppercase tracking-wider text-subtle">
            {title}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="-mr-1 -mt-1 rounded-control p-1.5 text-subtle hover:bg-card"
            aria-label="Закрыть"
          >
            <X className="size-5" />
          </button>
        </header>
        <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-5 pb-5 lg:px-7">{children}</div>
        {footer && (
          <footer className="flex flex-col-reverse gap-2 border-t border-line px-5 pb-[max(1rem,env(safe-area-inset-bottom))] pt-4 lg:px-7 lg:pb-7">
            {footer}
          </footer>
        )}
      </div>
    </div>,
    document.body,
  )
}
