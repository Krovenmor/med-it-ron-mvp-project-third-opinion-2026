import type { ErrorResponse } from './models'

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH'
  body?: unknown
  headers?: Record<string, string>
}

export async function request<T>(url: string, { method = 'GET', body, headers }: RequestOptions = {}): Promise<T> {
  const response = await fetch(url, {
    method,
    headers: {
      ...(body === undefined ? {} : { 'Content-Type': 'application/json' }),
      ...headers,
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as ErrorResponse | null
    throw new ApiError(response.status, payload?.error ?? response.statusText)
  }
  const text = await response.text()
  return (text === '' ? undefined : JSON.parse(text)) as T
}
