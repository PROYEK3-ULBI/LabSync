from fastapi import APIRouter
from app.config import settings
from app.schemas import HealthResponse

router = APIRouter(tags=["System"])

@router.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="healthy",
        service="face_service",
        model_name=settings.model_name,
        model_version=settings.model_version,
        is_calibrated=settings.is_calibrated,
    )
