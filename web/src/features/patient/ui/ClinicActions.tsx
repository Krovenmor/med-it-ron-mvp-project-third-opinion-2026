import clsx from 'clsx'
import { CalendarPlus, Copy, MapPin, Phone } from 'lucide-react'
import { toast } from 'sonner'

import type { Appointment, Clinic } from '@/api/models'
import { buttonStyles } from '@/shared/ui/buttonStyles'

import { downloadCalendarEvent } from '../calendar'
import { mapsHref, phoneHref } from '../models'

interface ClinicProps {
  clinic: Clinic
  className?: string
}

export function CallClinicButton({ clinic, className }: ClinicProps) {
  return (
    <a href={phoneHref(clinic)} className={clsx(buttonStyles('secondary'), className)}>
      <Phone className="size-4" aria-hidden />
      Позвонить в клинику
    </a>
  )
}

export function MapButton({ clinic, className }: ClinicProps) {
  return (
    <a href={mapsHref(clinic)} target="_blank" rel="noreferrer" className={clsx(buttonStyles('secondary'), className)}>
      <MapPin className="size-4" aria-hidden />
      Маршрут проезда
    </a>
  )
}

interface CalendarButtonProps extends ClinicProps {
  appointment: Appointment
}

export function CalendarButton({ appointment, clinic, className }: CalendarButtonProps) {
  return (
    <button
      type="button"
      onClick={() => downloadCalendarEvent(appointment, clinic)}
      className={clsx(buttonStyles('secondary'), className)}
    >
      <CalendarPlus className="size-4" aria-hidden />
      Добавить в календарь
    </button>
  )
}

export function ClinicPhone({ clinic, className }: ClinicProps) {
  const copy = () =>
    navigator.clipboard
      .writeText(clinic.phone)
      .then(() => toast.success('Номер скопирован'))
      .catch(() => toast.error('Не удалось скопировать номер'))

  return (
    <div className={clsx('flex items-center gap-2', className)}>
      <a href={phoneHref(clinic)} className="group flex min-w-0 items-center gap-3">
        <span className="flex size-10 shrink-0 items-center justify-center rounded-full bg-accent-tint text-accent-deep">
          <Phone className="size-4" aria-hidden />
        </span>
        <span>
          <span className="block font-medium text-accent-deep group-hover:underline">Позвонить в клинику</span>
          <span className="block text-sm tabular-nums text-subtle">{clinic.phone}</span>
        </span>
      </a>
      <button
        type="button"
        onClick={copy}
        className="rounded-control p-1.5 text-subtle hover:bg-card hover:text-ink"
        aria-label="Скопировать номер"
        title="Скопировать номер"
      >
        <Copy className="size-4" />
      </button>
    </div>
  )
}
