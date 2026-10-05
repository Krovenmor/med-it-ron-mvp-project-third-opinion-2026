import clsx from 'clsx'
import { useState, type ReactNode } from 'react'
import { toast } from 'sonner'

import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Layout'

import { Notice } from '../ui/Notice'
import { PageTitle } from '../ui/PageTitle'

const hours = Array.from({ length: 24 }, (_, hour) => `${hour.toString().padStart(2, '0')}:00`)

export function SettingsPage() {
  const [push, setPush] = useState(true)
  const [sms, setSms] = useState(true)
  const [consent, setConsent] = useState(true)
  const [quietFrom, setQuietFrom] = useState('22:00')
  const [quietTo, setQuietTo] = useState('08:00')

  const revoke = () => {
    setConsent(false)
    setPush(false)
    setSms(false)
    toast.info('Согласие отозвано: уведомления больше не придут')
  }

  return (
    <div className="flex max-w-2xl flex-col gap-5">
      <PageTitle main="Настройки" rest="уведомления" />
      <Notice tone="neutral">Это макет: в демо настройки не сохраняются.</Notice>

      <Card className="flex flex-col gap-1 p-5 sm:p-6">
        <h2 className="mb-2 text-lg font-medium">Каналы</h2>
        <SwitchRow label="Push-уведомления" hint="В приложении клиники" checked={push && consent} disabled={!consent} onChange={setPush} />
        <SwitchRow label="СМС" hint="Если push не дошёл" checked={sms && consent} disabled={!consent} onChange={setSms} />
      </Card>

      <Card className="flex flex-col gap-4 p-5 sm:p-6">
        <SwitchRow
          label="Согласие на уведомления"
          hint="Напоминания о шагах плана и записи"
          checked={consent}
          onChange={setConsent}
        />
        <div className="flex flex-wrap items-center gap-3">
          <span className="w-full text-subtle sm:w-auto">Тихие часы</span>
          <HourSelect label="с" value={quietFrom} onChange={setQuietFrom} disabled={!consent} />
          <HourSelect label="до" value={quietTo} onChange={setQuietTo} disabled={!consent} />
        </div>
      </Card>

      <Button variant="secondary" onClick={revoke} disabled={!consent} className="self-start">
        Отозвать согласие
      </Button>
    </div>
  )
}

interface SwitchRowProps {
  label: string
  hint: string
  checked: boolean
  disabled?: boolean
  onChange: (value: boolean) => void
}

function SwitchRow({ label, hint, checked, disabled = false, onChange }: SwitchRowProps) {
  return (
    <label className={clsx('flex items-center justify-between gap-4 py-2', disabled ? 'opacity-50' : 'cursor-pointer')}>
      <span>
        <span className="block font-medium">{label}</span>
        <span className="text-sm text-subtle">{hint}</span>
      </span>
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        aria-label={label}
        disabled={disabled}
        onClick={() => onChange(!checked)}
        className={clsx(
          'relative h-7 w-12 shrink-0 rounded-pill transition-colors disabled:cursor-not-allowed',
          checked ? 'bg-accent' : 'bg-line',
        )}
      >
        <span
          className={clsx(
            'absolute left-1 top-1 size-5 rounded-full bg-white shadow-sm transition-transform',
            checked && 'translate-x-5',
          )}
        />
      </button>
    </label>
  )
}

interface HourSelectProps {
  label: ReactNode
  value: string
  disabled: boolean
  onChange: (value: string) => void
}

function HourSelect({ label, value, disabled, onChange }: HourSelectProps) {
  return (
    <label className="flex items-center gap-2">
      <span className="text-subtle">{label}</span>
      <select
        value={value}
        disabled={disabled}
        onChange={(event) => onChange(event.target.value)}
        className="h-10 rounded-control border border-line bg-white px-3 outline-none focus:border-accent disabled:opacity-50"
      >
        {hours.map((hour) => (
          <option key={hour}>{hour}</option>
        ))}
      </select>
    </label>
  )
}
