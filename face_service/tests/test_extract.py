import pytest
from httpx import AsyncClient, ASGITransport

from main import app
from app.config import settings
from app.services.face_engine import MockFaceEngine, set_face_engine
from tests.test_image_validation import make_dummy_jpeg

@pytest.mark.asyncio
async def test_extract_success():
    set_face_engine(MockFaceEngine())
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg(200, 200)
        files = {"file": ("face.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)

        assert resp.status_code == 200
        data = resp.json()
        assert data["model_name"] == "ArcFace"
        assert data["pose"] == "FRONT"
        assert data["dim"] == 512
        assert len(data["embedding"]) == 512
        assert data["face_confidence"] > 0.0

@pytest.mark.asyncio
async def test_extract_no_face_detected():
    set_face_engine(MockFaceEngine(simulate_no_face=True))
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg(200, 200)
        files = {"file": ("noface.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "LEFT"}, headers=headers)

        assert resp.status_code == 422
        err = resp.json()["error"]
        assert err["code"] == "NO_FACE_DETECTED"

@pytest.mark.asyncio
async def test_extract_multiple_faces_detected():
    set_face_engine(MockFaceEngine(simulate_multi_face=True))
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg(200, 200)
        files = {"file": ("multiface.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "RIGHT"}, headers=headers)

        assert resp.status_code == 422
        err = resp.json()["error"]
        assert err["code"] == "MULTIPLE_FACES_DETECTED"

@pytest.mark.asyncio
async def test_extract_model_error():
    set_face_engine(MockFaceEngine(simulate_model_error=True))
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg(200, 200)
        files = {"file": ("model_err.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)

        assert resp.status_code == 500
        err = resp.json()["error"]
        assert err["code"] == "MODEL_ERROR"

@pytest.mark.asyncio
async def test_extract_invalid_pose():
    set_face_engine(MockFaceEngine())
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg(200, 200)
        files = {"file": ("face.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "TOP"}, headers=headers)

        assert resp.status_code == 400
        err = resp.json()["error"]
        assert err["code"] == "INVALID_POSE"
