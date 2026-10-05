import { useCallback } from 'react'
import { useLocation, useNavigate, useOutletContext, useSearchParams } from 'react-router'

import type { PatientContext, SheetMode } from './models'

interface SheetState {
  sheet?: boolean
}

export function usePatient(): PatientContext {
  return useOutletContext<PatientContext>()
}

export function useNodeSheet() {
  const [params, setParams] = useSearchParams()
  const location = useLocation()
  const navigate = useNavigate()
  const nodeId = params.get('step')
  const mode: SheetMode = params.get('book') === '1' ? 'book' : 'details'
  const fromApp = Boolean((location.state as SheetState | null)?.sheet)

  const open = useCallback(
    (id: string, next: SheetMode = 'details') => {
      const search = new URLSearchParams(params)
      search.set('step', id)
      if (next === 'book') {
        search.set('book', '1')
      } else {
        search.delete('book')
      }
      setParams(search, { replace: nodeId !== null, state: { sheet: true } satisfies SheetState })
    },
    [params, setParams, nodeId],
  )

  const close = useCallback(() => {
    if (fromApp) {
      navigate(-1)
      return
    }
    const search = new URLSearchParams(params)
    search.delete('step')
    search.delete('book')
    setParams(search, { replace: true })
  }, [fromApp, navigate, params, setParams])

  return { nodeId, mode, open, close }
}
