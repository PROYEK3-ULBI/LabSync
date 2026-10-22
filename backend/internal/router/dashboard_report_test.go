package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/config"
	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
	"labsync-backend/internal/service"
)

type dashReportTestContext struct {
	app            *fiber.App
	userRepo       *inMemoryUserRepo
	sessionRepo    *inMemorySessionRepo
	assetRepo      *inMemoryAssetRepo
	allocRepo      *inMemoryAllocationRepo
	maintRepo      *inMemoryMaintenanceRepo
	attendanceRepo *inMemoryAttendanceRepo
	auditRepo      *inMemoryAuditRepo
	helper         *AuthServiceTestHelper
}

func setupDashReportTestApp() *dashReportTestContext {
	cfg := &config.Config{
		AppEnv:                 "development",
		Port:                   "8000",
		CORSAllowOrigins:       "http://localhost:8443",
		BodyLimitMB:            10,
		RateLimitMax:           100,
		JWTSecret:              "test-super-secret-key-1234567890",
		JWTExpiryHours:         24,
		CookieSecure:           false,
		CookieSameSite:         "Lax",
		FaceServiceURL:         "http://localhost:8001",
		InternalAPIKey:         "labsync-internal-secret-key",
		FaceServiceTimeoutSec:  10,
		BiometricEncryptionKey: "01234567890123456789012345678901",
	}

	userRepo := newInMemoryUserRepo()
	auditRepo := &inMemoryAuditRepo{}
	healthRepo := repository.NewHealthRepository(nil)
	sessionRepo := newInMemorySessionRepo()
	faceRepo := newInMemoryFaceRepo()
	attendanceRepo := newInMemoryAttendanceRepo()
	assetRepo := newInMemoryAssetRepo()
	allocRepo := newInMemoryAllocationRepo(assetRepo)
	maintRepo := newInMemoryMaintenanceRepo(assetRepo)
	faceClient := &mockFaceClient{}

	authService := service.NewAuthService(userRepo, auditRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	userService := service.NewUserService(userRepo, auditRepo)
	sessionService := service.NewSessionService(sessionRepo, auditRepo)
	faceService := service.NewFaceService(userRepo, faceRepo, auditRepo, faceClient, cfg.BiometricEncryptionKey)
	attendanceService := service.NewAttendanceService(attendanceRepo, sessionRepo, faceRepo, userRepo, auditRepo, faceClient, cfg.BiometricEncryptionKey)
	assetService := service.NewAssetService(assetRepo, auditRepo)
	allocService := service.NewAllocationService(allocRepo, assetRepo, userRepo, sessionRepo, attendanceRepo, auditRepo)
	maintService := service.NewMaintenanceService(maintRepo, assetRepo, allocRepo, auditRepo)
	auditService := service.NewAuditService(auditRepo)
	dashService := service.NewDashboardService(userRepo, sessionRepo, attendanceRepo, assetRepo, allocRepo, maintRepo, auditRepo)
	reportService := service.NewReportService(attendanceRepo, assetRepo, allocRepo, maintRepo)
	healthService := service.NewHealthService(healthRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.CentralErrorHandler,
	})

	SetupWithServices(app, cfg, authService, userService, sessionService, faceService, attendanceService, assetService, allocService, maintService, auditService, dashService, reportService, healthService)

	helper := &AuthServiceTestHelper{
		authService:    authService,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		sessionRepo:    sessionRepo,
		sessionService: sessionService,
	}

	// Seed data pendukung:
	// 1. Presensi mahasiswa student-uuid-1 pada sess-1
	_ = attendanceRepo.Create(context.Background(), &model.Attendance{
		UserID:             "student-uuid-1",
		SessionID:          "sess-1",
		CheckedInAt:        time.Now().UTC(),
		VerificationStatus: "VERIFIED_FACE",
		SessionCode:        "PWB-01",
		SessionTitle:       "Praktikum Pemrograman Web",
		SessionDate:        "2026-09-18",
		UserIdentifier:     "TI240001",
		UserFullName:       "Budi Mahasiswa",
	})

	// 2. Alokasi ast-2 kepada student-uuid-1
	_ = allocRepo.CreateWithTransaction(context.Background(), &model.Allocation{
		AssetID:         "ast-1",
		UserID:          "student-uuid-1",
		SessionID:       "sess-1",
		AllocatedAt:     time.Now().UTC(),
		ConditionBefore: "GOOD",
		UserIdentifier:  "TI240001",
		UserFullName:    "Budi Mahasiswa",
		AssetCode:       "LAB-LT-001",
		AssetName:       "Dell Latitude 5490",
		SessionCode:     "PWB-01",
	})

	// 3. Tiket pemeliharaan ast-3
	_ = maintRepo.CreateWithTransaction(context.Background(), &model.Maintenance{
		AssetID:     "ast-3",
		ReportedBy:  "admin-uuid-1",
		IssueType:   "HARDWARE",
		Description: "Baterai drop",
		Priority:    model.MaintenancePriorityHigh,
	})

	// 4. Audit event
	_ = auditRepo.Create(context.Background(), &model.AuditLog{
		Action:     "INITIAL_SEED",
		EntityType: "SYSTEM",
		EntityID:   "00000000-0000-0000-0000-000000000000",
	})

	return &dashReportTestContext{
		app:            app,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		assetRepo:      assetRepo,
		allocRepo:      allocRepo,
		maintRepo:      maintRepo,
		attendanceRepo: attendanceRepo,
		auditRepo:      auditRepo,
		helper:         helper,
	}
}

