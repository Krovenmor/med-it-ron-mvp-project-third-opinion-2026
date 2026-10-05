import { Outlet } from 'react-router'

export function StaffLayout() {
  return (
    <>
      <header className="border-b border-line">
        <div className="mx-auto flex max-w-[1280px] items-baseline gap-3 px-8 py-5">
          <span className="text-lg font-semibold">Маршрут после находки</span>
          <span className="text-sm text-subtle">команда МедиДруны</span>
        </div>
      </header>
      <main className="mx-auto max-w-[1280px] px-8 py-8">
        <Outlet />
      </main>
    </>
  )
}
