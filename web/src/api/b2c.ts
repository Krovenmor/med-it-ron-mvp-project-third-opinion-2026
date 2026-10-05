import { request } from './http'
import type { Appointment, PatientDeclineRequest, PatientRoute, Slots } from './models'

const api = '/b2c/api/v1'
const recommendation = (patientId: string, recommendationId: string, action: string) =>
  `${api}/patients/${patientId}/recommendations/${recommendationId}/${action}`

export const b2c = {
  route: (patientId: string) => request<PatientRoute>(`${api}/patients/${patientId}/route`),

  slots: (patientId: string, recommendationId: string) =>
    request<Slots>(recommendation(patientId, recommendationId, 'slots')),

  book: (patientId: string, recommendationId: string, slotId: string) =>
    request<Appointment>(recommendation(patientId, recommendationId, 'appointments'), {
      method: 'POST',
      body: { slot_id: slotId },
    }),

  decline: (patientId: string, recommendationId: string, body: PatientDeclineRequest) =>
    request<void>(recommendation(patientId, recommendationId, 'decline'), { method: 'POST', body }),

  requestHelp: (patientId: string, recommendationId: string) =>
    request<void>(recommendation(patientId, recommendationId, 'help-requests'), { method: 'POST' }),
}
