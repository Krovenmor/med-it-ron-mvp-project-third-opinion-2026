import clsx from 'clsx'
import { Clock, PhoneCall, PhoneMissed, Stethoscope, UserCheck, X } from 'lucide-react'

import type { OperatorTask } from '@/api/models'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Layout'

interface CallActionsProps {
  task: OperatorTask
  comment: string
  busy: boolean
  onCommentChange: (value: string) => void
  onTake: () => void
  onNoAnswer: () => void
  onCallback: () => void
  onContacted: () => void
  onDecline: () => void
  onHandToDoctor: () => void
}

export function CallActions(props: CallActionsProps) {
  const { task, comment, busy } = props
  const closed = task.closed_at !== null
  const disabled = closed || busy

  return (
    <Card className="flex flex-col gap-4">
      <h2 className="text-lg font-medium">Результат звонка</h2>
      {!task.assignee && !closed && (
        <Button onClick={props.onTake} disabled={busy}>
          <UserCheck className="size-4" />
          Взять в работу
        </Button>
      )}
      <textarea
        value={comment}
        onChange={(event) => props.onCommentChange(event.target.value)}
        disabled={closed}
        placeholder="Комментарий к звонку"
        rows={3}
        className="w-full rounded-control border border-line bg-white px-4 py-3 outline-none focus:border-accent disabled:opacity-50"
      />
      <div className="flex flex-col gap-2 [&>button]:justify-start">
        <Button variant="secondary" onClick={props.onNoAnswer} disabled={disabled}>
          <PhoneMissed className="size-4" />
          Не дозвонился
        </Button>
        <Button variant="secondary" onClick={props.onCallback} disabled={disabled}>
          <Clock className="size-4" />
          Перезвонить позже
        </Button>
        <Button variant="secondary" onClick={props.onContacted} disabled={disabled}>
          <PhoneCall className="size-4" />
          Дозвонился, закрыть
        </Button>
        <Button variant="secondary" onClick={props.onDecline} disabled={disabled}>
          <X className="size-4" />
          Отказ
        </Button>
        <Button
          variant="secondary"
          onClick={props.onHandToDoctor}
          disabled={disabled}
          className={clsx(task.needs_doctor && 'border-overdue text-overdue')}
        >
          <Stethoscope className="size-4" />
          Передать врачу
        </Button>
      </div>
      <p className="text-sm text-subtle">Чтобы записать пациента, выберите слот в блоке «Что предложить».</p>
    </Card>
  )
}
