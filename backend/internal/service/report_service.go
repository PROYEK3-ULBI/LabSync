package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
	"labsync-backend/pkg/csvutil"
)

type ReportService struct {
	attendanceRepo repository.AttendanceRepository
	assetRepo      repository.AssetRepository
	allocRepo      repository.AllocationRepository
	maintRepo      repository.MaintenanceRepository
}

func NewReportService(
	attendanceRepo repository.AttendanceRepository,
	assetRepo repository.AssetRepository,
	allocRepo repository.AllocationRepository,
	maintRepo repository.MaintenanceRepository,
) *ReportService {
	return &ReportService{
		attendanceRepo: attendanceRepo,
		assetRepo:      assetRepo,
		allocRepo:      allocRepo,
		maintRepo:      maintRepo,
	}
}

func validateDateRange(startDate, endDate string) error {
	if startDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			return model.NewBadRequest("Format start_date harus YYYY-MM-DD")
		}
	}
	if endDate != "" {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			return model.NewBadRequest("Format end_date harus YYYY-MM-DD")
		}
	}
	if startDate != "" && endDate != "" {
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		if start.After(end) {
			return model.NewBadRequest("Tanggal awal (start_date) tidak boleh setelah tanggal akhir (end_date)")
		}
	}
	return nil
}

func (s *ReportService) AttendanceReport(ctx context.Context, q model.ReportFilterQuery) (*model.AttendanceReportResponse, []byte, error) {
	if err := validateDateRange(q.StartDate, q.EndDate); err != nil {
		return nil, nil, err
	}

	attQuery := model.AttendanceListQuery{
		Page:      1,
		Limit:     1000,
		SessionID: q.SessionID,
		UserID:    q.UserID,
		Date:      q.StartDate,
	}

	list, total, err := s.attendanceRepo.List(ctx, attQuery)
	if err != nil {
		return nil, nil, model.NewInternalError("Gagal mengambil data laporan presensi", err)
	}

	userSet := make(map[string]bool)
	sessSet := make(map[string]bool)
	data := make([]model.AttendanceResponse, len(list))

	for i, item := range list {
		data[i] = toAttendanceResponse(&item)
		userSet[item.UserID] = true
		sessSet[item.SessionID] = true
	}

	summary := model.AttendanceReportSummary{
		TotalAttendance: total,
		UniqueUsers:     len(userSet),
		UniqueSessions:  len(sessSet),
	}

	if strings.EqualFold(q.Format, "csv") {
		headers := []string{"ID", "NIM", "Nama Mahasiswa", "Kode Sesi", "Judul Sesi", "Tanggal", "Check-in", "Check-out", "Status Verifikasi"}
		rows := make([][]string, len(data))
		for i, d := range data {
			checkOut := "-"
			if d.CheckedOutAt != nil {
				checkOut = *d.CheckedOutAt
			}
			rows[i] = []string{
				d.ID,
				d.UserIdentifier,
				d.UserFullName,
				d.SessionCode,
				d.SessionTitle,
				d.SessionDate,
				d.CheckedInAt,
				checkOut,
				d.VerificationStatus,
			}
		}
		csvBytes, err := csvutil.GenerateCSV(headers, rows)
		if err != nil {
			return nil, nil, model.NewInternalError("Gagal membuat ekspor CSV presensi", err)
		}
		return nil, csvBytes, nil
	}

	return &model.AttendanceReportResponse{
		Summary: summary,
		Data:    data,
	}, nil, nil
}

func (s *ReportService) AssetReport(ctx context.Context, q model.ReportFilterQuery) (*model.AssetReportResponse, []byte, error) {
	assetQuery := model.AssetListQuery{
		Page:     1,
		Limit:    1000,
		Status:   q.Status,
		Category: q.Category,
	}

	list, total, err := s.assetRepo.List(ctx, assetQuery)
	if err != nil {
		return nil, nil, model.NewInternalError("Gagal mengambil data laporan aset", err)
	}

	byStatus := make(map[string]int)
	byCategory := make(map[string]int)
	data := make([]model.AssetResponse, len(list))

	for i, item := range list {
		data[i] = toAssetResponse(&item)
		byStatus[string(item.Status)]++
		byCategory[item.Category]++
	}

	summary := model.AssetReportSummary{
		TotalAssets: total,
		ByStatus:    byStatus,
		ByCategory:  byCategory,
	}

	if strings.EqualFold(q.Format, "csv") {
		headers := []string{"ID", "Kode Aset", "Nama Aset", "Kategori", "Merek", "Model", "Nomor Seri", "Lokasi", "Status"}
		rows := make([][]string, len(data))
		for i, d := range data {
			brand, modelName, sn := "-", "-", "-"
			if d.Brand != nil {
				brand = *d.Brand
			}
			if d.Model != nil {
				modelName = *d.Model
			}
			if d.SerialNumber != nil {
				sn = *d.SerialNumber
			}
			rows[i] = []string{
				d.ID,
				d.AssetCode,
				d.AssetName,
				d.Category,
				brand,
				modelName,
				sn,
				d.Location,
				string(d.Status),
			}
		}
		csvBytes, err := csvutil.GenerateCSV(headers, rows)
		if err != nil {
			return nil, nil, model.NewInternalError("Gagal membuat ekspor CSV aset", err)
		}
		return nil, csvBytes, nil
	}

	return &model.AssetReportResponse{
		Summary: summary,
		Data:    data,
	}, nil, nil
}

