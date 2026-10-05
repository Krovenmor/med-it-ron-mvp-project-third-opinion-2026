import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { b2b } from '@/api/b2b'
import type { CallbackRequest, DeclineRequest, OperatorBookingRequest, OutcomeRequest } from '@/api/models'
import { queryKeys } from '@/api/queries'

export function useCallActions(taskId: string) {
  const queryClient = useQueryClient()

  const refresh = () =>
    Promise.all([
      queryClient.invalidateQueries({ queryKey: queryKeys.operator }),
      queryClient.invalidateQueries({ queryKey: queryKeys.notifications }),
    ])
  const options = { onSuccess: () => void refresh(), onError: (error: Error) => toast.error(error.message) }

  return {
    take: useMutation({ mutationFn: () => b2b.takeTask(taskId), ...options }),
    noAnswer: useMutation({ mutationFn: (body: OutcomeRequest) => b2b.noAnswer(taskId, body), ...options }),
    callback: useMutation({ mutationFn: (body: CallbackRequest) => b2b.scheduleCallback(taskId, body), ...options }),
    contacted: useMutation({ mutationFn: (body: OutcomeRequest) => b2b.contacted(taskId, body), ...options }),
    decline: useMutation({ mutationFn: (body: DeclineRequest) => b2b.decline(taskId, body), ...options }),
    handToDoctor: useMutation({ mutationFn: (body: OutcomeRequest) => b2b.handToDoctor(taskId, body), ...options }),
    book: useMutation({ mutationFn: (body: OperatorBookingRequest) => b2b.bookByOperator(taskId, body), ...options }),
  }
}

export type CallActions = ReturnType<typeof useCallActions>
