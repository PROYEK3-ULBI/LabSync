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

type inMemoryAllocationRepo struct {
	allocations     map[string]*model.Allocation
	assetRepo       *inMemoryAssetRepo
	simulateDBError bool
}

func newInMemoryAllocationRepo(assetRepo *inMemoryAssetRepo) *inMemoryAllocationRepo {
	return &inMemoryAllocationRepo{
		allocations: make(map[string]*model.Allocation),
		assetRepo:   assetRepo,
	}
}

func (r *inMemoryAllocationRepo) CreateWithTransaction(ctx context.Context, a *model.Allocation) error {
	if r.simulateDBError {
		return errors.New("simulated database transaction error")
	}

	// 1. Cek status aset
	asset, ok := r.assetRepo.assets[a.AssetID]
	if !ok {
		return repository.ErrAssetNotFound
	}
	if asset.Status != model.AssetStatusAvailable {
		return repository.ErrAssetNotAvailable
	}

	// 2. Cek apakah ada alokasi aktif untuk aset ini
	for _, existing := range r.allocations {
		if existing.AssetID == a.AssetID && existing.ReturnedAt == nil {
			return repository.ErrDoubleAllocation
		}
	}

	// 3. Simpan alokasi dan perbarui status aset secara atomik
	a.ID = fmt.Sprintf("alloc-%d", len(r.allocations)+1)
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
	r.allocations[a.ID] = a

	asset.Status = model.AssetStatusAllocated
	asset.UpdatedAt = time.Now().UTC()
	r.assetRepo.allocations[a.AssetID] = true

	return nil
}

func (r *inMemoryAllocationRepo) ReturnWithTransaction(
	ctx context.Context,
	allocationID string,
	returnedAt time.Time,
	conditionAfter string,
	notes *string,
) (*model.Allocation, error) {
	if r.simulateDBError {
		return nil, errors.New("simulated database transaction error")
	}

	alloc, ok := r.allocations[allocationID]
	if !ok {
		return nil, repository.ErrAllocationNotFound
	}
	if alloc.ReturnedAt != nil {
		return nil, repository.ErrAlreadyReturned
	}

	asset, ok := r.assetRepo.assets[alloc.AssetID]
	if !ok {
		return nil, repository.ErrAssetNotFound
	}

	// Update alokasi
	alloc.ReturnedAt = &returnedAt
	alloc.ConditionAfter = &conditionAfter
	if notes != nil {
		alloc.Notes = notes
	}
	alloc.UpdatedAt = time.Now().UTC()

	// Update status aset
	if strings.EqualFold(conditionAfter, "DAMAGED") {
		asset.Status = model.AssetStatusMaintenance
		r.assetRepo.maintenances[asset.ID] = true
	} else {
		asset.Status = model.AssetStatusAvailable
	}
	r.assetRepo.allocations[asset.ID] = false
	asset.UpdatedAt = time.Now().UTC()

	return alloc, nil
}

func (r *inMemoryAllocationRepo) GetByID(ctx context.Context, id string) (*model.Allocation, error) {
	if a, ok := r.allocations[id]; ok {
		copyA := *a
		return &copyA, nil
	}
	return nil, repository.ErrAllocationNotFound
}

func (r *inMemoryAllocationRepo) List(ctx context.Context, q model.AllocationListQuery) ([]model.Allocation, int, error) {
	var list []model.Allocation
	for _, a := range r.allocations {
		if q.AssetID != "" && a.AssetID != q.AssetID {
			continue
		}
		if q.UserID != "" && a.UserID != q.UserID {
			continue
		}
		if q.SessionID != "" && a.SessionID != q.SessionID {
			continue
		}
		if strings.EqualFold(q.Status, "ACTIVE") && a.ReturnedAt != nil {
			continue
		}
		if strings.EqualFold(q.Status, "RETURNED") && a.ReturnedAt == nil {
			continue
		}
		list = append(list, *a)
	}
	return list, len(list), nil
}

