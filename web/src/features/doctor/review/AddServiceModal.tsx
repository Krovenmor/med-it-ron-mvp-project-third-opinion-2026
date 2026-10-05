import clsx from 'clsx'
import { Search } from 'lucide-react'
import { useDeferredValue, useState } from 'react'

import type { AcceptedMark, AddRecommendationRequest, CatalogService } from '@/api/models'
import { useServiceSearch } from '@/api/queries'
import { acceptedMarks, markLabels } from '@/shared/lib/labels'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

interface AddServiceModalProps {
  existingCodes: string[]
  pending: boolean
  onClose: () => void
  onSubmit: (body: AddRecommendationRequest) => void
}

export function AddServiceModal({ existingCodes, pending, onClose, onSubmit }: AddServiceModalProps) {
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<CatalogService | null>(null)
  const [patientText, setPatientText] = useState('')
  const [comment, setComment] = useState('')
  const [mark, setMark] = useState<AcceptedMark>('critical')
  const { data, isError } = useServiceSearch(useDeferredValue(query))

  const select = (service: CatalogService) => {
    setSelected(service)
    setPatientText(service.description)
  }

  const submit = () => {
    if (!selected) {
      return
    }
    onSubmit({
      service_code: selected.code,
      service_name: selected.name,
      patient_text: patientText.trim(),
      rationale: comment.trim(),
      mark,
    })
  }

  return (
    <Modal
      open
      title="Добавить исследование"
      onClose={onClose}
      footer={
        selected && (
          <>
            <Button variant="secondary" onClick={() => setSelected(null)}>
              Назад к справочнику
            </Button>
            <Button onClick={submit} disabled={pending || patientText.trim() === ''}>
              Добавить
            </Button>
          </>
        )
      }
    >
      {!selected ? (
        <div className="flex flex-col gap-4">
          <label className="flex h-12 items-center gap-3 rounded-control border border-line px-4 focus-within:border-accent">
            <Search className="size-5 text-muted" aria-hidden />
            <input
              autoFocus
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Поиск по справочнику услуг клиники"
              className="w-full outline-none"
            />
          </label>
          {isError && <div className="text-subtle">Справочник МИС недоступен, попробуйте позже.</div>}
          <ul className="flex max-h-80 flex-col gap-1 overflow-y-auto">
            {data?.services.map((service) => {
              const added = existingCodes.includes(service.code)
              return (
                <li key={service.code}>
                  <button
                    type="button"
                    disabled={added}
                    onClick={() => select(service)}
                    className="flex w-full items-center justify-between gap-4 rounded-control px-4 py-3 text-left enabled:hover:bg-card disabled:opacity-50"
                  >
                    <span>{service.name}</span>
                    <span className="text-sm text-subtle">{added ? 'уже в списке' : service.code}</span>
                  </button>
                </li>
              )
            })}
            {data?.services.length === 0 && <li className="px-4 py-3 text-subtle">Ничего не найдено</li>}
          </ul>
        </div>
      ) : (
        <div className="flex flex-col gap-5">
          <div className="font-medium">{selected.name}</div>
          <div className="flex gap-2">
            {acceptedMarks.map((value) => (
              <button
                key={value}
                type="button"
                onClick={() => setMark(value)}
                className={clsx(
                  'h-11 flex-1 rounded-control border text-sm font-medium',
                  mark === value
                    ? value === 'critical'
                      ? 'border-mark-critical-dark bg-mark-critical'
                      : 'border-mark-minor-dark bg-mark-minor'
                    : 'border-line text-subtle hover:bg-card',
                )}
              >
                {markLabels[value]}
              </button>
            ))}
          </div>
          <label className="flex flex-col gap-2">
            <span className="text-sm text-subtle">Текст для пациента</span>
            <textarea
              value={patientText}
              onChange={(event) => setPatientText(event.target.value)}
              rows={3}
              className="w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
            />
          </label>
          <label className="flex flex-col gap-2">
            <span className="text-sm text-subtle">Комментарий для коллег (необязательно)</span>
            <textarea
              value={comment}
              onChange={(event) => setComment(event.target.value)}
              rows={2}
              className="w-full rounded-control border border-line px-4 py-3 outline-none focus:border-accent"
            />
          </label>
        </div>
      )}
    </Modal>
  )
}
