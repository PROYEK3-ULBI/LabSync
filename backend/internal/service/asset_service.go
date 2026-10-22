package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type AssetService struct {
	assetRepo repository.AssetRepository
	auditRepo repository.AuditRepository
}

func NewAssetService(assetRepo repository.AssetRepository, auditRepo repository.AuditRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		auditRepo: auditRepo,
	}
}

func (s *AssetService) Create(ctx context.Context, adminID string, req model.CreateAssetRequest) (*model.AssetResponse, error) {
	req.AssetCode = strings.TrimSpace(req.AssetCode)
	req.AssetName = strings.TrimSpace(req.AssetName)

	if len(req.AssetCode) < 2 {
		return nil, model.NewBadRequest("Kode aset (asset_code) minimal 2 karakter")
	}
	if len(req.AssetName) < 2 {
		return nil, model.NewBadRequest("Nama aset (asset_name) minimal 2 karakter")
	}

	category := "Umum"
	if req.Category != nil && strings.TrimSpace(*req.Category) != "" {
		category = strings.TrimSpace(*req.Category)
	}

	location := "Lab Komputer"
	if req.Location != nil && strings.TrimSpace(*req.Location) != "" {
		location = strings.TrimSpace(*req.Location)
	}

	status := model.AssetStatusAvailable
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		st := strings.ToUpper(strings.TrimSpace(*req.Status))
		if !model.IsValidAssetStatus(st) {
			return nil, model.NewBadRequest("Status aset tidak valid, pilihan: AVAILABLE, ALLOCATED, MAINTENANCE, RETIRED")
		}
		if st == string(model.AssetStatusAllocated) {
			return nil, model.NewBadRequest("Status ALLOCATED hanya dapat diberikan melalui pembuatan alokasi")
		}
		status = model.AssetStatus(st)
	}

	// Cek keunikan kode aset
	if existing, _ := s.assetRepo.GetByCode(ctx, req.AssetCode); existing != nil {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "DUPLICATE_ASSET_CODE",
			Message:    "Kode aset sudah terdaftar dalam sistem",
		}
	}

	// Cek keunikan nomor seri jika diisi
	if req.SerialNumber != nil && strings.TrimSpace(*req.SerialNumber) != "" {
		sn := strings.TrimSpace(*req.SerialNumber)
		req.SerialNumber = &sn
		if existing, _ := s.assetRepo.GetBySerialNumber(ctx, sn); existing != nil {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_SERIAL_NUMBER",
				Message:    "Nomor seri sudah terdaftar pada aset lain",
			}
		}
	}

	asset := &model.Asset{
		AssetCode:    req.AssetCode,
		AssetName:    req.AssetName,
		Category:     category,
		Brand:        req.Brand,
		Model:        req.Model,
		SerialNumber: req.SerialNumber,
		Location:     location,
		Status:       status,
		Notes:        req.Notes,
	}

	if err := s.assetRepo.Create(ctx, asset); err != nil {
		if errors.Is(err, repository.ErrDuplicateAssetCode) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_ASSET_CODE",
				Message:    "Kode aset sudah terdaftar dalam sistem",
			}
		}
		if errors.Is(err, repository.ErrDuplicateSerialNumber) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_SERIAL_NUMBER",
				Message:    "Nomor seri sudah terdaftar pada aset lain",
			}
		}
		return nil, model.NewInternalError("Gagal menambahkan aset baru", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ASSET_CREATED",
		EntityType: "ASSET",
		EntityID:   asset.ID,
	})

	res := toAssetResponse(asset)
	return &res, nil
}

func (s *AssetService) GetByID(ctx context.Context, id string) (*model.AssetResponse, error) {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data aset", err)
	}

	res := toAssetResponse(asset)
	return &res, nil
}

