import clsx from 'clsx'

import type { Mark } from '@/api/models'
import { markLabels, urgencyOrder } from '@/shared/lib/labels'
import { UrgencyBadge } from '@/shared/ui/Badge'
import { Card, TwoToneHeading } from '@/shared/ui/Layout'

import { confirmationText } from '../review/models'

const urgencyDescriptions = {
  emergency: 'Критическая находка. Кейс вверху очереди, администратор сразу получает задачу связаться с пациентом. Проверить нужно за 30 минут.',
  priority: 'Пациенту сразу приходят пуш и СМС. Если не записался за сутки – задача администратору.',
  planned: 'Пуш сразу, напоминания на 3-й и 7-й день. Через 7 дней без записи – задача администратору.',
  normal: 'Отклонений нет. Пациент получит напоминание о плановом скрининге в срок.',
}

const markDescriptions: Record<Mark, string> = {
  critical: 'Пациенту важно пройти этот шаг в первую очередь.',
  minor: 'Шаг нужен, но не срочный.',
  rejected: 'Пациенту не покажется. Причина обязательна: «Противопоказано» или «Другое» с комментарием.',
}

const markSwatches: Record<Mark, string> = {
  critical: 'bg-mark-critical',
  minor: 'bg-mark-minor',
  rejected: 'bg-line',
}

export function HelpPage() {
  return (
    <div className="flex flex-col gap-5">
      <TwoToneHeading main="Справка" rest="для врача" />

      <Card>
        <h2 className="mb-4 text-lg font-medium">Уровни срочности</h2>
        <ul className="flex flex-col gap-4">
          {urgencyOrder.map((urgency) => (
            <li key={urgency} className="grid grid-cols-[180px_1fr] items-start gap-4">
              <div>
                <UrgencyBadge urgency={urgency} />
              </div>
              <span className="text-subtle">{urgencyDescriptions[urgency]}</span>
            </li>
          ))}
        </ul>
      </Card>

      <Card>
        <h2 className="mb-4 text-lg font-medium">Отметки рекомендаций</h2>
        <ul className="flex flex-col gap-4">
          {(Object.keys(markDescriptions) as Mark[]).map((mark) => (
            <li key={mark} className="grid grid-cols-[180px_1fr] items-center gap-4">
              <span className="flex items-center gap-2 font-medium">
                <span className={clsx('size-4 rounded-full', markSwatches[mark])} aria-hidden />
                {markLabels[mark]}
              </span>
              <span className="text-subtle">{markDescriptions[mark]}</span>
            </li>
          ))}
        </ul>
        <p className="mt-4 text-sm text-subtle">
          Подтвердить кейс можно, когда отмечена каждая строка. Понизить срочность можно только с указанием причины.
        </p>
      </Card>

      <Card>
        <h2 className="mb-4 text-lg font-medium">Дисклеймер</h2>
        <p className="leading-relaxed text-subtle">{confirmationText}</p>
        <p className="mt-4 leading-relaxed text-subtle">
          Пациент видит только подтверждённые вами рекомендации и понятный текст без диагноза: «Рекомендации подготовил
          врач на основе вашего заключения. Они не заменяют консультацию».
        </p>
      </Card>
    </div>
  )
}
