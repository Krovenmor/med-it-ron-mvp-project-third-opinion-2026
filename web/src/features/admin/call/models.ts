import type { Offer, OperatorCard, Slot } from '@/api/models'

export interface Selection {
  recommendationId: string
  slotId: string
}

export interface ScriptOffer {
  offer: Offer
  slot?: Slot
}

export function scriptOffer(card: OperatorCard, selection: Selection | null): ScriptOffer | null {
  const selected = selection && card.offers.find((offer) => offer.recommendation_id === selection.recommendationId)
  if (selected) {
    return { offer: selected, slot: selected.slots.find((slot) => slot.id === selection.slotId) }
  }
  const open = card.offers.find((offer) => !offer.booked)
  return open ? { offer: open, slot: open.slots[0] } : null
}
