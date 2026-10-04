import { request } from './http'
import type { DemoDeliveries, ServicesResponse } from './models'

export const mis = {
  searchServices: (query: string) =>
    request<ServicesResponse>(`/mis/api/v1/services?q=${encodeURIComponent(query)}`),

  sendDemoReports: () => request<DemoDeliveries>('/mis/demo/reports', { method: 'POST' }),

  resetDemo: () => request<void>('/mis/demo/reset', { method: 'POST' }),
}
