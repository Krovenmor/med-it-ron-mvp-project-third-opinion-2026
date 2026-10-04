import { createBrowserRouter, Navigate } from 'react-router'

import { DoctorLayout } from '@/features/doctor/DoctorLayout'
import { HelpPage } from '@/features/doctor/help/HelpPage'
import { QueuePage } from '@/features/doctor/queue/QueuePage'
import { ReviewPage } from '@/features/doctor/review/ReviewPage'
import { RoleInProgress } from '@/features/roles/RoleInProgress'

import { AppLayout } from './AppLayout'

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/doctor" replace /> },
      {
        path: 'doctor',
        element: <DoctorLayout />,
        children: [
          { index: true, element: <QueuePage /> },
          { path: 'help', element: <HelpPage /> },
        ],
      },
      { path: 'doctor/cases/:caseId', element: <ReviewPage /> },
      { path: 'patient', element: <RoleInProgress role="Пациент" /> },
      { path: 'admin', element: <RoleInProgress role="Администратор" /> },
      { path: 'manager', element: <RoleInProgress role="Руководитель" /> },
      { path: '*', element: <Navigate to="/doctor" replace /> },
    ],
  },
])