func (r *inMemoryAllocationRepo) GetActiveAllocationByAssetID(ctx context.Context, assetID string) (*model.Allocation, error) {
	for _, a := range r.allocations {
		if a.AssetID == assetID && a.ReturnedAt == nil {
			copyA := *a
			return &copyA, nil
		}
	}
	return nil, repository.ErrAllocationNotFound
}

func (r *inMemoryAllocationRepo) GetAssetTimeline(ctx context.Context, assetID string) ([]model.AssetTimelineEvent, error) {
	var timeline []model.AssetTimelineEvent
	for _, a := range r.allocations {
		if a.AssetID != assetID {
			continue
		}
		userID := a.UserID
		userName := a.UserFullName
		if userName == "" {
			userName = "Mahasiswa"
		}

		timeline = append(timeline, model.AssetTimelineEvent{
			ID:          fmt.Sprintf("alloc-%s", a.ID),
			EventType:   "ALLOCATION",
			Title:       "Aset Dipinjam",
			Description: fmt.Sprintf("Aset dialokasikan kepada %s", userName),
			Timestamp:   a.AllocatedAt.UTC().Format(time.RFC3339),
			ActorID:     &userID,
			ActorName:   &userName,
			Details: map[string]interface{}{
				"condition_before": a.ConditionBefore,
			},
		})

		if a.ReturnedAt != nil {
			condAfter := "GOOD"
			if a.ConditionAfter != nil {
				condAfter = *a.ConditionAfter
			}
			timeline = append(timeline, model.AssetTimelineEvent{
				ID:          fmt.Sprintf("return-%s", a.ID),
				EventType:   "RETURN",
				Title:       "Aset Dikembalikan",
				Description: fmt.Sprintf("Aset dikembalikan dengan kondisi %s", condAfter),
				Timestamp:   a.ReturnedAt.UTC().Format(time.RFC3339),
				ActorID:     &userID,
				ActorName:   &userName,
				Details: map[string]interface{}{
					"condition_after": condAfter,
				},
			})
		}
	}
	return timeline, nil
}

type allocationTestContext struct {
	app            *fiber.App
	userRepo       *inMemoryUserRepo
	sessionRepo    *inMemorySessionRepo
	assetRepo      *inMemoryAssetRepo
	attendanceRepo *inMemoryAttendanceRepo
	allocRepo      *inMemoryAllocationRepo
	helper         *AuthServiceTestHelper
}

func setupAllocationTestApp() *allocationTestContext {
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

	// Seed presensi masuk mahasiswa (student-uuid-1 pada sess-1 yang berstatus ACTIVE)
	_ = attendanceRepo.Create(context.Background(), &model.Attendance{
		UserID:             "student-uuid-1",
		SessionID:          "sess-1",
		CheckedInAt:        time.Now().UTC(),
		VerificationStatus: "VERIFIED_FACE",
	})

	return &allocationTestContext{
		app:            app,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		assetRepo:      assetRepo,
		attendanceRepo: attendanceRepo,
		allocRepo:      allocRepo,
		helper:         helper,
	}
}

