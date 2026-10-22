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

type AttendanceService struct {
	attendanceRepo repository.AttendanceRepository
	sessionRepo    repository.SessionRepository
	faceRepo       repository.FaceRepository
	userRepo       repository.UserRepository
	auditRepo      repository.AuditRepository
	faceClient     client.FaceClient
	encryptionKey  []byte
}

func NewAttendanceService(
	attendanceRepo repository.AttendanceRepository,
	sessionRepo repository.SessionRepository,
	faceRepo repository.FaceRepository,
	userRepo repository.UserRepository,
	auditRepo repository.AuditRepository,
	faceClient client.FaceClient,
	encryptionKeyStr string,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		sessionRepo:    sessionRepo,
		faceRepo:       faceRepo,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		faceClient:     faceClient,
		encryptionKey:  []byte(encryptionKeyStr[:32]),
	}
}

type CheckInInput struct {
	UserID    string
	SessionID string
	LiveImg   []byte
	Filename  string
	MimeType  string
	RequestID string
}

func (s *AttendanceService) CheckIn(ctx context.Context, in CheckInInput) (*model.AttendanceResponse, error) {
	// 1. Validasi akun user
	if s.userRepo != nil {
		u, err := s.userRepo.GetByID(ctx, in.UserID)
		if err != nil {
			return nil, model.NewUnauthorized("Pengguna tidak ditemukan")
		}
		if u.Status != "ACTIVE" {
			return nil, model.NewForbidden("Akun tidak aktif, tidak dapat melakukan presensi")
		}
	}

	// 2. Validasi sesi praktikum
	if strings.TrimSpace(in.SessionID) == "" {
		return nil, model.NewBadRequest("Session ID wajib disertakan")
	}

	sess, err := s.sessionRepo.GetByID(ctx, in.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, model.NewNotFound("Sesi praktikum tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal memeriksa sesi praktikum", err)
	}

	if sess.Status != model.SessionStatusActive {
		return nil, &model.AppError{
			StatusCode: 400,
			Code:       "SESSION_NOT_ACTIVE",
			Message:    "Sesi praktikum belum aktif atau sudah ditutup",
		}
	}

	// 3. Pre-check apakah user sudah presensi pada sesi ini
	existing, err := s.attendanceRepo.GetByUserAndSession(ctx, in.UserID, in.SessionID)
	if err != nil && !errors.Is(err, repository.ErrAttendanceNotFound) {
		return nil, model.NewInternalError("Gagal memeriksa riwayat presensi sesi", err)
	}
	if existing != nil {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ALREADY_CHECKED_IN",
			Message:    "Mahasiswa sudah melakukan presensi pada sesi ini",
		}
	}

	// 4. Validasi payload citra
	if len(in.LiveImg) == 0 {
		return nil, model.NewBadRequest("Berkas citra wajah live wajib diunggah")
	}
	if !isMimeAllowed(in.MimeType) {
		return nil, model.NewBadRequest("Format citra tidak didukung. Gunakan image/jpeg, image/png, atau image/webp")
	}

	// 5. Ambil sampel pendaftaran biometrik aktif milik pengguna (3 pose: FRONT, LEFT, RIGHT)
	_, samples, err := s.faceRepo.GetActiveTemplateByUserID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrFaceTemplateNotFound) {
			return nil, &model.AppError{
				StatusCode: 400,
				Code:       "FACE_NOT_ENROLLED",
				Message:    "Pengguna belum menyelesaikan pendaftaran wajah",
			}
		}
		return nil, model.NewInternalError("Gagal mengambil template wajah pengguna", err)
	}

	if len(samples) != 3 {
		return nil, &model.AppError{
			StatusCode: 400,
			Code:       "FACE_ENROLLMENT_INCOMPLETE",
			Message:    "Pendaftaran wajah pengguna belum lengkap (diperlukan 3 pose: FRONT, LEFT, RIGHT)",
		}
	}

	// 6. Dekripsi 3 embedding di memori
	enrolledPayloads := make([]model.EnrolledSamplePayload, len(samples))
	for i, smp := range samples {
		decrypted, decErr := crypto.DecryptEmbedding(s.encryptionKey, smp.EmbeddingCiphertext)
		if decErr != nil {
			return nil, model.NewInternalError("Gagal mendekripsi data biometrik pengguna", decErr)
		}
		enrolledPayloads[i] = model.EnrolledSamplePayload{
			Pose:      smp.Pose,
			Embedding: decrypted,
		}
	}

	// 7. Panggil microservice Python secara internal untuk verifikasi 1:1
	verifyRes, err := s.faceClient.Verify(
		ctx,
		in.LiveImg,
		in.Filename,
		in.MimeType,
		enrolledPayloads,
		nil, // Memakai threshold terkalibrasi dari face service
		"MIN_DISTANCE",
		in.RequestID,
	)

	// Pembersihan memori segera
	for i := range enrolledPayloads {
		enrolledPayloads[i].Embedding = nil
	}
	enrolledPayloads = nil

	if err != nil {
		return nil, err
	}

	// 8. Evaluasi keputusan verifikasi
	if !verifyRes.Verified {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			UserID:     &in.UserID,
			Action:     "ATTENDANCE_VERIFICATION_FAILED",
			EntityType: "ATTENDANCE",
			EntityID:   in.SessionID,
		})
		return nil, &model.AppError{
			StatusCode: 400,
			Code:       "VERIFICATION_FAILED",
			Message:    "Verifikasi wajah gagal. Wajah tidak cocok dengan data pendaftaran",
		}
	}

	// 9. Simpan presensi dengan timestamp server-authoritative
	now := time.Now().UTC()
	att := &model.Attendance{
		UserID:             in.UserID,
		SessionID:          in.SessionID,
		CheckedInAt:        now,
		VerificationStatus: "VERIFIED_FACE",
	}

	if err := s.attendanceRepo.Create(ctx, att); err != nil {
		if errors.Is(err, repository.ErrDuplicateAttendance) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "ALREADY_CHECKED_IN",
				Message:    "Mahasiswa sudah melakukan presensi pada sesi ini",
			}
		}
		return nil, model.NewInternalError("Gagal menyimpan presensi ke database", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &in.UserID,
		Action:     "ATTENDANCE_CHECK_IN",
		EntityType: "ATTENDANCE",
		EntityID:   att.ID,
	})

	// Lengkapi data relasi
	att.SessionCode = sess.Code
	att.SessionTitle = sess.Title
	att.SessionDate = sess.SessionDate

	res := toAttendanceResponse(att)
	return &res, nil
}

