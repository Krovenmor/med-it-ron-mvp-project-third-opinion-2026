import { useQueryClient } from '@tanstack/react-query'
import { ChevronLeft, Plus } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { toast } from 'sonner'

import { b2b } from '@/api/b2b'
import { ApiError } from '@/api/http'
import type { Mark, Recommendation } from '@/api/models'
import { queryKeys, useCaseDetails } from '@/api/queries'
import { markLabels } from '@/shared/lib/labels'
import { useServerNow } from '@/shared/lib/time'
import { Button } from '@/shared/ui/Button'
import { Collapsible } from '@/shared/ui/Layout'
import { ErrorState, LoadingState } from '@/shared/ui/States'

import { AddServiceModal } from './AddServiceModal'
import { CaseHeader, EmergencyBanner } from './CaseHeader'
import { ConfirmModal } from './ConfirmModal'
import { HistoryBlock } from './HistoryBlock'
import { countMarks } from './models'
import { RecommendationsTable } from './RecommendationsTable'
import { RejectModal } from './RejectModal'
import { UrgencyModal } from './UrgencyModal'
import { useReviewActions } from './useReviewActions'

export function ReviewPage() {
  const { caseId = '' } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const now = useServerNow()
  const { data: details, isPending, isError, error, refetch } = useCaseDetails(caseId)
  const actions = useReviewActions(caseId)
  const opened = useRef(false)

  const [rejecting, setRejecting] = useState<Recommendation | null>(null)
  const [adding, setAdding] = useState(false)
  const [changingUrgency, setChangingUrgency] = useState(false)
  const [confirming, setConfirming] = useState(false)

  useEffect(() => {
    if (details?.status !== 'in_review' || opened.current) {
      return
    }
    opened.current = true
    b2b
      .openCase(caseId)
      .then(() => queryClient.invalidateQueries({ queryKey: queryKeys.reviewQueue }))
      .catch((openError: Error) => toast.error(openError.message))
  }, [details?.status, caseId, queryClient])

  if (isPending) {
    return <LoadingState />
  }
  if (isError) {
    const notFound = error instanceof ApiError && error.status === 404
    return (
      <ErrorState
        title={notFound ? 'Кейс не найден' : 'Не удалось загрузить кейс'}
        message={notFound ? undefined : error.message}
        onRetry={notFound ? undefined : () => void refetch()}
      />
    )
  }

  const editable = details.status === 'in_review'
  const counts = countMarks(details.recommendations)
  const busy = actions.review.isPending || actions.add.isPending

  const mark = (recommendation: Recommendation, value: Mark) => {
    if (value === 'rejected') {
      setRejecting(recommendation)
      return
    }
    actions.review.mutate({ recommendationId: recommendation.id, body: { mark: value } })
  }

  const editText = (recommendation: Recommendation, text: string, done: () => void) =>
    actions.review.mutate({ recommendationId: recommendation.id, body: { patient_text: text } }, { onSuccess: done })

  const confirm = () =>
    actions.confirm.mutate(undefined, {
      onSuccess: () => {
        toast.success('Рекомендации подтверждены, пациент будет уведомлён')
        navigate('/doctor')
      },
    })

  return (
    <div className="flex flex-col gap-5">
      <Link to="/doctor" className="flex items-center gap-1 self-start font-medium text-accent-deep hover:underline">
        <ChevronLeft className="size-4" />
        Очередь на проверку
      </Link>

      <CaseHeader
        details={details}
        now={now}
        editable={editable}
        onChangeUrgency={() => setChangingUrgency(true)}
      />

      {editable && details.urgency === 'emergency' && <EmergencyBanner receivedAt={details.received_at} now={now} />}
      {!editable && (
        <div className="rounded-card bg-accent-tint px-6 py-4 font-medium text-accent-deep">
          Рекомендации подтверждены – кейс доступен только для чтения
        </div>
      )}

      <Collapsible title="Заключение" defaultOpen>
        <p className="leading-relaxed">{details.conclusion}</p>
      </Collapsible>

      <HistoryBlock caseId={caseId} />

      <section className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-medium">Рекомендации</h2>
          {editable && (
            <Button variant="secondary" size="sm" onClick={() => setAdding(true)}>
              <Plus className="size-4" />
              Добавить исследование
            </Button>
          )}
        </div>
        <RecommendationsTable
          recommendations={details.recommendations}
          editable={editable}
          busy={busy}
          onMark={mark}
          onEditText={editText}
        />
      </section>

      {editable && (
        <div className="sticky bottom-0 -mx-8 border-t border-line bg-white">
          <div className="flex items-center justify-between gap-6 px-8 py-4">
            <div className="flex flex-col">
              <span className="font-medium">
                {markLabels.critical} {counts.critical} · {markLabels.minor} {counts.minor} · {markLabels.rejected}{' '}
                {counts.rejected}
              </span>
              {counts.pending > 0 && (
                <span className="text-sm text-subtle">Не отмечено: {counts.pending} – отметьте каждую строку</span>
              )}
            </div>
            <Button onClick={() => setConfirming(true)} disabled={counts.pending > 0 || busy}>
              Подтвердить рекомендации
            </Button>
          </div>
        </div>
      )}

      {rejecting && (
        <RejectModal
          recommendation={rejecting}
          pending={actions.review.isPending}
          onClose={() => setRejecting(null)}
          onSubmit={(body) =>
            actions.review.mutate({ recommendationId: rejecting.id, body }, { onSuccess: () => setRejecting(null) })
          }
        />
      )}
      {adding && (
        <AddServiceModal
          existingCodes={details.recommendations.map((rec) => rec.service_code).filter(Boolean)}
          pending={actions.add.isPending}
          onClose={() => setAdding(false)}
          onSubmit={(body) => actions.add.mutate(body, { onSuccess: () => setAdding(false) })}
        />
      )}
      {changingUrgency && details.urgency && (
        <UrgencyModal
          current={details.urgency}
          pending={actions.changeUrgency.isPending}
          onClose={() => setChangingUrgency(false)}
          onSubmit={(body) => actions.changeUrgency.mutate(body, { onSuccess: () => setChangingUrgency(false) })}
        />
      )}
      {confirming && (
        <ConfirmModal
          noFurtherExamination={counts.critical + counts.minor === 0}
          pending={actions.confirm.isPending}
          onClose={() => setConfirming(false)}
          onConfirm={confirm}
        />
      )}
    </div>
  )
}
