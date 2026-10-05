import { Phone, Stethoscope, TriangleAlert } from 'lucide-react'

import type { OperatorCard } from '@/api/models'
import { formatDate } from '@/shared/lib/format'
import { Pill, UrgencyBadge } from '@/shared/ui/Badge'
import { Card, TwoToneHeading } from '@/shared/ui/Layout'

import { Deadline } from '../tasks/Deadline'
import { adminId, taskReasonLabels, taskStatusLabels, taskStatusStyles } from '../models'

interface CallHeaderProps {
  card: OperatorCard
  now: Date
}

export function CallHeader({ card, now }: CallHeaderProps) {
  const { task, patient } = card
  const emergency = card.case.urgency === 'emergency'
  return (
    <Card className="flex flex-wrap items-start justify-between gap-6">
      <div className="flex flex-col gap-2">
        <TwoToneHeading main={patient.full_name} rest={taskReasonLabels[task.reason].toLowerCase()} />
        <a href={`tel:${patient.phone}`} className="inline-flex items-center gap-2 text-lg font-medium text-accent-deep">
          <Phone className="size-5" aria-hidden />
          {patient.phone}
        </a>
        <div className="text-subtle">
          Записать не позднее {emergency ? 'сегодня' : formatDate(card.case.book_by)} · попыток: {task.attempts}
          {task.assignee && ` · ${task.assignee === adminId ? 'моя задача' : `у ${task.assignee}`}`}
        </div>
      </div>
      <div className="flex flex-col items-end gap-3">
        <div className="flex items-center gap-2">
          <UrgencyBadge urgency={card.case.urgency} />
          <Pill className={taskStatusStyles[task.status]}>{taskStatusLabels[task.status]}</Pill>
        </div>
        {!task.closed_at && <Deadline dueAt={task.due_at} now={now} />}
      </div>
    </Card>
  )
}

interface CallBannersProps {
  card: OperatorCard
}

export function CallBanners({ card }: CallBannersProps) {
  const { task } = card
  if (task.closed_at) {
    return (
      <div className="rounded-card bg-accent-tint px-6 py-4 font-medium text-accent-deep">
        Задача закрыта: {taskStatusLabels[task.status].toLowerCase()}
      </div>
    )
  }
  return (
    <>
      {card.case.urgency === 'emergency' && (
        <div className="flex items-center gap-3 rounded-card bg-emergency-bg px-6 py-4 font-medium text-emergency-fg">
          <TriangleAlert className="size-5" aria-hidden />
          Неотложный случай: позвоните пациенту сейчас
        </div>
      )}
      {task.needs_doctor && (
        <div className="flex items-center gap-3 rounded-card bg-white px-6 py-4 font-medium text-overdue ring-1 ring-inset ring-overdue">
          <Stethoscope className="size-5" aria-hidden />
          Три попытки без ответа – пациенту отправлено SMS. Передайте пациента врачу.
        </div>
      )}
    </>
  )
}
