import type { Role } from './models'

export const roles: Role[] = [
  { path: '/patient', label: 'Пациент', short: 'Пациент' },
  { path: '/doctor', label: 'Врач', short: 'Врач' },
  { path: '/admin', label: 'Администратор', short: 'Админ' },
  { path: '/manager', label: 'Руководитель', short: 'Руководитель' },
]
