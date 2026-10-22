package router

import (
	"bytes"
	"context"
	"encoding/json"
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

type inMemoryAssetRepo struct {
	assets       map[string]*model.Asset
	allocations  map[string]bool
	maintenances map[string]bool
}

func newInMemoryAssetRepo() *inMemoryAssetRepo {
	repo := &inMemoryAssetRepo{
		assets:       make(map[string]*model.Asset),
		allocations:  make(map[string]bool),
		maintenances: make(map[string]bool),
	}

	b1 := "Dell"
	m1 := "Latitude 5490"
	s1 := "SN-DELL-001"
	repo.assets["ast-1"] = &model.Asset{
		ID:           "ast-1",
		AssetCode:    "LAB-LT-001",
		AssetName:    "Dell Latitude 5490 Core i7",
		Category:     "Laptop",
		Brand:        &b1,
		Model:        &m1,
		SerialNumber: &s1,
		Location:     "Lab Komputer 1",
		Status:       model.AssetStatusAvailable,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	s2 := "SN-DELL-002"
	repo.assets["ast-2"] = &model.Asset{
		ID:           "ast-2",
		AssetCode:    "LAB-LT-002",
		AssetName:    "Dell Latitude 5490 Core i5",
		Category:     "Laptop",
		Brand:        &b1,
		Model:        &m1,
		SerialNumber: &s2,
		Location:     "Lab Komputer 1",
		Status:       model.AssetStatusAllocated,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	repo.allocations["ast-2"] = true

	s3 := "SN-DELL-003"
	repo.assets["ast-3"] = &model.Asset{
		ID:           "ast-3",
		AssetCode:    "LAB-LT-003",
		AssetName:    "Dell Latitude 5490 Rusak",
		Category:     "Laptop",
		Brand:        &b1,
		Model:        &m1,
		SerialNumber: &s3,
		Location:     "Lab Komputer 2",
		Status:       model.AssetStatusMaintenance,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	repo.maintenances["ast-3"] = true

	return repo
}

func (r *inMemoryAssetRepo) GetByID(ctx context.Context, id string) (*model.Asset, error) {
	if a, ok := r.assets[id]; ok {
		copyA := *a
		return &copyA, nil
	}
	return nil, repository.ErrAssetNotFound
}

func (r *inMemoryAssetRepo) GetByCode(ctx context.Context, code string) (*model.Asset, error) {
	for _, a := range r.assets {
		if strings.EqualFold(a.AssetCode, code) {
			copyA := *a
			return &copyA, nil
		}
	}
	return nil, repository.ErrAssetNotFound
}

func (r *inMemoryAssetRepo) GetBySerialNumber(ctx context.Context, sn string) (*model.Asset, error) {
	for _, a := range r.assets {
		if a.SerialNumber != nil && strings.EqualFold(*a.SerialNumber, sn) {
			copyA := *a
			return &copyA, nil
		}
	}
	return nil, repository.ErrAssetNotFound
}

func (r *inMemoryAssetRepo) Create(ctx context.Context, a *model.Asset) error {
	for _, existing := range r.assets {
		if strings.EqualFold(existing.AssetCode, a.AssetCode) {
			return repository.ErrDuplicateAssetCode
		}
		if a.SerialNumber != nil && existing.SerialNumber != nil && strings.EqualFold(*existing.SerialNumber, *a.SerialNumber) {
			return repository.ErrDuplicateSerialNumber
		}
	}
	a.ID = fmt.Sprintf("ast-%d", len(r.assets)+1)
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
	r.assets[a.ID] = a
	return nil
}

func (r *inMemoryAssetRepo) Update(ctx context.Context, a *model.Asset) error {
	if _, ok := r.assets[a.ID]; !ok {
		return repository.ErrAssetNotFound
	}
	for _, existing := range r.assets {
		if existing.ID != a.ID {
			if strings.EqualFold(existing.AssetCode, a.AssetCode) {
				return repository.ErrDuplicateAssetCode
			}
			if a.SerialNumber != nil && existing.SerialNumber != nil && strings.EqualFold(*existing.SerialNumber, *a.SerialNumber) {
				return repository.ErrDuplicateSerialNumber
			}
		}
	}
	a.UpdatedAt = time.Now().UTC()
	r.assets[a.ID] = a
	return nil
}

func (r *inMemoryAssetRepo) UpdateStatus(ctx context.Context, id string, status model.AssetStatus) error {
	a, ok := r.assets[id]
	if !ok {
		return repository.ErrAssetNotFound
	}
	a.Status = status
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *inMemoryAssetRepo) Delete(ctx context.Context, id string) error {
	if _, ok := r.assets[id]; !ok {
		return repository.ErrAssetNotFound
	}
	delete(r.assets, id)
	return nil
}

func (r *inMemoryAssetRepo) List(ctx context.Context, q model.AssetListQuery) ([]model.Asset, int, error) {
	var list []model.Asset
	for _, a := range r.assets {
		if q.Status != "" && !strings.EqualFold(string(a.Status), q.Status) {
			continue
		}
		if q.Category != "" && !strings.EqualFold(a.Category, q.Category) {
			continue
		}
		if q.Search != "" {
			s := strings.ToLower(q.Search)
			match := strings.Contains(strings.ToLower(a.AssetCode), s) ||
				strings.Contains(strings.ToLower(a.AssetName), s) ||
				strings.Contains(strings.ToLower(a.Location), s)
			if a.Brand != nil && strings.Contains(strings.ToLower(*a.Brand), s) {
				match = true
			}
			if a.SerialNumber != nil && strings.Contains(strings.ToLower(*a.SerialNumber), s) {
				match = true
			}
			if !match {
				continue
			}
		}
		list = append(list, *a)
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}

	start := (q.Page - 1) * q.Limit
	if start >= len(list) {
		return []model.Asset{}, len(list), nil
	}
	end := start + q.Limit
	if end > len(list) {
		end = len(list)
	}

	return list[start:end], len(list), nil
}

func (r *inMemoryAssetRepo) HasActiveAllocation(ctx context.Context, assetID string) (bool, error) {
	return r.allocations[assetID], nil
}

func (r *inMemoryAssetRepo) HasActiveMaintenance(ctx context.Context, assetID string) (bool, error) {
	return r.maintenances[assetID], nil
}

type assetTestContext struct {
	app       *fiber.App
	assetRepo *inMemoryAssetRepo
	helper    *AuthServiceTestHelper
}

func setupAssetTestApp() *assetTestContext {
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
	faceClient := &mockFaceClient{}

	maintRepo := newInMemoryMaintenanceRepo(assetRepo)

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

	return &assetTestContext{
		app:       app,
		assetRepo: assetRepo,
		helper:    helper,
	}
}

// 1. Test Admin Create Asset Success
func TestAdminCreateAssetSuccess(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	payload := `{
		"asset_code": "LAB-PC-010",
		"asset_name": "PC Workstation Lenovo ThinkCentre",
		"category": "PC",
		"brand": "Lenovo",
		"model": "ThinkCentre M70q",
		"serial_number": "SN-LENOVO-010",
		"location": "Lab Komputer 1",
		"notes": "Spesifikasi Core i7 RAM 16GB"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(payload))
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

	var res model.AssetResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	if res.AssetCode != "LAB-PC-010" {
		t.Errorf("expected asset_code LAB-PC-010, got %s", res.AssetCode)
	}
	if res.Status != model.AssetStatusAvailable {
		t.Errorf("expected default status AVAILABLE, got %s", res.Status)
	}
}

// 2. Test Create Asset Validation
func TestCreateAssetValidation(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Kode aset kosong/terlalu pendek
	payload1 := `{"asset_code":"A","asset_name":"Nama Aset"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(payload1))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for short asset code, got %d", resp1.StatusCode)
	}

	// b. Initial status ALLOCATED (tidak boleh dibuat langsung ALLOCATED)
	payload2 := `{"asset_code":"LAB-PC-099","asset_name":"PC Baru","status":"ALLOCATED"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(payload2))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for direct ALLOCATED on create, got %d", resp2.StatusCode)
	}
}

// 3. Test Duplicate Asset Code & Serial Number
func TestDuplicateAssetCodeAndSerialNumber(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Duplikat asset_code (LAB-LT-001 sudah ada)
	payloadDupCode := `{"asset_code":"LAB-LT-001","asset_name":"Laptop Duplikat"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(payloadDupCode))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate asset code, got %d", resp1.StatusCode)
	}

	// b. Duplikat serial_number (SN-DELL-001 sudah ada)
	payloadDupSN := `{"asset_code":"LAB-LT-999","asset_name":"Laptop Duplikat SN","serial_number":"SN-DELL-001"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(payloadDupSN))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate serial number, got %d", resp2.StatusCode)
	}
}

// 4. Test Student Cannot Manage Assets (Read-only Allowed)
func TestStudentCannotManageAssets(t *testing.T) {
	tc := setupAssetTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// a. Student POST /api/v1/assets -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewBufferString(`{"asset_code":"STU-1","asset_name":"Laptop"}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls POST /api/v1/assets, got %d", resp1.StatusCode)
	}

	// b. Student PUT /api/v1/assets/:id -> 403 Forbidden
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/assets/ast-1", bytes.NewBufferString(`{"asset_name":"Edit"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls PUT /api/v1/assets/:id, got %d", resp2.StatusCode)
	}

	// c. Student DELETE /api/v1/assets/:id -> 403 Forbidden
	req3 := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/ast-1", nil)
	req3.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp3, _ := tc.app.Test(req3)
	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls DELETE /api/v1/assets/:id, got %d", resp3.StatusCode)
	}

	// d. Student GET /api/v1/assets -> 200 OK (Katalog aset dapat dilihat mahasiswa)
	req4 := httptest.NewRequest(http.MethodGet, "/api/v1/assets", nil)
	req4.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp4, _ := tc.app.Test(req4)
	if resp4.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student reads assets catalog, got %d", resp4.StatusCode)
	}

	// e. Student GET /api/v1/assets/:id -> 200 OK
	req5 := httptest.NewRequest(http.MethodGet, "/api/v1/assets/ast-1", nil)
	req5.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp5, _ := tc.app.Test(req5)
	if resp5.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student reads asset detail, got %d", resp5.StatusCode)
	}
}

