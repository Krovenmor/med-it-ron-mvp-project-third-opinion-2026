import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

import { confirmationText } from './models'

interface ConfirmModalProps {
  noFurtherExamination: boolean
  pending: boolean
  onClose: () => void
  onConfirm: () => void
}

export function ConfirmModal({ noFurtherExamination, pending, onClose, onConfirm }: ConfirmModalProps) {
  return (
    <Modal
      open
      title="Подтвердить рекомендации"
      onClose={onClose}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>
            Вернуться к таблице
          </Button>
          <Button onClick={onConfirm} disabled={pending}>
            Подтвердить
          </Button>
        </>
      }
    >
      {noFurtherExamination && (
        <div className="mb-4 rounded-control bg-accent-tint px-4 py-3 font-medium text-accent-deep">
          Дообследование не требуется
        </div>
      )}
      <p className="leading-relaxed text-subtle">{confirmationText}</p>
    </Modal>
  )
}
