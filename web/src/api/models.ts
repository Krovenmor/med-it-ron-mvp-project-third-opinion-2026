export type Modality = 'CT' | 'DX' | 'MG'
export type Sex = 'male' | 'female'
export type Urgency = 'normal' | 'planned' | 'priority' | 'emergency'
export type CaseStatus =
  | 'received'
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

export type TaskReason = 'emergency' | 'no_booking' | 'help_request'
export type ActiveTaskStatus = 'new' | 'in_progress' | 'no_answer' | 'callback'
export type TaskStatus = ActiveTaskStatus | 'contacted' | 'booked' | 'declined' | 'handed_to_doctor' | 'cancelled'
export type CallOutcome = 'no_answer' | 'callback' | 'contacted' | 'declined' | 'booked' | 'handed_to_doctor'
export type DeclineReason = 'expensive' | 'far' | 'other_clinic' | 'not_needed' | 'other'
export type NotificationRecipient = 'patient' | 'duty_doctor' | 'admin_on_duty'
export type NotificationChannel = 'push' | 'sms' | 'staff'
export type NotificationKind =
  | 'plan_ready'
  | 'no_findings'
  | 'reminder'
  | 'urgent_contact'
  | 'unreachable'
  | 'review_escalation'
  | 'contact_escalation'
  | 'handed_to_doctor'

export interface OperatorTask {
  id: string
  case_id: string
  reason: TaskReason
  status: TaskStatus
  assignee: string
  due_at: string
  next_call_at: string | null
  attempts: number
  needs_doctor: boolean
  created_at: string
  closed_at: string | null
}

export interface OperatorCase {
  id: string
  status: CaseStatus
  urgency: Urgency
  modality: Modality
  book_by: string
}

export interface OperatorPatient {
  id: string
  full_name: string
  phone: string
}

export interface OperatorTaskItem {
  task: OperatorTask
  case: OperatorCase
  patient: OperatorPatient
  offer_service: string
}

export interface OperatorTasks {
  tasks: OperatorTaskItem[]
}

export interface CallAttempt {
  id: string
  outcome: CallOutcome
  decline_reason?: DeclineReason
  comment: string
  callback_at: string | null
  actor: string
  created_at: string
}

export interface Slot {
  id: string
  starts_at: string
}

export interface Offer {
  recommendation_id: string
  service_code: string
  service_name: string
  patient_text: string
  mark: AcceptedMark
  booked: boolean
  slots: Slot[]
  slots_unavailable: boolean
}

export interface OperatorCard {
  task: OperatorTask
  case: OperatorCase
  patient: OperatorPatient
  attempts: CallAttempt[]
  offers: Offer[]
}

export interface OutcomeRequest {
  comment?: string
}

export interface CallbackRequest extends OutcomeRequest {
  call_at: string
}

export interface DeclineRequest extends OutcomeRequest {
  reason: DeclineReason
}

export interface OperatorBookingRequest extends OutcomeRequest {
  recommendation_id: string
  slot_id: string
}

export interface Booking {
  booking_id: string
  case_id: string
  recommendation_id: string
  appointment_id: string
  case_status: CaseStatus
}

export interface Notification {
  id: string
  case_id: string
  recipient: NotificationRecipient
  channel: NotificationChannel
  kind: NotificationKind
  text: string
  urgency: Urgency
  patient: { id: string; full_name: string }
  created_at: string
}

export interface Notifications {
  notifications: Notification[]
}

export interface Ratio {
  on_time: number
  total: number
}

export interface DashboardFunnel {
  received: number
  recommended: number
  confirmed: number
  notified: number
  booked: number
  completed: number
}

export interface PeriodStats {
  funnel: DashboardFunnel
  doctor_sla: Ratio
  operator_sla: Ratio
  bookings: { self: number; operator: number }
}

export interface UrgencySlice {
  urgency: Urgency
  confirmed: number
  booked: number
  doctor_sla: Ratio
  operator_sla: Ratio
}

export interface ModalitySlice {
  modality: Modality
  confirmed: number
  booked: number
}

export interface DeclineCount {
  reason: DeclineReason
  count: number
}

export interface TrendPoint {
  from: string
  received: number
  booked: number
}

export interface DashboardPeriod {
  from: string
  to: string
  days: number
}

export interface Dashboard {
  period: DashboardPeriod
  previous_period: DashboardPeriod
  current: PeriodStats
  previous: PeriodStats
  by_urgency: UrgencySlice[]
  by_modality: ModalitySlice[]
  decline_reasons: DeclineCount[]
  trend_step: 'day' | 'week'
  trend: TrendPoint[]
}

export type TrunkKind = 'visit' | 'study'
export type StepKind = 'recommendation' | 'visit'
export type StepStatus = 'recommended' | 'booked' | 'done' | 'overdue' | 'declined'

export interface Clinic {
  name: string
  phone: string
  address: string
}

export interface Appointment {
  id: string
  service_name: string
  scheduled_at: string
}

export interface TrunkNode {
  id: string
  kind: TrunkKind
  title: string
  date: string
  ai_report: boolean
  pending: boolean
}

export interface RouteStep {
  id: string
  kind: StepKind
  origin_id: string
  case_id?: string
  title: string
  text?: string
  mark?: AcceptedMark
  status: StepStatus
  late: boolean
  assigned_at?: string
  due_at?: string
  date: string
  in_clinic: boolean
  bookable: boolean
  appointment: Appointment | null
  visited_at?: string
  decline_reason?: DeclineReason
}

export interface PatientRoute {
  now: string
  patient: { full_name: string }
  clinic: Clinic
  urgent_contact: boolean
  pending_review: boolean
  no_findings: boolean
  progress: number | null
  next_step_id: string
  trunk: TrunkNode[]
  steps: RouteStep[]
  appointments: Appointment[]
}

export interface Slots {
  slots: Slot[]
}

export interface PatientDeclineRequest {
  reason: DeclineReason
  comment: string
}

export interface ErrorResponse {
  error: string
}
