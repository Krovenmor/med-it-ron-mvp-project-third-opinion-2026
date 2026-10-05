import { ChevronLeft } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { toast } from 'sonner'

import { ApiError } from '@/api/http'
import { useOperatorCard } from '@/api/queries'
import { useServerNow } from '@/shared/lib/time'
import { ErrorState, LoadingState } from '@/shared/ui/States'

import { AttemptsHistory } from './AttemptsHistory'
import { CallActions } from './CallActions'
import { CallbackModal } from './CallbackModal'
import { CallBanners, CallHeader } from './CallHeader'
import { CallScript } from './CallScript'
import { DeclineModal } from './DeclineModal'
import type { Selection } from './models'
import { OffersBlock } from './OffersBlock'
import { useCallActions } from './useCallActions'

export function CallPage() {
  const { taskId = '' } = useParams()
  const navigate = useNavigate()
  const now = useServerNow()
  const { data: card, isPending, isError, error, refetch } = useOperatorCard(taskId)
  const actions = useCallActions(taskId)
  const [comment, setComment] = useState('')
  const [selection, setSelection] = useState<Selection | null>(null)
  const [declining, setDeclining] = useState(false)
  const [schedulingCallback, setSchedulingCallback] = useState(false)

  if (isPending) {
    return <LoadingState />
  }
  if (isError) {
    const notFound = error instanceof ApiError && error.status === 404
    return (
      <ErrorState
        title={notFound ? 'Задача не найдена' : 'Не удалось загрузить задачу'}
        message={notFound ? undefined : error.message}
        onRetry={notFound ? undefined : () => void refetch()}
      />
    )
  }

  const busy = Object.values(actions).some((mutation) => mutation.isPending)
  const body = { comment: comment.trim() }
  const finish = (message: string) => {
    toast.success(message)
    navigate('/admin')
  }
  const stay = (message: string) => {
    toast.success(message)
    setComment('')
  }

  return (
    <div className="flex flex-col gap-5">
      <Link to="/admin" className="flex items-center gap-1 self-start font-medium text-accent-deep hover:underline">
        <ChevronLeft className="size-4" />
        Позвонить
      </Link>

      <CallHeader card={card} now={now} />
      <CallBanners card={card} />

      <div className="grid grid-cols-[minmax(0,1fr)_380px] items-start gap-5">
        <div className="flex flex-col gap-5">
          <CallScript card={card} selection={selection} />
          <OffersBlock
            offers={card.offers}
            emergency={card.case.urgency === 'emergency'}
            selection={selection}
            disabled={card.task.closed_at !== null}
            booking={actions.book.isPending}
            onSelect={setSelection}
            onBook={(chosen) =>
              actions.book.mutate(
                { recommendation_id: chosen.recommendationId, slot_id: chosen.slotId, ...body },
                { onSuccess: () => finish('Пациент записан, задача закрыта') },
              )
            }
          />
        </div>
        <div className="flex flex-col gap-5">
          <CallActions
            task={card.task}
            comment={comment}
            busy={busy}
            onCommentChange={setComment}
            onTake={() => actions.take.mutate()}
            onNoAnswer={() =>
              actions.noAnswer.mutate(body, {
                onSuccess: (task) =>
                  stay(
                    task.needs_doctor
                      ? 'Третья попытка без ответа: пациенту отправлено SMS, передайте врачу'
                      : 'Попытка записана, повторный звонок запланирован',
                  ),
              })
            }
            onCallback={() => setSchedulingCallback(true)}
            onContacted={() => actions.contacted.mutate(body, { onSuccess: () => finish('Дозвонились, задача закрыта') })}
            onDecline={() => setDeclining(true)}
            onHandToDoctor={() =>
              actions.handToDoctor.mutate(body, { onSuccess: () => finish('Пациент передан дежурному врачу') })
            }
          />
          <AttemptsHistory attempts={card.attempts} />
        </div>
      </div>

      {declining && (
        <DeclineModal
          initialComment={comment}
          pending={actions.decline.isPending}
          onClose={() => setDeclining(false)}
          onSubmit={(request) => actions.decline.mutate(request, { onSuccess: () => finish('Отказ записан') })}
        />
      )}
      {schedulingCallback && (
        <CallbackModal
          now={now}
          comment={body.comment}
          pending={actions.callback.isPending}
          onClose={() => setSchedulingCallback(false)}
          onSubmit={(request) =>
            actions.callback.mutate(request, {
              onSuccess: () => {
                setSchedulingCallback(false)
                stay('Перезвон запланирован')
              },
            })
          }
        />
      )}
    </div>
  )
}
