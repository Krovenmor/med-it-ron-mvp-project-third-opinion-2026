import type { ReactNode } from 'react'

import type { OperatorCard } from '@/api/models'
import { formatDate, formatSlot } from '@/shared/lib/format'
import { Card } from '@/shared/ui/Layout'

import { clinicName } from '../models'

import { scriptOffer, type Selection } from './models'

function Value({ children }: { children: ReactNode }) {
  return <mark className="rounded-md bg-accent-soft px-1.5 py-0.5 text-ink">{children}</mark>
}

interface CallScriptProps {
  card: OperatorCard
  selection: Selection | null
}

export function CallScript({ card, selection }: CallScriptProps) {
  const emergency = card.case.urgency === 'emergency'
  const proposal = scriptOffer(card, selection)
  const slot = proposal?.slot ? <Value>{formatSlot(proposal.slot.starts_at)}</Value> : <Value>ближайшее время</Value>

  return (
    <Card className="flex flex-col gap-3">
      <div className="flex items-baseline justify-between gap-4">
        <h2 className="text-lg font-medium">Скрипт звонка</h2>
        <span className="text-sm text-subtle">без диагноза и находок</span>
      </div>
      <p className="text-lg leading-relaxed">
        Здравствуйте, это клиника <Value>{clinicName}</Value>.{' '}
        {card.task.reason === 'help_request' && <>Вы просили перезвонить вам по плану обследования. </>}
        {emergency ? (
          <>По результатам исследования врач просит вас срочно связаться с ним или приехать на консультацию.</>
        ) : proposal ? (
          <>
            Врач рекомендует вам <Value>{proposal.offer.service_name}</Value> до{' '}
            <Value>{formatDate(card.case.book_by)}</Value>. Могу записать вас сейчас, подойдёт {slot}?
          </>
        ) : (
          <>Врач подготовил для вас рекомендации. Удобно сейчас обсудить запись?</>
        )}
      </p>
      {emergency && proposal && (
        <p className="text-lg leading-relaxed">
          Могу сразу записать вас: <Value>{proposal.offer.service_name}</Value>, подойдёт {slot}?
        </p>
      )}
    </Card>
  )
}
