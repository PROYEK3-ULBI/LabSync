package handler

import (
	"github.com/gofiber/fiber/v2"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// OpenAPISpec melayani spesifikasi OpenAPI dalam format JSON
func (h *DocsHandler) OpenAPISpec(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "application/json; charset=utf-8")
	return c.SendString(openAPISpecJSON)
}

// SwaggerUI merender antarmuka interaktif dokumentasi API
func (h *DocsHandler) SwaggerUI(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(swaggerUIHTML)
}

const openAPISpecJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "LabSync API",
    "description": "Backend REST API Presensi dan Manajemen Aset Laboratorium Komputer",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "/api/v1",
      "description": "API v1 base URL"
    }
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health dan Readiness Probe",
        "tags": ["System"],
        "responses": {
          "200": {
            "description": "Layanan sehat",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/HealthResponse" }
              }
            }
          },
          "503": {
            "description": "Layanan degraded",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/HealthResponse" }
              }
            }
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login pengguna",
        "description": "Memvalidasi kredensial dan menerbitkan JWT via HttpOnly cookie",
        "tags": ["Auth"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/LoginRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login berhasil. Cookie token disetel.",
            "headers": {
              "Set-Cookie": {
                "schema": { "type": "string" },
                "description": "HttpOnly cookie JWT token"
              }
            },
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/LoginResponse" }
              }
            }
          },
          "401": {
            "description": "Kredensial tidak valid",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/ErrorResponse" }
              }
            }
          },
          "403": {
            "description": "Akun tidak aktif",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/ErrorResponse" }
              }
            }
          }
        }
      }
    },
    "/auth/logout": {
      "post": {
        "summary": "Logout pengguna",
        "description": "Menghapus cookie autentikasi dan mencatat log audit",
        "tags": ["Auth"],
        "security": [{ "cookieAuth": [] }],
        "responses": {
          "200": {
            "description": "Logout berhasil"
          },
          "401": {
            "description": "Tidak terautentikasi"
          }
        }
      }
    },
    "/auth/me": {
      "get": {
        "summary": "Profil pengguna saat ini",
        "description": "Mengambil informasi user login dari cookie JWT",
        "tags": ["Auth"],
        "security": [{ "cookieAuth": [] }],
        "responses": {
          "200": {
            "description": "Profil ditemukan",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserResponse" }
              }
            }
          },
          "401": {
            "description": "Sesi tidak valid atau telah berakhir"
          }
        }
      }
    },
    "/users": {
      "get": {
        "summary": "Daftar pengguna (Admin)",
        "tags": ["Users"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "search", "in": "query", "schema": { "type": "string" } },
          { "name": "role", "in": "query", "schema": { "type": "string" } },
          { "name": "status", "in": "query", "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Daftar pengguna",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserListResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya role ADMIN yang diizinkan"
          }
        }
      },
      "post": {
        "summary": "Buat pengguna baru (Admin)",
        "tags": ["Users"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateUserRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Pengguna berhasil dibuat",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserResponse" }
              }
            }
          },
          "409": {
            "description": "Identifier atau email duplikat"
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "summary": "Detail pengguna (Admin atau diri sendiri)",
        "tags": ["Users"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail pengguna",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserResponse" }
              }
            }
          },
          "403": {
            "description": "Tidak memiliki akses ke profil pengguna lain"
          },
          "404": {
            "description": "Pengguna tidak ditemukan"
          }
        }
      },
      "put": {
        "summary": "Perbarui pengguna (Admin atau diri sendiri)",
        "tags": ["Users"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateUserRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Pengguna berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserResponse" }
              }
            }
          },
          "403": {
            "description": "Tidak diizinkan mengubah pengguna lain atau field terbatas"
          },
          "409": {
            "description": "Identifier atau email duplikat"
          }
        }
      }
    },
    "/users/{id}/status": {
      "patch": {
        "summary": "Ubah status pengguna (Admin)",
        "tags": ["Users"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateUserStatusRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Status pengguna diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          }
        }
      }
    },
    "/sessions": {
      "get": {
        "summary": "Daftar sesi praktikum",
        "description": "Mengambil daftar sesi praktikum dengan filter dan paginasi (Tersedia untuk Admin dan Mahasiswa)",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "status", "in": "query", "schema": { "type": "string", "enum": ["SCHEDULED", "ACTIVE", "CLOSED", "CANCELLED"] } },
          { "name": "session_date", "in": "query", "schema": { "type": "string", "format": "date" } },
          { "name": "search", "in": "query", "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Daftar sesi praktikum",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/SessionListResponse" }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Buat sesi praktikum baru (Admin)",
        "description": "Membuat sesi praktikum baru dengan validasi rentang waktu dan lokasi",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateSessionRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Sesi praktikum berhasil dibuat",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/SessionResponse" }
              }
            }
          },
          "400": {
            "description": "Waktu selesai sebelum waktu mulai atau format tidak valid"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "409": {
            "description": "Kode sesi sudah digunakan"
          }
        }
      }
    },
    "/sessions/active": {
      "get": {
        "summary": "Daftar sesi praktikum yang sedang aktif",
        "description": "Mengembalikan semua sesi dengan status ACTIVE",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "responses": {
          "200": {
            "description": "Daftar sesi aktif",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/SessionResponse" }
                }
              }
            }
          }
        }
      }
    },
    "/sessions/{id}": {
      "get": {
        "summary": "Detail sesi praktikum",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail sesi",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/SessionResponse" }
              }
            }
          },
          "404": {
            "description": "Sesi tidak ditemukan"
          }
        }
      },
      "put": {
        "summary": "Perbarui sesi praktikum (Admin)",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateSessionRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Sesi berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/SessionResponse" }
              }
            }
          },
          "400": {
            "description": "Validasi waktu atau format gagal"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Sesi tidak ditemukan"
          },
          "409": {
            "description": "Kode sesi duplikat"
          }
        }
      },
      "delete": {
        "summary": "Hapus sesi praktikum (Admin)",
        "description": "Hanya sesi yang belum memiliki data presensi yang dapat dihapus secara fisik; sesi dengan presensi harus diubah statusnya menjadi CLOSED atau CANCELLED",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Sesi berhasil dihapus"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Sesi tidak ditemukan"
          },
          "409": {
            "description": "Sesi memiliki riwayat presensi dan tidak dapat dihapus fisik"
          }
        }
      }
    },
    "/sessions/{id}/status": {
      "patch": {
        "summary": "Ubah status sesi praktikum (Admin)",
        "description": "Mengubah status sesi (SCHEDULED, ACTIVE, CLOSED, CANCELLED)",
        "tags": ["Sessions"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateSessionStatusRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Status sesi berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/SessionResponse" }
              }
            }
          },
          "400": {
            "description": "Status tidak valid"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Sesi tidak ditemukan"
          }
        }
      }
    },
    "/face/enroll": {
      "post": {
        "summary": "Pendaftaran Wajah Multi-Angle (Mahasiswa)",
        "description": "Mengunggah 3 sampel pose wajah (FRONT, LEFT, RIGHT). Backend mengekstrak embedding via microservice Python, mengenkripsi ciphertext menggunakan AES-256-GCM, dan menyimpan template dalam satu transaksi. Raw image langsung dibuang dari memori.",
        "tags": ["Face Biometrics"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "required": ["front_image", "left_image", "right_image", "consent_agreed"],
                "properties": {
                  "front_image": { "type": "string", "format": "binary" },
                  "left_image": { "type": "string", "format": "binary" },
                  "right_image": { "type": "string", "format": "binary" },
                  "consent_agreed": { "type": "boolean", "example": true },
                  "re_enroll": { "type": "boolean", "example": false }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Pendaftaran wajah berhasil",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/FaceEnrollResponse" }
              }
            }
          },
          "400": {
            "description": "Berkas tidak lengkap, format tidak didukung, atau persetujuan biometrik belum disetujui"
          },
          "403": {
            "description": "Akun tidak aktif atau bukan peran mahasiswa"
          },
          "409": {
            "description": "Pengguna sudah memiliki template aktif (gunakan re_enroll=true untuk memperbarui)"
          },
          "422": {
            "description": "Tidak ada wajah atau terdeteksi lebih dari satu wajah"
          },
          "503": {
            "description": "Layanan biometrik internal tidak tersedia"
          },
          "504": {
            "description": "Waktu tunggu layanan biometrik terlampaui"
          }
        }
      }
    },
    "/face/status": {
      "get": {
        "summary": "Status Pendaftaran Wajah",
        "description": "Mengambil status enrollment dan daftar pose milik pengguna yang sedang login tanpa membocorkan embedding biometrik",
        "tags": ["Face Biometrics"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "user_id", "in": "query", "description": "ID pengguna (khusus Admin)", "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Status pendaftaran wajah",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/FaceStatusResponse" }
              }
            }
          },
          "401": {
            "description": "Tidak terautentikasi"
          }
        }
      }
    },
    "/face/template/{user_id}": {
      "delete": {
        "summary": "Cabut Template Wajah (Admin)",
        "description": "Mencabut status template wajah pengguna (ACTIVE -> REVOKED)",
        "tags": ["Face Biometrics"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "user_id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Template berhasil dicabut"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Template aktif tidak ditemukan"
          }
        }
      }
    },
    "/attendance/check-in": {
      "post": {
        "summary": "Presensi Masuk (Check-in) Berbasis Wajah (Mahasiswa)",
        "description": "Menerima session_id dan citra wajah live. Backend secara server-side memverifikasi wajah 1:1 terhadap 3 pose pendaftaran tersimpan milik mahasiswa yang sedang login. Check-in hanya tercatat jika verifikasi lolos.",
        "tags": ["Attendance"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "required": ["session_id", "image"],
                "properties": {
                  "session_id": { "type": "string", "example": "sess-uuid" },
                  "image": { "type": "string", "format": "binary" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Presensi masuk berhasil",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/CheckInResponse" }
              }
            }
          },
          "400": {
            "description": "Sesi tidak aktif, verifikasi biometrik gagal, atau pengguna belum terdaftar wajah"
          },
          "403": {
            "description": "Akun tidak aktif atau bukan mahasiswa"
          },
          "404": {
            "description": "Sesi tidak ditemukan"
          },
          "409": {
            "description": "Mahasiswa sudah melakukan presensi pada sesi ini"
          },
          "422": {
            "description": "Tidak ada wajah atau multi-wajah terdeteksi"
          },
          "503": {
            "description": "Layanan biometrik tidak dapat dihubungi"
          }
        }
      }
    },
    "/attendance/check-out": {
      "post": {
        "summary": "Presensi Keluar (Check-out) (Mahasiswa)",
        "description": "Mencatat waktu check-out server untuk sesi yang telah di-check-in sebelumnya",
        "tags": ["Attendance"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CheckOutRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Check-out berhasil",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/CheckOutResponse" }
              }
            }
          },
          "404": {
            "description": "Presensi masuk tidak ditemukan untuk sesi ini"
          },
          "409": {
            "description": "Sudah melakukan check-out sebelumnya"
          }
        }
      }
    },
    "/attendance/history": {
      "get": {
        "summary": "Riwayat Presensi Pribadi (Mahasiswa)",
        "tags": ["Attendance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } }
        ],
        "responses": {
          "200": {
            "description": "Daftar riwayat presensi",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AttendanceListResponse" }
              }
            }
          }
        }
      }
    },
    "/attendance": {
      "get": {
        "summary": "Monitoring Presensi Semua Mahasiswa (Admin)",
        "tags": ["Attendance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "session_id", "in": "query", "schema": { "type": "string" } },
          { "name": "user_id", "in": "query", "schema": { "type": "string" } },
          { "name": "date", "in": "query", "schema": { "type": "string", "format": "date" } }
        ],
        "responses": {
          "200": {
            "description": "Daftar presensi laboratorium",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AttendanceListResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          }
        }
      }
    },
    "/attendance/{id}": {
      "get": {
        "summary": "Detail Presensi",
        "tags": ["Attendance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail data presensi",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AttendanceResponse" }
              }
            }
          },
          "403": {
            "description": "Tidak diizinkan melihat presensi mahasiswa lain"
          },
          "404": {
            "description": "Presensi tidak ditemukan"
          }
        }
      }
    },
    "/assets": {
      "get": {
        "summary": "Daftar Aset Laboratorium (Katalog)",
        "description": "Mengambil daftar master aset komputer dengan filter status, kategori, kata kunci, dan paginasi (Dapat diakses Admin dan Mahasiswa)",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "status", "in": "query", "schema": { "type": "string", "enum": ["AVAILABLE", "ALLOCATED", "MAINTENANCE", "RETIRED"] } },
          { "name": "category", "in": "query", "schema": { "type": "string" } },
          { "name": "search", "in": "query", "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Daftar aset laboratorium",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetListResponse" }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Tambah Master Aset Baru (Admin)",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateAssetRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Aset berhasil ditambahkan",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetResponse" }
              }
            }
          },
          "400": {
            "description": "Validasi kode/nama aset gagal atau status awal tidak valid"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "409": {
            "description": "Kode aset atau nomor seri sudah terdaftar"
          }
        }
      }
    },
    "/assets/{id}": {
      "get": {
        "summary": "Detail Master Aset",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail data aset",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetResponse" }
              }
            }
          },
          "404": {
            "description": "Aset tidak ditemukan"
          }
        }
      },
      "put": {
        "summary": "Perbarui Data Master Aset (Admin)",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateAssetRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Data aset berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetResponse" }
              }
            }
          },
          "400": {
            "description": "Transisi status tidak valid"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Aset tidak ditemukan"
          },
          "409": {
            "description": "Konflik status: aset sedang dipinjam atau dalam perbaikan aktif"
          }
        }
      },
      "delete": {
        "summary": "Hapus Master Aset (Admin)",
        "description": "Menghapus master aset. Ditolak jika aset sedang dalam peminjaman aktif atau pemeliharaan terbuka.",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Aset berhasil dihapus"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Aset tidak ditemukan"
          },
          "409": {
            "description": "Aset sedang dipinjam atau dalam perbaikan dan tidak dapat dihapus"
          }
        }
      }
    },
    "/assets/{id}/status": {
      "patch": {
        "summary": "Perbarui Status Siklus Aset (Admin)",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateAssetStatusRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Status aset berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetResponse" }
              }
            }
          },
          "400": {
            "description": "Transisi status tidak diizinkan"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Aset tidak ditemukan"
          },
          "409": {
            "description": "Aset sedang dipinjam atau dalam perbaikan aktif"
          }
        }
      }
    },
    "/assets/{id}/history": {
      "get": {
        "summary": "Histori dan Timeline Siklus Aset",
        "description": "Mengambil urutan peristiwa peminjaman, pengembalian, kondisi, dan status aset",
        "tags": ["Assets"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Histori aset",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AssetHistoryResponse" }
              }
            }
          },
          "404": {
            "description": "Aset tidak ditemukan"
          }
        }
      }
    },
    "/assets/my-allocations": {
      "get": {
        "summary": "Daftar Aset yang Sedang Dipinjam Mahasiswa",
        "tags": ["Assets", "Allocations"],
        "security": [{ "cookieAuth": [] }],
        "responses": {
          "200": {
            "description": "Daftar alokasi milik mahasiswa login",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationListResponse" }
              }
            }
          }
        }
      }
    },
    "/allocations": {
      "get": {
        "summary": "Daftar Semua Alokasi Aset (Admin)",
        "tags": ["Allocations"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "asset_id", "in": "query", "schema": { "type": "string" } },
          { "name": "user_id", "in": "query", "schema": { "type": "string" } },
          { "name": "session_id", "in": "query", "schema": { "type": "string" } },
          { "name": "status", "in": "query", "schema": { "type": "string", "enum": ["ACTIVE", "RETURNED"] } }
        ],
        "responses": {
          "200": {
            "description": "Daftar alokasi aset",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationListResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          }
        }
      },
      "post": {
        "summary": "Alokasikan Aset ke Mahasiswa (Admin)",
        "description": "Mengalokasikan unit aset kepada mahasiswa pada sesi aktif. Memvalidasi ketersediaan aset dan prasyarat kehadiran mahasiswa (check-in aktif). Dijalankan dalam satu transaksi database.",
        "tags": ["Allocations"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateAllocationRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Alokasi aset berhasil dibuat",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationResponse" }
              }
            }
          },
          "400": {
            "description": "Mahasiswa belum presensi (ATTENDANCE_REQUIRED) atau sesi tidak aktif"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Aset, pengguna, atau sesi tidak ditemukan"
          },
          "409": {
            "description": "Aset sedang dialokasikan atau tidak berstatus AVAILABLE"
          }
        }
      }
    },
    "/allocations/my": {
      "get": {
        "summary": "Daftar Alokasi Pribadi (Mahasiswa)",
        "tags": ["Allocations"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } }
        ],
        "responses": {
          "200": {
            "description": "Daftar alokasi aset mahasiswa",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationListResponse" }
              }
            }
          }
        }
      }
    },
    "/allocations/{id}": {
      "get": {
        "summary": "Detail Alokasi Aset",
        "tags": ["Allocations"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail alokasi",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationResponse" }
              }
            }
          },
          "403": {
            "description": "Tidak memiliki akses ke alokasi mahasiswa lain"
          },
          "404": {
            "description": "Alokasi tidak ditemukan"
          }
        }
      }
    },
    "/allocations/{id}/return": {
      "post": {
        "summary": "Kembalikan Aset (Admin)",
        "description": "Mencatat pengembalian aset. Jika kondisi GOOD maka aset menjadi AVAILABLE; jika DAMAGED maka aset otomatis beralih ke MAINTENANCE.",
        "tags": ["Allocations"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/ReturnAllocationRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Aset berhasil dikembalikan",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AllocationResponse" }
              }
            }
          },
          "400": {
            "description": "Kondisi pengembalian tidak valid"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Alokasi tidak ditemukan"
          },
          "409": {
            "description": "Aset sudah dikembalikan sebelumnya (ALREADY_RETURNED)"
          }
        }
      }
    },
    "/maintenance": {
      "get": {
        "summary": "Daftar Tiket Pemeliharaan Aset (Admin)",
        "tags": ["Maintenance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "asset_id", "in": "query", "schema": { "type": "string" } },
          { "name": "status", "in": "query", "schema": { "type": "string", "enum": ["OPEN", "IN_PROGRESS", "COMPLETED", "CANCELLED"] } },
          { "name": "priority", "in": "query", "schema": { "type": "string", "enum": ["LOW", "MEDIUM", "HIGH"] } }
        ],
        "responses": {
          "200": {
            "description": "Daftar tiket pemeliharaan",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/MaintenanceListResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          }
        }
      },
      "post": {
        "summary": "Buat Tiket Pemeliharaan Aset (Admin)",
        "description": "Mendaftarkan tiket perbaikan atau pemeliharaan aset. Status aset otomatis beralih ke MAINTENANCE dalam transaksi. Aset yang sedang dialokasikan ditolak.",
        "tags": ["Maintenance"],
        "security": [{ "cookieAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateMaintenanceRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Tiket pemeliharaan berhasil dibuat",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/MaintenanceResponse" }
              }
            }
          },
          "400": {
            "description": "Validasi gagal atau aset RETIRED"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Aset tidak ditemukan"
          },
          "409": {
            "description": "Aset sedang dipinjam (ASSET_ALLOCATED)"
          }
        }
      }
    },
    "/maintenance/{id}": {
      "get": {
        "summary": "Detail Tiket Pemeliharaan (Admin)",
        "tags": ["Maintenance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Detail tiket pemeliharaan",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/MaintenanceResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Tiket tidak ditemukan"
          }
        }
      },
      "put": {
        "summary": "Perbarui Status/Resolusi Pemeliharaan (Admin)",
        "description": "Memperbarui status sesuai state machine (OPEN -> IN_PROGRESS -> COMPLETED/CANCELLED). Saat COMPLETED, resolusi wajib dicatat dan status aset kembali ke AVAILABLE atau RETIRED.",
        "tags": ["Maintenance"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateMaintenanceRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Tiket pemeliharaan berhasil diperbarui",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/MaintenanceResponse" }
              }
            }
          },
          "400": {
            "description": "Transisi status ilegal atau resolusi belum diisi"
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          },
          "404": {
            "description": "Tiket tidak ditemukan"
          }
        }
      }
    },
    "/audit-logs": {
      "get": {
        "summary": "Audit Trail Append-Only (Admin)",
        "description": "Menginspeksi log audit keamanan dan operasional sistem secara append-only (tidak ada endpoint edit/delete)",
        "tags": ["Audit Logs"],
        "security": [{ "cookieAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "action", "in": "query", "schema": { "type": "string" } },
          { "name": "entity_type", "in": "query", "schema": { "type": "string" } },
          { "name": "user_id", "in": "query", "schema": { "type": "string" } }
        ],
        "responses": {
          "200": {
            "description": "Daftar log audit",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AuditLogListResponse" }
              }
            }
          },
          "403": {
            "description": "Hanya admin yang diizinkan"
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "cookieAuth": {
        "type": "apiKey",
        "in": "cookie",
        "name": "token"
      }
    },
    "schemas": {
      "LoginRequest": {
        "type": "object",
        "required": ["identifier", "password"],
        "properties": {
          "identifier": { "type": "string", "example": "TI240001" },
          "password": { "type": "string", "example": "Rahasia123!" }
        }
      },
      "LoginResponse": {
        "type": "object",
        "properties": {
          "message": { "type": "string", "example": "Login berhasil" },
          "user": { "$ref": "#/components/schemas/UserResponse" }
        }
      },
      "CreateUserRequest": {
        "type": "object",
        "required": ["identifier", "email", "full_name", "password", "role"],
        "properties": {
          "identifier": { "type": "string", "example": "TI240002" },
          "email": { "type": "string", "example": "mhs@univ.ac.id" },
          "full_name": { "type": "string", "example": "Budi Santoso" },
          "password": { "type": "string", "example": "Password123!" },
          "role": { "type": "string", "enum": ["ADMIN", "STUDENT"], "example": "STUDENT" },
          "status": { "type": "string", "enum": ["ACTIVE", "INACTIVE"], "default": "ACTIVE" }
        }
      },
      "UpdateUserRequest": {
        "type": "object",
        "properties": {
          "identifier": { "type": "string" },
          "email": { "type": "string" },
          "full_name": { "type": "string" },
          "password": { "type": "string" },
          "role": { "type": "string", "enum": ["ADMIN", "STUDENT"] },
          "status": { "type": "string", "enum": ["ACTIVE", "INACTIVE"] }
        }
      },
      "UpdateUserStatusRequest": {
        "type": "object",
        "required": ["status"],
        "properties": {
          "status": { "type": "string", "enum": ["ACTIVE", "INACTIVE"], "example": "INACTIVE" }
        }
      },
      "UserResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "identifier": { "type": "string" },
          "email": { "type": "string" },
          "full_name": { "type": "string" },
          "role": { "type": "string" },
          "status": { "type": "string" },
          "face_enrolled": { "type": "boolean" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "UserListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/UserResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "CreateSessionRequest": {
        "type": "object",
        "required": ["code", "title", "location", "session_date", "start_at", "end_at"],
        "properties": {
          "code": { "type": "string", "example": "PWB-07" },
          "title": { "type": "string", "example": "Praktikum Pemrograman Web" },
          "course_code": { "type": "string", "example": "IF2101" },
          "location": { "type": "string", "example": "Lab Komputer 1" },
          "session_date": { "type": "string", "format": "date", "example": "2026-09-18" },
          "start_at": { "type": "string", "format": "date-time", "example": "2026-09-18T13:00:00Z" },
          "end_at": { "type": "string", "format": "date-time", "example": "2026-09-18T15:30:00Z" },
          "status": { "type": "string", "enum": ["SCHEDULED", "ACTIVE", "CLOSED", "CANCELLED"], "default": "SCHEDULED" }
        }
      },
      "UpdateSessionRequest": {
        "type": "object",
        "properties": {
          "code": { "type": "string" },
          "title": { "type": "string" },
          "course_code": { "type": "string" },
          "location": { "type": "string" },
          "session_date": { "type": "string", "format": "date" },
          "start_at": { "type": "string", "format": "date-time" },
          "end_at": { "type": "string", "format": "date-time" },
          "status": { "type": "string", "enum": ["SCHEDULED", "ACTIVE", "CLOSED", "CANCELLED"] }
        }
      },
      "UpdateSessionStatusRequest": {
        "type": "object",
        "required": ["status"],
        "properties": {
          "status": { "type": "string", "enum": ["SCHEDULED", "ACTIVE", "CLOSED", "CANCELLED"], "example": "ACTIVE" }
        }
      },
      "SessionResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "code": { "type": "string" },
          "title": { "type": "string" },
          "course_code": { "type": "string" },
          "location": { "type": "string" },
          "session_date": { "type": "string", "format": "date" },
          "start_at": { "type": "string", "format": "date-time" },
          "end_at": { "type": "string", "format": "date-time" },
          "status": { "type": "string" },
          "created_by": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "SessionListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/SessionResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "FaceEnrollResponse": {
        "type": "object",
        "properties": {
          "message": { "type": "string", "example": "Pendaftaran wajah berhasil" },
          "template_id": { "type": "string" },
          "model_name": { "type": "string", "example": "ArcFace" },
          "model_version": { "type": "string", "example": "1.0.0" },
          "sample_count": { "type": "integer", "example": 3 },
          "enrolled_at": { "type": "string", "format": "date-time" }
        }
      },
      "FaceStatusResponse": {
        "type": "object",
        "properties": {
          "enrolled": { "type": "boolean" },
          "status": { "type": "string", "example": "ACTIVE" },
          "model_name": { "type": "string", "example": "ArcFace" },
          "model_version": { "type": "string", "example": "1.0.0" },
          "consent_version": { "type": "string", "example": "v1.0" },
          "enrolled_at": { "type": "string", "format": "date-time" },
          "sample_count": { "type": "integer", "example": 3 },
          "poses": {
            "type": "array",
            "items": { "type": "string" },
            "example": ["FRONT", "LEFT", "RIGHT"]
          }
        }
      },
      "CheckInResponse": {
        "type": "object",
        "properties": {
          "message": { "type": "string", "example": "Presensi masuk (check-in) berhasil" },
          "attendance": { "$ref": "#/components/schemas/AttendanceResponse" }
        }
      },
      "CheckOutRequest": {
        "type": "object",
        "required": ["session_id"],
        "properties": {
          "session_id": { "type": "string", "example": "sess-uuid" }
        }
      },
      "CheckOutResponse": {
        "type": "object",
        "properties": {
          "message": { "type": "string", "example": "Presensi keluar (check-out) berhasil" },
          "attendance": { "$ref": "#/components/schemas/AttendanceResponse" }
        }
      },
      "AttendanceResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "user_id": { "type": "string" },
          "session_id": { "type": "string" },
          "checked_in_at": { "type": "string", "format": "date-time" },
          "checked_out_at": { "type": "string", "format": "date-time" },
          "verification_status": { "type": "string", "example": "VERIFIED_FACE" },
          "user_identifier": { "type": "string", "example": "TI240001" },
          "user_full_name": { "type": "string", "example": "Budi Santoso" },
          "session_code": { "type": "string", "example": "PWB-07" },
          "session_title": { "type": "string", "example": "Praktikum Pemrograman Web" },
          "session_date": { "type": "string", "format": "date", "example": "2026-09-18" }
        }
      },
      "AttendanceListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/AttendanceResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "CreateAssetRequest": {
        "type": "object",
        "required": ["asset_code", "asset_name"],
        "properties": {
          "asset_code": { "type": "string", "example": "LAB-LT-007" },
          "asset_name": { "type": "string", "example": "Dell Latitude 5490 Core i7" },
          "category": { "type": "string", "example": "Laptop" },
          "brand": { "type": "string", "example": "Dell" },
          "model": { "type": "string", "example": "Latitude 5490" },
          "serial_number": { "type": "string", "example": "SN-DELL-5490-007" },
          "location": { "type": "string", "example": "Lab Komputer 1" },
          "status": { "type": "string", "enum": ["AVAILABLE", "MAINTENANCE"], "default": "AVAILABLE" },
          "notes": { "type": "string" }
        }
      },
      "UpdateAssetRequest": {
        "type": "object",
        "properties": {
          "asset_code": { "type": "string" },
          "asset_name": { "type": "string" },
          "category": { "type": "string" },
          "brand": { "type": "string" },
          "model": { "type": "string" },
          "serial_number": { "type": "string" },
          "location": { "type": "string" },
          "status": { "type": "string", "enum": ["AVAILABLE", "MAINTENANCE", "RETIRED"] },
          "notes": { "type": "string" }
        }
      },
      "UpdateAssetStatusRequest": {
        "type": "object",
        "required": ["status"],
        "properties": {
          "status": { "type": "string", "enum": ["AVAILABLE", "MAINTENANCE", "RETIRED"], "example": "MAINTENANCE" }
        }
      },
      "AssetResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "asset_code": { "type": "string", "example": "LAB-LT-007" },
          "asset_name": { "type": "string", "example": "Dell Latitude 5490 Core i7" },
          "category": { "type": "string", "example": "Laptop" },
          "brand": { "type": "string", "example": "Dell" },
          "model": { "type": "string", "example": "Latitude 5490" },
          "serial_number": { "type": "string", "example": "SN-DELL-5490-007" },
          "location": { "type": "string", "example": "Lab Komputer 1" },
          "status": { "type": "string", "example": "AVAILABLE" },
          "notes": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "AssetListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/AssetResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "CreateAllocationRequest": {
        "type": "object",
        "required": ["asset_id", "user_id", "session_id"],
        "properties": {
          "asset_id": { "type": "string", "example": "ast-uuid" },
          "user_id": { "type": "string", "example": "user-uuid" },
          "session_id": { "type": "string", "example": "sess-uuid" },
          "condition_before": { "type": "string", "example": "GOOD" },
          "notes": { "type": "string" }
        }
      },
      "ReturnAllocationRequest": {
        "type": "object",
        "required": ["condition"],
        "properties": {
          "condition": { "type": "string", "enum": ["GOOD", "DAMAGED"], "example": "GOOD" },
          "notes": { "type": "string" }
        }
      },
      "AllocationResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "asset_id": { "type": "string" },
          "user_id": { "type": "string" },
          "session_id": { "type": "string" },
          "allocated_at": { "type": "string", "format": "date-time" },
          "returned_at": { "type": "string", "format": "date-time" },
          "condition_before": { "type": "string" },
          "condition_after": { "type": "string" },
          "notes": { "type": "string" },
          "asset_code": { "type": "string" },
          "asset_name": { "type": "string" },
          "user_identifier": { "type": "string" },
          "user_full_name": { "type": "string" },
          "session_code": { "type": "string" },
          "session_title": { "type": "string" }
        }
      },
      "AllocationListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/AllocationResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "AssetTimelineEvent": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "event_type": { "type": "string", "example": "ALLOCATION" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "timestamp": { "type": "string", "format": "date-time" },
          "actor_id": { "type": "string" },
          "actor_name": { "type": "string" }
        }
      },
      "AssetHistoryResponse": {
        "type": "object",
        "properties": {
          "asset": { "$ref": "#/components/schemas/AssetResponse" },
          "timeline": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/AssetTimelineEvent" }
          }
        }
      },
      "CreateMaintenanceRequest": {
        "type": "object",
        "required": ["asset_id", "description"],
        "properties": {
          "asset_id": { "type": "string", "example": "ast-uuid" },
          "issue_type": { "type": "string", "example": "HARDWARE" },
          "description": { "type": "string", "example": "Layar berkedip saat digerakkan" },
          "priority": { "type": "string", "enum": ["LOW", "MEDIUM", "HIGH"], "default": "MEDIUM" }
        }
      },
      "UpdateMaintenanceRequest": {
        "type": "object",
        "properties": {
          "status": { "type": "string", "enum": ["OPEN", "IN_PROGRESS", "COMPLETED", "CANCELLED"] },
          "priority": { "type": "string", "enum": ["LOW", "MEDIUM", "HIGH"] },
          "description": { "type": "string" },
          "resolution": { "type": "string" },
          "cost": { "type": "number" },
          "result_status": { "type": "string", "enum": ["AVAILABLE", "RETIRED"] }
        }
      },
      "MaintenanceResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "asset_id": { "type": "string" },
          "reported_by": { "type": "string" },
          "issue_type": { "type": "string" },
          "description": { "type": "string" },
          "priority": { "type": "string" },
          "status": { "type": "string" },
          "started_at": { "type": "string", "format": "date-time" },
          "completed_at": { "type": "string", "format": "date-time" },
          "resolution": { "type": "string" },
          "cost": { "type": "number" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" },
          "asset_code": { "type": "string" },
          "asset_name": { "type": "string" },
          "reporter_identifier": { "type": "string" },
          "reporter_name": { "type": "string" }
        }
      },
      "MaintenanceListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/MaintenanceResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "AuditLogResponse": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "user_id": { "type": "string" },
          "action": { "type": "string" },
          "entity_type": { "type": "string" },
          "entity_id": { "type": "string" },
          "status": { "type": "string" },
          "ip_address": { "type": "string" },
          "user_agent": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" },
          "user_identifier": { "type": "string" },
          "user_full_name": { "type": "string" }
        }
      },
      "AuditLogListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/AuditLogResponse" }
          },
          "meta": {
            "type": "object",
            "properties": {
              "page": { "type": "integer" },
              "limit": { "type": "integer" },
              "total": { "type": "integer" }
            }
          }
        }
      },
      "HealthResponse": {
        "type": "object",
        "properties": {
          "status": { "type": "string", "example": "healthy" },
          "services": {
            "type": "object",
            "properties": {
              "database": { "type": "string", "example": "connected" }
            }
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "properties": {
              "code": { "type": "string", "example": "UNAUTHORIZED" },
              "message": { "type": "string", "example": "Identifier atau password salah" },
              "request_id": { "type": "string", "example": "req-12345" }
            }
          }
        }
      }
    }
  }
}`

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>LabSync API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin="anonymous"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/openapi.json',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`
