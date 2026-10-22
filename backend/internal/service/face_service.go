package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"labsync-backend/internal/client"
	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
	"labsync-backend/pkg/crypto"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

type FaceService struct {
	userRepo      repository.UserRepository
	faceRepo      repository.FaceRepository
	auditRepo     repository.AuditRepository
	faceClient    client.FaceClient
	encryptionKey []byte
}

func NewFaceService(
	userRepo repository.UserRepository,
	faceRepo repository.FaceRepository,
	auditRepo repository.AuditRepository,
	faceClient client.FaceClient,
	encryptionKeyStr string,
) *FaceService {
	return &FaceService{
		userRepo:      userRepo,
		faceRepo:      faceRepo,
		auditRepo:     auditRepo,
		faceClient:    faceClient,
		encryptionKey: []byte(encryptionKeyStr[:32]),
	}
}

type EnrollInput struct {
	UserID        string
	RequestID     string
	FrontImg      []byte
	LeftImg       []byte
	RightImg      []byte
	FrontName     string
	LeftName      string
	RightName     string
	FrontMime     string
	LeftMime      string
	RightMime     string
	ConsentAgreed bool
	ReEnroll      bool
}

func (s *FaceService) Enroll(ctx context.Context, in EnrollInput) (*model.FaceEnrollResponse, error) {
	if s.userRepo != nil {
		u, err := s.userRepo.GetByID(ctx, in.UserID)
		if err != nil {
			return nil, model.NewUnauthorized("Pengguna tidak ditemukan")
		}
		if u.Status != "ACTIVE" {
			return nil, model.NewForbidden("Akun tidak aktif, tidak dapat melakukan pendaftaran wajah")
		}
	}

	if !in.ConsentAgreed {
		return nil, model.NewBadRequest("Persetujuan pemrosesan data biometrik (consent_agreed) wajib disetujui")
	}

	// 1. Validasi keberadaan ketiga file citra
	if len(in.FrontImg) == 0 || len(in.LeftImg) == 0 || len(in.RightImg) == 0 {
		return nil, model.NewBadRequest("Ketiga berkas citra wajah (front_image, left_image, right_image) wajib disertakan")
	}

	// 2. Validasi format MIME type
	if !isMimeAllowed(in.FrontMime) || !isMimeAllowed(in.LeftMime) || !isMimeAllowed(in.RightMime) {
		return nil, model.NewBadRequest("Format citra tidak didukung. Gunakan image/jpeg, image/png, atau image/webp")
	}

	// 3. Cek apakah user sudah memiliki template aktif
	existing, _, err := s.faceRepo.GetActiveTemplateByUserID(ctx, in.UserID)
	if err != nil && !errors.Is(err, repository.ErrFaceTemplateNotFound) {
		return nil, model.NewInternalError("Gagal memeriksa status template wajah", err)
	}

	if existing != nil && !in.ReEnroll {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "FACE_ALREADY_ENROLLED",
			Message:    "Pengguna sudah memiliki pendaftaran wajah aktif. Gunakan alur pendaftaran ulang (re_enroll=true) jika ingin memperbarui",
		}
	}

	// 4. Panggil Python Face Service untuk ekstraksi embedding masing-masing pose
	frontRes, err := s.faceClient.Extract(ctx, in.FrontImg, in.FrontName, in.FrontMime, "FRONT", in.RequestID)
	if err != nil {
		return nil, err
	}

	leftRes, err := s.faceClient.Extract(ctx, in.LeftImg, in.LeftName, in.LeftMime, "LEFT", in.RequestID)
	if err != nil {
		return nil, err
	}

	rightRes, err := s.faceClient.Extract(ctx, in.RightImg, in.RightName, in.RightMime, "RIGHT", in.RequestID)
	if err != nil {
		return nil, err
	}

	// 5. Enkripsi embedding dengan AES-256-GCM sebelum disimpan ke database
	frontCipher, err := crypto.EncryptEmbedding(s.encryptionKey, frontRes.Embedding)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengenkripsi embedding pose depan", err)
	}

	leftCipher, err := crypto.EncryptEmbedding(s.encryptionKey, leftRes.Embedding)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengenkripsi embedding pose kiri", err)
	}

	rightCipher, err := crypto.EncryptEmbedding(s.encryptionKey, rightRes.Embedding)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengenkripsi embedding pose kanan", err)
	}

	// Hapus embedding plaintext dari memori segera
	frontRes.Embedding = nil
	leftRes.Embedding = nil
	rightRes.Embedding = nil

	now := time.Now().UTC()
	template := &model.FaceTemplate{
		UserID:         in.UserID,
		ModelName:      frontRes.ModelName,
		ModelVersion:   frontRes.ModelVersion,
		ConsentVersion: "v1.0",
		ConsentedAt:    now,
		Status:         "ACTIVE",
		EnrolledAt:     now,
		UpdatedAt:      now,
	}

	samples := []model.FaceTemplateSample{
		{
			Pose:                "FRONT",
			EmbeddingCiphertext: frontCipher,
			QualityStatus:       "VALID",
		},
		{
			Pose:                "LEFT",
			EmbeddingCiphertext: leftCipher,
			QualityStatus:       "VALID",
		},
		{
			Pose:                "RIGHT",
			EmbeddingCiphertext: rightCipher,
			QualityStatus:       "VALID",
		},
	}

	// 6. Simpan dalam satu transaksi atomik di database
	if err := s.faceRepo.SaveEnrollment(ctx, template, samples, in.ReEnroll); err != nil {
		if errors.Is(err, repository.ErrFaceAlreadyEnrolled) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "FACE_ALREADY_ENROLLED",
				Message:    "Pengguna sudah memiliki pendaftaran wajah aktif",
			}
		}
		return nil, model.NewInternalError("Gagal menyimpan data biometrik terenkripsi", err)
	}

	// 7. Catat log audit tanpa membocorkan data biometrik
	action := "FACE_ENROLLED"
	if in.ReEnroll {
		action = "FACE_RE_ENROLLED"
	}
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &in.UserID,
		Action:     action,
		EntityType: "FACE_TEMPLATE",
		EntityID:   template.ID,
	})

	return &model.FaceEnrollResponse{
		Message:      "Pendaftaran wajah berhasil",
		TemplateID:   template.ID,
		ModelName:    template.ModelName,
		ModelVersion: template.ModelVersion,
		SampleCount:  len(samples),
		EnrolledAt:   template.EnrolledAt.Format(time.RFC3339),
	}, nil
}

