import json
import pytest
from httpx import AsyncClient, ASGITransport

from main import app
from app.config import settings
from app.services.face_engine import MockFaceEngine, set_face_engine
from tests.test_image_validation import make_dummy_jpeg

def sample_vector(val: float = 0.1) -> list[float]:
    vec = [val] * 512
    norm = sum(x * x for x in vec) ** 0.5
    return [x / norm for x in vec]

@pytest.mark.asyncio
async def test_verify_success_calibration_pending():
    # Model returns deterministic embedding
    fixed_vec = sample_vector(0.1)
    set_face_engine(MockFaceEngine(fixed_embedding=fixed_vec))

    # Pastikan default is_calibrated = False
    settings.is_calibrated = False

    enrolled = [
        {"pose": "FRONT", "embedding": fixed_vec},         # Distance ~ 0.0 (identical)
        {"pose": "LEFT", "embedding": sample_vector(0.3)},  # Distance ~ 0.0 (same direction)
        {"pose": "RIGHT", "embedding": sample_vector(-0.1)}, # Distance ~ 2.0 (opposite direction)
    ]

    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg()
        files = {"file": ("live.jpg", img_bytes, "image/jpeg")}
        data = {
            "enrolled_samples": json.dumps(enrolled),
            "threshold": "0.40",
            "strategy": "MIN_DISTANCE",
        }
        headers = {"X-Internal-API-Key": settings.internal_api_key}

        resp = await client.post("/api/v1/internal/face/verify", files=files, data=data, headers=headers)
        assert resp.status_code == 200
        res = resp.json()

        # Karena is_calibrated=False, verified HARUS False untuk mencegah false acceptance
        assert res["verified"] is False
        assert res["is_calibrated"] is False
        assert "CALIBRATION_PENDING" in res["decision_rule"]
        assert res["best_matched_pose"] == "FRONT"
        assert res["distance"] <= 0.01
        assert "FRONT" in res["distances_by_pose"]
        assert "LEFT" in res["distances_by_pose"]
        assert "RIGHT" in res["distances_by_pose"]

@pytest.mark.asyncio
async def test_verify_calibrated_decision():
    fixed_vec = sample_vector(0.1)
    set_face_engine(MockFaceEngine(fixed_embedding=fixed_vec))

    # Simulasikan kalibrasi telah selesai
    settings.is_calibrated = True

    enrolled = [
        {"pose": "FRONT", "embedding": fixed_vec},         # Cocok identik -> dist ~ 0.0
        {"pose": "LEFT", "embedding": sample_vector(0.5)},
        {"pose": "RIGHT", "embedding": sample_vector(-0.5)},
    ]

    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg()
        files = {"file": ("live.jpg", img_bytes, "image/jpeg")}
        data = {
            "enrolled_samples": json.dumps(enrolled),
            "threshold": "0.40",
            "strategy": "MIN_DISTANCE",
        }
        headers = {"X-Internal-API-Key": settings.internal_api_key}

        resp = await client.post("/api/v1/internal/face/verify", files=files, data=data, headers=headers)
        assert resp.status_code == 200
        res = resp.json()

        # Bila kalibrasi aktif dan distance <= threshold -> verified True
        assert res["verified"] is True
        assert res["is_calibrated"] is True
        assert "THRESHOLD_PASSED" in res["decision_rule"]

@pytest.mark.asyncio
async def test_verify_incomplete_samples():
    set_face_engine(MockFaceEngine())
    # Hanya kirim 2 sampel (kurang dari 3 pose FRONT, LEFT, RIGHT)
    enrolled = [
        {"pose": "FRONT", "embedding": sample_vector(0.1)},
        {"pose": "LEFT", "embedding": sample_vector(0.2)},
    ]

    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg()
        files = {"file": ("live.jpg", img_bytes, "image/jpeg")}
        data = {"enrolled_samples": json.dumps(enrolled)}
        headers = {"X-Internal-API-Key": settings.internal_api_key}

        resp = await client.post("/api/v1/internal/face/verify", files=files, data=data, headers=headers)
        assert resp.status_code == 400
        err = resp.json()["error"]
        assert err["code"] == "INCOMPLETE_SAMPLES"

@pytest.mark.asyncio
async def test_verify_invalid_samples_payload():
    set_face_engine(MockFaceEngine())
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg()
        files = {"file": ("live.jpg", img_bytes, "image/jpeg")}
        data = {"enrolled_samples": "not-a-valid-json"}
        headers = {"X-Internal-API-Key": settings.internal_api_key}

        resp = await client.post("/api/v1/internal/face/verify", files=files, data=data, headers=headers)
        assert resp.status_code == 400
        err = resp.json()["error"]
        assert err["code"] == "INVALID_SAMPLES_PAYLOAD"
