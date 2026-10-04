export type Modality = 'CT' | 'DX' | 'MG'
export type Sex = 'male' | 'female'
export type Urgency = 'normal' | 'planned' | 'priority' | 'emergency'
export type CaseStatus =
  | 'draft'
  | 'in_review'
  | 'confirmed'
  | 'notified'
  | 'booked'
  | 'completed'
  | 'declined'
  | 'unreachable'
export type Mark = 'critical' | 'minor' | 'rejected'
export type AcceptedMark = Exclude<Mark, 'rejected'>
export type RejectReason = 'contraindicated' | 'other'
export type RecommendationSource = 'ai' | 'doctor'
export type Importance = 'high' | 'low'

export interface Patient {
  id: string
  full_name: string
  birth_date: string
  sex: Sex
}

export interface QueueCase {
  case_id: string
  urgency: Urgency
  modality: Modality
  performed_at: string
  received_at: string
  patient: Patient
  recommendations_total: number
  recommendations_reviewed: number
  opened_at?: string
}

export interface ReviewQueue {
  cases: QueueCase[]
}

export interface Review {
  mark: Mark
  reject_reason?: RejectReason
  reject_comment?: string
  reviewed_by: string
  reviewed_at: string
}

export interface Recommendation {
  id: string
  position: number
  source: RecommendationSource
  service_code: string
  service_name: string
  importance: Importance
  rationale: string
  guideline_ref: string
  patient_text: string
  already_booked: boolean
  review: Review | null
}

export interface CaseDetails {
  id: string
  patient_id: string
  source_system: string
  study_id: string
  modality: Modality
  body_site: string
  performed_at: string
  conclusion: string
  status: CaseStatus
  urgency?: Urgency
  catalog_version?: string
  guidelines_version?: string
  received_at: string
  updated_at: string
  patient: Patient
  recommendations: Recommendation[]
}

export interface HistoryVisit {
  service_code: string
  service_name: string
  visited_at: string
}

export interface HistoryAppointment {
  service_code: string
  service_name: string
  scheduled_at: string
}

export interface CaseHistory {
  visits: HistoryVisit[]
  appointments: HistoryAppointment[]
}

export interface CaseRef {
  case_id: string
  status: CaseStatus
}

export interface Clock {
  now: string
}

export interface ReviewRecommendationRequest {
  mark?: Mark
  reject_reason?: RejectReason
  reject_comment?: string
  patient_text?: string
}

export interface AddRecommendationRequest {
  service_code: string
  service_name: string
  patient_text: string
  rationale: string
  mark: AcceptedMark
}

export interface ChangeUrgencyRequest {
  urgency: Urgency
  reason: string
}

export interface CatalogService {
  code: string
  name: string
  description: string
}

export interface ServicesResponse {
  services: CatalogService[]
}

export interface DemoDelivery {
  study_id: string
  status: number
  response: CaseRef | ErrorResponse
}

export interface DemoDeliveries {
  delivered: DemoDelivery[]
  error?: string
}

export interface ErrorResponse {
  error: string
}
