import pytest

DEEPFACE_AVAILABLE = False
try:
    import deepface
    DEEPFACE_AVAILABLE = True
except ImportError:
    DEEPFACE_AVAILABLE = False

@pytest.mark.skipif(
    not DEEPFACE_AVAILABLE,
    reason="DeepFace / TensorFlow tidak terpasang di runtime lokal ini (memerlukan Python 3.10-3.12 dan download weights)",
)
def test_deepface_integration_runtime():
    """
    Test integrasi nyata yang memuat DeepFace dan kandidat model ArcFace.
    Dijalankan secara terisolasi pada environment CI / container ber-GPU / runtime Python 3.10-3.12.
    """
    from app.services.face_engine import DeepFaceEngine
    engine = DeepFaceEngine()
    assert engine.deepface is not None