// 1. Test Admin Create Allocation Success
func TestCreateAllocationSuccess(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// ast-1 berstatus AVAILABLE, student-uuid-1 punya presensi pada sess-1
	payload := `{
		"asset_id": "ast-1",
		"user_id": "student-uuid-1",
		"session_id": "sess-1",
		"condition_before": "GOOD"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
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

	var res model.AllocationResponse
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &res)

	if res.AssetID != "ast-1" || res.UserID != "student-uuid-1" {
		t.Errorf("mismatched allocation details: %+v", res)
	}

	// Pastikan status aset diperbarui menjadi ALLOCATED dalam transaksi
	asset, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if asset.Status != model.AssetStatusAllocated {
		t.Errorf("expected asset status ALLOCATED, got %s", asset.Status)
	}
}

// 2. Test Double Allocation & Concurrency Prevention
func TestDoubleAllocationPrevention(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Alokasi pertama -> Sukses
	payload := `{"asset_id":"ast-1","user_id":"student-uuid-1","session_id":"sess-1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("first allocation failed: %d", resp1.StatusCode)
	}

	// Alokasi kedua pada aset yang sama sebelum dikembalikan -> 409 Conflict
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)

	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict on double allocation, got %d", resp2.StatusCode)
	}
	b2, _ := io.ReadAll(resp2.Body)
	var errResp2 model.ErrorResponse
	_ = json.Unmarshal(b2, &errResp2)
	if errResp2.Error.Code != "ASSET_ALREADY_ALLOCATED" && errResp2.Error.Code != "ASSET_NOT_AVAILABLE" {
		t.Errorf("expected conflict error code, got %s", errResp2.Error.Code)
	}
}

// 3. Test Allocation Attendance Prerequisite
func TestAllocationAttendancePrerequisite(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// student-uuid-2 belum memiliki presensi masuk pada sess-1
	_ = tc.userRepo.UpdateStatus(context.Background(), "student-uuid-2", "ACTIVE")

	payload := `{"asset_id":"ast-1","user_id":"student-uuid-2","session_id":"sess-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request when attendance prerequisite is missing, got %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(b, &errResp)
	if errResp.Error.Code != "ATTENDANCE_REQUIRED" {
		t.Errorf("expected ATTENDANCE_REQUIRED, got %s", errResp.Error.Code)
	}
}

// 4. Test Allocation Blocked for Asset in Maintenance or Retired
func TestAllocationAssetMaintenanceBlocked(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// ast-3 berstatus MAINTENANCE
	payload := `{"asset_id":"ast-3","user_id":"student-uuid-1","session_id":"sess-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict when allocating asset in maintenance, got %d", resp.StatusCode)
	}
}

// 5. Test Return Allocation Success (Good vs Damaged)
func TestReturnAllocationSuccess(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Buat alokasi ast-1
	payload := `{"asset_id":"ast-1","user_id":"student-uuid-1","session_id":"sess-1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	var created model.AllocationResponse
	b1, _ := io.ReadAll(resp1.Body)
	_ = json.Unmarshal(b1, &created)

	// a. Return dengan condition GOOD -> Aset kembali AVAILABLE
	returnPayload := `{"condition":"GOOD","notes":"Semua tombol berfungsi normal"}`
	reqReturn := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/allocations/%s/return", created.ID), bytes.NewBufferString(returnPayload))
	reqReturn.Header.Set("Content-Type", "application/json")
	reqReturn.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respReturn, _ := tc.app.Test(reqReturn)

	if respReturn.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on return, got %d", respReturn.StatusCode)
	}
	var resReturn model.AllocationResponse
	bRet, _ := io.ReadAll(respReturn.Body)
	_ = json.Unmarshal(bRet, &resReturn)

	if resReturn.ReturnedAt == nil || *resReturn.ReturnedAt == "" {
		t.Errorf("expected returned_at timestamp")
	}

	assetAfter, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if assetAfter.Status != model.AssetStatusAvailable {
		t.Errorf("expected asset status to revert to AVAILABLE, got %s", assetAfter.Status)
	}
}

// 6. Test Return Allocation Twice (Double Return)
func TestReturnAllocationTwice(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Buat alokasi
	payload := `{"asset_id":"ast-1","user_id":"student-uuid-1","session_id":"sess-1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	var created model.AllocationResponse
	b1, _ := io.ReadAll(resp1.Body)
	_ = json.Unmarshal(b1, &created)

	// Return pertama -> 200 OK
	returnPayload := `{"condition":"GOOD"}`
	reqReturn1 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/allocations/%s/return", created.ID), bytes.NewBufferString(returnPayload))
	reqReturn1.Header.Set("Content-Type", "application/json")
	reqReturn1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respReturn1, _ := tc.app.Test(reqReturn1)
	if respReturn1.StatusCode != http.StatusOK {
		t.Fatalf("first return failed: %d", respReturn1.StatusCode)
	}

	// Return kedua -> 409 Conflict (ALREADY_RETURNED)
	reqReturn2 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/allocations/%s/return", created.ID), bytes.NewBufferString(returnPayload))
	reqReturn2.Header.Set("Content-Type", "application/json")
	reqReturn2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respReturn2, _ := tc.app.Test(reqReturn2)

	if respReturn2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict on repeated return, got %d", respReturn2.StatusCode)
	}
}

