import { request } from './http'
import type { Clock, DemoDeliveries, ServicesResponse } from './models'

export const mis = {
  searchServices: (query: string) =>
    request<ServicesResponse>(`/mis/api/v1/services?q=${encodeURIComponent(query)}`),

  sendDemoReports: () => request<DemoDeliveries>('/mis/demo/reports', { method: 'POST' }),

  advanceClock: (by: string) => request<Clock>('/mis/demo/clock/advance', { method: 'POST', body: { by } }),

  resetDemo: () => request<void>('/mis/demo/reset', { method: 'POST' }),
}
