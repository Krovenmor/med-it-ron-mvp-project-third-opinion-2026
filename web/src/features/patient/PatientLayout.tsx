import clsx from 'clsx'
import { RefreshCw, Sprout } from 'lucide-react'
import { useMemo } from 'react'
import { NavLink, Outlet, useLocation, useNavigate, useParams } from 'react-router'

import { ApiError } from '@/api/http'
import { usePatientRoute } from '@/api/queries'
import { formatDateTime } from '@/shared/lib/format'
import { EmptyState, ErrorState, LoadingState } from '@/shared/ui/States'
import { Tabs } from '@/shared/ui/Tabs'

import { NodeSheet } from './card/NodeSheet'
import { defaultPatientId, demoPatients, patientTabs, type PatientContext } from './models'
import { useNodeSheet } from './usePatient'

export function PatientLayout() {
  const { patientId = defaultPatientId } = useParams()
  const query = usePatientRoute(patientId)
  const sheet = useNodeSheet()
  const route = query.data
  const notFound = query.error instanceof ApiError && query.error.status === 404

  const context = useMemo<PatientContext | null>(
    () => (route ? { patientId, route, now: new Date(route.now), openNode: sheet.open } : null),
    [patientId, route, sheet.open],
  )

  const demo = demoPatients.find((patient) => patient.id === patientId)
  const base = `/patient/${patientId}`

  return (
    <div className="bg-white">
      <header className="border-b border-line lg:border-b-0">
        <div className="mx-auto flex max-w-[1280px] flex-wrap items-end gap-x-6 gap-y-3 px-4 py-4 md:px-8 lg:pb-0 lg:pt-7">
          <div className="min-w-0 flex-1">
            <div className="text-xs font-medium uppercase tracking-wider text-accent-deep">
              {route?.clinic.name ?? 'Клиника'} · ваш маршрут
            </div>
            <h1 className="mt-1 truncate text-2xl font-medium tracking-tight lg:text-3xl">
              {route?.patient.full_name ?? demo?.name ?? patientId}
            </h1>
            <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-subtle">
              {route && <span>обновлено {formatDateTime(route.now)}</span>}
              <button
                type="button"
                onClick={() => void query.refetch()}
                disabled={query.isFetching}
                className="inline-flex items-center gap-1.5 font-medium text-accent-deep hover:underline disabled:opacity-60"
              >
                <RefreshCw className={clsx('size-3.5', query.isFetching && 'animate-spin')} aria-hidden />
                Обновить
              </button>
            </div>
          </div>
          <PatientSwitcher patientId={patientId} />
          <div className="hidden w-full lg:block">
            <Tabs
              tabs={patientTabs.map((tab) => ({
                to: tab.path ? `${base}/${tab.path}` : base,
                label: tab.label,
                end: tab.path === '',
              }))}
            />
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-[1280px] px-4 pb-28 pt-5 md:px-8 lg:pb-12 lg:pt-8">
        {context ? (
          <>
            {query.isError && (
              <div className="mb-5 flex flex-wrap items-center gap-3 rounded-card bg-priority-bg px-5 py-3 text-sm text-priority-fg">
                Не удалось обновить данные, показываем последние.
                <button type="button" onClick={() => void query.refetch()} className="font-semibold underline">
                  Повторить
                </button>
              </div>
            )}
            <Outlet context={context} />
          </>
        ) : query.isPending ? (
          <LoadingState text="Загружаем ваш маршрут" />
        ) : notFound ? (
          <EmptyState
            icon={Sprout}
            title="Маршрута пока нет"
            text="Он появится после обследования: врач проверит результаты и подготовит план. Мы пришлём уведомление."
          />
        ) : (
          <ErrorState
            title="Не удалось обновить данные"
            message="Проверьте соединение и попробуйте ещё раз."
            onRetry={() => void query.refetch()}
          />
        )}
      </main>

      <BottomNav base={base} />
      {context && sheet.nodeId && (
        <NodeSheet context={context} nodeId={sheet.nodeId} mode={sheet.mode} onClose={sheet.close} />
      )}
    </div>
  )
}

interface PatientSwitcherProps {
  patientId: string
}

function PatientSwitcher({ patientId }: PatientSwitcherProps) {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  const tab = pathname.split('/').slice(3).join('/')

  return (
    <label className="flex w-full items-center gap-2 text-sm sm:w-auto lg:mb-1 lg:self-start">
      <span className="shrink-0 text-subtle">Демо-пациент</span>
      <select
        value={patientId}
        onChange={(event) => navigate(`/patient/${event.target.value}${tab ? `/${tab}` : ''}`)}
        className="h-10 min-w-0 flex-1 rounded-control border border-line bg-white px-3 text-ink outline-none focus:border-accent sm:flex-none"
      >
        {demoPatients.map((patient) => (
          <option key={patient.id} value={patient.id}>
            {patient.name} · {patient.story}
          </option>
        ))}
      </select>
    </label>
  )
}

interface BottomNavProps {
  base: string
}

function BottomNav({ base }: BottomNavProps) {
  return (
    <nav className="fixed inset-x-0 bottom-0 z-30 border-t border-line bg-white/95 pb-[env(safe-area-inset-bottom)] backdrop-blur lg:hidden">
      <div className="mx-auto grid max-w-lg grid-cols-5">
        {patientTabs.map((tab) => (
          <NavLink
            key={tab.label}
            to={tab.path ? `${base}/${tab.path}` : base}
            end={tab.path === ''}
            className={({ isActive }) =>
              clsx(
                'flex flex-col items-center gap-0.5 pb-2 pt-1.5 text-[11px] font-medium',
                isActive ? 'text-accent-deep' : 'text-subtle',
              )
            }
          >
            {({ isActive }) => (
              <>
                <span
                  className={clsx(
                    'flex h-7 w-12 items-center justify-center rounded-pill transition-colors',
                    isActive && 'bg-accent-soft',
                  )}
                >
                  <tab.icon className="size-5" aria-hidden />
                </span>
                {tab.short}
              </>
            )}
          </NavLink>
        ))}
      </div>
    </nav>
  )
}
