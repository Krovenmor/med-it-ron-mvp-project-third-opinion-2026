/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_DOCTOR_ID: string
  readonly VITE_ADMIN_ID: string
  readonly VITE_CLINIC_NAME: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
