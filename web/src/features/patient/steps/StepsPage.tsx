import { ListChecks } from 'lucide-react'

import type { RouteStep } from '@/api/models'
import { Button } from '@/shared/ui/Button'
import { EmptyState } from '@/shared/ui/States'

import { isOpen, recommendations, type SheetMode } from '../models'
import { stepCaption } from '../tree/captions'
import { stepLook } from '../tree/looks'
import { NodeGlyph } from '../tree/marks'
import { PageTitle } from '../ui/PageTitle'
import { MarkBadge } from '../ui/StepBadges'
import { usePatient } from '../usePatient'

interface StepGroup {
  title: string
  steps: RouteStep[]
}

export function StepsPage() {
  const { route, openNode } = usePatient()
  const plan = recommendations(route)
  const byDate = (a: RouteStep, b: RouteStep) => Date.parse(a.date) - Date.parse(b.date)

  const groups: StepGroup[] = [
    {
      title: 'Нужно записаться',
      steps: plan
        .filter(isOpen)
        .sort((a, b) => Number(b.status === 'overdue') - Number(a.status === 'overdue') || byDate(a, b)),
    },
    { title: 'Вы записаны', steps: plan.filter((step) => step.status === 'booked').sort(byDate) },
    { title: 'Пройдено', steps: plan.filter((step) => step.status === 'done').sort((a, b) => byDate(b, a)) },
    { title: 'Вы отказались', steps: plan.filter((step) => step.status === 'declined').sort(byDate) },
  ].filter((group) => group.steps.length > 0)

  return (
    <div className="flex max-w-3xl flex-col gap-6">
      <PageTitle main="Мои шаги" rest="по сроку" />
      {groups.length === 0 ? (
        <EmptyState
          icon={ListChecks}
          title={route.pending_review ? 'Врач готовит ваш план' : 'Шагов пока нет'}
          text="Здесь появятся назначения врача: что пройти и до какого срока."
        />
      ) : (
        groups.map((group) => (
          <section key={group.title}>
            <h2 className="mb-1 text-sm font-medium text-subtle">
              {group.title} · {group.steps.length}
            </h2>
            <ul className="divide-y divide-line border-y border-line">
              {group.steps.map((step) => (
                <StepRow key={step.id} step={step} onOpen={openNode} />
              ))}
            </ul>
          </section>
        ))
      )}
    </div>
  )
}

interface StepRowProps {
  step: RouteStep
  onOpen: (id: string, mode?: SheetMode) => void
}

function StepRow({ step, onOpen }: StepRowProps) {
  const caption = stepCaption(step)
  return (
    <li className="flex flex-col gap-3 py-4 sm:flex-row sm:items-center">
      <button type="button" onClick={() => onOpen(step.id)} className="flex min-w-0 flex-1 items-start gap-3 text-left">
        <NodeGlyph look={stepLook(step)} size={26} className="mt-0.5" />
        <span className="min-w-0 flex-1">
          <span className="block font-medium hover:underline">{step.title}</span>
          <span className={caption.alert ? 'text-sm font-semibold text-overdue' : 'text-sm text-subtle'}>
            {caption.text}
          </span>
        </span>
      </button>
      <div className="flex items-center gap-3 pl-[38px] sm:pl-0">
        {step.mark && <MarkBadge mark={step.mark} />}
        {step.bookable && (
          <Button size="sm" onClick={() => onOpen(step.id, 'book')} className="ml-auto sm:ml-0">
            Записаться
          </Button>
        )}
      </div>
    </li>
  )
}
