import clsx from 'clsx'
import { NavLink } from 'react-router'

export interface Tab {
  to: string
  label: string
  end: boolean
}

interface TabsProps {
  tabs: Tab[]
}

export function Tabs({ tabs }: TabsProps) {
  return (
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
  )
}
