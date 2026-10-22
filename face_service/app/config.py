import os
from dataclasses import dataclass


def load_dotenv(path: str = ".env") -> None:
    if not os.path.exists(path):
        return
    with open(path, encoding="utf-8") as file:
        for line in file:
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, value = line.split("=", 1)
            os.environ.setdefault(key.strip(), value.strip().strip('"').strip("'"))


load_dotenv()

@dataclass
class Settings:
    app_env: str = "development"
    port: int = 8001
    internal_api_key: str = "labsync-internal-secret-key"
    model_name: str = "ArcFace"
    model_version: str = "1.0.0"
    detector_backend: str = "opencv"
    max_file_size_mb: int = 5
    min_image_width: int = 160
    min_image_height: int = 160
    default_similarity_threshold: float = 0.40
    is_calibrated: bool = False
    default_strategy: str = "MIN_DISTANCE"

    @classmethod
    def load(cls) -> "Settings":
        app_env = os.getenv("APP_ENV", "development")
        port = int(os.getenv("PORT", "8001"))
        api_key = os.getenv("INTERNAL_API_KEY", "labsync-internal-secret-key")
        model_name = os.getenv("MODEL_NAME", "ArcFace")
        model_version = os.getenv("MODEL_VERSION", "1.0.0")
        detector_backend = os.getenv("DETECTOR_BACKEND", "opencv")
        max_file_size_mb = int(os.getenv("MAX_FILE_SIZE_MB", "5"))
        min_width = int(os.getenv("MIN_IMAGE_WIDTH", "160"))
        min_height = int(os.getenv("MIN_IMAGE_HEIGHT", "160"))
        threshold = float(os.getenv("DEFAULT_SIMILARITY_THRESHOLD", "0.40"))
        is_calibrated = os.getenv("IS_CALIBRATED", "false").lower() in ("true", "1")
        default_strategy = os.getenv("DEFAULT_STRATEGY", "MIN_DISTANCE")

        return cls(
            app_env=app_env,
            port=port,
            internal_api_key=api_key,
            model_name=model_name,
            model_version=model_version,
            detector_backend=detector_backend,
            max_file_size_mb=max_file_size_mb,
            min_image_width=min_width,
            min_image_height=min_height,
            default_similarity_threshold=threshold,
            is_calibrated=is_calibrated,
            default_strategy=default_strategy,
        )

settings = Settings.load()
