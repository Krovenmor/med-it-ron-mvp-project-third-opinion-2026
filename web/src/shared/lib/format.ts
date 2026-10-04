const dateFormat = new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
const timeFormat = new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit' })

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

export function formatAge(age: number): string {
  const lastTwo = age % 100
  const last = age % 10
  if (lastTwo >= 11 && lastTwo <= 14) {
    return `${age} лет`
  }
  if (last === 1) {
    return `${age} год`
  }
  if (last >= 2 && last <= 4) {
    return `${age} года`
  }
  return `${age} лет`
}
