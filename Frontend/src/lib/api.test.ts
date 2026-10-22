import test from "node:test"
import assert from "node:assert"
import { ApiError, getApiBaseUrl, request, api, unwrapList } from "./api.ts"

test("getApiBaseUrl fallback", () => {
  const url = getApiBaseUrl()
  assert.ok(url.length > 0)
  assert.ok(!url.endsWith("/"))
})

test("ApiError properties", () => {
  const err = new ApiError("Kredensial salah", 401, "UNAUTHORIZED", "req-123")
  assert.strictEqual(err.message, "Kredensial salah")
  assert.strictEqual(err.statusCode, 401)
  assert.strictEqual(err.code, "UNAUTHORIZED")
  assert.strictEqual(err.requestId, "req-123")
  assert.strictEqual(err.name, "ApiError")
})

test("request parses JSON success with credentials include", async () => {
  const originalFetch = globalThis.fetch
  let capturedOptions: RequestInit | undefined

  globalThis.fetch = ((async (
    url: string | URL | Request,
    options?: RequestInit,
  ) => {
    capturedOptions = options
    return new Response(
      JSON.stringify({
        message: "Success",
        user: { id: "u-1", identifier: "TI240001" },
      }),
      {
        status: 200,
        headers: {
          "Content-Type": "application/json",
          "X-Request-ID": "req-test-99",
        },
      },
    )
  }) as typeof fetch)

  try {
    const res = await api.auth.login({
      identifier: "TI240001",
      password: "secretpassword",
    })
    assert.strictEqual(res.message, "Success")
    assert.strictEqual(res.user.identifier, "TI240001")
    assert.strictEqual(capturedOptions?.credentials, "include")
    assert.ok(capturedOptions?.headers)
    const headers = capturedOptions.headers as Record<string, string>
    assert.ok(headers["X-Request-ID"])
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("request parses structured error response", async () => {
  const originalFetch = globalThis.fetch

  globalThis.fetch = ((async () => {
    return new Response(
      JSON.stringify({
        error: {
          code: "UNAUTHORIZED",
          message: "Identifier atau password salah",
          request_id: "req-err-401",
        },
      }),
      {
        status: 401,
        headers: {
          "Content-Type": "application/json",
          "X-Request-ID": "req-err-401",
        },
      },
    )
  }) as typeof fetch)

  try {
    await assert.rejects(
      async () => {
        await api.auth.login({ identifier: "TI240001", password: "wrong" })
      },
      (err: unknown) => {
        assert.ok(err instanceof ApiError)
        assert.strictEqual(err.statusCode, 401)
        assert.strictEqual(err.code, "UNAUTHORIZED")
        assert.strictEqual(err.message, "Identifier atau password salah")
        assert.strictEqual(err.requestId, "req-err-401")
        return true
      },
    )
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("request handles 204 no content", async () => {
  const originalFetch = globalThis.fetch

  globalThis.fetch = ((async () => {
    return new Response(null, {
      status: 204,
    })
  }) as typeof fetch)

  try {
    const res = await request("/test/empty")
    assert.strictEqual(res, undefined)
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("request keeps browser multipart content-type boundary", async () => {
  const originalFetch = globalThis.fetch
  let capturedHeaders: HeadersInit | undefined
  globalThis.fetch = ((async (_url, options) => {
    capturedHeaders = options?.headers
    return new Response(JSON.stringify({ message: "ok" }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    })
  }) as typeof fetch)

  try {
    const form = new FormData()
    form.append("session_id", "s-1")
    form.append("image", new Blob(["x"], { type: "image/jpeg" }), "face.jpg")
    await api.attendance.checkIn(form)
    const headers = capturedHeaders as Record<string, string>
    assert.strictEqual(headers["Content-Type"], undefined)
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("unwrapList handles backend list shapes", () => {
  assert.deepStrictEqual(unwrapList([{ id: "1" }]), [{ id: "1" }])
  assert.deepStrictEqual(unwrapList({ data: [{ id: "2" }] }), [{ id: "2" }])
  assert.deepStrictEqual(unwrapList({ sessions: [{ id: "3" }] }), [{ id: "3" }])
  assert.deepStrictEqual(unwrapList({ allocations: [{ id: "4" }] }), [
    { id: "4" },
  ])
  assert.deepStrictEqual(unwrapList({ users: [{ id: "5" }] }), [{ id: "5" }])
  assert.deepStrictEqual(unwrapList({ audit_logs: [{ id: "6" }] }), [
    { id: "6" },
  ])
})

test("admin endpoints use REST paths", async () => {
  const originalFetch = globalThis.fetch
  const calls: string[] = []
  globalThis.fetch = (async (url) => {
    calls.push(String(url))
    return new Response(JSON.stringify({ data: [] }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    })
  }) as typeof fetch

  try {
    await api.users.list()
    await api.assets.list()
    await api.allocations.list()
    await api.maintenance.list()
    await api.audit.list()
    assert.ok(calls.some((url) => url.endsWith("/users")))
    assert.ok(calls.some((url) => url.endsWith("/assets")))
    assert.ok(calls.some((url) => url.endsWith("/allocations")))
    assert.ok(calls.some((url) => url.endsWith("/maintenance")))
    assert.ok(calls.some((url) => url.endsWith("/audit-logs")))
  } finally {
    globalThis.fetch = originalFetch
  }
})
