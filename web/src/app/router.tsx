import { createBrowserRouter, Navigate } from 'react-router'

import { AdminLayout } from '@/features/admin/AdminLayout'
import { CallPage } from '@/features/admin/call/CallPage'
import { JournalPage } from '@/features/admin/journal/JournalPage'
import { TasksPage } from '@/features/admin/tasks/TasksPage'
import { DoctorLayout } from '@/features/doctor/DoctorLayout'
import { HelpPage } from '@/features/doctor/help/HelpPage'
import { QueuePage } from '@/features/doctor/queue/QueuePage'
import { ReviewPage } from '@/features/doctor/review/ReviewPage'
import { defaultPatientId } from '@/features/patient/models'
import { LoadingState } from '@/shared/ui/States'

import { AppLayout } from './AppLayout'
import { StaffLayout } from './StaffLayout'

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    hydrateFallbackElement: <LoadingState />,
    children: [
      { index: true, element: <Navigate to="/doctor" replace /> },
      { path: 'patient', element: <Navigate to={`/patient/${defaultPatientId}`} replace /> },
      {
        path: 'patient/:patientId',
        lazy: async () => ({ Component: (await import('@/features/patient/PatientLayout')).PatientLayout }),
        children: [
          { index: true, lazy: async () => ({ Component: (await import('@/features/patient/tree/TreePage')).TreePage }) },
          {
            path: 'steps',
            lazy: async () => ({ Component: (await import('@/features/patient/steps/StepsPage')).StepsPage }),
          },
          {
            path: 'appointments',
            lazy: async () => ({
              Component: (await import('@/features/patient/appointments/AppointmentsPage')).AppointmentsPage,
            }),
          },
          {
            path: 'history',
            lazy: async () => ({ Component: (await import('@/features/patient/history/HistoryPage')).HistoryPage }),
          },
          {
            path: 'settings',
            lazy: async () => ({ Component: (await import('@/features/patient/settings/SettingsPage')).SettingsPage }),
          },
        ],
      },
      {
        element: <StaffLayout />,
        children: [
          {
            path: 'doctor',
            element: <DoctorLayout />,
            children: [
              { index: true, element: <QueuePage /> },
              { path: 'help', element: <HelpPage /> },
            ],
          },
          { path: 'doctor/cases/:caseId', element: <ReviewPage /> },
          {
            path: 'admin',
            element: <AdminLayout />,
            children: [
              { index: true, element: <TasksPage /> },
              { path: 'notifications', element: <JournalPage /> },
            ],
          },
          { path: 'admin/tasks/:taskId', element: <CallPage /> },
          {
            path: 'manager',
            lazy: async () => ({ Component: (await import('@/features/manager/ManagerPage')).ManagerPage }),
          },
        ],
      },
      { path: '*', element: <Navigate to="/doctor" replace /> },
    ],
  },
])
