class FaceServiceException(Exception):
    def __init__(self, status_code: int, code: str, message: str):
        self.status_code = status_code
        self.code = code
        self.message = message
        super().__init__(message)

class UnauthorizedInternalException(FaceServiceException):
    def __init__(self, message: str = "Akses ditolak: kunci API internal tidak valid"):
        super().__init__(401, "UNAUTHORIZED_INTERNAL", message)

class InvalidImageException(FaceServiceException):
    def __init__(self, message: str = "Berkas citra tidak valid atau rusak"):
        super().__init__(400, "INVALID_IMAGE", message)

class FileTooLargeException(FaceServiceException):
    def __init__(self, message: str = "Ukuran berkas melebihi batas yang diizinkan"):
        super().__init__(413, "FILE_TOO_LARGE", message)

class NoFaceDetectedException(FaceServiceException):
    def __init__(self, message: str = "Tidak ada wajah terdeteksi pada citra"):
        super().__init__(422, "NO_FACE_DETECTED", message)

class MultipleFacesDetectedException(FaceServiceException):
    def __init__(self, message: str = "Terdeteksi lebih dari satu wajah pada citra"):
        super().__init__(422, "MULTIPLE_FACES_DETECTED", message)

class ModelInferenceException(FaceServiceException):
    def __init__(self, message: str = "Terjadi kegagalan saat inferensi model biometrik"):
        super().__init__(500, "MODEL_ERROR", message)
