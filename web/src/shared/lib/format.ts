const dateFormat = new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
const timeFormat = new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit' })
const slotFormat = new Intl.DateTimeFormat('ru-RU', {
  weekday: 'short',
  day: 'numeric',
  month: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

const shortDateTimeFormat = new Intl.DateTimeFormat('ru-RU', {
  day: '2-digit',
  month: '2-digit',
  hour: '2-digit',
  minute: '2-digit',
})
const longDateFormat = new Intl.DateTimeFormat('ru-RU', { weekday: 'long', day: 'numeric', month: 'long' })

const minute = 60_000
const hour = 60 * minute
const day = 24 * hour

export function formatDate(value: string | Date): string {
  return dateFormat.format(new Date(value))
}

export function formatDateTime(value: string | Date): string {
  const date = new Date(value)
  return `${dateFormat.format(date)} ${timeFormat.format(date)}`
}

export function formatSlot(value: string | Date): string {
  return slotFormat.format(new Date(value))
}

export function formatShortDateTime(value: string | Date): string {
  return shortDateTimeFormat.format(new Date(value))
}

export function formatLongDate(value: string | Date): string {
  return longDateFormat.format(new Date(value))
}

export function capitalize(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1)
}

export function formatTime(value: string | Date): string {
  return timeFormat.format(new Date(value))
}

export function toDateTimeLocal(date: Date): string {
  const pad = (value: number) => value.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function formatDuration(ms: number): string {
  const total = Math.max(ms, 0)
  if (total < minute) {
    return 'меньше минуты'
  }
  if (total < hour) {
    return `${Math.floor(total / minute)} мин`
  }
  if (total < day) {
    return `${Math.floor(total / hour)} ч ${Math.floor((total % hour) / minute)} мин`
  }
  return `${Math.floor(total / day)} дн ${Math.floor((total % day) / hour)} ч`
}

export function formatCountdown(ms: number): string {
  const totalSeconds = Math.floor(Math.abs(ms) / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

export function ageAt(birthDate: string, now: Date): number {
  const [year, month, dayOfMonth] = birthDate.split('-').map(Number)
  let age = now.getFullYear() - year
  const beforeBirthday = now.getMonth() + 1 < month || (now.getMonth() + 1 === month && now.getDate() < dayOfMonth)
  if (beforeBirthday) {
    age--
  }
  return age
}

export function pluralize(count: number, one: string, few: string, many: string): string {
  const lastTwo = count % 100
  const last = count % 10
  if (lastTwo >= 11 && lastTwo <= 14) {
    return many
  }
  if (last === 1) {
    return one
  }
  if (last >= 2 && last <= 4) {
    return few
  }
  return many
}

export function formatAge(age: number): string {
  return `${age} ${pluralize(age, 'год', 'года', 'лет')}`
}