// 1. Test Student Dashboard
func TestStudentDashboard(t *testing.T) {
	tc := setupDashReportTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/student", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: studentToken})

	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(b))
	}

	var res model.StudentDashboardResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	if res.User.Identifier != "TI240001" {
		t.Errorf("expected student NIM TI240001, got %s", res.User.Identifier)
	}
	if res.TotalAttended != 1 {
		t.Errorf("expected total attended 1, got %d", res.TotalAttended)
	}
	if res.ActiveSession == nil || res.ActiveSession.Code != "PWB-01" {
		t.Errorf("expected active session PWB-01")
	}
	if res.CurrentAllocation == nil || res.CurrentAllocation.AssetID != "ast-1" {
		t.Errorf("expected current allocation for ast-1")
	}
}

// 2. Test Admin Dashboard
func TestAdminDashboard(t *testing.T) {
	tc := setupDashReportTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/admin", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})

	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(b))
	}

	var res model.AdminDashboardResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	// Verifikasi konsistensi agregasi
	if res.ActiveSessionsCount != len(res.ActiveSessions) {
		t.Errorf("active sessions count %d does not match length %d", res.ActiveSessionsCount, len(res.ActiveSessions))
	}
	if res.ActiveAllocationsCount < 1 {
		t.Errorf("expected at least 1 active allocation, got %d", res.ActiveAllocationsCount)
	}
	if res.PendingMaintenanceCount < 1 {
		t.Errorf("expected at least 1 pending maintenance, got %d", res.PendingMaintenanceCount)
	}
	if len(res.RecentAuditLogs) == 0 {
		t.Errorf("expected non-empty recent audit logs")
	}
}

// 3. Test Dashboard Authorization
func TestDashboardAuthorization(t *testing.T) {
	tc := setupDashReportTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Student calls admin dashboard -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/admin", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls admin dashboard, got %d", resp1.StatusCode)
	}

	// b. Admin calls student dashboard -> 403 Forbidden (student role required)
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/student", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when admin calls student dashboard, got %d", resp2.StatusCode)
	}
}

// 4. Test Reports Authorization (Admin Only)
func TestReportsAuthorization(t *testing.T) {
	tc := setupDashReportTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	endpoints := []string{
		"/api/v1/reports/attendance",
		"/api/v1/reports/assets",
		"/api/v1/reports/allocations",
		"/api/v1/reports/maintenance",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
		resp, _ := tc.app.Test(req)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403 when student calls %s, got %d", ep, resp.StatusCode)
		}
	}
}

