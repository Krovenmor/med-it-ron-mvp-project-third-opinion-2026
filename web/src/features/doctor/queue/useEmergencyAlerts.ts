import { useEffect, useRef } from 'react'
import { toast } from 'sonner'

import type { QueueCase } from '@/api/models'
import { playAlertSound, showBrowserNotification } from '@/shared/lib/alerts'

export function useEmergencyAlerts(cases: QueueCase[] | undefined) {
  const seen = useRef<Set<string> | null>(null)

  useEffect(() => {
    if (!cases) {
      return
    }
    const emergencies = cases.filter((item) => item.urgency === 'emergency')
    if (seen.current === null) {
      seen.current = new Set(emergencies.map((item) => item.case_id))
      return
    }

    const known = seen.current
    const fresh = emergencies.filter((item) => !known.has(item.case_id))
    if (fresh.length === 0) {
      return
    }
    fresh.forEach((item) => known.add(item.case_id))

    const patients = fresh.map((item) => item.patient.id).join(', ')
    playAlertSound()
    toast.warning('Неотложный случай в очереди', { description: `Пациент ${patients}` })
    showBrowserNotification('Неотложный случай', `Пациент ${patients} ожидает проверки`)
  }, [cases])
}
