import { Hourglass, Sparkles } from 'lucide-react'

import type { TrunkNode } from '@/api/models'
import { formatDate } from '@/shared/lib/format'
import { Button } from '@/shared/ui/Button'
import { Sheet } from '@/shared/ui/Sheet'

import { findNode, type PatientContext, type SheetMode } from '../models'
import { stepCaption } from '../tree/captions'
import { stepLook } from '../tree/looks'
import { NodeGlyph } from '../tree/marks'

import { StepSheet } from './StepSheet'

interface NodeSheetProps {
  context: PatientContext
  nodeId: string
  mode: SheetMode
  onClose: () => void
}

export function NodeSheet({ context, nodeId, mode, onClose }: NodeSheetProps) {
  const found = findNode(context.route, nodeId)
  if (!found) {
    return null
  }
  if (found.type === 'step') {
    return <StepSheet key={`${nodeId}:${mode}`} context={context} step={found.step} mode={mode} onClose={onClose} />
  }
  return <TrunkSheet context={context} node={found.node} onClose={onClose} />
}

interface TrunkSheetProps {
  context: PatientContext
  node: TrunkNode
  onClose: () => void
}

function TrunkSheet({ context, node, onClose }: TrunkSheetProps) {
  const branches = context.route.steps.filter((step) => step.origin_id === node.id)

  return (
    <Sheet
      title={node.kind === 'study' ? 'Исследование' : 'Приём'}
      onClose={onClose}
      footer={
        <Button variant="secondary" onClick={onClose}>
          Закрыть
        </Button>
      }
    >
      <div className="flex flex-col gap-5">
        <div>
          <h3 className="text-2xl font-medium tracking-tight">{node.title}</h3>
          <p className="mt-1 text-subtle">{formatDate(node.date)}</p>
        </div>

        {node.pending ? (
          <div className="flex gap-3 rounded-card bg-card p-4">
            <Hourglass className="mt-0.5 size-5 shrink-0 text-subtle" aria-hidden />
            <p>Врач готовит ваш план по этому исследованию. Мы пришлём уведомление, когда он будет готов.</p>
          </div>
        ) : (
          node.ai_report && (
            <div className="flex gap-3 rounded-card bg-accent-tint p-4">
              <Sparkles className="mt-0.5 size-5 shrink-0 text-accent-deep" aria-hidden />
              <p>
                По заключению этого исследования ИИ-сервис подготовил рекомендации, а врач проверил их и составил ваш
                план.
              </p>
            </div>
          )
        )}

        {branches.length > 0 && (
          <section>
            <h4 className="mb-2 text-sm font-medium text-subtle">Отсюда начинаются шаги</h4>
            <ul className="flex flex-col">
              {branches.map((step) => {
                const caption = stepCaption(step)
                return (
                  <li key={step.id}>
                    <button
                      type="button"
                      onClick={() => context.openNode(step.id)}
                      className="flex w-full items-center gap-3 rounded-control px-2 py-2.5 text-left hover:bg-card"
                    >
                      <NodeGlyph look={stepLook(step)} size={22} />
                      <span className="min-w-0 flex-1">
                        <span className="block truncate font-medium">{step.title}</span>
                        <span className={caption.alert ? 'text-sm font-medium text-overdue' : 'text-sm text-subtle'}>
                          {caption.text}
                        </span>
                      </span>
                    </button>
                  </li>
                )
              })}
            </ul>
          </section>
        )}
      </div>
    </Sheet>
  )
}
