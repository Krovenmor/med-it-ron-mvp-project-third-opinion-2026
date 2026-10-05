import clsx from 'clsx'
import { CalendarCheck } from 'lucide-react'

import type { Offer } from '@/api/models'
import { formatSlot } from '@/shared/lib/format'
import { markLabels } from '@/shared/lib/labels'
import { Pill } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Layout'

import type { Selection } from './models'

const markStyles = {
  critical: 'bg-mark-critical text-ink',
  minor: 'bg-mark-minor text-ink',
}

interface OffersBlockProps {
  offers: Offer[]
  emergency: boolean
  selection: Selection | null
  disabled: boolean
  booking: boolean
  onSelect: (selection: Selection) => void
  onBook: (selection: Selection) => void
}

export function OffersBlock({ offers, emergency, selection, disabled, booking, onSelect, onBook }: OffersBlockProps) {
  return (
    <Card className="flex flex-col gap-4">
      <h2 className="text-lg font-medium">Что предложить</h2>
      {offers.length === 0 && (
        <p className="text-subtle">
          {emergency
            ? 'Врач ещё не подтвердил рекомендации. Пригласите пациента на срочную консультацию или передайте врачу.'
            : 'Рекомендаций для записи нет.'}
        </p>
      )}
      {offers.map((offer) => {
        const selectedSlot = selection?.recommendationId === offer.recommendation_id ? selection.slotId : null
        const slot = offer.slots.find((item) => item.id === selectedSlot)
        return (
          <div key={offer.recommendation_id} className="flex flex-col gap-3 rounded-card bg-white p-5">
            <div className="flex flex-wrap items-center gap-2">
              <span className="font-medium">{offer.service_name}</span>
              <Pill className={markStyles[offer.mark]}>{markLabels[offer.mark]}</Pill>
              {offer.booked && (
                <Pill className="bg-normal-bg text-normal-fg">
                  <CalendarCheck className="size-4" aria-hidden />
                  записан
                </Pill>
              )}
            </div>
            <p className="text-sm text-subtle">{offer.patient_text}</p>
            {!offer.booked && offer.service_code === '' && (
              <p className="text-sm text-subtle">Услуги нет в клинике – запись здесь невозможна.</p>
            )}
            {offer.slots_unavailable && <p className="text-sm text-overdue">МИС не ответила – слоты недоступны.</p>}
            {offer.slots.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {offer.slots.map((item) => (
                  <button
                    key={item.id}
                    type="button"
                    disabled={disabled}
                    onClick={() => onSelect({ recommendationId: offer.recommendation_id, slotId: item.id })}
                    className={clsx(
                      'rounded-pill border px-3 py-1.5 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50',
                      item.id === selectedSlot
                        ? 'border-accent-deep bg-accent text-ink'
                        : 'border-line bg-white hover:bg-card',
                    )}
                  >
                    {formatSlot(item.starts_at)}
                  </button>
                ))}
              </div>
            )}
            {slot && (
              <Button
                className="self-start"
                disabled={disabled || booking}
                onClick={() => onBook({ recommendationId: offer.recommendation_id, slotId: slot.id })}
              >
                <CalendarCheck className="size-4" />
                Записал на {formatSlot(slot.starts_at)}
              </Button>
            )}
          </div>
        )
      })}
    </Card>
  )
}
