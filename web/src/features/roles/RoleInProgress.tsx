import { Hourglass } from 'lucide-react'

import { EmptyState } from '@/shared/ui/States'

interface RoleInProgressProps {
  role: string
}

export function RoleInProgress({ role }: RoleInProgressProps) {
  return (
    <EmptyState
      icon={Hourglass}
      title={`Раздел «${role}» в разработке`}
      text="В демо уже доступна роль «Врач»: очередь на проверку и карточка проверки рекомендаций."
    />
  )
}
