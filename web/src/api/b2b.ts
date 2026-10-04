import { request } from './http'
import type {
  AddRecommendationRequest,
  CaseDetails,
  CaseHistory,
  CaseRef,
  ChangeUrgencyRequest,
  Clock,
  Recommendation,
  ReviewQueue,
  ReviewRecommendationRequest,
} from './models'

const api = '/b2b/api/v1'
const demo = '/b2b/demo'
const doctor = { 'X-User-ID': import.meta.env.VITE_DOCTOR_ID }

export const b2b = {
  clock: () => request<Clock>(`${api}/clock`),

  reviewQueue: () => request<ReviewQueue>(`${api}/review/queue`),

  caseDetails: (caseId: string) => request<CaseDetails>(`${api}/cases/${caseId}`),

  caseHistory: (caseId: string) => request<CaseHistory>(`${api}/cases/${caseId}/history`),

  openCase: (caseId: string) => request<void>(`${api}/cases/${caseId}/open`, { method: 'POST', headers: doctor }),

  reviewRecommendation: (caseId: string, recommendationId: string, body: ReviewRecommendationRequest) =>
    request<Recommendation>(`${api}/cases/${caseId}/recommendations/${recommendationId}`, {
      method: 'PATCH',
      body,
      headers: doctor,
    }),

  addRecommendation: (caseId: string, body: AddRecommendationRequest) =>
    request<Recommendation>(`${api}/cases/${caseId}/recommendations`, { method: 'POST', body, headers: doctor }),

  changeUrgency: (caseId: string, body: ChangeUrgencyRequest) =>
    request<void>(`${api}/cases/${caseId}/urgency`, { method: 'PUT', body, headers: doctor }),

  confirm: (caseId: string) => request<CaseRef>(`${api}/cases/${caseId}/confirm`, { method: 'POST', headers: doctor }),

  advanceClock: (by: string) => request<Clock>(`${demo}/clock/advance`, { method: 'POST', body: { by } }),

  resetDemo: () => request<void>(`${demo}/reset`, { method: 'POST' }),
}
