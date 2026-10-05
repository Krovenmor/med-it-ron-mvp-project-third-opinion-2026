import clsx from 'clsx'
import { useState } from 'react'

import type { DeclineReason, DeclineRequest } from '@/api/models'
import { declineReasonLabels, declineReasons } from '@/shared/lib/labels'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

interface DeclineModalProps {
  initialComment: string
  pending: boolean
  onClose: () => void
  onSubmit: (body: DeclineRequest) => void
}

export function DeclineModal({ initialComment, pending, onClose, onSubmit }: DeclineModalProps) {
  const [reason, setReason] = useState<DeclineReason | null>(null)
  const [comment, setComment] = useState(initialComment)
  const valid = reason !== null && (reason !== 'other' || comment.trim() !== '')

  return (
    <Modal
      open
      title="Почему пациент отказался?"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Отмена
          </Button>
          <Button onClick={() => reason && onSubmit({ reason, comment: comment.trim() })} disabled={!valid || pending}>
            Записать отказ
          </Button>
        </>
      }
    >
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
              name="decline-reason"
              checked={reason === value}
              onChange={() => setReason(value)}
              className="accent-accent-deep"
            />
            {declineReasonLabels[value]}
          </label>
        ))}
      </div>
      <textarea
        value={comment}
        onChange={(event) => setComment(event.target.value)}
        placeholder={reason === 'other' ? 'Укажите причину' : 'Комментарий, если нужен'}
        rows={3}
        className="mt-4 w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
      />
    </Modal>
  )
}
