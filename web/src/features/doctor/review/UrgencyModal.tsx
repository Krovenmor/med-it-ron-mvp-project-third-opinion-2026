import clsx from 'clsx'
import { useState } from 'react'

import type { ChangeUrgencyRequest, Urgency } from '@/api/models'
import { urgencyOrder, urgencyRank } from '@/shared/lib/labels'
import { UrgencyBadge } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

interface UrgencyModalProps {
  current: Urgency
  pending: boolean
  onClose: () => void
  onSubmit: (body: ChangeUrgencyRequest) => void
}

export function UrgencyModal({ current, pending, onClose, onSubmit }: UrgencyModalProps) {
  const [urgency, setUrgency] = useState<Urgency>(current)
  const [reason, setReason] = useState('')
  const lowering = urgencyRank(urgency) < urgencyRank(current)
  const valid = urgency !== current && (!lowering || reason.trim() !== '')

  return (
    <Modal
      open
      title="Изменить срочность"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Отмена
          </Button>
          <Button onClick={() => onSubmit({ urgency, reason: reason.trim() })} disabled={!valid || pending}>
            Сохранить
          </Button>
        </>
      }
    >
      <div className="flex flex-col gap-2">
        {urgencyOrder.map((value) => (
          <label
            key={value}
            className={clsx(
              'flex cursor-pointer items-center justify-between rounded-control border px-4 py-3',
              urgency === value ? 'border-accent bg-accent-tint' : 'border-line hover:bg-card',
            )}
          >
            <span className="flex items-center gap-3">
              <input
                type="radio"
                name="urgency"
                checked={urgency === value}
                onChange={() => setUrgency(value)}
                className="accent-accent-deep"
              />
              <UrgencyBadge urgency={value} />
            </span>
            {value === current && <span className="text-sm text-subtle">сейчас</span>}
          </label>
        ))}
      </div>
      <label className="mt-4 flex flex-col gap-2">
        <span className="text-sm text-subtle">
          Причина {lowering ? '(обязательна при понижении)' : '(необязательно)'}
        </span>
        <textarea
          value={reason}
          onChange={(event) => setReason(event.target.value)}
          rows={3}
          className="w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
        />
      </label>
    </Modal>
  )
}
