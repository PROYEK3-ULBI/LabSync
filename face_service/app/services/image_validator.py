import io
from typing import Tuple
import numpy as np
from PIL import Image

from app.config import settings
from app.exceptions import FaceServiceException, InvalidImageException, FileTooLargeException

ALLOWED_MIME_TYPES = {
    "image/jpeg",
    "image/jpg",
    "image/png",
    "image/webp",
}

def validate_and_load_image(
    content: bytes,
    content_type: str | None = None,
) -> Tuple[np.ndarray, int, int]:
    # 1. Validasi ukuran file
    max_bytes = settings.max_file_size_mb * 1024 * 1024
    if len(content) > max_bytes:
        raise FileTooLargeException(
            f"Ukuran berkas ({len(content)} bytes) melebihi batas maksimum {settings.max_file_size_mb} MB"
        )

    if len(content) == 0:
        raise InvalidImageException("Berkas citra kosong")

    # 2. Validasi MIME Type jika tersedia
    if content_type and content_type.lower() not in ALLOWED_MIME_TYPES:
        raise InvalidImageException(
            f"Tipe konten '{content_type}' tidak diizinkan. Gunakan image/jpeg, image/png, atau image/webp"
        )

    # 3. Validasi dekode citra di memori
    try:
        with Image.open(io.BytesIO(content)) as img:
            img.verify()
    except Exception as e:
        raise InvalidImageException(f"Berkas rusak atau bukan format citra valid: {str(e)}")

    # 4. Baca array RGB dan cek dimensi minimum
    try:
        with Image.open(io.BytesIO(content)) as img:
            rgb_img = img.convert("RGB")
            width, height = rgb_img.size
            if width < settings.min_image_width or height < settings.min_image_height:
                raise InvalidImageException(
                    f"Dimensi citra {width}x{height} terlalu kecil (minimal {settings.min_image_width}x{settings.min_image_height})"
                )
            arr = np.array(rgb_img)
            return arr, width, height
    except FaceServiceException:
        raise
    except Exception as e:
        raise InvalidImageException(f"Gagal memproses citra di memori: {str(e)}")
