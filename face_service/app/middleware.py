import logging
import time
import uuid
from typing import Callable
from fastapi import Request, Response
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.config import settings
from app.exceptions import FaceServiceException, UnauthorizedInternalException

logger = logging.getLogger("face_service")
logging.basicConfig(
    level=logging.INFO,
    format='{"time":"%(asctime)s","level":"%(levelname)s","message":"%(message)s"}',
)

class RequestIDMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        req_id = request.headers.get("X-Request-ID")
        if not req_id:
            req_id = str(uuid.uuid4())
        request.state.request_id = req_id

        response = await call_next(request)
        response.headers["X-Request-ID"] = req_id
        return response

class SafeLoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        start_time = time.perf_counter()
        req_id = getattr(request.state, "request_id", "unknown")

        response = await call_next(request)
        latency_ms = round((time.perf_counter() - start_time) * 1000, 2)

        # Logging terstruktur aman tanpa merekam isi body, binary citra, atau token rahasia
        logger.info(
            f"http_request method={request.method} path={request.url.path} "
            f"status={response.status_code} latency_ms={latency_ms} request_id={req_id}"
        )
        return response

async def verify_internal_api_key(request: Request):
    api_key = request.headers.get("X-Internal-API-Key")
    if not api_key or api_key != settings.internal_api_key:
        raise UnauthorizedInternalException()

async def face_service_exception_handler(request: Request, exc: FaceServiceException) -> JSONResponse:
    req_id = getattr(request.state, "request_id", None)
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "error": {
                "code": exc.code,
                "message": exc.message,
                "request_id": req_id,
            }
        },
    )

async def unhandled_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    req_id = getattr(request.state, "request_id", None)
    logger.error(f"unhandled_error error={str(exc)} request_id={req_id}")
    return JSONResponse(
        status_code=500,
        content={
            "error": {
                "code": "INTERNAL_SERVER_ERROR",
                "message": "Terjadi kesalahan internal pada layanan biometrik",
                "request_id": req_id,
            }
        },
    )