// 7. Test Transaction Rollback on Failure
func TestAllocationTransactionRollback(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")
	tc.allocRepo.simulateDBError = true

	payload := `{"asset_id":"ast-1","user_id":"student-uuid-1","session_id":"sess-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 on transaction rollback, got %d", resp.StatusCode)
	}

	// Aset harus tetap AVAILABLE (tidak menjadi ALLOCATED)
	asset, _ := tc.assetRepo.GetByID(context.Background(), "ast-1")
	if asset.Status != model.AssetStatusAvailable {
		t.Errorf("asset status must remain AVAILABLE after transaction rollback, got %s", asset.Status)
	}
}

// 8. Test Ownership Scoping & Student Access
func TestAllocationOwnershipScoping(t *testing.T) {
	tc := setupAllocationTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Student calls POST /api/v1/allocations -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(`{}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls POST /api/v1/allocations, got %d", resp1.StatusCode)
	}

	// b. Student calls GET /api/v1/allocations (Admin only list) -> 403 Forbidden
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/allocations", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls GET /api/v1/allocations, got %d", resp2.StatusCode)
	}

	// c. Student calls GET /api/v1/allocations/my -> 200 OK
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/allocations/my", nil)
	req3.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp3, _ := tc.app.Test(req3)
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student reads their own allocations, got %d", resp3.StatusCode)
	}

	// d. Admin calls GET /api/v1/allocations -> 200 OK
	reqAdmin := httptest.NewRequest(http.MethodGet, "/api/v1/allocations", nil)
	reqAdmin.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respAdmin, _ := tc.app.Test(reqAdmin)
	if respAdmin.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when admin calls GET /api/v1/allocations, got %d", respAdmin.StatusCode)
	}
}

// 9. Test Asset History Timeline
func TestAssetTimelineHistory(t *testing.T) {
	tc := setupAllocationTestApp()
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// Buat alokasi ast-1 lalu kembalikan
	payload := `{"asset_id":"ast-1","user_id":"student-uuid-1","session_id":"sess-1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/allocations", bytes.NewBufferString(payload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := tc.app.Test(req1)
	var created model.AllocationResponse
	b1, _ := io.ReadAll(resp1.Body)
	_ = json.Unmarshal(b1, &created)

	returnPayload := `{"condition":"DAMAGED","notes":"Layar flickering"}`
	reqReturn := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/allocations/%s/return", created.ID), bytes.NewBufferString(returnPayload))
	reqReturn.Header.Set("Content-Type", "application/json")
	reqReturn.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	_, _ = tc.app.Test(reqReturn)

	// Ambil histori timeline ast-1
	reqHistory := httptest.NewRequest(http.MethodGet, "/api/v1/assets/ast-1/history", nil)
	reqHistory.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respHistory, _ := tc.app.Test(reqHistory)

	if respHistory.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on asset history, got %d", respHistory.StatusCode)
	}

	var history model.AssetHistoryResponse
	bHist, _ := io.ReadAll(respHistory.Body)
	_ = json.Unmarshal(bHist, &history)

	if len(history.Timeline) < 2 {
		t.Errorf("expected at least 2 events in timeline (ALLOCATION and RETURN), got %d", len(history.Timeline))
	}
}
