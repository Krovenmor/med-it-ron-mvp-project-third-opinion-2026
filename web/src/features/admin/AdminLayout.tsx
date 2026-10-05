import { Outlet } from 'react-router'

import { Tabs, type Tab } from '@/shared/ui/Tabs'

const tabs: Tab[] = [
  { to: '/admin', label: 'Позвонить', end: true },
  { to: '/admin/notifications', label: 'Журнал уведомлений', end: false },
]

export function AdminLayout() {
  return (
    <div className="flex flex-col gap-8">
      <Tabs tabs={tabs} />
      <Outlet />
    </div>
  )
}
