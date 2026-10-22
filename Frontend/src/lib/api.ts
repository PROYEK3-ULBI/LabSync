// Centralized API Client for LabSync REST API (/api/v1)

export class ApiError extends Error {
  statusCode: number
  code: string
  requestId?: string

  constructor(
    message: string,
    statusCode: number,
    code: string,
    requestId?: string,
  ) {
    super(message)
    this.name = "ApiError"
    this.statusCode = statusCode
    this.code = code
    this.requestId = requestId
  }
}

export type UserRole = "STUDENT" | "ADMIN"

export type AuthUser = {
  id: string
  identifier: string
  email: string
  full_name: string
  role: UserRole
  status: string
  face_enrolled?: boolean
  created_at?: string
  updated_at?: string
}

export type UserPayload = {
  identifier: string
  email: string
  full_name: string
  password?: string
  role?: UserRole
  status?: string
}

export type LoginCredentials = {
  identifier: string
  password: string
}

export type LoginResponse = {
  message: string
  user: AuthUser
}

export type FaceStatus = {
  enrolled: boolean
  status?: string
  model_name?: string
  model_version?: string
  enrolled_at?: string
  sample_count?: number
  poses?: string[]
}

export type LabSession = {
  id: string
  code: string
  title: string
  course_code?: string
  location: string
  session_date: string
  start_at: string
  end_at: string
  status: string
  attendance_status?: string
}

export type AttendanceRecord = {
  id: string
  session_id: string
  user_id?: string
  checked_in_at?: string
  checked_out_at?: string
  verification_status?: string
  status?: string
  session?: LabSession
}

export type AttendanceResponse = {
  message: string
  attendance: AttendanceRecord
}

export type AssetSummary = {
  id?: string
  asset_code: string
  asset_name: string
  category?: string
  brand?: string
  model?: string
  location?: string
  status?: string
}

