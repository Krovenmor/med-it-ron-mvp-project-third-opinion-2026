import type { Appointment, Clinic } from '@/api/models'

const visitMs = 30 * 60_000

export function downloadCalendarEvent(appointment: Appointment, clinic: Clinic): void {
  const start = new Date(appointment.scheduled_at)
  const lines = [
    'BEGIN:VCALENDAR',
    'VERSION:2.0',
    'PRODID:-//Meditron//Patient route//RU',
    'CALSCALE:GREGORIAN',
    'METHOD:PUBLISH',
    'BEGIN:VEVENT',
    `UID:${appointment.id}@meditron`,
    `DTSTAMP:${icsTime(new Date())}`,
    `DTSTART:${icsTime(start)}`,
    `DTEND:${icsTime(new Date(start.getTime() + visitMs))}`,
    `SUMMARY:${icsText(appointment.service_name)}`,
    `LOCATION:${icsText(`${clinic.name}, ${clinic.address}`)}`,
    `DESCRIPTION:${icsText(`Телефон клиники: ${clinic.phone}`)}`,
    'BEGIN:VALARM',
    'TRIGGER:-PT2H',
    'ACTION:DISPLAY',
    `DESCRIPTION:${icsText(appointment.service_name)}`,
    'END:VALARM',
    'END:VEVENT',
    'END:VCALENDAR',
  ]
  const url = URL.createObjectURL(new Blob([lines.join('\r\n')], { type: 'text/calendar;charset=utf-8' }))
  const link = document.createElement('a')
  link.href = url
  link.download = `${appointment.id}.ics`
  link.click()
  URL.revokeObjectURL(url)
}

function icsTime(date: Date): string {
  return date.toISOString().replace(/[-:]/g, '').replace(/\.\d{3}/, '')
}

function icsText(value: string): string {
  return value.replace(/\\/g, '\\\\').replace(/;/g, '\\;').replace(/,/g, '\\,').replace(/\n/g, '\\n')
}
