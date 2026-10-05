import { TriangleAlert } from 'lucide-react'

import type { CaseDetails } from '@/api/models'
import { ageAt, formatAge, formatDate, formatDateTime } from '@/shared/lib/format'
import { modalityLabels, sexLabels } from '@/shared/lib/labels'
import { UrgencyBadge } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'
import { Card, TwoToneHeading } from '@/shared/ui/Layout'

import { SlaTimer } from '../SlaTimer'

interface CaseHeaderProps {
  details: CaseDetails
  now: Date
  editable: boolean
  onChangeUrgency: () => void
}

export function CaseHeader({ details, now, editable, onChangeUrgency }: CaseHeaderProps) {
  const { patient } = details
  return (
    <Card className="flex flex-wrap items-start justify-between gap-6">
      <div className="flex flex-col gap-2">
        <TwoToneHeading
          main={`Пациент ${patient.id}`}
          rest={`${formatAge(ageAt(patient.birth_date, now))}, ${sexLabels[patient.sex]}`}
        />
        <div className="text-subtle">
          {modalityLabels[details.modality]} от {formatDate(details.performed_at)} · поступило{' '}
          {formatDateTime(details.received_at)}
        </div>
      </div>
      <div className="flex items-center gap-3">
        {details.urgency && <UrgencyBadge urgency={details.urgency} />}
        {editable && (
          <Button variant="secondary" size="sm" onClick={onChangeUrgency}>
            Изменить срочность
          </Button>
        )}
      </div>
    </Card>
  )
}

interface EmergencyBannerProps {
  receivedAt: string
  now: Date
}

export function EmergencyBanner({ receivedAt, now }: EmergencyBannerProps) {
  return (
    <div className="flex items-center justify-between gap-4 rounded-card bg-emergency-bg px-6 py-4 text-emergency-fg">
      <div className="flex items-center gap-3">
        <TriangleAlert className="size-5" aria-hidden />
        <span className="font-medium">Неотложный случай: передано администратору</span>
      </div>
      <SlaTimer receivedAt={receivedAt} now={now} />
    </div>
  )
}
