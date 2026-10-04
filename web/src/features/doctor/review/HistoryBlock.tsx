import { useCaseHistory } from '@/api/queries'
import { formatDate, formatDateTime } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { Collapsible } from '@/shared/ui/Layout'

interface HistoryBlockProps {
  caseId: string
}

export function HistoryBlock({ caseId }: HistoryBlockProps) {
  const { data, isPending, isError, refetch } = useCaseHistory(caseId)

  return (
    <Collapsible
      title="История приёмов и обследований"
      aside={isError && <span className="text-overdue">история не загружена</span>}
    >
      {isPending && <div className="text-subtle">Загружаем историю из МИС</div>}
      {isError && (
        <div className="flex items-center gap-4">
          <span className="text-subtle">МИС не ответила, история не загружена.</span>
          <Button variant="secondary" size="sm" onClick={() => void refetch()}>
            Повторить
          </Button>
        </div>
      )}
      {data && data.visits.length === 0 && data.appointments.length === 0 && (
        <div className="text-subtle">В МИС нет посещений и записей пациента.</div>
      )}
      {data && (data.visits.length > 0 || data.appointments.length > 0) && (
        <div className="grid grid-cols-2 gap-5">
          <div className="rounded-card bg-card-nested p-5">
            <div className="mb-3 text-sm font-medium text-subtle">Посещения</div>
            {data.visits.length === 0 && <div className="text-subtle">Нет</div>}
            <ul className="flex flex-col gap-2">
              {data.visits.map((visit) => (
                <li key={`${visit.service_code}-${visit.visited_at}`} className="flex justify-between gap-4">
                  <span>{visit.service_name}</span>
                  <span className="text-subtle">{formatDate(visit.visited_at)}</span>
                </li>
              ))}
            </ul>
          </div>
          <div className="rounded-card bg-card-nested p-5">
            <div className="mb-3 text-sm font-medium text-subtle">Записи</div>
            {data.appointments.length === 0 && <div className="text-subtle">Нет</div>}
            <ul className="flex flex-col gap-2">
              {data.appointments.map((appointment) => (
                <li
                  key={`${appointment.service_code}-${appointment.scheduled_at}`}
                  className="flex justify-between gap-4"
                >
                  <span>{appointment.service_name}</span>
                  <span className="text-subtle">{formatDateTime(appointment.scheduled_at)}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      )}
    </Collapsible>
  )
}
