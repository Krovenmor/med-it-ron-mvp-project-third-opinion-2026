import { CalendarX2 } from 'lucide-react'

import type { Appointment, Clinic, RouteStep } from '@/api/models'
import { formatTime } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Layout'
import { EmptyState } from '@/shared/ui/States'

import { CalendarButton, ClinicPhone, MapButton } from '../ui/ClinicActions'
import { PageTitle } from '../ui/PageTitle'
import { usePatient } from '../usePatient'

const weekdayFormat = new Intl.DateTimeFormat('ru-RU', { weekday: 'short' })
const monthFormat = new Intl.DateTimeFormat('ru-RU', { month: 'short' })

export function AppointmentsPage() {
  const { route, openNode } = usePatient()

  return (
    <div className="flex max-w-3xl flex-col gap-5">
      <PageTitle main="Записи" rest="предстоящие приёмы" />
      {route.appointments.length === 0 ? (
        <EmptyState
          icon={CalendarX2}
          title="Предстоящих записей нет"
          text="Записаться можно из карточки шага на дереве или во вкладке «Мои шаги»."
        />
      ) : (
        route.appointments.map((appointment) => (
          <AppointmentCard
            key={appointment.id}
            appointment={appointment}
            clinic={route.clinic}
            step={route.steps.find((step) => step.appointment?.id === appointment.id)}
            onOpen={openNode}
          />
        ))
      )}
      <div className="flex flex-col gap-2 pt-2">
        <p className="text-sm text-subtle">Перенести или отменить запись можно по телефону клиники.</p>
        <ClinicPhone clinic={route.clinic} />
      </div>
    </div>
  )
}

interface AppointmentCardProps {
  appointment: Appointment
  clinic: Clinic
  step: RouteStep | undefined
  onOpen: (id: string) => void
}

function AppointmentCard({ appointment, clinic, step, onOpen }: AppointmentCardProps) {
  const date = new Date(appointment.scheduled_at)
  return (
    <Card className="flex flex-col gap-5 p-5 sm:p-6">
      <div className="flex gap-4">
        <div className="flex w-16 shrink-0 flex-col items-center justify-center rounded-control bg-white py-2">
          <span className="text-xs uppercase text-subtle">{weekdayFormat.format(date)}</span>
          <span className="text-2xl font-medium leading-tight">{date.getDate()}</span>
          <span className="text-xs text-subtle">{monthFormat.format(date)}</span>
        </div>
        <div className="min-w-0 flex-1">
          <div className="text-lg font-medium">{appointment.service_name}</div>
          <div className="mt-0.5">в {formatTime(date)}</div>
          <div className="mt-0.5 text-sm text-subtle">
            {clinic.name}, {clinic.address}
          </div>
        </div>
      </div>
      <div className="grid gap-2 sm:flex sm:flex-wrap">
        <MapButton clinic={clinic} />
        <CalendarButton appointment={appointment} clinic={clinic} />
        {step && (
          <Button variant="ghost" onClick={() => onOpen(step.id)}>
            Подробнее
          </Button>
        )}
      </div>
    </Card>
  )
}
