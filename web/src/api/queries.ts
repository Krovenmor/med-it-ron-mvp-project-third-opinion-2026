import { useQuery } from '@tanstack/react-query'

import { b2b } from './b2b'
import { mis } from './mis'

const queuePollingMs = 5_000
const clockSyncMs = 30_000

export const queryKeys = {
  clockOffset: ['clock-offset'] as const,
  reviewQueue: ['review-queue'] as const,
  caseDetails: (caseId: string) => ['case', caseId] as const,
  caseHistory: (caseId: string) => ['case', caseId, 'history'] as const,
  services: (query: string) => ['services', query] as const,
}

export function useClockOffset() {
  return useQuery({
    queryKey: queryKeys.clockOffset,
    queryFn: async () => {
      const sentAt = Date.now()
      const { now } = await b2b.clock()
      return Date.parse(now) - sentAt
    },
    refetchInterval: clockSyncMs,
  })
}

export function useReviewQueue() {
  return useQuery({
    queryKey: queryKeys.reviewQueue,
    queryFn: b2b.reviewQueue,
    refetchInterval: queuePollingMs,
  })
}

export function useCaseDetails(caseId: string) {
  return useQuery({
    queryKey: queryKeys.caseDetails(caseId),
    queryFn: () => b2b.caseDetails(caseId),
  })
}

export function useCaseHistory(caseId: string) {
  return useQuery({
    queryKey: queryKeys.caseHistory(caseId),
    queryFn: () => b2b.caseHistory(caseId),
    retry: false,
  })
}

export function useServiceSearch(query: string) {
  return useQuery({
    queryKey: queryKeys.services(query),
    queryFn: () => mis.searchServices(query),
    placeholderData: (previous) => previous,
  })
}
