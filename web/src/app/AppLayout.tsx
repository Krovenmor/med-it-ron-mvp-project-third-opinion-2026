import { Outlet } from 'react-router'

import { DemoPanel } from '@/features/demo/DemoPanel'

export function AppLayout() {
  return (
    <div className="min-h-screen">
      <DemoPanel />
      <Outlet />
    </div>
  )
}