export type AssetRecord = AssetSummary & {
  id: string
  serial_number?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export type PaginationMeta = {
  page?: number
  total_pages?: number
  total?: number
}

export type AssetAllocation = {
  id: string
  asset_id?: string
  user_id?: string
  session_id?: string
  allocated_at?: string
  returned_at?: string
  condition?: string
  status?: string
  asset?: AssetSummary
  session?: LabSession
}

export type MaintenanceRecord = {
  id: string
  asset_id?: string
  issue_type?: string
  description?: string
  priority?: string
  status?: string
  resolution?: string
  cost?: number
  created_at?: string
  completed_at?: string
  asset?: AssetSummary
}

export type AuditRecord = {
  id: string
  action: string
  entity_type: string
  entity_id: string
  status: string
  created_at?: string
  user?: AuthUser
}

export type Paginated = {
  data?: unknown[]
  items?: unknown[]
  sessions?: unknown[]
  attendance?: unknown[]
  allocations?: unknown[]
  users?: unknown[]
  assets?: unknown[]
  maintenance?: unknown[]
  audit_logs?: unknown[]
  meta?: PaginationMeta
  pagination?: PaginationMeta
}

export type ErrorPayload = {
  error: {
    code: string
    message: string
    request_id?: string
  }
}

// Menentukan Base URL dari environment variable dengan fallback yang aman
export function getApiBaseUrl(): string {
  // Dukung VITE_API_URL maupun VITE_API_BASE_URL
  const envUrl =
    (typeof import.meta !== "undefined" && import.meta.env
      ? import.meta.env.VITE_API_URL || import.meta.env.VITE_API_BASE_URL
      : undefined) || ""

  if (envUrl) {
    return envUrl.replace(/\/+$/, "")
  }

  // Fallback default: relatif /api/v1 saat menggunakan Vite proxy, atau localhost:8000
  if (
    typeof window !== "undefined" &&
    window.location &&
    window.location.origin
  ) {
    return "/api/v1"
  }
  return "http://localhost:8000/api/v1"
}

type RequestOptions = {
  method?: string
  headers?: Record<string, string>
  body?: unknown
  timeoutMs?: number
  signal?: AbortSignal
}

function generateRequestId(): string {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  return "req-" + Math.random().toString(36).substring(2, 11)
}

export async function request<T>(
  endpoint: string,
  options: RequestOptions = {},
): Promise<T> {
  const baseUrl = getApiBaseUrl()
  const cleanEndpoint = endpoint.startsWith("/") ? endpoint : `/${endpoint}`
  const url = `${baseUrl}${cleanEndpoint}`

  const headers: Record<string, string> = {
    Accept: "application/json",
    "X-Request-ID": generateRequestId(),
    ...(options.headers || {}),
  }

  const isFormData =
    typeof FormData !== "undefined" && options.body instanceof FormData
  let reqBody: BodyInit | null | undefined = undefined

  if (options.body != null) {
    if (isFormData) {
      // Biarkan browser menetapkan Content-Type beserta multipart boundary secara otomatis
      reqBody = (options.body as FormData)
      delete headers["Content-Type"]
    } else if (typeof options.body === "string") {
      headers["Content-Type"] = "application/json"
      reqBody = options.body
    } else {
      headers["Content-Type"] = "application/json"
      reqBody = JSON.stringify(options.body)
    }
  }

  const timeoutMs = options.timeoutMs ?? 15000
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)

  try {
    const res = await fetch(url, {
      method: options.method || "GET",
      headers,
      body: reqBody,
      // INVARIANT: selalu sertakan credentials agar HttpOnly cookie JWT terkirim ke backend
      credentials: "include",
      signal: options.signal || controller.signal,
    })

    const contentType = res.headers.get("content-type") || ""
    const isJson = contentType.includes("application/json")

    if (!res.ok) {
      let errCode = "HTTP_ERROR"
      let errMsg = `Permintaan gagal dengan kode status ${res.status}`
      let requestId: string | undefined =
        res.headers.get("x-request-id") || undefined

      if (isJson) {
        try {
          const errData = (await res.json()) as ErrorPayload
          if (errData && errData.error) {
            errCode = errData.error.code || errCode
            errMsg = errData.error.message || errMsg
            requestId = errData.error.request_id || requestId
          }
        } catch {
          // Gagal parse json error body
        }
      } else {
        const text = await res.text()
        if (text) errMsg = text
      }

      throw new ApiError(errMsg, res.status, errCode, requestId)
    }

    if (res.status === 204) {
      return undefined as unknown as T
    }

    if (isJson) {
      return (await res.json()) as T
    }

    return (await res.text()) as unknown as T
  } catch (err: unknown) {
    if (err instanceof ApiError) {
      throw err
    }
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new ApiError(
        "Waktu tunggu koneksi ke server terlampaui (timeout)",
        408,
        "REQUEST_TIMEOUT",
      )
    }
    const message =
      err instanceof Error ? err.message : "Gagal menghubungi server"
    throw new ApiError(
      `Gagal terhubung ke server backend: ${message}`,
      0,
      "NETWORK_ERROR",
    )
  } finally {
    clearTimeout(timeoutId)
  }
}

