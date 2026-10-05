import clsx from 'clsx'

export type Variant = 'primary' | 'secondary' | 'ghost'
export type Size = 'md' | 'sm'

const variants: Record<Variant, string> = {
  primary: 'bg-accent text-ink hover:bg-accent-dark disabled:hover:bg-accent',
  secondary: 'border border-line bg-white text-ink hover:bg-card disabled:hover:bg-white',
  ghost: 'text-accent-deep hover:bg-accent-tint disabled:hover:bg-transparent',
}

const sizes: Record<Size, string> = {
  md: 'h-11 px-5 text-base',
  sm: 'h-9 px-3 text-sm',
}

export function buttonStyles(variant: Variant = 'primary', size: Size = 'md'): string {
  return clsx(
    'inline-flex shrink-0 items-center justify-center gap-2 rounded-control font-semibold transition-colors',
    'disabled:cursor-not-allowed disabled:opacity-50',
    variants[variant],
    sizes[size],
  )
}
