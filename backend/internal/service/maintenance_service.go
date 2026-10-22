package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type MaintenanceService struct {
	maintRepo repository.MaintenanceRepository
	assetRepo repository.AssetRepository
	allocRepo repository.AllocationRepository
	auditRepo repository.AuditRepository
}

func NewMaintenanceService(
	maintRepo repository.MaintenanceRepository,
	assetRepo repository.AssetRepository,
	allocRepo repository.AllocationRepository,
	auditRepo repository.AuditRepository,
) *MaintenanceService {
	return &MaintenanceService{
		maintRepo: maintRepo,
		assetRepo: assetRepo,
		allocRepo: allocRepo,
		auditRepo: auditRepo,
	}
}

func (s *MaintenanceService) Create(ctx context.Context, adminID string, req model.CreateMaintenanceRequest) (*model.MaintenanceResponse, error) {
	req.AssetID = strings.TrimSpace(req.AssetID)
	req.Description = strings.TrimSpace(req.Description)

	if req.AssetID == "" {
		return nil, model.NewBadRequest("Parameter asset_id wajib diisi")
	}
	if len(req.Description) < 3 {
		return nil, model.NewBadRequest("Deskripsi kendala pemeliharaan minimal 3 karakter")
	}

	// 1. Validasi keberadaan aset
	asset, err := s.assetRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal memeriksa data aset", err)
	}

	if asset.Status == model.AssetStatusRetired {
		return nil, model.NewBadRequest("Aset berstatus RETIRED tidak dapat diajukan pemeliharaan")
	}

	// 2. Invariant: Aset yang memiliki alokasi aktif tidak boleh diajukan ke maintenance
	hasAlloc, _ := s.assetRepo.HasActiveAllocation(ctx, req.AssetID)
	if hasAlloc || asset.Status == model.AssetStatusAllocated {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_ALLOCATED",
			Message:    "Aset sedang dipinjam dan tidak dapat diajukan pemeliharaan sebelum pengembalian aset diselesaikan",
		}
	}

	issueType := "HARDWARE"
	if req.IssueType != nil && strings.TrimSpace(*req.IssueType) != "" {
		issueType = strings.ToUpper(strings.TrimSpace(*req.IssueType))
	}

	priority := model.MaintenancePriorityMedium
	if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
		p := strings.ToUpper(strings.TrimSpace(*req.Priority))
		if !model.IsValidMaintenancePriority(p) {
			return nil, model.NewBadRequest("Prioritas tidak valid, pilihan: LOW, MEDIUM, HIGH")
		}
		priority = model.MaintenancePriority(p)
	}

	m := &model.Maintenance{
		AssetID:     req.AssetID,
		ReportedBy:  adminID,
		IssueType:   issueType,
		Description: req.Description,
		Priority:    priority,
		Status:      model.MaintenanceStatusOpen,
	}

	if err := s.maintRepo.CreateWithTransaction(ctx, m); err != nil {
		if errors.Is(err, repository.ErrAssetHasActiveAllocation) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "ASSET_ALLOCATED",
				Message:    "Aset sedang dipinjam dan tidak dapat diajukan pemeliharaan",
			}
		}
		return nil, model.NewInternalError("Gagal membuat tiket pemeliharaan aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "MAINTENANCE_CREATED",
		EntityType: "MAINTENANCE",
		EntityID:   m.ID,
	})

	m.AssetCode = asset.AssetCode
	m.AssetName = asset.AssetName

	res := toMaintenanceResponse(m)
	return &res, nil
}