// 5. Test Reports Date Range Validation
func TestReportsDateRangeValidation(t *testing.T) {
	tc := setupDashReportTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Tanggal awal setelah tanggal akhir (start > end) -> 400 Bad Request
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/reports/attendance?start_date=2026-10-20&end_date=2026-10-10", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)

	if resp1.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request when start_date > end_date, got %d", resp1.StatusCode)
	}
	b1, _ := io.ReadAll(resp1.Body)
	var errResp1 model.ErrorResponse
	_ = json.Unmarshal(b1, &errResp1)
	if !strings.Contains(errResp1.Error.Message, "setelah") {
		t.Errorf("expected date range error message, got: %s", errResp1.Error.Message)
	}

	// b. Format tanggal salah -> 400 Bad Request
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/reports/attendance?start_date=20-10-2026", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on invalid date format, got %d", resp2.StatusCode)
	}
}

// 6. Test Reports JSON and CSV Export
func TestReportsJSONAndCSVExport(t *testing.T) {
	tc := setupDashReportTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Attendance Report (JSON)
	reqAttJSON := httptest.NewRequest(http.MethodGet, "/api/v1/reports/attendance", nil)
	reqAttJSON.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAttJSON, _ := tc.app.Test(reqAttJSON)
	if respAttJSON.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on attendance report json, got %d", respAttJSON.StatusCode)
	}
	var attRes model.AttendanceReportResponse
	bAJ, _ := io.ReadAll(respAttJSON.Body)
	_ = json.Unmarshal(bAJ, &attRes)
	if attRes.Summary.TotalAttendance < 1 {
		t.Errorf("expected total attendance >= 1, got %d", attRes.Summary.TotalAttendance)
	}

	// b. Attendance Report (CSV)
	reqAttCSV := httptest.NewRequest(http.MethodGet, "/api/v1/reports/attendance?format=csv", nil)
	reqAttCSV.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAttCSV, _ := tc.app.Test(reqAttCSV)
	if respAttCSV.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on attendance report csv, got %d", respAttCSV.StatusCode)
	}
	if !strings.Contains(respAttCSV.Header.Get("Content-Type"), "text/csv") {
		t.Errorf("expected Content-Type text/csv, got %s", respAttCSV.Header.Get("Content-Type"))
	}
	bAC, _ := io.ReadAll(respAttCSV.Body)
	if !strings.Contains(string(bAC), "TI240001") {
		t.Errorf("expected student NIM in attendance CSV: %s", string(bAC))
	}

	// c. Assets Report (CSV)
	reqAssetCSV := httptest.NewRequest(http.MethodGet, "/api/v1/reports/assets?format=csv", nil)
	reqAssetCSV.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAssetCSV, _ := tc.app.Test(reqAssetCSV)
	if respAssetCSV.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK on assets report csv, got %d", respAssetCSV.StatusCode)
	}

	// d. Allocations Report (CSV)
	reqAllocCSV := httptest.NewRequest(http.MethodGet, "/api/v1/reports/allocations?format=csv", nil)
	reqAllocCSV.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAllocCSV, _ := tc.app.Test(reqAllocCSV)
	if respAllocCSV.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK on allocations report csv, got %d", respAllocCSV.StatusCode)
	}

	// e. Maintenance Report (CSV)
	reqMaintCSV := httptest.NewRequest(http.MethodGet, "/api/v1/reports/maintenance?format=csv", nil)
	reqMaintCSV.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respMaintCSV, _ := tc.app.Test(reqMaintCSV)
	if respMaintCSV.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK on maintenance report csv, got %d", respMaintCSV.StatusCode)
	}
}

// 7. Test Reports Empty States
func TestReportsEmptyStates(t *testing.T) {
	tc := setupDashReportTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Filter session_id yang tidak ada -> harus 200 OK dengan slice kosong, bukan crash
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/attendance?session_id=non-existent-session-id", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on empty report state, got %d", resp.StatusCode)
	}

	var res model.AttendanceReportResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	if res.Summary.TotalAttendance != 0 {
		t.Errorf("expected 0 total attendance, got %d", res.Summary.TotalAttendance)
	}
	if len(res.Data) != 0 {
		t.Errorf("expected empty data array, got %d", len(res.Data))
	}
}
