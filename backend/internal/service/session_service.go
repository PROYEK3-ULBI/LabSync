package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type SessionService struct {
	sessionRepo repository.SessionRepository
	auditRepo   repository.AuditRepository
}

func NewSessionService(sessionRepo repository.SessionRepository, auditRepo repository.AuditRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		auditRepo:   auditRepo,
	}
}

func (s *SessionService) Create(ctx context.Context, adminID string, req model.CreateSessionRequest) (*model.SessionResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Title = strings.TrimSpace(req.Title)
	req.Location = strings.TrimSpace(req.Location)
	req.SessionDate = strings.TrimSpace(req.SessionDate)

	if len(req.Code) < 2 {
		return nil, model.NewBadRequest("Kode sesi (code) minimal 2 karakter")
	}
	if len(req.Title) < 2 {
		return nil, model.NewBadRequest("Judul sesi (title) minimal 2 karakter")
	}
	if len(req.Location) < 2 {
		return nil, model.NewBadRequest("Lokasi lab (location) minimal 2 karakter")
	}

	// Validasi format tanggal YYYY-MM-DD
	if _, err := time.Parse("2006-01-02", req.SessionDate); err != nil {
		return nil, model.NewBadRequest("Format tanggal (session_date) harus YYYY-MM-DD (contoh: 2026-09-18)")
	}

	// Validasi format waktu RFC3339 dengan semantik UTC
	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return nil, model.NewBadRequest("Format waktu mulai (start_at) tidak valid, gunakan RFC3339 (contoh: 2026-09-18T13:00:00Z)")
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return nil, model.NewBadRequest("Format waktu selesai (end_at) tidak valid, gunakan RFC3339 (contoh: 2026-09-18T15:30:00Z)")
	}

	if !startAt.Before(endAt) {
		return nil, model.NewBadRequest("Waktu selesai (end_at) harus setelah waktu mulai (start_at)")
	}

	status := model.SessionStatusScheduled
	if req.Status != nil && *req.Status != "" {
		st := strings.ToUpper(strings.TrimSpace(*req.Status))
		if !model.IsValidSessionStatus(st) {
			return nil, model.NewBadRequest("Status sesi tidak valid, pilihan: SCHEDULED, ACTIVE, CLOSED, CANCELLED")
		}
		status = model.SessionStatus(st)
	}

	// Cek duplikasi kode sesi
	if existing, _ := s.sessionRepo.GetByCode(ctx, req.Code); existing != nil {
		return nil, model.NewConflict("Kode sesi sudah digunakan")
	}

	var courseCode *string
	if req.CourseCode != nil && strings.TrimSpace(*req.CourseCode) != "" {
		c := strings.TrimSpace(*req.CourseCode)
		courseCode = &c
	}

	newSession := &model.LabSession{
		Code:        req.Code,
		Title:       req.Title,
		CourseCode:  courseCode,
		Location:    req.Location,
		SessionDate: req.SessionDate,
		StartAt:     startAt.UTC(),
		EndAt:       endAt.UTC(),
		Status:      status,
		CreatedBy:   &adminID,
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		if errors.Is(err, repository.ErrDuplicateSessionCode) {
			return nil, model.NewConflict("Kode sesi sudah digunakan")
		}
		return nil, model.NewInternalError("Gagal membuat sesi praktikum", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "SESSION_CREATED",
		EntityType: "LAB_SESSION",
		EntityID:   newSession.ID,
	})

	res := toSessionResponse(*newSession)
	return &res, nil
}

func (s *SessionService) List(ctx context.Context, q model.SessionListQuery) (*model.SessionListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	sessions, total, err := s.sessionRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar sesi praktikum", err)
	}

	data := make([]model.SessionResponse, len(sessions))
	for i, sess := range sessions {
		data[i] = toSessionResponse(sess)
	}

	return &model.SessionListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *SessionService) GetActive(ctx context.Context) ([]model.SessionResponse, error) {
	sessions, err := s.sessionRepo.GetActive(ctx)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar sesi aktif", err)
	}

	data := make([]model.SessionResponse, len(sessions))
	for i, sess := range sessions {
		data[i] = toSessionResponse(sess)
	}

	return data, nil
}

func (s *SessionService) GetByID(ctx context.Context, id string) (*model.SessionResponse, error) {
	sess, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, model.NewNotFound("Sesi praktikum tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil detail sesi praktikum", err)
	}

	res := toSessionResponse(*sess)
	return &res, nil
}

