import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { b2b } from '@/api/b2b'
import type { AddRecommendationRequest, CaseDetails, ChangeUrgencyRequest, Recommendation } from '@/api/models'
import { queryKeys } from '@/api/queries'

import type { ReviewVariables } from './models'

export function useReviewActions(caseId: string) {
  const queryClient = useQueryClient()
  const caseKey = queryKeys.caseDetails(caseId)

  const updateRecommendations = (update: (current: Recommendation[]) => Recommendation[]) =>
    queryClient.setQueryData<CaseDetails>(caseKey, (current) =>
      current ? { ...current, recommendations: update(current.recommendations) } : current,
    )
  const refreshQueue = () => queryClient.invalidateQueries({ queryKey: queryKeys.reviewQueue })
  const showError = (error: Error) => toast.error(error.message)

  const review = useMutation({
    mutationFn: ({ recommendationId, body }: ReviewVariables) =>
      b2b.reviewRecommendation(caseId, recommendationId, body),
    onSuccess: (saved) => {
      updateRecommendations((current) => current.map((rec) => (rec.id === saved.id ? saved : rec)))
      void refreshQueue()
    },
    onError: showError,
  })

  const add = useMutation({
    mutationFn: (body: AddRecommendationRequest) => b2b.addRecommendation(caseId, body),
    onSuccess: (created) => {
      updateRecommendations((current) => [...current, created])
      void refreshQueue()
    },
    onError: showError,
  })

  const changeUrgency = useMutation({
    mutationFn: (body: ChangeUrgencyRequest) => b2b.changeUrgency(caseId, body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: caseKey })
      void refreshQueue()
    },
    onError: showError,
  })

  const confirm = useMutation({
    mutationFn: () => b2b.confirm(caseId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: caseKey })
      void refreshQueue()
    },
    onError: showError,
  })

  return { review, add, changeUrgency, confirm }
}
