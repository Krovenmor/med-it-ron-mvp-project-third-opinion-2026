import clsx from 'clsx'
import { useState } from 'react'

import type { CallbackRequest } from '@/api/models'
import { formatDateTime, toDateTimeLocal } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

const hour = 60 * 60_000

interface Preset {
  label: string
  at: (now: Date) => Date
}

const presets: Preset[] = [
  { label: 'Через час', at: (now) => new Date(now.getTime() + hour) },
  { label: 'Через 3 часа', at: (now) => new Date(now.getTime() + 3 * hour) },
  {
    label: 'Завтра в 10:00',
    at: (now) => {
      const tomorrow = new Date(now)
      tomorrow.setDate(tomorrow.getDate() + 1)
      tomorrow.setHours(10, 0, 0, 0)
      return tomorrow
    },
  },
]

interface CallbackModalProps {
  now: Date
  comment: string
  pending: boolean
  onClose: () => void
  onSubmit: (body: CallbackRequest) => void
}

export function CallbackModal({ now, comment, pending, onClose, onSubmit }: CallbackModalProps) {
  const [value, setValue] = useState(() => toDateTimeLocal(presets[0].at(now)))
  const at = new Date(value)
  const valid = !Number.isNaN(at.getTime()) && at.getTime() > now.getTime()

  return (
    <Modal
      open
      title="Когда перезвонить?"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Отмена
          </Button>
          <Button onClick={() => onSubmit({ call_at: at.toISOString(), comment })} disabled={!valid || pending}>
            Перезвонить {valid ? formatDateTime(at) : ''}
          </Button>
        </>
      }
    >
      <div className="flex flex-wrap gap-2">
        {presets.map((preset) => {
          const presetValue = toDateTimeLocal(preset.at(now))
          return (
            <button
              key={preset.label}
              type="button"
              onClick={() => setValue(presetValue)}
              className={clsx(
                'rounded-pill border px-3 py-1.5 text-sm font-medium',
                value === presetValue ? 'border-accent-deep bg-accent' : 'border-line hover:bg-card',
              )}
            >
              {preset.label}
            </button>
          )
        })}
      </div>
      <input
        type="datetime-local"
        value={value}
        onChange={(event) => setValue(event.target.value)}
        className="mt-4 w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
      />
      {!valid && <p className="mt-2 text-sm text-overdue">Время звонка должно быть в будущем</p>}
    </Modal>
  )
}
