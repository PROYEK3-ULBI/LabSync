import io
import pytest
from PIL import Image
from httpx import AsyncClient, ASGITransport

from main import app
from app.config import settings

def make_dummy_jpeg(width: int = 200, height: int = 200) -> bytes:
    img = Image.new("RGB", (width, height), color=(100, 150, 200))
    buf = io.BytesIO()
    img.save(buf, format="JPEG")
    return buf.getvalue()

@pytest.mark.asyncio
async def test_auth_rejection_without_api_key():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        # Request tanpa header X-Internal-API-Key
        img_bytes = make_dummy_jpeg()
        files = {"file": ("test.jpg", img_bytes, "image/jpeg")}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"})
        assert resp.status_code == 401
        err = resp.json()["error"]
        assert err["code"] == "UNAUTHORIZED_INTERNAL"

@pytest.mark.asyncio
async def test_auth_rejection_with_wrong_key():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        img_bytes = make_dummy_jpeg()
        files = {"file": ("test.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": "wrong-key"}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)
        assert resp.status_code == 401
        err = resp.json()["error"]
        assert err["code"] == "UNAUTHORIZED_INTERNAL"

@pytest.mark.asyncio
async def test_invalid_image_non_image_file():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        files = {"file": ("test.txt", b"this is not an image", "text/plain")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)
        assert resp.status_code == 400
        err = resp.json()["error"]
        assert err["code"] == "INVALID_IMAGE"

@pytest.mark.asyncio
async def test_invalid_image_resolution_too_small():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        # Buat citra 50x50 (kurang dari minimum 160x160)
        img_bytes = make_dummy_jpeg(width=50, height=50)
        files = {"file": ("small.jpg", img_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)
        assert resp.status_code == 400
        err = resp.json()["error"]
        assert err["code"] == "INVALID_IMAGE"
        assert "terlalu kecil" in err["message"]

@pytest.mark.asyncio
async def test_file_too_large():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        # Mock file melebihi MAX_FILE_SIZE_MB (misal 5MB + 1024 bytes)
        oversized_bytes = b"0" * (settings.max_file_size_mb * 1024 * 1024 + 1024)
        files = {"file": ("large.jpg", oversized_bytes, "image/jpeg")}
        headers = {"X-Internal-API-Key": settings.internal_api_key}
        resp = await client.post("/api/v1/internal/face/extract", files=files, data={"pose": "FRONT"}, headers=headers)
        assert resp.status_code == 413
        err = resp.json()["error"]
        assert err["code"] == "FILE_TOO_LARGE"
