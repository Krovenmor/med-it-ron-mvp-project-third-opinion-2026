import type { Role } from './models'

export const roles: Role[] = [
  { path: '/patient', label: 'Пациент' },
  { path: '/doctor', label: 'Врач' },
  { path: '/admin', label: 'Администратор' },
  { path: '/manager', label: 'Руководитель' },
]
