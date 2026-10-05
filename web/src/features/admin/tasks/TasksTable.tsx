import clsx from 'clsx'
import { ChevronRight, Phone } from 'lucide-react'
import { useNavigate } from 'react-router'

import type { OperatorTaskItem } from '@/api/models'
import { formatDate, formatDateTime } from '@/shared/lib/format'
import { Pill, UrgencyBadge } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'

import { adminId, taskReasonLabels, taskStatusLabels, taskStatusStyles } from '../models'

import { Deadline } from './Deadline'

interface TasksTableProps {
  items: OperatorTaskItem[]
  now: Date
  taking: string | null
  onTake: (taskId: string) => void
}

export function TasksTable({ items, now, taking, onTake }: TasksTableProps) {
  const navigate = useNavigate()

  return (
    <table className="w-full table-fixed border-collapse text-left">
      <thead>
        <tr className="border-b border-line text-sm font-medium text-subtle">
          <th className="w-[14%] py-3 pr-4 pl-2 font-medium">Уровень</th>
          <th className="w-[15%] py-3 pr-4 font-medium">Пациент</th>
          <th className="w-[17%] py-3 pr-4 font-medium">Что предложить</th>
          <th className="w-[11%] py-3 pr-4 font-medium">Причина звонка</th>
          <th className="w-[12%] py-3 pr-4 font-medium">Дедлайн</th>
          <th className="w-[6%] py-3 pr-4 font-medium">Попытки</th>
          <th className="py-3 pr-4 font-medium">Статус</th>
          <th className="w-40 py-3 pr-2" />
        </tr>
      </thead>
      <tbody>
        {items.map(({ task, case: caseInfo, patient, offer_service }) => {
          const emergency = caseInfo.urgency === 'emergency'
          return (
            <tr
              key={task.id}
              onClick={() => navigate(`/admin/tasks/${task.id}`)}
              className={clsx(
                'cursor-pointer border-b border-line align-top transition-colors',
                emergency ? 'bg-emergency-bg/60 hover:bg-emergency-bg' : 'hover:bg-card-nested',
              )}
            >
              <td className="py-4 pr-4 pl-2">
                <UrgencyBadge urgency={caseInfo.urgency} />
                {emergency && <div className="mt-1.5 text-sm font-medium text-emergency-fg">позвонить сейчас</div>}
              </td>
              <td className="py-4 pr-4">
                <div className="font-medium">{patient.full_name}</div>
                <a
                  href={`tel:${patient.phone}`}
                  onClick={(event) => event.stopPropagation()}
                  className="mt-0.5 inline-flex items-center gap-1.5 text-sm text-accent-deep hover:underline"
                >
                  <Phone className="size-3.5" aria-hidden />
                  {patient.phone}
                </a>
              </td>
              <td className="py-4 pr-4">
                <div>{offer_service || (emergency ? 'Срочная консультация врача' : 'Ждёт подтверждения врача')}</div>
                <div className="text-sm text-subtle">
                  не позднее {emergency ? 'сегодня' : formatDate(caseInfo.book_by)}
                </div>
              </td>
              <td className="py-4 pr-4">{taskReasonLabels[task.reason]}</td>
              <td className="py-4 pr-4">
                <Deadline dueAt={task.due_at} now={now} />
              </td>
              <td className="py-4 pr-4">{task.attempts}</td>
              <td className="py-4 pr-4">
                <div className="flex flex-wrap gap-1.5">
                  <Pill className={taskStatusStyles[task.status]}>{taskStatusLabels[task.status]}</Pill>
                  {task.needs_doctor && (
                    <Pill className="bg-white text-overdue ring-1 ring-inset ring-overdue">Передать врачу</Pill>
                  )}
                </div>
                {task.next_call_at && (
                  <div className="mt-1 text-sm text-subtle">перезвонить {formatDateTime(task.next_call_at)}</div>
                )}
                {task.assignee && (
                  <div className="mt-1 text-sm text-subtle">
                    {task.assignee === adminId ? 'моя задача' : `у ${task.assignee}`}
                  </div>
                )}
              </td>
              <td className="py-4 pr-2 text-right">
                {task.assignee ? (
                  <ChevronRight className="ml-auto size-5 text-muted" aria-hidden />
                ) : (
                  <Button
                    size="sm"
                    className="whitespace-nowrap"
                    disabled={taking === task.id}
                    onClick={(event) => {
                      event.stopPropagation()
                      onTake(task.id)
                    }}
                  >
                    Взять в работу
                  </Button>
                )}
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}
