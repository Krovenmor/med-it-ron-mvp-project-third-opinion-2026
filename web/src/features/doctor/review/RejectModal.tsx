import clsx from 'clsx'
import { useState } from 'react'

import type { Recommendation, RejectReason, ReviewRecommendationRequest } from '@/api/models'
import { rejectReasonLabels } from '@/shared/lib/labels'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

const reasons: RejectReason[] = ['contraindicated', 'other']

interface RejectModalProps {
  recommendation: Recommendation
  pending: boolean
  onClose: () => void
  onSubmit: (body: ReviewRecommendationRequest) => void
}

export function RejectModal({ recommendation, pending, onClose, onSubmit }: RejectModalProps) {
  const [reason, setReason] = useState<RejectReason | null>(null)
  const [comment, setComment] = useState('')
  const valid = reason === 'contraindicated' || (reason === 'other' && comment.trim() !== '')

  const submit = () => {
    if (!reason) {
      return
    }
    onSubmit({
      mark: 'rejected',
      reject_reason: reason,
      reject_comment: reason === 'other' ? comment.trim() : undefined,
    })
  }

  return (
    <Modal
      open
      title="Почему услуга не нужна?"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Отмена
          </Button>
          <Button onClick={submit} disabled={!valid || pending}>
            Отметить «Не нужна»
          </Button>
        </>
      }
    >
      <p className="mb-4 text-subtle">{recommendation.service_name}</p>
      <div className="flex flex-col gap-2">
        {reasons.map((value) => (
          <label
            key={value}
            className={clsx(
              'flex cursor-pointer items-center gap-3 rounded-control border px-4 py-3',
              reason === value ? 'border-accent bg-accent-tint' : 'border-line hover:bg-card',
            )}
          >
            <input
              type="radio"
              name="reject-reason"
              checked={reason === value}
              onChange={() => setReason(value)}
              className="accent-accent-deep"
            />
            {rejectReasonLabels[value]}
          </label>
        ))}
      </div>
      {reason === 'other' && (
        <textarea
          autoFocus
          value={comment}
          onChange={(event) => setComment(event.target.value)}
          placeholder="Укажите причину"
          rows={3}
          className="mt-4 w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
        />
      )}
    </Modal>
  )
}