func (s *MaintenanceService) Update(ctx context.Context, adminID string, id string, req model.UpdateMaintenanceRequest) (*model.MaintenanceResponse, error) {
	m, err := s.maintRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMaintenanceNotFound) {
			return nil, model.NewNotFound("Tiket pemeliharaan tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data pemeliharaan", err)
	}

	if req.Description != nil && strings.TrimSpace(*req.Description) != "" {
		m.Description = strings.TrimSpace(*req.Description)
	}
	if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
		p := strings.ToUpper(strings.TrimSpace(*req.Priority))
		if !model.IsValidMaintenancePriority(p) {
			return nil, model.NewBadRequest("Prioritas tidak valid, pilihan: LOW, MEDIUM, HIGH")
		}
		m.Priority = model.MaintenancePriority(p)
	}
	if req.Resolution != nil && strings.TrimSpace(*req.Resolution) != "" {
		resStr := strings.TrimSpace(*req.Resolution)
		m.Resolution = &resStr
	}
	if req.Cost != nil {
		m.Cost = req.Cost
	}

	var resultAssetStatus *model.AssetStatus
	auditAction := "MAINTENANCE_UPDATED"

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		newStatus := model.MaintenanceStatus(strings.ToUpper(strings.TrimSpace(*req.Status)))
		if !model.IsValidMaintenanceStatus(string(newStatus)) {
			return nil, model.NewBadRequest("Status pemeliharaan tidak valid, pilihan: OPEN, IN_PROGRESS, COMPLETED, CANCELLED")
		}

		if newStatus == model.MaintenanceStatusCompleted {
			if m.Resolution == nil || len(strings.TrimSpace(*m.Resolution)) == 0 {
				return nil, model.NewBadRequest("Resolusi perbaikan (resolution) wajib dicatat saat menyelesaikan tiket pemeliharaan")
			}
			targetStatus := model.AssetStatusAvailable
			if req.ResultStatus != nil && strings.EqualFold(*req.ResultStatus, "RETIRED") {
				targetStatus = model.AssetStatusRetired
			}
			resultAssetStatus = &targetStatus
			auditAction = "MAINTENANCE_COMPLETED"
		} else if newStatus == model.MaintenanceStatusCancelled {
			auditAction = "MAINTENANCE_CANCELLED"
		} else if newStatus == model.MaintenanceStatusInProgress {
			auditAction = "MAINTENANCE_IN_PROGRESS"
		}

		m.Status = newStatus
	}

	if err := s.maintRepo.UpdateWithTransaction(ctx, m, resultAssetStatus); err != nil {
		if errors.Is(err, repository.ErrInvalidStatusTransition) {
			return nil, &model.AppError{
				StatusCode: 400,
				Code:       "INVALID_STATUS_TRANSITION",
				Message:    fmt.Sprintf("Transisi status pemeliharaan ke %s tidak diizinkan", m.Status),
			}
		}
		return nil, model.NewInternalError("Gagal memperbarui tiket pemeliharaan", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     auditAction,
		EntityType: "MAINTENANCE",
		EntityID:   id,
	})

	res := toMaintenanceResponse(m)
	return &res, nil
}

func (s *MaintenanceService) GetByID(ctx context.Context, id string) (*model.MaintenanceResponse, error) {
	m, err := s.maintRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMaintenanceNotFound) {
			return nil, model.NewNotFound("Tiket pemeliharaan tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data pemeliharaan", err)
	}

	res := toMaintenanceResponse(m)
	return &res, nil
}

func (s *MaintenanceService) List(ctx context.Context, q model.MaintenanceListQuery) (*model.MaintenanceListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	list, total, err := s.maintRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar pemeliharaan", err)
	}

	data := make([]model.MaintenanceResponse, len(list))
	for i, item := range list {
		data[i] = toMaintenanceResponse(&item)
	}

	return &model.MaintenanceListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func toMaintenanceResponse(m *model.Maintenance) model.MaintenanceResponse {
	var startedAtStr, completedAtStr *string
	if m.StartedAt != nil {
		t := m.StartedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		startedAtStr = &t
	}
	if m.CompletedAt != nil {
		t := m.CompletedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		completedAtStr = &t
	}

	return model.MaintenanceResponse{
		ID:                 m.ID,
		AssetID:            m.AssetID,
		ReportedBy:         m.ReportedBy,
		IssueType:          m.IssueType,
		Description:        m.Description,
		Priority:           m.Priority,
		Status:             m.Status,
		StartedAt:          startedAtStr,
		CompletedAt:        completedAtStr,
		Resolution:         m.Resolution,
		Cost:               m.Cost,
		CreatedAt:          m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          m.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		AssetCode:          m.AssetCode,
		AssetName:          m.AssetName,
		ReporterIdentifier: m.ReporterIdentifier,
		ReporterName:       m.ReporterName,
	}
}
