import clsx from 'clsx'

export interface Option<T extends string> {
  value: T
  label: string
}

interface SegmentedProps<T extends string> {
  label?: string
  options: Option<T>[]
  value: T
  onChange: (value: T) => void
}

export function Segmented<T extends string>({ label, options, value, onChange }: SegmentedProps<T>) {
  return (
    <div className="flex items-center gap-3">
      {label && <span className="text-sm text-subtle">{label}</span>}
      <div className="flex gap-1 rounded-pill bg-card p-1">
        {options.map((option) => (
          <button
            key={option.value}
            type="button"
            onClick={() => onChange(option.value)}
            className={clsx(
              'rounded-pill px-3 py-1 text-sm font-medium transition-colors',
              option.value === value ? 'bg-white text-ink shadow-sm' : 'text-subtle hover:text-ink',
            )}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  )
}