func (s *FaceService) GetStatus(ctx context.Context, callerID, callerRole, targetUserID string) (*model.FaceStatusResponse, error) {
	if callerRole != "ADMIN" && callerID != targetUserID {
		return nil, model.NewForbidden("Akses ditolak: tidak dapat melihat status biometrik pengguna lain")
	}

	template, samples, err := s.faceRepo.GetActiveTemplateByUserID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, repository.ErrFaceTemplateNotFound) {
			return &model.FaceStatusResponse{
				Enrolled:    false,
				Status:      "NOT_ENROLLED",
				SampleCount: 0,
				Poses:       []string{},
			}, nil
		}
		return nil, model.NewInternalError("Gagal mengambil status pendaftaran wajah", err)
	}

	poses := make([]string, len(samples))
	for i, smp := range samples {
		poses[i] = smp.Pose
	}

	return &model.FaceStatusResponse{
		Enrolled:       true,
		Status:         template.Status,
		ModelName:      template.ModelName,
		ModelVersion:   template.ModelVersion,
		ConsentVersion: template.ConsentVersion,
		EnrolledAt:     template.EnrolledAt.Format(time.RFC3339),
		SampleCount:    len(samples),
		Poses:          poses,
	}, nil
}

func (s *FaceService) Revoke(ctx context.Context, adminID, targetUserID string) error {
	if err := s.faceRepo.RevokeTemplate(ctx, targetUserID); err != nil {
		if errors.Is(err, repository.ErrFaceTemplateNotFound) {
			return model.NewNotFound("Template wajah aktif tidak ditemukan")
		}
		return model.NewInternalError("Gagal mencabut pendaftaran wajah", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "FACE_TEMPLATE_REVOKED",
		EntityType: "FACE_TEMPLATE",
		EntityID:   targetUserID,
	})

	return nil
}

func isMimeAllowed(mime string) bool {
	clean := strings.ToLower(strings.TrimSpace(strings.Split(mime, ";")[0]))
	return allowedMimeTypes[clean]
}
