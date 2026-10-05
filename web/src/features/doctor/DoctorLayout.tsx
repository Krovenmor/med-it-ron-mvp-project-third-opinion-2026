import { Outlet } from 'react-router'

import { Tabs, type Tab } from '@/shared/ui/Tabs'

const tabs: Tab[] = [
  { to: '/doctor', label: 'Очередь на проверку', end: true },
  { to: '/doctor/help', label: 'Справка', end: false },
]

export function DoctorLayout() {
  return (
    <div className="flex flex-col gap-8">
      <Tabs tabs={tabs} />
      <Outlet />
    </div>
  )
}
