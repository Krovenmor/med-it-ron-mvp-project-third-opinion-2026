import { useMutation, useQueryClient } from '@tanstack/react-query'
import { CalendarCheck, LifeBuoy, PhoneCall, TriangleAlert } from 'lucide-react'
import { useState, type ReactNode } from 'react'
import { toast } from 'sonner'

import { b2c } from '@/api/b2c'
import { ApiError } from '@/api/http'
import type { Appointment, RouteStep, Slot } from '@/api/models'
import { queryKeys } from '@/api/queries'
import { capitalize, formatDate, formatLongDate, formatTime } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { Sheet } from '@/shared/ui/Sheet'

import { isOpen, patientDeclineReasonLabels, type PatientContext, type SheetMode } from '../models'
import { CalendarButton, CallClinicButton, MapButton } from '../ui/ClinicActions'
import { Notice } from '../ui/Notice'
import { MarkBadge, StatusPill } from '../ui/StepBadges'

import { DeclineDialog } from './DeclineDialog'
import { SlotPicker } from './SlotPicker'

type View = 'details' | 'slots' | 'confirm' | 'booked'

const titles: Record<View, string> = {
  details: 'Шаг плана',
  slots: 'Запись',
  confirm: 'Подтверждение записи',
  booked: 'Запись',
}

interface StepSheetProps {
  context: PatientContext
  step: RouteStep
  mode: SheetMode
  onClose: () => void
}

export function StepSheet({ context, step, mode, onClose }: StepSheetProps) {
  const { patientId, route } = context
  const queryClient = useQueryClient()
  const [view, setView] = useState<View>(mode === 'book' && step.bookable ? 'slots' : 'details')
  const [slot, setSlot] = useState<Slot | null>(null)
  const [appointment, setAppointment] = useState<Appointment | null>(null)
  const [helpSent, setHelpSent] = useState(false)
  const [declining, setDeclining] = useState(false)

  const refreshRoute = () => queryClient.invalidateQueries({ queryKey: queryKeys.patientRoute(patientId) })

  const book = useMutation({
    mutationFn: (slotId: string) => b2c.book(patientId, step.id, slotId),
    onSuccess: (created) => {
      setAppointment(created)
      setView('booked')
      void refreshRoute()
    },
    onError: (error) => {
      if (error instanceof ApiError && error.status === 409) {
        toast.error('Это время уже недоступно – выберите другое')
        setSlot(null)
        setView('slots')
        void queryClient.invalidateQueries({ queryKey: queryKeys.slots(patientId, step.id) })
        void refreshRoute()
        return
      }
      toast.error('Не удалось записаться. Попробуйте ещё раз или позвоните в клинику.')
    },
  })

  const help = useMutation({
    mutationFn: () => b2c.requestHelp(patientId, step.id),
    onSuccess: () => setHelpSent(true),
    onError: () => toast.error('Не удалось отправить просьбу. Позвоните в клинику.'),
  })

  const decline = useMutation({
    mutationFn: (body: Parameters<typeof b2c.decline>[2]) => b2c.decline(patientId, step.id, body),
    onSuccess: () => {
      setDeclining(false)
      toast.success('Спасибо, мы учли ваш ответ')
      void refreshRoute()
    },
    onError: () => toast.error('Не удалось отправить ответ. Попробуйте ещё раз.'),
  })

  const current: View = view === 'slots' || view === 'confirm' ? (step.bookable ? view : 'details') : view
  const title = current === 'details' && step.kind === 'visit' ? 'Приём' : titles[current]

  const footer = (() => {
    switch (current) {
      case 'slots':
        return (
          <>
            <Button variant="secondary" onClick={() => setView('details')}>
              Назад
            </Button>
            <Button onClick={() => setView('confirm')} disabled={!slot}>
              Далее
            </Button>
          </>
        )
      case 'confirm':
        return (
          <>
            <Button variant="secondary" onClick={() => setView('slots')} disabled={book.isPending}>
              Назад
            </Button>
            <Button onClick={() => slot && book.mutate(slot.id)} disabled={!slot || book.isPending}>
              Записаться
            </Button>
          </>
        )
      case 'booked':
        return (
          <>
            {appointment && <CalendarButton appointment={appointment} clinic={route.clinic} />}
            <MapButton clinic={route.clinic} />
            <Button onClick={onClose}>Готово</Button>
          </>
        )
      default:
        return <DetailsFooter context={context} step={step} helpSent={helpSent} helpPending={help.isPending} onBook={() => setView('slots')} onHelp={() => help.mutate()} onDecline={() => setDeclining(true)} onClose={onClose} />
    }
  })()

  return (
    <Sheet title={title} onClose={onClose} footer={footer}>
      {current === 'slots' && (
        <>
          <h3 className="mb-4 text-xl font-medium">{step.title}</h3>
          <SlotPicker patientId={patientId} recommendationId={step.id} selected={slot} onSelect={setSlot} />
        </>
      )}
      {current === 'confirm' && slot && (
        <div className="flex flex-col gap-4">
          <h3 className="text-2xl font-medium tracking-tight">{step.title}</h3>
          <dl className="flex flex-col gap-3 rounded-card bg-card p-5">
            <Fact label="Когда">
              {capitalize(formatLongDate(slot.starts_at))}, {formatTime(slot.starts_at)}
            </Fact>
            <Fact label="Где">
              {route.clinic.name}, {route.clinic.address}
            </Fact>
          </dl>
        </div>
      )}
      {current === 'booked' && appointment && (
        <div className="flex flex-col items-center gap-3 py-6 text-center">
          <span className="flex size-16 items-center justify-center rounded-full bg-accent-soft">
            <CalendarCheck className="size-8 text-accent-deep" aria-hidden />
          </span>
          <h3 className="text-2xl font-medium tracking-tight">Вы записаны</h3>
          <p className="text-lg">
            {appointment.service_name}
            <br />
            {capitalize(formatLongDate(appointment.scheduled_at))}, {formatTime(appointment.scheduled_at)}
          </p>
          <p className="text-subtle">
            {route.clinic.name}, {route.clinic.address}
          </p>
        </div>
      )}
      {current === 'details' && <StepDetails context={context} step={step} helpSent={helpSent} />}

      {declining && (
        <DeclineDialog pending={decline.isPending} onClose={() => setDeclining(false)} onSubmit={(body) => decline.mutate(body)} />
      )}
    </Sheet>
  )
}

