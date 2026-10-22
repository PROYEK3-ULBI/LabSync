package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type AllocationService struct {
	allocRepo      repository.AllocationRepository
	assetRepo      repository.AssetRepository
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	attendanceRepo repository.AttendanceRepository
	auditRepo      repository.AuditRepository
}

func NewAllocationService(
	allocRepo repository.AllocationRepository,
	assetRepo repository.AssetRepository,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	attendanceRepo repository.AttendanceRepository,
	auditRepo repository.AuditRepository,
) *AllocationService {
	return &AllocationService{
		allocRepo:      allocRepo,
		assetRepo:      assetRepo,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		attendanceRepo: attendanceRepo,
		auditRepo:      auditRepo,
	}
}

func (s *AllocationService) Create(ctx context.Context, adminID string, req model.CreateAllocationRequest) (*model.AllocationResponse, error) {
	req.AssetID = strings.TrimSpace(req.AssetID)
	req.UserID = strings.TrimSpace(req.UserID)
	req.SessionID = strings.TrimSpace(req.SessionID)

	if req.AssetID == "" || req.UserID == "" || req.SessionID == "" {
		return nil, model.NewBadRequest("Parameter asset_id, user_id, dan session_id wajib diisi")
	}

	// 1. Validasi pengguna aktif
	if s.userRepo != nil {
		u, err := s.userRepo.GetByID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, model.NewNotFound("Pengguna tidak ditemukan")
			}
			return nil, model.NewInternalError("Gagal memeriksa data pengguna", err)
		}
		if u.Status != "ACTIVE" {
			return nil, model.NewForbidden("Akun pengguna tidak aktif")
		}
	}

	// 2. Validasi sesi aktif
	if s.sessionRepo != nil {
		sess, err := s.sessionRepo.GetByID(ctx, req.SessionID)
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
	}

	// 3. Prasyarat Presensi (Attendance Prerequisite):
	// Mahasiswa wajib memiliki catatan presensi masuk aktif pada sesi ini sebelum menerima alokasi aset
	if s.attendanceRepo != nil {
		att, err := s.attendanceRepo.GetByUserAndSession(ctx, req.UserID, req.SessionID)
		if err != nil {
			if errors.Is(err, repository.ErrAttendanceNotFound) {
				return nil, &model.AppError{
					StatusCode: 400,
					Code:       "ATTENDANCE_REQUIRED",
					Message:    "Mahasiswa wajib melakukan presensi masuk (check-in) pada sesi ini sebelum aset dapat dialokasikan",
				}
			}
			return nil, model.NewInternalError("Gagal memeriksa data presensi mahasiswa", err)
		}
		if att.CheckedOutAt != nil {
			return nil, &model.AppError{
				StatusCode: 400,
				Code:       "ATTENDANCE_ALREADY_CHECKED_OUT",
				Message:    "Mahasiswa sudah melakukan check-out pada sesi ini",
			}
		}
	}

	// 4. Validasi ketersediaan aset
	asset, err := s.assetRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal memeriksa data aset", err)
	}

	if asset.Status != model.AssetStatusAvailable {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_NOT_AVAILABLE",
			Message:    "Aset tidak berstatus AVAILABLE dan tidak dapat dialokasikan",
		}
	}

	condBefore := "GOOD"
	if req.ConditionBefore != nil && strings.TrimSpace(*req.ConditionBefore) != "" {
		condBefore = strings.ToUpper(strings.TrimSpace(*req.ConditionBefore))
	}

	now := time.Now().UTC()
	alloc := &model.Allocation{
		AssetID:         req.AssetID,
		UserID:          req.UserID,
		SessionID:       req.SessionID,
		AllocatedAt:     now,
		ConditionBefore: condBefore,
		Notes:           req.Notes,
	}

	// 5. Eksekusi alokasi dan pembaruan status aset dalam satu transaksi database atomik
	if err := s.allocRepo.CreateWithTransaction(ctx, alloc); err != nil {
		if errors.Is(err, repository.ErrDoubleAllocation) || errors.Is(err, repository.ErrAssetNotAvailable) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "ASSET_ALREADY_ALLOCATED",
				Message:    "Aset sedang dialokasikan dan belum dikembalikan",
			}
		}
		return nil, model.NewInternalError("Gagal memproses transaksi alokasi aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ALLOCATION_CREATED",
		EntityType: "ALLOCATION",
		EntityID:   alloc.ID,
	})

	// Lengkapi data relasional untuk respons
	alloc.AssetCode = asset.AssetCode
	alloc.AssetName = asset.AssetName

	res := toAllocationResponse(alloc)
	return &res, nil
}

