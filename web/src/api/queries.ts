import { useQuery } from '@tanstack/react-query'

import { b2b } from './b2b'
import { b2c } from './b2c'
import { ApiError } from './http'
import { mis } from './mis'

const queuePollingMs = 5_000
const cardPollingMs = 10_000
const clockSyncMs = 30_000

export const queryKeys = {
  clockOffset: ['clock-offset'] as const,
  reviewQueue: ['review-queue'] as const,
  caseDetails: (caseId: string) => ['case', caseId] as const,
  caseHistory: (caseId: string) => ['case', caseId, 'history'] as const,
  services: (query: string) => ['services', query] as const,
  operator: ['operator'] as const,
  operatorTasks: (mine: boolean) => ['operator', 'tasks', mine] as const,
  operatorCard: (taskId: string) => ['operator', 'task', taskId] as const,
  notifications: ['notifications'] as const,
  dashboard: (days: number) => ['dashboard', days] as const,
  patientRoute: (patientId: string) => ['patient', patientId, 'route'] as const,
  slots: (patientId: string, recommendationId: string) => ['patient', patientId, 'slots', recommendationId] as const,
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

export function useOperatorTasks(mine: boolean) {
  return useQuery({
    queryKey: queryKeys.operatorTasks(mine),
    queryFn: () => b2b.operatorTasks(mine),
    refetchInterval: queuePollingMs,
  })
}

export function useOperatorCard(taskId: string) {
  return useQuery({
    queryKey: queryKeys.operatorCard(taskId),
    queryFn: () => b2b.operatorCard(taskId),
    refetchInterval: cardPollingMs,
  })
}

export function useNotifications() {
  return useQuery({
    queryKey: queryKeys.notifications,
    queryFn: b2b.notifications,
    refetchInterval: queuePollingMs,
  })
}

export function useDashboard(days: number) {
  return useQuery({
    queryKey: queryKeys.dashboard(days),
    queryFn: () => b2b.dashboard(days),
    refetchInterval: cardPollingMs,
    placeholderData: (previous) => previous,
  })
}

export function usePatientRoute(patientId: string) {
  return useQuery({
    queryKey: queryKeys.patientRoute(patientId),
    queryFn: () => b2c.route(patientId),
    refetchInterval: cardPollingMs,
    retry: (failures, error) => !(error instanceof ApiError && error.status === 404) && failures < 1,
  })
}

export function useSlots(patientId: string, recommendationId: string) {
  return useQuery({
    queryKey: queryKeys.slots(patientId, recommendationId),
    queryFn: () => b2c.slots(patientId, recommendationId),
    staleTime: 0,
  })
}