interface DetailsFooterProps {
  context: PatientContext
  step: RouteStep
  helpSent: boolean
  helpPending: boolean
  onBook: () => void
  onHelp: () => void
  onDecline: () => void
  onClose: () => void
}

function DetailsFooter({ context, step, helpSent, helpPending, onBook, onHelp, onDecline, onClose }: DetailsFooterProps) {
  const { clinic } = context.route
  if (step.status === 'booked' && step.appointment) {
    return (
      <>
        <CallClinicButton clinic={clinic} />
        <MapButton clinic={clinic} />
        <CalendarButton appointment={step.appointment} clinic={clinic} />
      </>
    )
  }
  if (isOpen(step) && step.bookable) {
    return (
      <>
        {step.status === 'recommended' && (
          <Button variant="ghost" onClick={onDecline}>
            Не буду записываться
          </Button>
        )}
        {!helpSent && (
          <Button variant="secondary" onClick={onHelp} disabled={helpPending}>
            <LifeBuoy className="size-4" aria-hidden />
            Нужна помощь
          </Button>
        )}
        <Button onClick={onBook}>Записаться</Button>
      </>
    )
  }
  if (isOpen(step) || step.status === 'declined') {
    return <CallClinicButton clinic={clinic} />
  }
  return (
    <Button variant="secondary" onClick={onClose}>
      Закрыть
    </Button>
  )
}

interface StepDetailsProps {
  context: PatientContext
  step: RouteStep
  helpSent: boolean
}

function StepDetails({ context, step, helpSent }: StepDetailsProps) {
  const { route } = context
  const recommendation = step.kind === 'recommendation'

  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-3">
        <div className="flex flex-wrap gap-2">
          <StatusPill status={step.status} />
          {step.mark && <MarkBadge mark={step.mark} />}
        </div>
        <h3 className="text-2xl font-medium tracking-tight">{step.title}</h3>
      </div>

      {step.status === 'overdue' && step.bookable && (
        <Notice tone="overdue" icon={<TriangleAlert className="size-5" />}>
          Срок прошёл, запишитесь. Если что-то мешает, нажмите «Нужна помощь» – администратор перезвонит.
        </Notice>
      )}
      {isOpen(step) && route.urgent_contact && (
        <Notice tone="urgent" icon={<PhoneCall className="size-5" />}>
          Врач просит сначала связаться с клиникой по телефону {route.clinic.phone}. Администратор поможет с записью.
        </Notice>
      )}
      {recommendation && !step.in_clinic && step.status !== 'done' && (
        <Notice tone="neutral">
          Этой услуги нет в нашей клинике. Обратитесь к профильному специалисту – администратор подскажет, куда.
        </Notice>
      )}
      {helpSent && (
        <Notice tone="accent" icon={<LifeBuoy className="size-5" />}>
          Мы передали просьбу администратору клиники – вам перезвонят и помогут с записью.
        </Notice>
      )}

      <dl className="flex flex-col gap-3">
        {isOpen(step) && step.due_at && (
          <Fact label="Срок">
            {step.status === 'overdue' ? 'был до' : 'не позднее'} {formatDate(step.due_at)}
          </Fact>
        )}
        {step.status === 'booked' && step.appointment && (
          <>
            <Fact label="Когда">
              {capitalize(formatLongDate(step.appointment.scheduled_at))}, {formatTime(step.appointment.scheduled_at)}
            </Fact>
            <Fact label="Где">
              {route.clinic.name}, {route.clinic.address}
            </Fact>
          </>
        )}
        {step.status === 'done' && step.visited_at && (
          <Fact label="Когда">
            {formatDate(step.visited_at)}
            {step.late && <span className="text-subtle"> · позже срока</span>}
          </Fact>
        )}
        {step.status === 'declined' && step.decline_reason && (
          <Fact label="Ваш ответ">{patientDeclineReasonLabels[step.decline_reason]}</Fact>
        )}
      </dl>

      {step.text && (
        <section>
          <h4 className="mb-1 text-sm font-medium text-subtle">Почему назначено</h4>
          <p className="leading-relaxed">{step.text}</p>
        </section>
      )}

      {step.status === 'booked' && (
        <p className="text-sm text-subtle">Перенести или отменить запись можно по телефону клиники {route.clinic.phone}.</p>
      )}
      {step.status === 'declined' && (
        <p className="text-sm text-subtle">Передумали? Позвоните в клинику – вас запишут.</p>
      )}
    </div>
  )
}

interface FactProps {
  label: string
  children: ReactNode
}

function Fact({ label, children }: FactProps) {
  return (
    <div className="grid grid-cols-[88px_minmax(0,1fr)] gap-3">
      <dt className="text-subtle">{label}</dt>
      <dd>{children}</dd>
    </div>
  )
}
