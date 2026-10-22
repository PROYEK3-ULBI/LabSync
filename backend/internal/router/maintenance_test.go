package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type inMemoryMaintenanceRepo struct {
	maintenances    map[string]*model.Maintenance
	assetRepo       *inMemoryAssetRepo
	simulateDBError bool
}

func newInMemoryMaintenanceRepo(assetRepo *inMemoryAssetRepo) *inMemoryMaintenanceRepo {
	return &inMemoryMaintenanceRepo{
		maintenances: make(map[string]*model.Maintenance),
		assetRepo:    assetRepo,
	}
}

func (r *inMemoryMaintenanceRepo) CreateWithTransaction(ctx context.Context, m *model.Maintenance) error {
	if r.simulateDBError {
		return errors.New("simulated database transaction error")
	}

	asset, ok := r.assetRepo.assets[m.AssetID]
	if !ok {
		return repository.ErrAssetNotFound
	}

	if asset.Status == model.AssetStatusAllocated || r.assetRepo.allocations[m.AssetID] {
		return repository.ErrAssetHasActiveAllocation
	}
	if asset.Status == model.AssetStatusRetired {
		return errors.New("aset berstatus RETIRED tidak dapat diajukan pemeliharaan")
	}

	m.ID = fmt.Sprintf("maint-%d", len(r.maintenances)+1)
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now().UTC()
	if m.Status == "" {
		m.Status = model.MaintenanceStatusOpen
	}
	r.maintenances[m.ID] = m

	asset.Status = model.AssetStatusMaintenance
	asset.UpdatedAt = time.Now().UTC()
	r.assetRepo.maintenances[m.AssetID] = true

	return nil
}

func (r *inMemoryMaintenanceRepo) UpdateWithTransaction(ctx context.Context, m *model.Maintenance, resultAssetStatus *model.AssetStatus) error {
	if r.simulateDBError {
		return errors.New("simulated database transaction error")
	}

	existing, ok := r.maintenances[m.ID]
	if !ok {
		return repository.ErrMaintenanceNotFound
	}

	if !repository.IsValidMaintenanceTransition(existing.Status, m.Status) {
		return repository.ErrInvalidStatusTransition
	}

	asset, ok := r.assetRepo.assets[existing.AssetID]
	if !ok {
		return repository.ErrAssetNotFound
	}

	now := time.Now().UTC()
	if m.Status == model.MaintenanceStatusInProgress && existing.StartedAt == nil {
		existing.StartedAt = &now
	}

	if m.Status == model.MaintenanceStatusCompleted {
		existing.CompletedAt = &now
		existing.Resolution = m.Resolution
		existing.Cost = m.Cost

		target := model.AssetStatusAvailable
		if resultAssetStatus != nil && *resultAssetStatus == model.AssetStatusRetired {
			target = model.AssetStatusRetired
		}
		asset.Status = target
		asset.UpdatedAt = now
		r.assetRepo.maintenances[asset.ID] = false
	} else if m.Status == model.MaintenanceStatusCancelled {
		asset.Status = model.AssetStatusAvailable
		asset.UpdatedAt = now
		r.assetRepo.maintenances[asset.ID] = false
	}

	existing.Status = m.Status
	existing.Priority = m.Priority
	existing.Description = m.Description
	existing.UpdatedAt = now
	r.maintenances[m.ID] = existing

	return nil
}

func (r *inMemoryMaintenanceRepo) GetByID(ctx context.Context, id string) (*model.Maintenance, error) {
	if m, ok := r.maintenances[id]; ok {
		copyM := *m
		return &copyM, nil
	}
	return nil, repository.ErrMaintenanceNotFound
}

func (r *inMemoryMaintenanceRepo) List(ctx context.Context, q model.MaintenanceListQuery) ([]model.Maintenance, int, error) {
	var list []model.Maintenance
	for _, m := range r.maintenances {
		if q.AssetID != "" && m.AssetID != q.AssetID {
			continue
		}
		if q.Status != "" && !strings.EqualFold(string(m.Status), q.Status) {
			continue
		}
		if q.Priority != "" && !strings.EqualFold(string(m.Priority), q.Priority) {
			continue
		}
		list = append(list, *m)
	}
	return list, len(list), nil
}