func (s *AssetService) List(ctx context.Context, q model.AssetListQuery) (*model.AssetListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	list, total, err := s.assetRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar aset", err)
	}

	data := make([]model.AssetResponse, len(list))
	for i, item := range list {
		data[i] = toAssetResponse(&item)
	}

	return &model.AssetListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *AssetService) Update(ctx context.Context, adminID string, id string, req model.UpdateAssetRequest) (*model.AssetResponse, error) {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data aset", err)
	}

	// 1. Guard proteksi alokasi aktif
	hasAlloc, _ := s.assetRepo.HasActiveAllocation(ctx, id)
	if hasAlloc && req.Status != nil && *req.Status != string(asset.Status) {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_ALLOCATED",
			Message:    "Aset sedang dipinjam/dialokasikan dan tidak dapat diubah statusnya secara manual",
		}
	}

	// 2. Guard proteksi perbaikan/pemeliharaan aktif
	hasMaint, _ := s.assetRepo.HasActiveMaintenance(ctx, id)
	if hasMaint && req.Status != nil && *req.Status != string(asset.Status) {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_IN_MAINTENANCE",
			Message:    "Aset sedang dalam proses perbaikan/pemeliharaan aktif",
		}
	}

	// 3. Validasi perubahan kode aset
	if req.AssetCode != "" && req.AssetCode != asset.AssetCode {
		req.AssetCode = strings.TrimSpace(req.AssetCode)
		if len(req.AssetCode) < 2 {
			return nil, model.NewBadRequest("Kode aset minimal 2 karakter")
		}
		if existing, _ := s.assetRepo.GetByCode(ctx, req.AssetCode); existing != nil && existing.ID != id {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_ASSET_CODE",
				Message:    "Kode aset sudah digunakan oleh aset lain",
			}
		}
		asset.AssetCode = req.AssetCode
	}

	if req.AssetName != "" {
		req.AssetName = strings.TrimSpace(req.AssetName)
		if len(req.AssetName) < 2 {
			return nil, model.NewBadRequest("Nama aset minimal 2 karakter")
		}
		asset.AssetName = req.AssetName
	}

	if req.Category != nil {
		asset.Category = strings.TrimSpace(*req.Category)
	}
	if req.Brand != nil {
		asset.Brand = req.Brand
	}
	if req.Model != nil {
		asset.Model = req.Model
	}
	if req.Location != nil {
		asset.Location = strings.TrimSpace(*req.Location)
	}
	if req.Notes != nil {
		asset.Notes = req.Notes
	}

	// 4. Validasi nomor seri
	if req.SerialNumber != nil && *req.SerialNumber != "" {
		sn := strings.TrimSpace(*req.SerialNumber)
		if existing, _ := s.assetRepo.GetBySerialNumber(ctx, sn); existing != nil && existing.ID != id {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_SERIAL_NUMBER",
				Message:    "Nomor seri sudah terdaftar pada aset lain",
			}
		}
		asset.SerialNumber = &sn
	}

	// 5. Validasi transisi status
	if req.Status != nil && *req.Status != "" {
		toStatus := model.AssetStatus(strings.ToUpper(strings.TrimSpace(*req.Status)))
		if !model.IsValidAssetStatus(string(toStatus)) {
			return nil, model.NewBadRequest("Status tidak valid, pilihan: AVAILABLE, ALLOCATED, MAINTENANCE, RETIRED")
		}
		if !isValidTransition(asset.Status, toStatus) {
			return nil, &model.AppError{
				StatusCode: 400,
				Code:       "INVALID_STATUS_TRANSITION",
				Message:    fmt.Sprintf("Transisi status dari %s ke %s tidak diizinkan", asset.Status, toStatus),
			}
		}
		if toStatus == model.AssetStatusAllocated && asset.Status != model.AssetStatusAllocated {
			return nil, model.NewBadRequest("Status ALLOCATED hanya dapat diberikan melalui pembuatan alokasi")
		}
		asset.Status = toStatus
	}

	if err := s.assetRepo.Update(ctx, asset); err != nil {
		if errors.Is(err, repository.ErrDuplicateAssetCode) {
			return nil, &model.AppError{
				StatusCode: 409,
				Code:       "DUPLICATE_ASSET_CODE",
				Message:    "Kode aset sudah digunakan oleh aset lain",
			}
		}
		return nil, model.NewInternalError("Gagal memperbarui aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ASSET_UPDATED",
		EntityType: "ASSET",
		EntityID:   id,
	})

	res := toAssetResponse(asset)
	return &res, nil
}

