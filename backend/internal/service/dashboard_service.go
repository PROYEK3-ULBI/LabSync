package service

import (
	"context"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type DashboardService struct {
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	attendanceRepo repository.AttendanceRepository
	assetRepo      repository.AssetRepository
	allocRepo      repository.AllocationRepository
	maintRepo      repository.MaintenanceRepository
	auditRepo      repository.AuditRepository
}

func NewDashboardService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	attendanceRepo repository.AttendanceRepository,
	assetRepo repository.AssetRepository,
	allocRepo repository.AllocationRepository,
	maintRepo repository.MaintenanceRepository,
	auditRepo repository.AuditRepository,
) *DashboardService {
	return &DashboardService{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		attendanceRepo: attendanceRepo,
		assetRepo:      assetRepo,
		allocRepo:      allocRepo,
		maintRepo:      maintRepo,
		auditRepo:      auditRepo,
	}
}

func (s *DashboardService) GetStudentDashboard(ctx context.Context, studentID string) (*model.StudentDashboardResponse, error) {
	// 1. Profil pengguna
	user, err := s.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, model.NewNotFound("Pengguna tidak ditemukan")
	}

	faceEnrolled, _ := s.userRepo.HasFaceTemplate(ctx, studentID)

	userRes := model.UserResponse{
		ID:           user.ID,
		Identifier:   user.Identifier,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         user.RoleCode,
		Status:       user.Status,
		FaceEnrolled: &faceEnrolled,
	}

	// 2. Sesi aktif saat ini
	var activeSessionRes *model.SessionResponse
	activeSessions, _ := s.sessionRepo.GetActive(ctx)
	if len(activeSessions) > 0 {
		res := toSessionResponse(activeSessions[0])
		activeSessionRes = &res
	}

	// 3. Sesi praktikum mendatang (SCHEDULED)
	upcomingList, _, _ := s.sessionRepo.List(ctx, model.SessionListQuery{
		Status: string(model.SessionStatusScheduled),
		Limit:  5,
	})
	upcomingRes := make([]model.SessionResponse, len(upcomingList))
	for i, sess := range upcomingList {
		upcomingRes[i] = toSessionResponse(sess)
	}

	// 4. Riwayat presensi terbaru (5 terakhir) & total kehadiran
	attList, totalAttended, _ := s.attendanceRepo.List(ctx, model.AttendanceListQuery{
		UserID: studentID,
		Limit:  5,
	})
	recentAtt := make([]model.AttendanceResponse, len(attList))
	for i, a := range attList {
		recentAtt[i] = toAttendanceResponse(&a)
	}

	// 5. Alokasi aset aktif saat ini & total alokasi
	var currentAllocRes *model.AllocationResponse
	activeAllocs, _, _ := s.allocRepo.List(ctx, model.AllocationListQuery{
		UserID: studentID,
		Status: "ACTIVE",
		Limit:  1,
	})
	if len(activeAllocs) > 0 {
		res := toAllocationResponse(&activeAllocs[0])
		currentAllocRes = &res
	}

	_, totalAllocations, _ := s.allocRepo.List(ctx, model.AllocationListQuery{
		UserID: studentID,
	})

	return &model.StudentDashboardResponse{
		User:              userRes,
		FaceEnrolled:      faceEnrolled,
		ActiveSession:     activeSessionRes,
		UpcomingSessions:  upcomingRes,
		RecentAttendance:  recentAtt,
		TotalAttended:     totalAttended,
		CurrentAllocation: currentAllocRes,
		TotalAllocations:  totalAllocations,
	}, nil
}

func (s *DashboardService) GetAdminDashboard(ctx context.Context) (*model.AdminDashboardResponse, error) {
	// 1. Sesi aktif
	activeSessions, _ := s.sessionRepo.GetActive(ctx)
	activeSessionsRes := make([]model.SessionResponse, len(activeSessions))
	for i, sess := range activeSessions {
		activeSessionsRes[i] = toSessionResponse(sess)
	}

	// 2. Total pengguna & mahasiswa
	_, totalUsers, _ := s.userRepo.List(ctx, model.UserListQuery{Limit: 1})
	_, totalStudents, _ := s.userRepo.List(ctx, model.UserListQuery{Role: "STUDENT", Limit: 1})

	// 3. Ringkasan aset per status
	assetStatusCounts := make(map[string]int)
	_, availCount, _ := s.assetRepo.List(ctx, model.AssetListQuery{Status: string(model.AssetStatusAvailable), Limit: 1})
	_, allocCount, _ := s.assetRepo.List(ctx, model.AssetListQuery{Status: string(model.AssetStatusAllocated), Limit: 1})
	_, maintCount, _ := s.assetRepo.List(ctx, model.AssetListQuery{Status: string(model.AssetStatusMaintenance), Limit: 1})
	_, retCount, _ := s.assetRepo.List(ctx, model.AssetListQuery{Status: string(model.AssetStatusRetired), Limit: 1})
	_, totalAssets, _ := s.assetRepo.List(ctx, model.AssetListQuery{Limit: 1})

	assetStatusCounts["AVAILABLE"] = availCount
	assetStatusCounts["ALLOCATED"] = allocCount
	assetStatusCounts["MAINTENANCE"] = maintCount
	assetStatusCounts["RETIRED"] = retCount
	assetStatusCounts["TOTAL"] = totalAssets

	// 4. Alokasi aktif
	_, activeAllocsCount, _ := s.allocRepo.List(ctx, model.AllocationListQuery{Status: "ACTIVE", Limit: 1})

	// 5. Pemeliharaan yang butuh penanganan (OPEN + IN_PROGRESS)
	_, openMaintCount, _ := s.maintRepo.List(ctx, model.MaintenanceListQuery{Status: string(model.MaintenanceStatusOpen), Limit: 1})
	_, inProgMaintCount, _ := s.maintRepo.List(ctx, model.MaintenanceListQuery{Status: string(model.MaintenanceStatusInProgress), Limit: 1})
	pendingMaint := openMaintCount + inProgMaintCount

	// 6. Total presensi
	_, totalAttendance, _ := s.attendanceRepo.List(ctx, model.AttendanceListQuery{Limit: 1})

	// 7. Audit log terbaru (10 event)
	auditLogs, _, _ := s.auditRepo.List(ctx, model.AuditLogListQuery{Limit: 10})
	recentAudit := make([]model.AuditLogResponse, len(auditLogs))
	for i, l := range auditLogs {
		recentAudit[i] = model.AuditLogResponse{
			ID:             l.ID,
			UserID:         l.UserID,
			Action:         l.Action,
			EntityType:     l.EntityType,
			EntityID:       l.EntityID,
			Status:         l.Status,
			IPAddress:      l.IPAddress,
			UserAgent:      l.UserAgent,
			CreatedAt:      l.CreatedAt.UTC().Format(time.RFC3339),
			UserIdentifier: l.UserIdentifier,
			UserFullName:   l.UserFullName,
		}
	}

	return &model.AdminDashboardResponse{
		ActiveSessionsCount:     len(activeSessions),
		ActiveSessions:          activeSessionsRes,
		TotalUsers:              totalUsers,
		TotalStudents:           totalStudents,
		AssetStatusCounts:       assetStatusCounts,
		ActiveAllocationsCount:  activeAllocsCount,
		PendingMaintenanceCount: pendingMaint,
		TotalAttendanceCount:    totalAttendance,
		RecentAuditLogs:         recentAudit,
	}, nil
}
