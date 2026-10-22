from typing import Dict, List, Optional
from pydantic import BaseModel, Field

class HealthResponse(BaseModel):
    status: str
    service: str
    model_name: str
    model_version: str
    is_calibrated: bool

class ErrorDetail(BaseModel):
    code: str
    message: str
    request_id: Optional[str] = None

class ErrorResponse(BaseModel):
    error: ErrorDetail

class ExtractResponse(BaseModel):
    model_name: str
    model_version: str
    pose: str
    embedding: List[float]
    dim: int
    face_confidence: float

class EnrolledSample(BaseModel):
    pose: str
    embedding: List[float]

class VerifyResponse(BaseModel):
    verified: bool
    distance: float
    threshold: float
    strategy: str
    best_matched_pose: str
    distances_by_pose: Dict[str, float]
    is_calibrated: bool
    decision_rule: str
    model_name: str
