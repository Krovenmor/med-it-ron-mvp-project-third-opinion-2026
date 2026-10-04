import { useMutation, useQueryClient } from '@tanstack/react-query'
import clsx from 'clsx'
import { FastForward, RotateCcw, Send } from 'lucide-react'
import { useState } from 'react'
import { NavLink } from 'react-router'
import { toast } from 'sonner'

import { b2b } from '@/api/b2b'
import { mis } from '@/api/mis'
import { formatDateTime } from '@/shared/lib/format'
import { useServerNow } from '@/shared/lib/time'
import { Button } from '@/shared/ui/Button'
import { Modal } from '@/shared/ui/Modal'

import type { SimulationResult } from './models'
import { roles } from './roles'
import { createdCaseIds, waitUntilInQueue } from './simulation'

const simulationToast = 'demo-simulation'

export function DemoPanel() {
  const queryClient = useQueryClient()
  const now = useServerNow()
  const [resetOpen, setResetOpen] = useState(false)

  const refresh = () => queryClient.invalidateQueries()
  const fail = (error: Error) => toast.error(error.message)

  const simulate = useMutation({
    mutationFn: async (): Promise<SimulationResult> => {
      const { delivered } = await mis.sendDemoReports()
      const created = createdCaseIds(delivered)
      if (created.length === 0) {
        return { created: 0, ready: true }
      }
      toast.loading(`Поступило заключений: ${created.length}. ИИ оценивает…`, { id: simulationToast })
      return { created: created.length, ready: await waitUntilInQueue(queryClient, created) }
    },
    onSuccess: ({ created, ready }) => {
      if (created === 0) {
        toast.info('Новых заключений нет – все уже приняты', { id: simulationToast })
      } else if (ready) {
        toast.success(`В очереди новых кейсов: ${created}`, { id: simulationToast })
      } else {
        toast.warning('ИИ ещё оценивает часть заключений – они появятся в очереди позже', { id: simulationToast })
      }
      void refresh()
    },
    onError: (error) => toast.error(error.message, { id: simulationToast }),
  })

  const advance = useMutation({
    mutationFn: b2b.advanceClock,
    onSuccess: (_, by) => {
      toast.success(by === '30m' ? 'Время перемотано на 30 минут' : 'Время перемотано на 1 день')
      void refresh()
    },
    onError: fail,
  })

  const reset = useMutation({
    mutationFn: async () => {
      await b2b.resetDemo()
      await mis.resetDemo()
    },
    onSuccess: () => {
      setResetOpen(false)
      toast.success('Демо сброшено')
      void refresh()
    },
    onError: fail,
  })

  return (
    <div className="bg-ink text-white">
      <div className="mx-auto flex max-w-[1280px] flex-wrap items-center gap-x-6 gap-y-3 px-8 py-3">
        <span className="text-xs font-medium uppercase tracking-wider text-muted">Демо</span>
        <nav className="flex gap-1">
          {roles.map((role) => (
            <NavLink
              key={role.path}
              to={role.path}
              className={({ isActive }) =>
                clsx(
                  'rounded-pill px-3 py-1 text-sm font-medium transition-colors',
                  isActive ? 'bg-accent text-ink' : 'text-white/80 hover:bg-white/10',
                )
              }
            >
              {role.label}
            </NavLink>
          ))}
        </nav>
        <div className="ml-auto flex flex-wrap items-center gap-2">
          <span className="mr-2 text-sm text-white/70">Время демо: {formatDateTime(now)}</span>
          <Button size="sm" onClick={() => simulate.mutate()} disabled={simulate.isPending}>
            <Send className="size-4" />
            Симулировать поступление из МИС
          </Button>
          <Button size="sm" variant="secondary" onClick={() => advance.mutate('30m')} disabled={advance.isPending}>
            <FastForward className="size-4" />
            +30 мин
          </Button>
          <Button size="sm" variant="secondary" onClick={() => advance.mutate('24h')} disabled={advance.isPending}>
            <FastForward className="size-4" />
            +1 день
          </Button>
          <Button size="sm" variant="secondary" onClick={() => setResetOpen(true)}>
            <RotateCcw className="size-4" />
            Сбросить демо
          </Button>
        </div>
      </div>

      <Modal
        open={resetOpen}
        title="Сбросить демо?"
        onClose={() => setResetOpen(false)}
        footer={
          <>
            <Button variant="secondary" onClick={() => setResetOpen(false)}>
              Отмена
            </Button>
            <Button onClick={() => reset.mutate()} disabled={reset.isPending}>
              Сбросить
            </Button>
          </>
        }
      >
        <p className="text-subtle">
          Удалятся все кейсы, отметки врача и записи пациентов, время вернётся к текущему. Демо-заключения можно
          будет отправить заново.
        </p>
      </Modal>
    </div>
  )
}
