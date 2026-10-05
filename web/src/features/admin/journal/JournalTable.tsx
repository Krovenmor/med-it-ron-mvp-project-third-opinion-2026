import clsx from 'clsx'
import { ChevronDown } from 'lucide-react'
import { useState } from 'react'

import type { Notification } from '@/api/models'
import { formatDateTime } from '@/shared/lib/format'
import { UrgencyBadge } from '@/shared/ui/Badge'

import { channelLabels, escalationKinds, notificationKindLabels, recipientLabels } from '../models'

interface JournalTableProps {
  notifications: Notification[]
}

export function JournalTable({ notifications }: JournalTableProps) {
  const [expanded, setExpanded] = useState<ReadonlySet<string>>(new Set())

  const toggle = (id: string) =>
    setExpanded((current) => {
      const next = new Set(current)
      if (!next.delete(id)) {
        next.add(id)
      }
      return next
    })

  return (
    <table className="w-full table-fixed border-collapse text-left">
      <thead>
        <tr className="border-b border-line text-sm font-medium text-subtle">
          <th className="w-[14%] py-3 pr-4 pl-2 font-medium">Время</th>
          <th className="w-[18%] py-3 pr-4 font-medium">Пациент</th>
          <th className="w-[15%] py-3 pr-4 font-medium">Уровень</th>
          <th className="w-[17%] py-3 pr-4 font-medium">Кому и канал</th>
          <th className="py-3 pr-2 font-medium">Сообщение</th>
        </tr>
      </thead>
      <tbody>
        {notifications.map((notification) => {
          const open = expanded.has(notification.id)
          return (
            <tr
              key={notification.id}
              onClick={() => toggle(notification.id)}
              className="cursor-pointer border-b border-line align-top transition-colors hover:bg-card-nested"
              aria-expanded={open}
            >
              <td className="py-4 pr-4 pl-2">{formatDateTime(notification.created_at)}</td>
              <td className="py-4 pr-4">
                <div className="font-medium">{notification.patient.full_name}</div>
                <div className="text-sm text-subtle">{notification.patient.id}</div>
              </td>
              <td className="py-4 pr-4">{notification.urgency && <UrgencyBadge urgency={notification.urgency} />}</td>
              <td className="py-4 pr-4">
                <div>{recipientLabels[notification.recipient]}</div>
                <div className="text-sm text-subtle">{channelLabels[notification.channel]}</div>
              </td>
              <td className="py-4 pr-2">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div
                      className={clsx(
                        'font-medium',
                        escalationKinds.includes(notification.kind) && 'text-overdue',
                      )}
                    >
                      {notificationKindLabels[notification.kind]}
                    </div>
                    <div className={clsx('text-sm text-subtle', !open && 'truncate')}>{notification.text}</div>
                  </div>
                  <ChevronDown
                    className={clsx('mt-0.5 size-4 shrink-0 text-muted transition-transform', open && 'rotate-180')}
                    aria-hidden
                  />
                </div>
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}
