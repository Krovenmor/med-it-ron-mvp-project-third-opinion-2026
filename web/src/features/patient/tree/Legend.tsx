import type { ReactNode } from 'react'

import { Modal } from '@/shared/ui/Modal'

import type { BranchTone } from './layout'
import type { NodeLook } from './looks'
import { AiBadge, BranchGlyph, NodeGlyph } from './marks'

const nodes: { look: NodeLook; label: string }[] = [
  { look: 'done', label: 'Пройдено: приём или услуга уже были' },
  { look: 'booked-minor', label: 'Вы записаны, приём впереди' },
  { look: 'planned', label: 'Назначено, но вы ещё не записались' },
  { look: 'overdue', label: 'Срок прошёл: запишитесь' },
]

const branches: { tone: BranchTone; label: string }[] = [
  { tone: 'dashed', label: 'Пунктир: шаг назначен, записи нет' },
  { tone: 'green', label: 'Зелёная: шаг пройден' },
  { tone: 'minor', label: 'Жёлтая: записаны, отметка врача «Вторично»' },
  { tone: 'critical', label: 'Коричневая: записаны, отметка врача «Приоритетно»' },
]

interface LegendProps {
  open: boolean
  onClose: () => void
}

export function Legend({ open, onClose }: LegendProps) {
  return (
    <Modal open={open} title="Как читать дерево" onClose={onClose}>
      <div className="flex flex-col gap-6">
        <p className="text-subtle">
          Время идёт снизу вверх. Ствол – ваши приёмы у терапевта и исследования, ветки – шаги, которые назначил врач.
          Толщина ветки ничего не значит, важность показывает цвет.
        </p>
        <LegendGroup title="Узлы">
          {nodes.map((item) => (
            <LegendRow key={item.look} sample={<NodeGlyph look={item.look} size={26} />} label={item.label} />
          ))}
        </LegendGroup>
        <LegendGroup title="Ветки">
          <LegendRow
            sample={
              <svg width={44} height={22} aria-hidden>
                <path d="M22 1 Q18 1 18 6 L16 21 L22 21 Z" className="fill-accent" />
                <path d="M22 1 Q26 1 26 6 L28 21 L22 21 Z" className="fill-accent-dark" />
              </svg>
            }
            label="Ствол: приёмы у терапевта и исследования"
          />
          {branches.map((item) => (
            <LegendRow key={item.tone} sample={<BranchGlyph tone={item.tone} />} label={item.label} />
          ))}
        </LegendGroup>
        <LegendGroup title="Другое">
          <LegendRow
            sample={
              <svg width={44} height={16} aria-hidden>
                <line x1={2} x2={42} y1={8} y2={8} strokeWidth={1.5} className="stroke-today" />
                <circle cx={4} cy={8} r={3.5} className="fill-today" />
              </svg>
            }
            label="Сегодня"
          />
          <LegendRow
            sample={
              <svg width={44} height={22} viewBox="-22 -11 44 22" aria-hidden>
                <AiBadge x={0} y={0} />
              </svg>
            }
            label="Исследование с заключением ИИ: с него начался ваш план"
          />
          <LegendRow sample={<span className="w-11 text-center text-lg leading-none text-muted">⋯</span>} label="Сжатый период без событий" />
        </LegendGroup>
      </div>
    </Modal>
  )
}

interface LegendGroupProps {
  title: string
  children: ReactNode
}

function LegendGroup({ title, children }: LegendGroupProps) {
  return (
    <section>
      <h3 className="mb-2 text-sm font-medium text-subtle">{title}</h3>
      <ul className="flex flex-col gap-2.5">{children}</ul>
    </section>
  )
}

interface LegendRowProps {
  sample: ReactNode
  label: string
}

function LegendRow({ sample, label }: LegendRowProps) {
  return (
    <li className="flex items-center gap-3">
      <span className="flex w-11 shrink-0 justify-center">{sample}</span>
      {label}
    </li>
  )
}
