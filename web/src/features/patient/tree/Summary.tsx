import clsx from 'clsx'
import { CircleCheck, Hourglass, Phone } from 'lucide-react'

import type { PatientRoute, RouteStep } from '@/api/models'
import { capitalize, formatDate, formatLongDate, formatTime } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { buttonStyles } from '@/shared/ui/buttonStyles'
import { Card } from '@/shared/ui/Layout'

import { phoneHref, type SheetMode } from '../models'
import { Notice } from '../ui/Notice'

interface StatusBannersProps {
  route: PatientRoute
}

export function StatusBanners({ route }: StatusBannersProps) {
  return (
    <>
      {route.urgent_contact && (
        <Notice tone="urgent" icon={<Phone className="size-5" />} className="p-5">
          <h2 className="text-lg font-medium">Врач просит срочно связаться с клиникой</h2>
          <p className="mt-1">
            Позвоните, пожалуйста, сегодня: администратор ответит на вопросы и запишет вас. Телефон{' '}
            <span className="whitespace-nowrap font-medium">{route.clinic.phone}</span>
          </p>
          <a href={phoneHref(route.clinic)} className={clsx(buttonStyles(), 'mt-4 w-full sm:w-auto')}>
            <Phone className="size-4" aria-hidden />
            Позвонить в клинику
          </a>
        </Notice>
      )}
      {route.pending_review && (
        <Notice tone="neutral" icon={<Hourglass className="size-5 text-subtle" />} className="p-5">
          <h2 className="text-lg font-medium">Врач готовит ваш план</h2>
          <p className="mt-1 text-subtle">
            Мы пришлём уведомление, когда план будет готов. Пока на дереве – ваши приёмы и исследование.
          </p>
        </Notice>
      )}
      {route.no_findings && (
        <Notice tone="accent" icon={<CircleCheck className="size-5" />} className="p-5">
          <h2 className="text-lg font-medium">Отклонений не найдено</h2>
          <p className="mt-1 text-subtle">
            Врач проверил результаты исследования. Дальше – плановый профилактический контроль, мы напомним о нём.
          </p>
        </Notice>
      )}
    </>
  )
}

interface ProgressCardProps {
  progress: number | null
}

export function ProgressCard({ progress }: ProgressCardProps) {
  return (
    <Card className="flex flex-col gap-3 p-5 sm:p-6">
      <div className="flex items-baseline justify-between gap-4">
        <h2 className="text-lg font-medium">Прогресс плана</h2>
        <span className="text-[28px] font-medium leading-none tabular-nums">{progress === null ? '–' : `${progress}%`}</span>
      </div>
      <div
        role="progressbar"
        aria-label="Прогресс плана"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={progress ?? undefined}
        className="h-2.5 overflow-hidden rounded-pill bg-white"
      >
        <div className="h-full rounded-pill bg-accent transition-[width] duration-700" style={{ width: `${progress ?? 0}%` }} />
      </div>
      <p className="text-sm text-subtle">
        Выполнение плана, а не состояние здоровья.
        {progress === null && ' Появится, когда наступит срок первого шага.'}
      </p>
    </Card>
  )
}

interface NextStepCardProps {
  step: RouteStep
  overdueCount: number
  onOpen: (id: string, mode?: SheetMode) => void
}

export function NextStepCard({ step, overdueCount, onOpen }: NextStepCardProps) {
  const overdue = step.status === 'overdue'
  const booked = step.status === 'booked' && step.appointment

  return (
    <section className={clsx('flex flex-col gap-4 rounded-card p-5 sm:p-6', overdue ? 'bg-overdue-tint' : 'bg-accent-tint')}>
      <div>
        <div className={clsx('text-sm font-semibold', overdue ? 'text-overdue' : 'text-accent-deep')}>
          {overdue ? 'Срок прошёл, запишитесь' : booked ? 'Ближайший приём' : 'Ближайший шаг'}
        </div>
        <h2 className="mt-1 text-xl font-medium">{step.title}</h2>
        <p className="mt-1 text-subtle">
          {booked
            ? `${capitalize(formatLongDate(booked.scheduled_at))}, ${formatTime(booked.scheduled_at)}`
            : step.due_at && `${overdue ? 'срок был' : 'не позднее'} ${formatDate(step.due_at)}`}
        </p>
        {overdueCount > 1 && <p className="mt-1 text-sm text-overdue">Ещё просроченных шагов: {overdueCount - 1}</p>}
      </div>
      <div className="flex gap-2">
        {step.bookable && (
          <Button onClick={() => onOpen(step.id, 'book')} className="flex-1 sm:flex-none">
            Записаться
          </Button>
        )}
        <Button variant="secondary" onClick={() => onOpen(step.id)} className="flex-1 sm:flex-none">
          Подробнее
        </Button>
      </div>
    </section>
  )
}

export function Disclaimer() {
  return (
    <p className="text-sm text-subtle">
      Рекомендации подготовил врач на основе вашего заключения. Они не заменяют консультацию.
    </p>
  )
}
