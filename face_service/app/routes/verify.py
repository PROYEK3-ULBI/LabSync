import json
from typing import List, Optional
from fastapi import APIRouter, Depends, File, Form, UploadFile

from app.config import settings
from app.exceptions import FaceServiceException
from app.middleware import verify_internal_api_key
from app.schemas import EnrolledSample, VerifyResponse
from app.services.image_validator import validate_and_load_image
from app.services.face_engine import get_face_engine

router = APIRouter(prefix="/api/v1/internal/face", tags=["Face Internal"])

ALLOWED_STRATEGIES = {"MIN_DISTANCE", "AVERAGE_DISTANCE"}

@router.post("/verify", response_model=VerifyResponse, dependencies=[Depends(verify_internal_api_key)])
async def verify_face(
    file: UploadFile = File(...),
    enrolled_samples: str = Form(..., description="JSON string array of 3 enrolled samples"),
    threshold: Optional[float] = Form(None),
    strategy: Optional[str] = Form(None),
):
    # 1. Parse dan validasi sampel enrollment (FRONT, LEFT, RIGHT)
    try:
        raw_samples = json.loads(enrolled_samples)
        if not isinstance(raw_samples, list) or len(raw_samples) == 0:
            raise ValueError("Harus berupa array data sampel")
        samples: List[EnrolledSample] = [EnrolledSample(**item) for item in raw_samples]
    except Exception as e:
        raise FaceServiceException(400, "INVALID_SAMPLES_PAYLOAD", f"Format data enrolled_samples tidak valid: {str(e)}")

    if len(samples) != 3:
        raise FaceServiceException(
            400,
            "INCOMPLETE_SAMPLES",
            f"Jumlah sampel enrollment harus tepat 3 pose (FRONT, LEFT, RIGHT). Diterima: {len(samples)} sampel",
        )

    # 2. Tentukan threshold dan strategi
    eff_threshold = threshold if threshold is not None else settings.default_similarity_threshold
    eff_strategy = (strategy.strip().upper() if strategy else settings.default_strategy)
    if eff_strategy not in ALLOWED_STRATEGIES:
        raise FaceServiceException(400, "INVALID_STRATEGY", "Strategi harus MIN_DISTANCE atau AVERAGE_DISTANCE")

    # 3. Baca dan validasi live frame
    content = await file.read()
    try:
        img_arr, _, _ = validate_and_load_image(content, file.content_type)
        engine = get_face_engine()
        live_embedding, _ = engine.extract_embedding(img_arr, "LIVE")
    finally:
        del content
        if "img_arr" in locals():
            del img_arr

    # 4. Bandingkan 1:1 live embedding dengan 3 sampel enrollment milik user terpilih
    result = engine.compare_1_to_3(
        live_embedding=live_embedding,
        enrolled_samples=samples,
        threshold=eff_threshold,
        strategy=eff_strategy,
    )

    return VerifyResponse(**result)
