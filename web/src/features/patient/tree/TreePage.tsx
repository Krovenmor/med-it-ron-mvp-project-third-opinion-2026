import { Info, Sprout } from 'lucide-react'
import { useState } from 'react'
import { useSearchParams } from 'react-router'

import { Button } from '@/shared/ui/Button'
import { EmptyState } from '@/shared/ui/States'

import { recommendations } from '../models'
import { ClinicPhone } from '../ui/ClinicActions'
import { usePatient } from '../usePatient'

import { HealthTree } from './HealthTree'
import { Legend } from './Legend'
import { Disclaimer, NextStepCard, ProgressCard, StatusBanners } from './Summary'

export function TreePage() {
  const { patientId, route, openNode } = usePatient()
  const [params] = useSearchParams()
  const [legendOpen, setLegendOpen] = useState(false)
  const plan = recommendations(route)
  const next = route.steps.find((step) => step.id === route.next_step_id)
  const overdueCount = plan.filter((step) => step.status === 'overdue').length

  const contacts = (
    <div className="flex flex-col gap-3">
      <ClinicPhone clinic={route.clinic} />
      <Disclaimer />
    </div>
  )

  return (
    <div className="grid gap-5 lg:grid-cols-[minmax(0,380px)_minmax(0,1fr)] lg:items-start lg:gap-8">
      <div className="flex flex-col gap-5 lg:sticky lg:top-6">
        <StatusBanners route={route} />
        {!route.urgent_contact && plan.length > 0 && !(route.no_findings && route.progress === null) && (
          <ProgressCard progress={route.progress} />
        )}
        {!route.urgent_contact && next && <NextStepCard step={next} overdueCount={overdueCount} onOpen={openNode} />}
        <div className="hidden lg:block">{contacts}</div>
      </div>

      <div className="flex flex-col gap-5">
        <section className="-mx-4 border-y border-line sm:mx-0 sm:rounded-card sm:border">
          <header className="flex items-center justify-between gap-3 px-4 pb-1 pt-4 sm:px-5">
            <div>
              <h2 className="text-lg font-medium">Дерево маршрута</h2>
              <p className="text-sm text-subtle">Время идёт снизу вверх</p>
            </div>
            <Button size="sm" variant="ghost" onClick={() => setLegendOpen(true)}>
              <Info className="size-4" aria-hidden />
              Легенда
            </Button>
          </header>
          {route.trunk.length === 0 ? (
            <EmptyState icon={Sprout} title="Дерево появится после первого приёма" />
          ) : (
            <div className="px-2 pb-2 sm:px-3">
              <HealthTree key={patientId} route={route} selectedId={params.get('step')} onSelect={(id) => openNode(id)} />
            </div>
          )}
        </section>
        <div className="lg:hidden">{contacts}</div>
      </div>

      <Legend open={legendOpen} onClose={() => setLegendOpen(false)} />
    </div>
  )
}