func (s *AllocationService) Return(ctx context.Context, adminID string, id string, req model.ReturnAllocationRequest) (*model.AllocationResponse, error) {
	cond := strings.ToUpper(strings.TrimSpace(req.Condition))
	if cond != "GOOD" && cond != "DAMAGED" {
		return nil, model.NewBadRequest("Kondisi pengembalian (condition) harus GOOD atau DAMAGED")
	}

	now := time.Now().UTC()
	updatedAlloc, err := s.allocRepo.ReturnWithTransaction(ctx, id, now, cond, req.Notes)
	if err != nil {
		if errors.Is(err, repository.ErrAllocationNotFound) {
			return nil, model.NewNotFound("Alokasi aset tidak ditemukan")
		}
		if errors.Is(err, repository.ErrAlreadyReturned) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "ALREADY_RETURNED",
				Message:    "Alokasi aset sudah dikembalikan sebelumnya",
			}
		}
		return nil, model.NewInternalError("Gagal memproses pengembalian aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ALLOCATION_RETURNED",
		EntityType: "ALLOCATION",
		EntityID:   id,
	})

	res := toAllocationResponse(updatedAlloc)
	return &res, nil
}

func (s *AllocationService) GetByID(ctx context.Context, callerID, callerRole, id string) (*model.AllocationResponse, error) {
	alloc, err := s.allocRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAllocationNotFound) {
			return nil, model.NewNotFound("Alokasi aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data alokasi", err)
	}

	if callerRole != "ADMIN" && alloc.UserID != callerID {
		return nil, model.NewForbidden("Akses ditolak: tidak dapat melihat alokasi pengguna lain")
	}

	res := toAllocationResponse(alloc)
	return &res, nil
}

func (s *AllocationService) List(ctx context.Context, q model.AllocationListQuery) (*model.AllocationListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	list, total, err := s.allocRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar alokasi", err)
	}

	data := make([]model.AllocationResponse, len(list))
	for i, item := range list {
		data[i] = toAllocationResponse(&item)
	}

	return &model.AllocationListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *AllocationService) GetMyAllocations(ctx context.Context, studentID string, q model.AllocationListQuery) (*model.AllocationListResponse, error) {
	q.UserID = studentID
	return s.List(ctx, q)
}

func (s *AllocationService) GetAssetHistory(ctx context.Context, callerID, callerRole, assetID string) (*model.AssetHistoryResponse, error) {
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data aset", err)
	}

	timeline, err := s.allocRepo.GetAssetTimeline(ctx, assetID)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil histori timeline aset", err)
	}

	// Filter timeline jika diakses oleh mahasiswa (hanya melihat event miliknya)
	var scopedTimeline []model.AssetTimelineEvent
	if callerRole != "ADMIN" {
		for _, ev := range timeline {
			if ev.ActorID != nil && *ev.ActorID == callerID {
				scopedTimeline = append(scopedTimeline, ev)
			}
		}
	} else {
		scopedTimeline = timeline
	}

	assetRes := toAssetResponse(asset)
	return &model.AssetHistoryResponse{
		Asset:    assetRes,
		Timeline: scopedTimeline,
	}, nil
}

func toAllocationResponse(a *model.Allocation) model.AllocationResponse {
	var returnedAtStr *string
	if a.ReturnedAt != nil {
		t := a.ReturnedAt.UTC().Format(time.RFC3339)
		returnedAtStr = &t
	}

	return model.AllocationResponse{
		ID:              a.ID,
		AssetID:         a.AssetID,
		UserID:          a.UserID,
		SessionID:       a.SessionID,
		AllocatedAt:     a.AllocatedAt.UTC().Format(time.RFC3339),
		ReturnedAt:      returnedAtStr,
		ConditionBefore: a.ConditionBefore,
		ConditionAfter:  a.ConditionAfter,
		Notes:           a.Notes,
		AssetCode:       a.AssetCode,
		AssetName:       a.AssetName,
		UserIdentifier:  a.UserIdentifier,
		UserFullName:    a.UserFullName,
		SessionCode:     a.SessionCode,
		SessionTitle:    a.SessionTitle,
	}
}
