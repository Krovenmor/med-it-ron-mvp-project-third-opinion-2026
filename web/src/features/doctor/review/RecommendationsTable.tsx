import clsx from 'clsx'
import { Check, Pencil, Sparkles, Stethoscope } from 'lucide-react'
import { useState } from 'react'

import type { Mark, Recommendation } from '@/api/models'
import { markLabels, rejectReasonLabels } from '@/shared/lib/labels'
import { Pill } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'

const markColumns: Mark[] = ['critical', 'minor', 'rejected']

const selectedStyles: Record<Mark, string> = {
  critical: 'border-mark-critical-dark bg-mark-critical text-ink',
  minor: 'border-mark-minor-dark bg-mark-minor text-ink',
  rejected: 'border-subtle bg-line text-ink',
}

interface RecommendationsTableProps {
  recommendations: Recommendation[]
  editable: boolean
  busy: boolean
  onMark: (recommendation: Recommendation, mark: Mark) => void
  onEditText: (recommendation: Recommendation, text: string, done: () => void) => void
}

export function RecommendationsTable({
  recommendations,
  editable,
  busy,
  onMark,
  onEditText,
}: RecommendationsTableProps) {
  return (
    <table className="w-full border-collapse text-left">
      <thead>
        <tr className="border-b border-line text-sm font-medium text-subtle">
          <th className="w-[30%] py-3 pr-4 font-medium">Услуга</th>
          <th className="py-3 pr-4 font-medium">Обоснование</th>
          {markColumns.map((mark) => (
            <th key={mark} className="w-36 py-3 text-center font-medium">
              {markLabels[mark]}
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {recommendations.map((recommendation) => (
          <tr key={recommendation.id} className="border-b border-line align-top">
            <td className="py-4 pr-4">
              <ServiceCell recommendation={recommendation} />
            </td>
            <td className="py-4 pr-4">
              <div className="flex flex-col gap-3">
                <RationaleCell recommendation={recommendation} />
                <PatientText
                  recommendation={recommendation}
                  editable={editable}
                  busy={busy}
                  onSave={(text, done) => onEditText(recommendation, text, done)}
                />
              </div>
            </td>
            {markColumns.map((mark) => {
              const selected = recommendation.review?.mark === mark
              return (
                <td key={mark} className="px-1 py-4">
                  <button
                    type="button"
                    disabled={!editable || busy}
                    onClick={() => !selected && onMark(recommendation, mark)}
                    aria-pressed={selected}
                    className={clsx(
                      'flex h-11 w-full items-center justify-center gap-1.5 rounded-control border text-sm font-medium transition-colors',
                      selected
                        ? selectedStyles[mark]
                        : 'border-line bg-white text-subtle enabled:hover:bg-card',
                      !editable && !selected && 'opacity-40',
                    )}
                  >
                    {selected && <Check className="size-4" aria-hidden />}
                    {markLabels[mark]}
                  </button>
                </td>
              )
            })}
          </tr>
        ))}
      </tbody>
    </table>
  )
}

interface CellProps {
  recommendation: Recommendation
}

function ServiceCell({ recommendation }: CellProps) {
  const { review } = recommendation
  return (
    <div className="flex flex-col gap-2">
      <span className="font-medium">{recommendation.service_name}</span>
      <div className="flex flex-wrap gap-1.5">
        {recommendation.source === 'ai' ? (
          <Pill className="bg-accent-soft text-accent-deep">
            <Sparkles className="size-3.5" aria-hidden />
            предложено ИИ
          </Pill>
        ) : (
          <Pill className="bg-card text-subtle">
            <Stethoscope className="size-3.5" aria-hidden />
            добавлено врачом
          </Pill>
        )}
        {recommendation.already_booked && <Pill className="bg-normal-bg text-normal-fg">уже записан</Pill>}
        {recommendation.service_code === '' && <Pill className="bg-planned-bg text-planned-fg">нет в клинике</Pill>}
      </div>
      {review?.mark === 'rejected' && review.reject_reason && (
        <span className="text-sm text-subtle">
          Не нужна: {rejectReasonLabels[review.reject_reason]}
          {review.reject_comment && ` – ${review.reject_comment}`}
        </span>
      )}
    </div>
  )
}

function RationaleCell({ recommendation }: CellProps) {
  const [expanded, setExpanded] = useState(false)
  const hasDetails = recommendation.guideline_ref !== '' || recommendation.rationale.length > 120
  return (
    <div className="flex flex-col gap-2 text-sm">
      <span className={clsx(!expanded && 'line-clamp-2')}>{recommendation.rationale || 'Обоснование не указано'}</span>
      {expanded && recommendation.guideline_ref && (
        <span className="text-subtle">Клинические рекомендации: {recommendation.guideline_ref}</span>
      )}
      {hasDetails && (
        <button
          type="button"
          onClick={() => setExpanded((value) => !value)}
          className="self-start font-medium text-accent-deep hover:underline"
        >
          {expanded ? 'Свернуть' : 'Подробнее ›'}
        </button>
      )}
    </div>
  )
}

interface PatientTextProps {
  recommendation: Recommendation
  editable: boolean
  busy: boolean
  onSave: (text: string, done: () => void) => void
}

function PatientText({ recommendation, editable, busy, onSave }: PatientTextProps) {
  const [draft, setDraft] = useState<string | null>(null)

  if (draft !== null) {
    const text = draft.trim()
    return (
      <div className="flex flex-col gap-2 text-sm">
        <span className="text-subtle">Текст для пациента – без диагноза и пугающих формулировок</span>
        <textarea
          autoFocus
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          rows={3}
          className="w-full rounded-control border border-line bg-white px-3 py-2 outline-none focus:border-accent"
        />
        <div className="flex gap-2">
          <Button
            size="sm"
            disabled={busy || text === '' || text === recommendation.patient_text}
            onClick={() => onSave(text, () => setDraft(null))}
          >
            Сохранить
          </Button>
          <Button size="sm" variant="secondary" onClick={() => setDraft(null)}>
            Отмена
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="rounded-control bg-card-nested px-3 py-2 text-sm">
      <div className="mb-1 flex items-center justify-between gap-3 text-subtle">
        <span>Пациент увидит</span>
        {editable && (
          <button
            type="button"
            onClick={() => setDraft(recommendation.patient_text)}
            className="flex items-center gap-1 font-medium text-accent-deep hover:underline"
          >
            <Pencil className="size-3.5" aria-hidden />
            Изменить
          </button>
        )}
      </div>
      <p>{recommendation.patient_text}</p>
    </div>
  )
}
