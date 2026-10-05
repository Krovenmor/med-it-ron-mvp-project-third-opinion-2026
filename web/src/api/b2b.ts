import { request } from './http'
import type {
  AddRecommendationRequest,
  Booking,
  CallbackRequest,
  CaseDetails,
  CaseHistory,
  CaseRef,
  ChangeUrgencyRequest,
  Clock,
  Dashboard,
  DeclineRequest,
  Notifications,
  OperatorBookingRequest,
  OperatorCard,
  OperatorTask,
  OperatorTasks,
  OutcomeRequest,
  Recommendation,
  ReviewQueue,
  ReviewRecommendationRequest,
} from './models'

const api = '/b2b/api/v1'
const demo = '/b2b/demo'
const doctor = { 'X-User-ID': import.meta.env.VITE_DOCTOR_ID }
const admin = { 'X-User-ID': import.meta.env.VITE_ADMIN_ID }
const task = (taskId: string, action = '') => `${api}/operator/tasks/${taskId}${action && `/${action}`}`

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

  operatorTasks: (mine: boolean) =>
    request<OperatorTasks>(`${api}/operator/tasks${mine ? '?mine=true' : ''}`, { headers: admin }),

  operatorCard: (taskId: string) => request<OperatorCard>(task(taskId)),

  takeTask: (taskId: string) => request<OperatorTask>(task(taskId, 'take'), { method: 'POST', headers: admin }),

  noAnswer: (taskId: string, body: OutcomeRequest) =>
    request<OperatorTask>(task(taskId, 'no-answer'), { method: 'POST', body, headers: admin }),

  scheduleCallback: (taskId: string, body: CallbackRequest) =>
    request<OperatorTask>(task(taskId, 'callback'), { method: 'POST', body, headers: admin }),

  contacted: (taskId: string, body: OutcomeRequest) =>
    request<OperatorTask>(task(taskId, 'contacted'), { method: 'POST', body, headers: admin }),

  decline: (taskId: string, body: DeclineRequest) =>
    request<OperatorTask>(task(taskId, 'decline'), { method: 'POST', body, headers: admin }),

  handToDoctor: (taskId: string, body: OutcomeRequest) =>
    request<OperatorTask>(task(taskId, 'hand-to-doctor'), { method: 'POST', body, headers: admin }),

  bookByOperator: (taskId: string, body: OperatorBookingRequest) =>
    request<Booking>(task(taskId, 'bookings'), { method: 'POST', body, headers: admin }),

  notifications: () => request<Notifications>(`${api}/notifications`),

  dashboard: (days: number) => request<Dashboard>(`${api}/dashboard?days=${days}`),

  advanceClock: (by: string) => request<Clock>(`${demo}/clock/advance`, { method: 'POST', body: { by } }),

  resetDemo: () => request<void>(`${demo}/reset`, { method: 'POST' }),
}
