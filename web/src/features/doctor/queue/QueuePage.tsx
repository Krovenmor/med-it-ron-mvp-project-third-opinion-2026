import { BellRing } from 'lucide-react'
import { useState } from 'react'

import { useReviewQueue } from '@/api/queries'
import { canAskNotificationPermission } from '@/shared/lib/alerts'
import { useServerNow } from '@/shared/lib/time'
import { Button } from '@/shared/ui/Button'
import { Metric, TwoToneHeading } from '@/shared/ui/Layout'
import { EmptyState, ErrorState, LoadingState } from '@/shared/ui/States'

import { matchesFilters, queueStatusOf, type StatusFilter, type UrgencyFilter } from './models'
import { QueueFilters } from './QueueFilters'
import { QueueTable } from './QueueTable'
import { useEmergencyAlerts } from './useEmergencyAlerts'

export function QueuePage() {
  const { data, isPending, isError, error, refetch } = useReviewQueue()
  const now = useServerNow()
  const [urgency, setUrgency] = useState<UrgencyFilter>('all')
  const [status, setStatus] = useState<StatusFilter>('all')
  const [askNotifications, setAskNotifications] = useState(canAskNotificationPermission)

  useEmergencyAlerts(data?.cases)

  const cases = data?.cases ?? []
  const visible = cases.filter((item) => matchesFilters(item, urgency, status))

  const enableNotifications = async () => {
    await Notification.requestPermission()
    setAskNotifications(canAskNotificationPermission())
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-end justify-between gap-4">
        <TwoToneHeading main="Очередь" rest="на проверку" />
        {askNotifications && (
          <Button variant="secondary" size="sm" onClick={enableNotifications}>
            <BellRing className="size-4" />
            Включить уведомления о неотложных
          </Button>
        )}
      </div>

      <div className="grid grid-cols-3 gap-5">
        <Metric value={cases.length} label="ожидают проверки" />
        <Metric value={cases.filter((item) => item.urgency === 'emergency').length} label="неотложных" />
        <Metric value={cases.filter((item) => queueStatusOf(item) === 'new').length} label="ещё не открыты" />
      </div>

      <QueueFilters urgency={urgency} status={status} onUrgencyChange={setUrgency} onStatusChange={setStatus} />

      {isPending && <LoadingState />}
      {isError && <ErrorState message={error.message} onRetry={() => void refetch()} />}
      {data && cases.length === 0 && (
        <EmptyState title="Очередь пуста" text="Новые заключения появятся здесь после оценки ИИ." />
      )}
      {data && cases.length > 0 && visible.length === 0 && <EmptyState title="Нет кейсов под выбранные фильтры" />}
      {visible.length > 0 && <QueueTable cases={visible} now={now} />}
    </div>
  )
}
