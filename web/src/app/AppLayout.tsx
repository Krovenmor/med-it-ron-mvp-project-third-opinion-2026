import { Outlet } from 'react-router'

import { DemoPanel } from '@/features/demo/DemoPanel'

export function AppLayout() {
  return (
    <div className="min-h-screen">
      <DemoPanel />
      <header className="border-b border-line">
        <div className="mx-auto flex max-w-[1280px] items-baseline gap-3 px-8 py-5">
          <span className="text-lg font-semibold">Маршрут после находки</span>
          <span className="text-sm text-subtle">команда МедиДруны</span>
        </div>
      </header>
      <main className="mx-auto max-w-[1280px] px-8 py-8">
        <Outlet />
      </main>
    </div>
  )
}
