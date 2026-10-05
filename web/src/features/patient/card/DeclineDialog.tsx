import clsx from 'clsx'
import { useState } from 'react'

import type { DeclineReason, PatientDeclineRequest } from '@/api/models'
import { declineReasons } from '@/shared/lib/labels'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

import { patientDeclineReasonLabels } from '../models'

interface DeclineDialogProps {
  pending: boolean
  onClose: () => void
  onSubmit: (body: PatientDeclineRequest) => void
}

export function DeclineDialog({ pending, onClose, onSubmit }: DeclineDialogProps) {
  const [reason, setReason] = useState<DeclineReason | null>(null)
  const [comment, setComment] = useState('')
  const valid = reason !== null && (reason !== 'other' || comment.trim() !== '')

  return (
    <Modal
      open
      title="Почему не будете записываться?"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Отмена
          </Button>
          <Button onClick={() => reason && onSubmit({ reason, comment: comment.trim() })} disabled={!valid || pending}>
            Отправить
          </Button>
        </>
      }
    >
      <p className="mb-4 text-subtle">Ответ увидят врач и клиника. Если передумаете, позвоните в клинику.</p>
      <div className="flex flex-col gap-2">
        {declineReasons.map((value) => (
          <label
            key={value}
            className={clsx(
              'flex cursor-pointer items-center gap-3 rounded-control border px-4 py-3',
              reason === value ? 'border-accent bg-accent-tint' : 'border-line hover:bg-card',
            )}
          >
            <input
              type="radio"
              name="patient-decline-reason"
              checked={reason === value}
              onChange={() => setReason(value)}
              className="accent-accent-deep"
            />
            {patientDeclineReasonLabels[value]}
          </label>
        ))}
      </div>
      <textarea
        value={comment}
        onChange={(event) => setComment(event.target.value)}
        placeholder={reason === 'other' ? 'Расскажите, что мешает' : 'Комментарий, если хотите'}
        rows={3}
        className="mt-4 w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
      />
    </Modal>
  )
}
