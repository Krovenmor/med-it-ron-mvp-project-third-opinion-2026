import clsx from 'clsx'
import { useNavigate } from 'react-router'

import type { QueueCase } from '@/api/models'
import { ageAt, formatAge, formatDate, formatDuration } from '@/shared/lib/format'
import { modalityLabels, sexLabels } from '@/shared/lib/labels'
import { Pill, UrgencyBadge } from '@/shared/ui/Badge'

import { SlaTimer } from '../SlaTimer'

import { queueStatusLabels, queueStatusOf } from './models'

interface QueueTableProps {
  cases: QueueCase[]
  now: Date
}

export function QueueTable({ cases, now }: QueueTableProps) {
  const navigate = useNavigate()

  return (
    <table className="w-full table-fixed border-collapse text-left">
      <thead>
        <tr className="border-b border-line text-sm font-medium text-subtle">
          <th className="w-[18%] py-3 pr-4 pl-2 font-medium">Срочность</th>
          <th className="w-[20%] py-3 pr-4 font-medium">Пациент</th>
          <th className="w-[20%] py-3 pr-4 font-medium">Исследование</th>
          <th className="w-[22%] py-3 pr-4 font-medium">Ожидает</th>
          <th className="py-3 font-medium">Статус</th>
        </tr>
      </thead>
      <tbody>
        {cases.map((item) => {
          const emergency = item.urgency === 'emergency'
          const status = queueStatusOf(item)
          return (
            <tr
              key={item.case_id}
              onClick={() => navigate(`/doctor/cases/${item.case_id}`)}
              className={clsx(
                'cursor-pointer border-b border-line transition-colors',
                emergency ? 'bg-emergency-bg/60 hover:bg-emergency-bg' : 'hover:bg-card-nested',
              )}
            >
              <td className="py-4 pr-4 pl-2">
                <UrgencyBadge urgency={item.urgency} />
              </td>
              <td className="py-4 pr-4">
                <div className="font-medium">{item.patient.id}</div>
                <div className="text-sm text-subtle">
                  {formatAge(ageAt(item.patient.birth_date, now))}, {sexLabels[item.patient.sex]}
                </div>
              </td>
              <td className="py-4 pr-4">
                <div>{modalityLabels[item.modality]}</div>
                <div className="text-sm text-subtle">{formatDate(item.performed_at)}</div>
              </td>
              <td className="py-4 pr-4">
                <div>{formatDuration(now.getTime() - Date.parse(item.received_at))}</div>
                {emergency && <SlaTimer receivedAt={item.received_at} now={now} />}
              </td>
              <td className="py-4">
                <Pill className={status === 'new' ? 'bg-accent-soft text-accent-deep' : 'bg-card text-subtle'}>
                  {queueStatusLabels[status]}
                </Pill>
                <div className="mt-1 text-sm text-subtle">
                  Отмечено {item.recommendations_reviewed} из {item.recommendations_total}
                </div>
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}