type maintenanceTestContext struct {
	app       *fiber.App
	userRepo  *inMemoryUserRepo
	assetRepo *inMemoryAssetRepo
	allocRepo *inMemoryAllocationRepo
	maintRepo *inMemoryMaintenanceRepo
	auditRepo *inMemoryAuditRepo
	helper    *AuthServiceTestHelper
}

func setupMaintenanceTestApp() *maintenanceTestContext {
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

	return &maintenanceTestContext{
		app:       app,
		userRepo:  userRepo,
		assetRepo: assetRepo,
		allocRepo: allocRepo,
		maintRepo: maintRepo,
		auditRepo: auditRepo,
		helper:    helper,
	}
}

// 1. Test Admin Create Maintenance Success
func TestCreateMaintenanceSuccess(t *testing.T) {
	tc := setupMaintenanceTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	payload := `{
		"asset_id": "ast-1",
		"issue_type": "HARDWARE",
		"description": "Layar LCD berkedip dan port HDMI longgar",
		"priority": "HIGH"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})

	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201 Created, got %d: %s", resp.StatusCode, string(b))
	}

	var res model.MaintenanceResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	if res.Status != model.MaintenanceStatusOpen {
		t.Errorf("expected initial status OPEN, got %s", res.Status)
	}

	// Status aset harus beralih ke MAINTENANCE
	asset, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if asset.Status != model.AssetStatusMaintenance {
		t.Errorf("expected asset status MAINTENANCE, got %s", asset.Status)
	}
}

// 2. Test Create Maintenance Blocked on Allocated Asset
func TestCreateMaintenanceAllocatedAssetBlocked(t *testing.T) {
	tc := setupMaintenanceTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// ast-2 sedang dialokasikan
	payload := `{
		"asset_id": "ast-2",
		"description": "Kendala kipas laptop berbunyi bising",
		"priority": "MEDIUM"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict when maintaining allocated asset, got %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(b, &errResp)
	if errResp.Error.Code != "ASSET_ALLOCATED" {
		t.Errorf("expected code ASSET_ALLOCATED, got %s", errResp.Error.Code)
	}
}

// 3. Test Maintenance State Machine Transitions & Resolution
func TestMaintenanceStateMachineTransitions(t *testing.T) {
	tc := setupMaintenanceTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Buat tiket pemeliharaan -> OPEN
	createPayload := `{"asset_id":"ast-1","description":"Pembersihan debu dan pasta termal"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(createPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respCreate, _ := tc.app.Test(reqCreate)
	var created model.MaintenanceResponse
	bC, _ := io.ReadAll(respCreate.Body)
	_ = json.Unmarshal(bC, &created)

	// b. Transisi Ilegal: OPEN langsung ke COMPLETED (ditolak)
	reqIllegal := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"COMPLETED","resolution":"Selesai"}`))
	reqIllegal.Header.Set("Content-Type", "application/json")
	reqIllegal.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respIllegal, _ := tc.app.Test(reqIllegal)
	if respIllegal.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for illegal transition from OPEN to COMPLETED, got %d", respIllegal.StatusCode)
	}

	// c. Transisi Legal: OPEN -> IN_PROGRESS
	reqProgress := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"IN_PROGRESS"}`))
	reqProgress.Header.Set("Content-Type", "application/json")
	reqProgress.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respProgress, _ := tc.app.Test(reqProgress)
	if respProgress.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on transition to IN_PROGRESS, got %d", respProgress.StatusCode)
	}

	// d. Selesai tanpa resolusi -> Ditolak 400 Bad Request
	reqNoRes := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"COMPLETED"}`))
	reqNoRes.Header.Set("Content-Type", "application/json")
	reqNoRes.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respNoRes, _ := tc.app.Test(reqNoRes)
	if respNoRes.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when completing without resolution, got %d", respNoRes.StatusCode)
	}

	// e. Transisi Legal: IN_PROGRESS -> COMPLETED (dengan resolusi)
	reqComplete := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"COMPLETED","resolution":"Thermal paste diganti dan kipas dibersihkan"}`))
	reqComplete.Header.Set("Content-Type", "application/json")
	reqComplete.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respComplete, _ := tc.app.Test(reqComplete)
	if respComplete.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on valid completion, got %d", respComplete.StatusCode)
	}

	// Pastikan status aset kembali ke AVAILABLE
	assetAfter, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if assetAfter.Status != model.AssetStatusAvailable {
		t.Errorf("expected asset status AVAILABLE after completion, got %s", assetAfter.Status)
	}

	// f. Terminal State: COMPLETED tidak boleh diubah lagi
	reqTerminal := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"IN_PROGRESS"}`))
	reqTerminal.Header.Set("Content-Type", "application/json")
	reqTerminal.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respTerminal, _ := tc.app.Test(reqTerminal)
	if respTerminal.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when mutating terminal COMPLETED maintenance, got %d", respTerminal.StatusCode)
	}
}

