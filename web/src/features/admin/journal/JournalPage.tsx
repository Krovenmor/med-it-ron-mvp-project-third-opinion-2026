import { BellOff } from 'lucide-react'
import { useState } from 'react'

import type { Notification } from '@/api/models'
import { useNotifications } from '@/api/queries'
import { TwoToneHeading } from '@/shared/ui/Layout'
import { Segmented, type Option } from '@/shared/ui/Segmented'
import { EmptyState, ErrorState, LoadingState } from '@/shared/ui/States'

import type { RecipientFilter } from '../models'

import { JournalTable } from './JournalTable'

const recipientOptions: Option<RecipientFilter>[] = [
  { value: 'all', label: 'Все' },
  { value: 'patient', label: 'Пациентам' },
  { value: 'staff', label: 'Сотрудникам' },
]

function matchesRecipient(notification: Notification, filter: RecipientFilter): boolean {
  if (filter === 'all') {
    return true
  }
  return (notification.recipient === 'patient') === (filter === 'patient')
}

export function JournalPage() {
  const { data, isPending, isError, error, refetch } = useNotifications()
  const [recipient, setRecipient] = useState<RecipientFilter>('all')

  const notifications = data?.notifications ?? []
  const visible = notifications.filter((notification) => matchesRecipient(notification, recipient))

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-2">
        <TwoToneHeading main="Журнал" rest="уведомлений" />
        <p className="text-subtle">В MVP уведомления не отправляются – здесь видно, что получили бы пациенты и сотрудники.</p>
      </div>

      <Segmented label="Кому" options={recipientOptions} value={recipient} onChange={setRecipient} />

      {isPending && <LoadingState />}
      {isError && <ErrorState message={error.message} onRetry={() => void refetch()} />}
      {data && notifications.length === 0 && (
        <EmptyState
          icon={BellOff}
          title="Уведомлений пока нет"
          text="Они появятся после подтверждения врачом, по таймерам протокола и при эскалациях."
        />
      )}
      {data && notifications.length > 0 && visible.length === 0 && <EmptyState title="Нет уведомлений под фильтр" />}
      {visible.length > 0 && <JournalTable notifications={visible} />}
    </div>
  )
}
