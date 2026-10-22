# LabSync Face Verification Service

Microservice internal berbasis Python + FastAPI + DeepFace untuk ekstraksi dan verifikasi biometrik 1:1 multi-angle (`FRONT`, `LEFT`, `RIGHT`).

> **PENTING:** Service ini adalah layanan internal privat. Frontend **TIDAK PERNAH** memanggil service ini secara langsung. Seluruh request diproksi dan diotentikasi oleh backend Go menggunakan header `X-Internal-API-Key`.

---

## 1. Spesifikasi Model & Keputusan Biometrik

| Komponen | Spesifikasi / Kandidat |
| :--- | :--- |
| **Framework** | DeepFace 0.0.93+ |
| **Kandidat Model** | **ArcFace** (dimensi output: 512 float) |
| **Detektor Wajah** | OpenCV (default ringan) / RetinaFace (opsional presisi tinggi) |
| **Metrik Jarak** | Cosine Distance ($d = 1 - \frac{u \cdot v}{\|u\|_2 \|v\|_2}$) |
| **Mekanisme** | Verifikasi 1:1 terhadap 3 embedding tersimpan milik user (bukan pencarian 1:N) |
| **Threshold Awal** | `0.40` (Kandidat default Cosine distance ArcFace) |

### Aturan Keputusan Multi-Angle & Kalibrasi
- **Strategi Pembandingan:** Menggunakan strategi `MIN_DISTANCE`, yaitu membandingkan vektor live frame terhadap ketiga sampel pendaftaran (`FRONT`, `LEFT`, `RIGHT`), kemudian mengambil jarak terkecil (pose yang paling bersesuaian dengan sudut hadap kamera).
- **Kebijakan Kalibrasi:** Secara default `IS_CALIBRATED=false`. Sebelum kalibrasi empiris (pengujian False Acceptance Rate / False Rejection Rate pada populasi mahasiswa kampus) selesai dilakukan, flag `verified` tidak diizinkan meloloskan presensi otomatis untuk mencegah salah verifikasi di lingkungan produksi.

---

## 2. Kebutuhan Sumber Daya & Kompatibilitas Runtime

- **Versi Python:** **Python 3.10 – 3.12** direkomendasikan.
  *(Catatan: TensorFlow/tf-keras saat ini belum menyediakan wheel biner untuk Python 3.14+).*
- **Kebutuhan CPU:** Minimal 2 vCPU dengan dukungan instruksi **AVX2** untuk inferensi optimal.
- **Kebutuhan Memori (RAM):** Minimal **1.5 GB – 2.0 GB** saat memuat model ArcFace ke dalam memori.
- **Kebutuhan Disk:** Sekitar **150 MB** untuk file bobot ArcFace (`arcface_weights.h5`) dan ~20 MB untuk detektor, diunduh otomatis pada eksekusi pertama ke folder `~/.deepface/weights/`.
- **Catatan Lisensi:** Pustaka DeepFace berlisensi MIT. Bobot pretrained ArcFace (InsightFace) umumnya dirilis untuk keperluan riset/non-komersial; wajib diverifikasi legalitas institusi sebelum rilis produksi.

---

## 3. Menjalankan Layanan Lokal

### A. Persiapan Virtual Environment
```bash
python3.11 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### B. Konfigurasi Environment
```bash
cp .env.example .env
```
Pastikan `INTERNAL_API_KEY` sama dengan konfigurasi `backend/.env`.

### C. Menjalankan Server
```bash
uvicorn main:app --host 0.0.0.0 --port 8001 --reload
```
Layanan akan berjalan di `http://localhost:8001`.

---

## 4. Menjalankan Pengujian

```bash
pytest -v
```

- **Unit Tests (15 tests):** Menguji validasi format citra, MIME type, resolusi minimum (160x160), batas ukuran file (5MB), deteksi tanpa wajah, deteksi multi-wajah, kegagalan model, dan kalkulasi jarak 1:1 multi-pose menggunakan engine mock.
- **Integration Test:** `test_integration.py` menguji pemuatan bobot model DeepFace nyata secara terisolasi dan akan otomatis di-skip (`SKIPPED`) jika pustaka DeepFace atau bobot model belum terunduh di mesin lokal.

---

## 5. Kontrak Endpoint Internal

### `GET /health`
Liveness probe publik internal.
- **Response `200 OK`:**
  ```json
  {
    "status": "healthy",
    "service": "face_service",
    "model_name": "ArcFace",
    "model_version": "1.0.0",
    "is_calibrated": false
  }
  ```

### `POST /api/v1/internal/face/extract`
Mengekstrak embedding vektor dari 1 pose enrollment (`FRONT`, `LEFT`, atau `RIGHT`).
- **Header:** `X-Internal-API-Key: <INTERNAL_API_KEY>`
- **Request Form:** `file` (binary citra), `pose` ("FRONT" | "LEFT" | "RIGHT")
- **Response `200 OK`:**
  ```json
  {
    "model_name": "ArcFace",
    "model_version": "1.0.0",
    "pose": "FRONT",
    "embedding": [0.0123, -0.0456, ...],
    "dim": 512,
    "face_confidence": 0.985
  }
  ```

### `POST /api/v1/internal/face/verify`
Membandingkan 1 live frame terhadap 3 vektor enrollment milik pengguna.
- **Header:** `X-Internal-API-Key: <INTERNAL_API_KEY>`
- **Request Form:**
  - `file`: live frame binary
  - `enrolled_samples`: JSON string berisi 3 sampel pose
  - `threshold`: optional float (default `0.40`)
  - `strategy`: optional string (`MIN_DISTANCE` | `AVERAGE_DISTANCE`)
- **Response `200 OK`:**
  ```json
  {
    "verified": false,
    "distance": 0.28,
    "threshold": 0.40,
    "strategy": "MIN_DISTANCE",
    "best_matched_pose": "FRONT",
    "distances_by_pose": {
      "FRONT": 0.28,
      "LEFT": 0.54,
      "RIGHT": 0.61
    },
    "is_calibrated": false,
    "decision_rule": "MIN_DISTANCE_UNVERIFIED_CALIBRATION_PENDING",
    "model_name": "ArcFace"
  }
  ```
