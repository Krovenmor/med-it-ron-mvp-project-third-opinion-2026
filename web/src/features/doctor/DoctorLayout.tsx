import clsx from 'clsx'
import { NavLink, Outlet } from 'react-router'

const tabs = [
  { to: '/doctor', label: 'Очередь на проверку', end: true },
  { to: '/doctor/help', label: 'Справка', end: false },
]

export function DoctorLayout() {
  return (
    <div className="flex flex-col gap-8">
      <nav className="flex gap-8 border-b border-line">
        {tabs.map((tab) => (
          <NavLink
            key={tab.to}
            to={tab.to}
            end={tab.end}
            className={({ isActive }) =>
              clsx(
                '-mb-px border-b-2 pb-3 text-base font-medium',
                isActive ? 'border-accent text-ink' : 'border-transparent text-subtle hover:text-ink',
              )
            }
          >
            {tab.label}
          </NavLink>
        ))}
      </nav>
      <Outlet />
    </div>
  )
}
