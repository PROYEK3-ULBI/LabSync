from fastapi import APIRouter, Depends, File, Form, UploadFile

from app.config import settings
from app.exceptions import FaceServiceException
from app.middleware import verify_internal_api_key
from app.schemas import ExtractResponse
from app.services.image_validator import validate_and_load_image
from app.services.face_engine import get_face_engine

router = APIRouter(prefix="/api/v1/internal/face", tags=["Face Internal"])

ALLOWED_POSES = {"FRONT", "LEFT", "RIGHT"}

@router.post("/extract", response_model=ExtractResponse, dependencies=[Depends(verify_internal_api_key)])
async def extract_face_embedding(
    file: UploadFile = File(...),
    pose: str = Form("FRONT"),
):
    pose_upper = pose.strip().upper()
    if pose_upper not in ALLOWED_POSES:
        raise FaceServiceException(400, "INVALID_POSE", "Pose harus salah satu dari FRONT, LEFT, atau RIGHT")

    content = await file.read()
    try:
        img_arr, _, _ = validate_and_load_image(content, file.content_type)
        engine = get_face_engine()
        embedding, confidence = engine.extract_embedding(img_arr, pose_upper)
    finally:
        # Pembersihan memori eksplisit: raw binary dan numpy array dilepaskan
        del content
        if "img_arr" in locals():
            del img_arr

    return ExtractResponse(
        model_name=settings.model_name,
        model_version=settings.model_version,
        pose=pose_upper,
        embedding=embedding,
        dim=len(embedding),
        face_confidence=round(confidence, 4),
    )
