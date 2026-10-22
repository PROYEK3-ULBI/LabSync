import os
from contextlib import asynccontextmanager
from fastapi import FastAPI

from app.config import settings
from app.exceptions import FaceServiceException
from app.middleware import (
    RequestIDMiddleware,
    SafeLoggingMiddleware,
    face_service_exception_handler,
    unhandled_exception_handler,
)
from app.routes import health, extract, verify
from app.services.face_engine import DeepFaceEngine, MockFaceEngine, set_face_engine

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Inisialisasi engine saat startup
    try:
        engine = DeepFaceEngine()
        set_face_engine(engine)
    except Exception as e:
        # Fallback ke Mock engine jika DeepFace belum terpasang (misal pada environment testing / dev awal)
        set_face_engine(MockFaceEngine())

    yield

app = FastAPI(
    title="LabSync Face Verification Internal Service",
    description="Microservice internal biometrik 1:1 face verification",
    version="1.0.0",
    docs_url="/docs" if settings.app_env != "production" else None,
    redoc_url=None,
    lifespan=lifespan,
)

# Registrasi Middleware
app.add_middleware(SafeLoggingMiddleware)
app.add_middleware(RequestIDMiddleware)

# Registrasi Exception Handlers
app.add_exception_handler(FaceServiceException, face_service_exception_handler)
app.add_exception_handler(Exception, unhandled_exception_handler)

# Registrasi Rute
app.include_router(health.router)
app.include_router(extract.router)
app.include_router(verify.router)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=settings.port, reload=True)
