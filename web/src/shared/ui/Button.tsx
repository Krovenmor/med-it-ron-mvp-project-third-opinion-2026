import clsx from 'clsx'
import type { ButtonHTMLAttributes } from 'react'

import { buttonStyles, type Size, type Variant } from './buttonStyles'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  size?: Size
}

export function Button({ variant = 'primary', size = 'md', className, type = 'button', ...props }: ButtonProps) {
  return <button type={type} className={clsx(buttonStyles(variant, size), className)} {...props} />
}
