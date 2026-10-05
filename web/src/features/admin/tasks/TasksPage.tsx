import { useMutation, useQueryClient } from '@tanstack/react-query'
import { PhoneOff } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'

import { b2b } from '@/api/b2b'
import { queryKeys, useOperatorTasks } from '@/api/queries'
import { useServerNow } from '@/shared/lib/time'
import { Metric, TwoToneHeading } from '@/shared/ui/Layout'
import { EmptyState, ErrorState, LoadingState } from '@/shared/ui/States'

import { isOverdue, matchesFilters, type OwnerFilter, type StatusFilter, type UrgencyFilter } from '../models'

import { TaskFilters } from './TaskFilters'
import { TasksTable } from './TasksTable'

export function TasksPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const now = useServerNow()
  const [owner, setOwner] = useState<OwnerFilter>('all')
  const [urgency, setUrgency] = useState<UrgencyFilter>('all')
  const [status, setStatus] = useState<StatusFilter>('all')
  const { data, isPending, isError, error, refetch } = useOperatorTasks(owner === 'mine')

  const take = useMutation({
    mutationFn: b2b.takeTask,
    onSuccess: (task) => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.operator })
      navigate(`/admin/tasks/${task.id}`)
    },
    onError: (takeError) => toast.error(takeError.message),
  })

  const items = data?.tasks ?? []
  const visible = items.filter((item) => matchesFilters(item, urgency, status))

  return (
    <div className="flex flex-col gap-6">
      <TwoToneHeading main="Позвонить" rest="пациентам" />

      <div className="grid grid-cols-3 gap-5">
        <Metric value={items.length} label="задач к звонку" />
        <Metric value={items.filter((item) => item.case.urgency === 'emergency').length} label="неотложных" />
        <Metric value={items.filter((item) => isOverdue(item.task.due_at, now)).length} label="просрочено" />
      </div>

      <TaskFilters
        owner={owner}
        urgency={urgency}
        status={status}
        onOwnerChange={setOwner}
        onUrgencyChange={setUrgency}
        onStatusChange={setStatus}
      />

      {isPending && <LoadingState />}
      {isError && <ErrorState message={error.message} onRetry={() => void refetch()} />}
      {data && items.length === 0 && (
        <EmptyState
          icon={PhoneOff}
          title="Сейчас звонить некому"
          text="Задачи появятся, когда пациент не запишется вовремя или ИИ отметит неотложный случай."
        />
      )}
      {data && items.length > 0 && visible.length === 0 && <EmptyState title="Нет задач под выбранные фильтры" />}
      {visible.length > 0 && (
        <TasksTable
          items={visible}
          now={now}
          taking={take.isPending ? (take.variables ?? null) : null}
          onTake={(taskId) => take.mutate(taskId)}
        />
      )}
    </div>
  )
}