func (s *AssetService) UpdateStatus(ctx context.Context, adminID string, id string, req model.UpdateAssetStatusRequest) (*model.AssetResponse, error) {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return nil, model.NewNotFound("Aset tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data aset", err)
	}

	// Guard proteksi relasi aktif
	hasAlloc, _ := s.assetRepo.HasActiveAllocation(ctx, id)
	if hasAlloc {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_ALLOCATED",
			Message:    "Aset sedang dipinjam/dialokasikan dan tidak dapat diubah statusnya secara manual",
		}
	}

	hasMaint, _ := s.assetRepo.HasActiveMaintenance(ctx, id)
	if hasMaint {
		return nil, &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_IN_MAINTENANCE",
			Message:    "Aset sedang dalam proses perbaikan/pemeliharaan aktif",
		}
	}

	toStatus := model.AssetStatus(strings.ToUpper(strings.TrimSpace(req.Status)))
	if !model.IsValidAssetStatus(string(toStatus)) {
		return nil, model.NewBadRequest("Status tidak valid, pilihan: AVAILABLE, ALLOCATED, MAINTENANCE, RETIRED")
	}

	if !isValidTransition(asset.Status, toStatus) {
		return nil, &model.AppError{
			StatusCode: 400,
			Code:       "INVALID_STATUS_TRANSITION",
			Message:    fmt.Sprintf("Transisi status dari %s ke %s tidak diizinkan", asset.Status, toStatus),
		}
	}

	if toStatus == model.AssetStatusAllocated && asset.Status != model.AssetStatusAllocated {
		return nil, model.NewBadRequest("Status ALLOCATED hanya dapat diberikan melalui pembuatan alokasi")
	}

	if err := s.assetRepo.UpdateStatus(ctx, id, toStatus); err != nil {
		return nil, model.NewInternalError("Gagal memperbarui status aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ASSET_STATUS_UPDATED",
		EntityType: "ASSET",
		EntityID:   id,
	})

	asset.Status = toStatus
	res := toAssetResponse(asset)
	return &res, nil
}

func (s *AssetService) Delete(ctx context.Context, adminID string, id string) error {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return model.NewNotFound("Aset tidak ditemukan")
		}
		return model.NewInternalError("Gagal mengambil data aset", err)
	}

	// Guard: tidak boleh menghapus aset yang sedang dipinjam
	hasAlloc, _ := s.assetRepo.HasActiveAllocation(ctx, id)
	if hasAlloc || asset.Status == model.AssetStatusAllocated {
		return &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_ALLOCATED",
			Message:    "Aset sedang dipinjam dan tidak dapat dihapus",
		}
	}

	// Guard: tidak boleh menghapus aset yang sedang dalam pemeliharaan
	hasMaint, _ := s.assetRepo.HasActiveMaintenance(ctx, id)
	if hasMaint || asset.Status == model.AssetStatusMaintenance {
		return &model.AppError{
			StatusCode: 409,
			Code:       "ASSET_IN_MAINTENANCE",
			Message:    "Aset sedang dalam proses perbaikan/pemeliharaan dan tidak dapat dihapus",
		}
	}

	if err := s.assetRepo.Delete(ctx, id); err != nil {
		return model.NewInternalError("Gagal menghapus aset", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "ASSET_DELETED",
		EntityType: "ASSET",
		EntityID:   id,
	})

	return nil
}

func isValidTransition(from, to model.AssetStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case model.AssetStatusAvailable:
		return to == model.AssetStatusAllocated || to == model.AssetStatusMaintenance || to == model.AssetStatusRetired
	case model.AssetStatusAllocated:
		return to == model.AssetStatusAvailable || to == model.AssetStatusMaintenance
	case model.AssetStatusMaintenance:
		return to == model.AssetStatusAvailable || to == model.AssetStatusRetired
	case model.AssetStatusRetired:
		return false
	default:
		return false
	}
}

func toAssetResponse(a *model.Asset) model.AssetResponse {
	return model.AssetResponse{
		ID:           a.ID,
		AssetCode:    a.AssetCode,
		AssetName:    a.AssetName,
		Category:     a.Category,
		Brand:        a.Brand,
		Model:        a.Model,
		SerialNumber: a.SerialNumber,
		Location:     a.Location,
		Status:       a.Status,
		Notes:        a.Notes,
		CreatedAt:    a.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    a.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
