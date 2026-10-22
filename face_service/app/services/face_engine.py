import math
from typing import Any, Dict, List, Tuple
import numpy as np

from app.config import settings
from app.exceptions import (
    NoFaceDetectedException,
    MultipleFacesDetectedException,
    ModelInferenceException,
)
from app.schemas import EnrolledSample

def cosine_distance(u: List[float], v: List[float]) -> float:
    if len(u) != len(v) or len(u) == 0:
        raise ValueError("Dimensi vektor embedding tidak cocok atau kosong")

    dot = sum(a * b for a, b in zip(u, v))
    norm_u = math.sqrt(sum(a * a for a in u))
    norm_v = math.sqrt(sum(b * b for b in v))

    if norm_u == 0.0 or norm_v == 0.0:
        return 1.0

    similarity = dot / (norm_u * norm_v)
    similarity = max(-1.0, min(1.0, similarity))
    return float(1.0 - similarity)

class BaseFaceEngine:
    def extract_embedding(self, img_arr: np.ndarray, pose: str) -> Tuple[List[float], float]:
        raise NotImplementedError

    def compare_1_to_3(
        self,
        live_embedding: List[float],
        enrolled_samples: List[EnrolledSample],
        threshold: float,
        strategy: str,
    ) -> Dict[str, Any]:
        if not enrolled_samples:
            raise ValueError("Daftar sampel enrollment tidak boleh kosong")

        distances_by_pose: Dict[str, float] = {}
        for sample in enrolled_samples:
            dist = cosine_distance(live_embedding, sample.embedding)
            distances_by_pose[sample.pose] = round(float(dist), 4)

        if strategy == "AVERAGE_DISTANCE":
            computed_distance = round(sum(distances_by_pose.values()) / len(distances_by_pose), 4)
            best_pose = "AVERAGE"
        else:
            # Default: MIN_DISTANCE (kecocokan pose terbaik)
            best_pose = min(distances_by_pose, key=distances_by_pose.get)
            computed_distance = distances_by_pose[best_pose]

        # Kebijakan Kalibrasi (Requirement 6):
        # Jika kalibrasi empiris belum selesai (is_calibrated = False),
        # hasil verified dinyatakan False untuk mencegah false acceptance pada presensi sebelum kalibrasi.
        is_calibrated = settings.is_calibrated
        if not is_calibrated:
            verified = False
            rule_desc = f"{strategy}_UNVERIFIED_CALIBRATION_PENDING"
        else:
            verified = computed_distance <= threshold
            rule_desc = f"{strategy}_THRESHOLD_{'PASSED' if verified else 'EXCEEDED'}"

        return {
            "verified": verified,
            "distance": computed_distance,
            "threshold": threshold,
            "strategy": strategy,
            "best_matched_pose": best_pose,
            "distances_by_pose": distances_by_pose,
            "is_calibrated": is_calibrated,
            "decision_rule": rule_desc,
            "model_name": settings.model_name,
        }

class MockFaceEngine(BaseFaceEngine):
    def __init__(
        self,
        simulate_no_face: bool = False,
        simulate_multi_face: bool = False,
        simulate_model_error: bool = False,
        fixed_embedding: List[float] | None = None,
    ):
        self.simulate_no_face = simulate_no_face
        self.simulate_multi_face = simulate_multi_face
        self.simulate_model_error = simulate_model_error
        self.fixed_embedding = fixed_embedding

    def extract_embedding(self, img_arr: np.ndarray, pose: str) -> Tuple[List[float], float]:
        if self.simulate_model_error:
            raise ModelInferenceException("Simulasi model failure")
        if self.simulate_no_face:
            raise NoFaceDetectedException("Tidak ada wajah terdeteksi pada citra")
        if self.simulate_multi_face:
            raise MultipleFacesDetectedException("Terdeteksi lebih dari satu wajah")

        if self.fixed_embedding is not None:
            return self.fixed_embedding, 0.99

        # Hasilkan embedding 512 dimensi sintetis terdeterministik berbasis ukuran gambar
        dim = 512
        seed_val = float(np.sum(img_arr) % 1000) / 1000.0
        vec = [math.sin(seed_val + float(i)) for i in range(dim)]
        norm = math.sqrt(sum(x * x for x in vec))
        normalized = [x / norm for x in vec]
        return normalized, 0.98

class DeepFaceEngine(BaseFaceEngine):
    def __init__(self):
        try:
            from deepface import DeepFace  # type: ignore
            self.deepface = DeepFace
        except ImportError as e:
            raise RuntimeError(
                f"DeepFace tidak dapat diimpor pada runtime ini ({str(e)}). "
                "Pastikan dependensi diinstal pada Python 3.10-3.12."
            )

    def extract_embedding(self, img_arr: np.ndarray, pose: str) -> Tuple[List[float], float]:
        try:
            results = self.deepface.represent(
                img_path=img_arr,
                model_name=settings.model_name,
                detector_backend=settings.detector_backend,
                enforce_detection=True,
                align=True,
            )
        except Exception as e:
            err_msg = str(e).lower()
            if "no face" in err_msg or "face could not be detected" in err_msg:
                raise NoFaceDetectedException("Tidak ada wajah terdeteksi pada citra")
            raise ModelInferenceException(f"Inferensi model DeepFace gagal: {str(e)}")

        if not results or len(results) == 0:
            raise NoFaceDetectedException("Tidak ada wajah terdeteksi pada citra")
        if len(results) > 1:
            raise MultipleFacesDetectedException(
                f"Terdeteksi {len(results)} wajah pada citra, pastikan hanya ada 1 wajah"
            )

        face_res = results[0]
        embedding = face_res.get("embedding", [])
        confidence = float(face_res.get("face_confidence", 1.0))
        return embedding, confidence

# Singleton / Global Engine State
_engine: BaseFaceEngine = MockFaceEngine()

def get_face_engine() -> BaseFaceEngine:
    global _engine
    return _engine

def set_face_engine(engine: BaseFaceEngine) -> None:
    global _engine
    _engine = engine
