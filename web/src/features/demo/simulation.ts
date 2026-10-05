import type { QueryClient } from '@tanstack/react-query'

import { b2b } from '@/api/b2b'
import type { DemoDelivery } from '@/api/models'
import { queryKeys } from '@/api/queries'

const pollMs = 500
const timeoutMs = 20_000

export function createdCaseIds(delivered: DemoDelivery[]): string[] {
  return delivered.flatMap((delivery) =>
    delivery.status === 201 && 'case_id' in delivery.response ? [delivery.response.case_id] : [],
  )
}

export async function waitUntilInQueue(queryClient: QueryClient, caseIds: string[]): Promise<boolean> {
  const deadline = Date.now() + timeoutMs
  for (;;) {
    const queue = await queryClient.fetchQuery({
      queryKey: queryKeys.reviewQueue,
      queryFn: b2b.reviewQueue,
      staleTime: 0,
    })
    const inQueue = new Set(queue.cases.map((item) => item.case_id))
    if (caseIds.every((id) => inQueue.has(id))) {
      return true
    }
    if (Date.now() > deadline) {
      return false
    }
    await new Promise((resolve) => setTimeout(resolve, pollMs))
  }
}