// 5. Test Status Transitions
func TestAssetStatusTransitions(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. AVAILABLE -> MAINTENANCE (Valid: dilaporkan rusak)
	req1 := httptest.NewRequest(http.MethodPatch, "/api/v1/assets/ast-1/status", bytes.NewBufferString(`{"status":"MAINTENANCE"}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on transition to MAINTENANCE, got %d", resp1.StatusCode)
	}

	// b. MAINTENANCE -> RETIRED (Valid: rusak permanen)
	req2 := httptest.NewRequest(http.MethodPatch, "/api/v1/assets/ast-1/status", bytes.NewBufferString(`{"status":"RETIRED"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on transition to RETIRED, got %d", resp2.StatusCode)
	}

	// c. RETIRED -> AVAILABLE (Illegal: status RETIRED adalah terminal)
	req3 := httptest.NewRequest(http.MethodPatch, "/api/v1/assets/ast-1/status", bytes.NewBufferString(`{"status":"AVAILABLE"}`))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp3, _ := tc.app.Test(req3)
	if resp3.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for illegal transition from RETIRED to AVAILABLE, got %d", resp3.StatusCode)
	}
}

// 6. Test Guards on Active Allocation & Maintenance
func TestAssetGuardsOnActiveAllocationAndMaintenance(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. ast-2 sedang dialokasikan (allocations["ast-2"] = true)
	// Coba ubah status secara manual -> 409 Conflict (ASSET_ALLOCATED)
	req1 := httptest.NewRequest(http.MethodPatch, "/api/v1/assets/ast-2/status", bytes.NewBufferString(`{"status":"AVAILABLE"}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict when editing status of allocated asset, got %d", resp1.StatusCode)
	}

	// Coba hapus aset yang sedang dialokasikan -> 409 Conflict
	reqDelAlloc := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/ast-2", nil)
	reqDelAlloc.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respDelAlloc, _ := tc.app.Test(reqDelAlloc)
	if respDelAlloc.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict when deleting allocated asset, got %d", respDelAlloc.StatusCode)
	}

	// b. ast-3 sedang dalam maintenance aktif (maintenances["ast-3"] = true)
	// Coba ubah status secara manual -> 409 Conflict (ASSET_IN_MAINTENANCE)
	req2 := httptest.NewRequest(http.MethodPatch, "/api/v1/assets/ast-3/status", bytes.NewBufferString(`{"status":"AVAILABLE"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict when editing status of asset in maintenance, got %d", resp2.StatusCode)
	}

	// Coba hapus aset yang sedang dalam maintenance -> 409 Conflict
	reqDelMaint := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/ast-3", nil)
	reqDelMaint.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respDelMaint, _ := tc.app.Test(reqDelMaint)
	if respDelMaint.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict when deleting asset in maintenance, got %d", respDelMaint.StatusCode)
	}
}

// 7. Test Asset Not Found
func TestAssetNotFound(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/assets/unknown-id", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for unknown asset, got %d", resp1.StatusCode)
	}

	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/unknown-id", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 on deleting unknown asset, got %d", resp2.StatusCode)
	}
}

// 8. Test Asset List Filtering and Pagination
func TestAssetListFilteringAndPagination(t *testing.T) {
	tc := setupAssetTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Filter status AVAILABLE
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/assets?status=AVAILABLE", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	var res1 model.AssetListResponse
	b1, _ := io.ReadAll(resp1.Body)
	_ = json.Unmarshal(b1, &res1)

	if len(res1.Data) != 1 || res1.Data[0].AssetCode != "LAB-LT-001" {
		t.Errorf("expected 1 AVAILABLE asset, got %d", len(res1.Data))
	}

	// b. Paging limit=2
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/assets?page=1&limit=2", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	var res2 model.AssetListResponse
	b2, _ := io.ReadAll(resp2.Body)
	_ = json.Unmarshal(b2, &res2)

	if len(res2.Data) != 2 {
		t.Errorf("expected 2 assets in page 1, got %d", len(res2.Data))
	}
	if res2.Meta.Total != 3 {
		t.Errorf("expected total 3, got %d", res2.Meta.Total)
	}
}
