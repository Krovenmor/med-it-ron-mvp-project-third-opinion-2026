/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_DOCTOR_ID: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