func (s *SessionService) Update(ctx context.Context, adminID string, id string, req model.UpdateSessionRequest) (*model.SessionResponse, error) {
	sess, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, model.NewNotFound("Sesi praktikum tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data sesi praktikum", err)
	}

	if req.Code != "" && req.Code != sess.Code {
		req.Code = strings.TrimSpace(req.Code)
		if len(req.Code) < 2 {
			return nil, model.NewBadRequest("Kode sesi minimal 2 karakter")
		}
		if existing, _ := s.sessionRepo.GetByCode(ctx, req.Code); existing != nil && existing.ID != id {
			return nil, model.NewConflict("Kode sesi sudah digunakan oleh sesi lain")
		}
		sess.Code = req.Code
	}

	if req.Title != "" {
		req.Title = strings.TrimSpace(req.Title)
		if len(req.Title) < 2 {
			return nil, model.NewBadRequest("Judul sesi minimal 2 karakter")
		}
		sess.Title = req.Title
	}

	if req.Location != "" {
		req.Location = strings.TrimSpace(req.Location)
		if len(req.Location) < 2 {
			return nil, model.NewBadRequest("Lokasi lab minimal 2 karakter")
		}
		sess.Location = req.Location
	}

	if req.CourseCode != nil {
		c := strings.TrimSpace(*req.CourseCode)
		if c == "" {
			sess.CourseCode = nil
		} else {
			sess.CourseCode = &c
		}
	}

	if req.SessionDate != "" {
		req.SessionDate = strings.TrimSpace(req.SessionDate)
		if _, err := time.Parse("2006-01-02", req.SessionDate); err != nil {
			return nil, model.NewBadRequest("Format tanggal harus YYYY-MM-DD")
		}
		sess.SessionDate = req.SessionDate
	}

	startChanged := req.StartAt != ""
	endChanged := req.EndAt != ""

	startAt := sess.StartAt
	endAt := sess.EndAt

	if startChanged {
		parsedStart, err := time.Parse(time.RFC3339, req.StartAt)
		if err != nil {
			return nil, model.NewBadRequest("Format waktu mulai tidak valid (gunakan RFC3339)")
		}
		startAt = parsedStart.UTC()
	}

	if endChanged {
		parsedEnd, err := time.Parse(time.RFC3339, req.EndAt)
		if err != nil {
			return nil, model.NewBadRequest("Format waktu selesai tidak valid (gunakan RFC3339)")
		}
		endAt = parsedEnd.UTC()
	}

	if !startAt.Before(endAt) {
		return nil, model.NewBadRequest("Waktu selesai (end_at) harus setelah waktu mulai (start_at)")
	}

	sess.StartAt = startAt
	sess.EndAt = endAt

	if req.Status != nil && *req.Status != "" {
		st := strings.ToUpper(strings.TrimSpace(*req.Status))
		if !model.IsValidSessionStatus(st) {
			return nil, model.NewBadRequest("Status tidak valid, pilihan: SCHEDULED, ACTIVE, CLOSED, CANCELLED")
		}
		sess.Status = model.SessionStatus(st)
	}

	if err := s.sessionRepo.Update(ctx, sess); err != nil {
		if errors.Is(err, repository.ErrDuplicateSessionCode) {
			return nil, model.NewConflict("Kode sesi sudah digunakan oleh sesi lain")
		}
		return nil, model.NewInternalError("Gagal memperbarui sesi praktikum", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "SESSION_UPDATED",
		EntityType: "LAB_SESSION",
		EntityID:   id,
	})

	res := toSessionResponse(*sess)
	return &res, nil
}

func (s *SessionService) UpdateStatus(ctx context.Context, adminID string, id string, req model.UpdateSessionStatusRequest) (*model.SessionResponse, error) {
	st := strings.ToUpper(strings.TrimSpace(req.Status))
	if !model.IsValidSessionStatus(st) {
		return nil, model.NewBadRequest("Status tidak valid, pilihan: SCHEDULED, ACTIVE, CLOSED, CANCELLED")
	}

	status := model.SessionStatus(st)
	if err := s.sessionRepo.UpdateStatus(ctx, id, status); err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, model.NewNotFound("Sesi praktikum tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal memperbarui status sesi praktikum", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "SESSION_STATUS_UPDATED",
		EntityType: "LAB_SESSION",
		EntityID:   id,
	})

	return s.GetByID(ctx, id)
}

func (s *SessionService) Delete(ctx context.Context, adminID string, id string) error {
	hasAtt, err := s.sessionRepo.HasAttendance(ctx, id)
	if err != nil {
		return model.NewInternalError("Gagal memeriksa data presensi sesi", err)
	}
	if hasAtt {
		return model.NewConflict("Sesi praktikum telah memiliki data presensi dan tidak dapat dihapus secara fisik. Silakan ubah status menjadi CANCELLED atau CLOSED")
	}

	if err := s.sessionRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return model.NewNotFound("Sesi praktikum tidak ditemukan")
		}
		if errors.Is(err, repository.ErrSessionHasAttendance) {
			return model.NewConflict("Sesi praktikum telah memiliki data presensi dan tidak dapat dihapus secara fisik. Silakan ubah status menjadi CANCELLED atau CLOSED")
		}
		return model.NewInternalError("Gagal menghapus sesi praktikum", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "SESSION_DELETED",
		EntityType: "LAB_SESSION",
		EntityID:   id,
	})

	return nil
}

func toSessionResponse(s model.LabSession) model.SessionResponse {
	return model.SessionResponse{
		ID:          s.ID,
		Code:        s.Code,
		Title:       s.Title,
		CourseCode:  s.CourseCode,
		Location:    s.Location,
		SessionDate: s.SessionDate,
		StartAt:     s.StartAt.UTC().Format(time.RFC3339),
		EndAt:       s.EndAt.UTC().Format(time.RFC3339),
		Status:      s.Status,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
