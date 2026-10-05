import { useEffect, useState } from 'react'

import { useClockOffset } from '@/api/queries'

export const emergencyReviewSlaMs = 30 * 60_000

export function useServerNow(): Date {
  const { data: offset = 0 } = useClockOffset()
  const [tick, setTick] = useState(() => Date.now())

  useEffect(() => {
    const timer = setInterval(() => setTick(Date.now()), 1_000)
    return () => clearInterval(timer)
  }, [])

  return new Date(tick + offset)
}