// API Endpoints
export const api = {
  auth: {
    login: (credentials: LoginCredentials) =>
      request<LoginResponse>("/auth/login", {
        method: "POST",
        body: credentials,
      }),
    logout: () =>
      request<{ message: string }>("/auth/logout", {
        method: "POST",
      }),
    me: () => request<AuthUser>("/auth/me"),
  },
  face: {
    status: () => request<FaceStatus>("/face/status"),
    enroll: (form: FormData) =>
      request<{ message: string }>("/face/enroll", {
        method: "POST",
        body: form,
        timeoutMs: 90000,
      }),
  },
  sessions: {
    list: () => request<Paginated | LabSession[]>("/sessions"),
    active: () => request<Paginated | LabSession[]>("/sessions/active"),
    detail: (id: string) => request<LabSession>(`/sessions/${id}`),
    create: (body: Partial<LabSession>) =>
      request<LabSession>("/sessions", { method: "POST", body }),
    update: (id: string, body: Partial<LabSession>) =>
      request<LabSession>(`/sessions/${id}`, { method: "PUT", body }),
    updateStatus: (id: string, status: string) =>
      request<LabSession>(`/sessions/${id}/status`, {
        method: "PATCH",
        body: { status },
      }),
  },
  attendance: {
    history: () =>
      request<Paginated | AttendanceRecord[]>("/attendance/history"),
    checkIn: (form: FormData) =>
      request<AttendanceResponse>("/attendance/check-in", {
        method: "POST",
        body: form,
        timeoutMs: 90000,
      }),
    checkOut: (sessionId: string) =>
      request<AttendanceResponse>("/attendance/check-out", {
        method: "POST",
        body: { session_id: sessionId },
      }),
    list: () => request<Paginated | AttendanceRecord[]>("/attendance"),
  },
  assets: {
    list: () => request<Paginated | AssetRecord[]>("/assets"),
    create: (body: Partial<AssetRecord>) =>
      request<AssetRecord>("/assets", { method: "POST", body }),
    update: (id: string, body: Partial<AssetRecord>) =>
      request<AssetRecord>(`/assets/${id}`, { method: "PUT", body }),
    updateStatus: (id: string, status: string) =>
      request<AssetRecord>(`/assets/${id}/status`, {
        method: "PATCH",
        body: { status },
      }),
    myAllocations: () =>
      request<Paginated | AssetAllocation[]>("/assets/my-allocations"),
    history: (assetId: string) =>
      request<Paginated | AssetAllocation[]>(`/assets/${assetId}/history`),
  },
  users: {
    list: () => request<Paginated | AuthUser[]>("/users"),
    create: (body: UserPayload) =>
      request<AuthUser>("/users", { method: "POST", body }),
    update: (id: string, body: Partial<UserPayload>) =>
      request<AuthUser>(`/users/${id}`, { method: "PUT", body }),
    updateStatus: (id: string, status: string) =>
      request<AuthUser>(`/users/${id}/status`, {
        method: "PATCH",
        body: { status },
      }),
  },
  allocations: {
    list: () => request<Paginated | AssetAllocation[]>("/allocations"),
    create: (body: { asset_id: string; user_id: string; session_id: string }) =>
      request<AssetAllocation>("/allocations", { method: "POST", body }),
    returnAsset: (id: string, condition: string) =>
      request<AssetAllocation>(`/allocations/${id}/return`, {
        method: "POST",
        body: { condition },
      }),
  },
  maintenance: {
    list: () => request<Paginated | MaintenanceRecord[]>("/maintenance"),
    create: (body: Partial<MaintenanceRecord>) =>
      request<MaintenanceRecord>("/maintenance", { method: "POST", body }),
    update: (id: string, body: Partial<MaintenanceRecord>) =>
      request<MaintenanceRecord>(`/maintenance/${id}`, {
        method: "PUT",
        body,
      }),
  },
  dashboard: {
    admin: () => request<Record<string, unknown>>("/dashboard/admin"),
  },
  reports: {
    attendance: () => request<unknown>("/reports/attendance"),
    assets: () => request<unknown>("/reports/assets"),
    allocations: () => request<unknown>("/reports/allocations"),
    maintenance: () => request<unknown>("/reports/maintenance"),
    csvUrl: (name: string) => `${getApiBaseUrl()}/reports/${name}?format=csv`,
  },
  audit: {
    list: () => request<Paginated | AuditRecord[]>("/audit-logs"),
  },
}

export function unwrapList<T>(payload: Paginated | T[]): T[] {
  if (Array.isArray(payload)) return payload
  return (payload.data ||
    payload.items ||
    payload.sessions ||
    payload.attendance ||
    payload.allocations ||
    payload.users ||
    payload.assets ||
    payload.maintenance ||
    payload.audit_logs ||
    []) as T[]
}
