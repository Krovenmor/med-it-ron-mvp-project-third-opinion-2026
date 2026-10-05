import type { Recommendation, ReviewRecommendationRequest } from '@/api/models'

export interface ReviewVariables {
  recommendationId: string
  body: ReviewRecommendationRequest
}

export interface MarkCounts {
  critical: number
  minor: number
  rejected: number
  pending: number
}

export function countMarks(recommendations: Recommendation[]): MarkCounts {
  const counts: MarkCounts = { critical: 0, minor: 0, rejected: 0, pending: 0 }
  for (const recommendation of recommendations) {
    if (recommendation.review) {
      counts[recommendation.review.mark]++
    } else {
      counts.pending++
    }
  }
  return counts
}

export const confirmationText =
  'Подтверждая рекомендации, я подтверждаю, что ознакомился с заключением и данными пациента, проверил предложенный перечень исследований, консультаций и услуг, внёс необходимые изменения и согласен с выставленными приоритетами. Рекомендации ИИ носят вспомогательный характер и не заменяют клиническую оценку врача; решение о назначениях принимает врач и несёт за него ответственность.'