func (s *AttendanceService) CheckOut(ctx context.Context, userID, sessionID string) (*model.AttendanceResponse, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, model.NewBadRequest("Session ID wajib diisi")
	}

	// 1. Ambil data presensi yang ada
	att, err := s.attendanceRepo.GetByUserAndSession(ctx, userID, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrAttendanceNotFound) {
			return nil, &model.AppError{
				StatusCode: 404,
				Code:       "ATTENDANCE_NOT_FOUND",
				Message:    "Data presensi masuk belum ditemukan untuk sesi ini",
			}
		}
		return nil, model.NewInternalError("Gagal mengambil data presensi", err)
	}

	// 2. Cegah double check-out
	if att.CheckedOutAt != nil {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ALREADY_CHECKED_OUT",
			Message:    "Mahasiswa sudah melakukan check-out sebelumnya",
		}
	}

	// 3. Update stempel waktu check-out dari server
	now := time.Now().UTC()
	if err := s.attendanceRepo.CheckOut(ctx, att.ID, now); err != nil {
		if errors.Is(err, repository.ErrAlreadyCheckedOut) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "ALREADY_CHECKED_OUT",
				Message:    "Mahasiswa sudah melakukan check-out sebelumnya",
			}
		}
		return nil, model.NewInternalError("Gagal menyimpan waktu check-out", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "ATTENDANCE_CHECK_OUT",
		EntityType: "ATTENDANCE",
		EntityID:   att.ID,
	})

	att.CheckedOutAt = &now
	res := toAttendanceResponse(att)
	return &res, nil
}

func (s *AttendanceService) GetHistory(ctx context.Context, studentID string, page, limit int) (*model.AttendanceListResponse, error) {
	q := model.AttendanceListQuery{
		Page:   page,
		Limit:  limit,
		UserID: studentID,
	}

	list, total, err := s.attendanceRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil riwayat presensi mahasiswa", err)
	}

	data := make([]model.AttendanceResponse, len(list))
	for i, item := range list {
		data[i] = toAttendanceResponse(&item)
	}

	return &model.AttendanceListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *AttendanceService) List(ctx context.Context, q model.AttendanceListQuery) (*model.AttendanceListResponse, error) {
	list, total, err := s.attendanceRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar presensi", err)
	}

	data := make([]model.AttendanceResponse, len(list))
	for i, item := range list {
		data[i] = toAttendanceResponse(&item)
	}

	return &model.AttendanceListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *AttendanceService) GetByID(ctx context.Context, callerID, callerRole, id string) (*model.AttendanceResponse, error) {
	att, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAttendanceNotFound) {
			return nil, model.NewNotFound("Data presensi tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil detail presensi", err)
	}

	if callerRole != "ADMIN" && att.UserID != callerID {
		return nil, model.NewForbidden("Akses ditolak: tidak dapat melihat presensi pengguna lain")
	}

	res := toAttendanceResponse(att)
	return &res, nil
}

func toAttendanceResponse(a *model.Attendance) model.AttendanceResponse {
	var checkedOutStr *string
	if a.CheckedOutAt != nil {
		t := a.CheckedOutAt.UTC().Format(time.RFC3339)
		checkedOutStr = &t
	}

	return model.AttendanceResponse{
		ID:                 a.ID,
		UserID:             a.UserID,
		SessionID:          a.SessionID,
		CheckedInAt:        a.CheckedInAt.UTC().Format(time.RFC3339),
		CheckedOutAt:       checkedOutStr,
		VerificationStatus: a.VerificationStatus,
		UserIdentifier:     a.UserIdentifier,
		UserFullName:       a.UserFullName,
		SessionCode:        a.SessionCode,
		SessionTitle:       a.SessionTitle,
		SessionDate:        a.SessionDate,
	}
}