// 4. Test Complete Maintenance Retiring Asset
func TestCompleteMaintenanceRetiringAsset(t *testing.T) {
	tc := setupMaintenanceTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Buat tiket
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(`{"asset_id":"ast-1","description":"Motherboard terbakar"}`))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respCreate, _ := tc.app.Test(reqCreate)
	var created model.MaintenanceResponse
	bC, _ := io.ReadAll(respCreate.Body)
	_ = json.Unmarshal(bC, &created)

	// IN_PROGRESS
	reqProgress := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(`{"status":"IN_PROGRESS"}`))
	reqProgress.Header.Set("Content-Type", "application/json")
	reqProgress.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	_, _ = tc.app.Test(reqProgress)

	// COMPLETED dengan result_status RETIRED
	completePayload := `{"status":"COMPLETED","resolution":"Kerusakan fatal tidak dapat diperbaiki","result_status":"RETIRED"}`
	reqComplete := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/maintenance/%s", created.ID), bytes.NewBufferString(completePayload))
	reqComplete.Header.Set("Content-Type", "application/json")
	reqComplete.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respComplete, _ := tc.app.Test(reqComplete)

	if respComplete.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on retiring maintenance, got %d", respComplete.StatusCode)
	}

	assetAfter, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if assetAfter.Status != model.AssetStatusRetired {
		t.Errorf("expected asset status RETIRED, got %s", assetAfter.Status)
	}
}

// 5. Test Maintenance Permissions (Admin Only)
func TestMaintenancePermissions(t *testing.T) {
	tc := setupMaintenanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// a. Student POST /api/v1/maintenance -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(`{"asset_id":"ast-1","description":"Rusak"}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls POST /api/v1/maintenance, got %d", resp1.StatusCode)
	}

	// b. Student GET /api/v1/maintenance -> 403 Forbidden
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/maintenance", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls GET /api/v1/maintenance, got %d", resp2.StatusCode)
	}
}

// 6. Test Audit Logs Append-Only & Inspection
func TestAuditLogsAppendOnlyAndInspection(t *testing.T) {
	tc := setupMaintenanceTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// Buat maintenance untuk men-trigger event audit
	reqMaint := httptest.NewRequest(http.MethodPost, "/api/v1/maintenance", bytes.NewBufferString(`{"asset_id":"ast-1","description":"Kendala keyboard"}`))
	reqMaint.Header.Set("Content-Type", "application/json")
	reqMaint.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	_, _ = tc.app.Test(reqMaint)

	// a. Student GET /api/v1/audit-logs -> 403 Forbidden
	reqStudent := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs", nil)
	reqStudent.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	respStudent, _ := tc.app.Test(reqStudent)
	if respStudent.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student inspects audit logs, got %d", respStudent.StatusCode)
	}

	// b. Admin GET /api/v1/audit-logs -> 200 OK
	reqAdmin := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs", nil)
	reqAdmin.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAdmin, _ := tc.app.Test(reqAdmin)
	if respAdmin.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK when admin reads audit logs, got %d", respAdmin.StatusCode)
	}

	var auditRes model.AuditLogListResponse
	bAdmin, _ := io.ReadAll(respAdmin.Body)
	_ = json.Unmarshal(bAdmin, &auditRes)

	if len(auditRes.Data) == 0 {
		t.Errorf("expected audit logs to contain events")
	}

	// Pastikan metadata audit log tidak mengandung rahasia (password, token, embedding)
	for _, l := range auditRes.Data {
		metaStr, _ := json.Marshal(l.Metadata)
		s := string(metaStr)
		if strings.Contains(s, "password") || strings.Contains(s, "token") || strings.Contains(s, "embedding") {
			t.Errorf("audit log leaked sensitive data: %s", s)
		}
	}
}