func (s *ReportService) AllocationReport(ctx context.Context, q model.ReportFilterQuery) (*model.AllocationReportResponse, []byte, error) {
	if err := validateDateRange(q.StartDate, q.EndDate); err != nil {
		return nil, nil, err
	}

	allocQuery := model.AllocationListQuery{
		Page:      1,
		Limit:     1000,
		AssetID:   q.AssetID,
		UserID:    q.UserID,
		SessionID: q.SessionID,
		Status:    q.Status,
	}

	list, total, err := s.allocRepo.List(ctx, allocQuery)
	if err != nil {
		return nil, nil, model.NewInternalError("Gagal mengambil data laporan alokasi", err)
	}

	activeCount := 0
	returnedCount := 0
	damagedCount := 0
	data := make([]model.AllocationResponse, len(list))

	for i, item := range list {
		data[i] = toAllocationResponse(&item)
		if item.ReturnedAt == nil {
			activeCount++
		} else {
			returnedCount++
			if item.ConditionAfter != nil && strings.EqualFold(*item.ConditionAfter, "DAMAGED") {
				damagedCount++
			}
		}
	}

	summary := model.AllocationReportSummary{
		TotalAllocations: total,
		ActiveCount:      activeCount,
		ReturnedCount:    returnedCount,
		DamagedCount:     damagedCount,
	}

	if strings.EqualFold(q.Format, "csv") {
		headers := []string{"ID", "Kode Aset", "Nama Aset", "NIM", "Nama Mahasiswa", "Kode Sesi", "Waktu Pinjam", "Waktu Kembali", "Kondisi Sebelum", "Kondisi Sesudah"}
		rows := make([][]string, len(data))
		for i, d := range data {
			returnedAt := "-"
			if d.ReturnedAt != nil {
				returnedAt = *d.ReturnedAt
			}
			condAfter := "-"
			if d.ConditionAfter != nil {
				condAfter = *d.ConditionAfter
			}
			rows[i] = []string{
				d.ID,
				d.AssetCode,
				d.AssetName,
				d.UserIdentifier,
				d.UserFullName,
				d.SessionCode,
				d.AllocatedAt,
				returnedAt,
				d.ConditionBefore,
				condAfter,
			}
		}
		csvBytes, err := csvutil.GenerateCSV(headers, rows)
		if err != nil {
			return nil, nil, model.NewInternalError("Gagal membuat ekspor CSV alokasi", err)
		}
		return nil, csvBytes, nil
	}

	return &model.AllocationReportResponse{
		Summary: summary,
		Data:    data,
	}, nil, nil
}

func (s *ReportService) MaintenanceReport(ctx context.Context, q model.ReportFilterQuery) (*model.MaintenanceReportResponse, []byte, error) {
	if err := validateDateRange(q.StartDate, q.EndDate); err != nil {
		return nil, nil, err
	}

	maintQuery := model.MaintenanceListQuery{
		Page:     1,
		Limit:    1000,
		AssetID:  q.AssetID,
		Status:   q.Status,
		Priority: q.Priority,
	}

	list, total, err := s.maintRepo.List(ctx, maintQuery)
	if err != nil {
		return nil, nil, model.NewInternalError("Gagal mengambil data laporan pemeliharaan", err)
	}

	openCount, inProgCount, compCount, cancelCount := 0, 0, 0, 0
	byPriority := make(map[string]int)
	var totalCost float64
	data := make([]model.MaintenanceResponse, len(list))

	for i, item := range list {
		data[i] = toMaintenanceResponse(&item)
		byPriority[string(item.Priority)]++
		switch item.Status {
		case model.MaintenanceStatusOpen:
			openCount++
		case model.MaintenanceStatusInProgress:
			inProgCount++
		case model.MaintenanceStatusCompleted:
			compCount++
		case model.MaintenanceStatusCancelled:
			cancelCount++
		}
		if item.Cost != nil {
			totalCost += *item.Cost
		}
	}

	summary := model.MaintenanceReportSummary{
		TotalTickets:    total,
		OpenCount:       openCount,
		InProgressCount: inProgCount,
		CompletedCount:  compCount,
		CancelledCount:  cancelCount,
		ByPriority:      byPriority,
		TotalCost:       totalCost,
	}

	if strings.EqualFold(q.Format, "csv") {
		headers := []string{"ID", "Kode Aset", "Nama Aset", "Pelapor", "Jenis Kendala", "Deskripsi", "Prioritas", "Status", "Waktu Selesai", "Resolusi", "Biaya"}
		rows := make([][]string, len(data))
		for i, d := range data {
			completedAt := "-"
			if d.CompletedAt != nil {
				completedAt = *d.CompletedAt
			}
			resStr := "-"
			if d.Resolution != nil {
				resStr = *d.Resolution
			}
			costStr := "0"
			if d.Cost != nil {
				costStr = fmt.Sprintf("%.2f", *d.Cost)
			}
			rows[i] = []string{
				d.ID,
				d.AssetCode,
				d.AssetName,
				d.ReporterName,
				d.IssueType,
				d.Description,
				string(d.Priority),
				string(d.Status),
				completedAt,
				resStr,
				costStr,
			}
		}
		csvBytes, err := csvutil.GenerateCSV(headers, rows)
		if err != nil {
			return nil, nil, model.NewInternalError("Gagal membuat ekspor CSV pemeliharaan", err)
		}
		return nil, csvBytes, nil
	}

	return &model.MaintenanceReportResponse{
		Summary: summary,
		Data:    data,
	}, nil, nil
}
