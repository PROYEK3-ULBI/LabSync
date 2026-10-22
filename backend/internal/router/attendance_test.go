package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/config"
	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
	"labsync-backend/internal/service"
	"labsync-backend/pkg/crypto"
)

type inMemoryAttendanceRepo struct {
	attendances     map[string]*model.Attendance
	simulateDBError bool
}

func newInMemoryAttendanceRepo() *inMemoryAttendanceRepo {
	return &inMemoryAttendanceRepo{
		attendances: make(map[string]*model.Attendance),
	}
}

func (r *inMemoryAttendanceRepo) GetByUserAndSession(ctx context.Context, userID, sessionID string) (*model.Attendance, error) {
	for _, a := range r.attendances {
		if a.UserID == userID && a.SessionID == sessionID {
			copyA := *a
			return &copyA, nil
		}
	}
	return nil, repository.ErrAttendanceNotFound
}

func (r *inMemoryAttendanceRepo) GetByID(ctx context.Context, id string) (*model.Attendance, error) {
	if a, ok := r.attendances[id]; ok {
		copyA := *a
		return &copyA, nil
	}
	return nil, repository.ErrAttendanceNotFound
}

func (r *inMemoryAttendanceRepo) Create(ctx context.Context, a *model.Attendance) error {
	if r.simulateDBError {
		return errors.New("simulated database failure")
	}
	for _, existing := range r.attendances {
		if existing.UserID == a.UserID && existing.SessionID == a.SessionID {
			return repository.ErrDuplicateAttendance
		}
	}
	a.ID = fmt.Sprintf("att-%d", len(r.attendances)+1)
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
	r.attendances[a.ID] = a
	return nil
}

func (r *inMemoryAttendanceRepo) CheckOut(ctx context.Context, id string, checkedOutAt time.Time) error {
	if r.simulateDBError {
		return errors.New("simulated database failure")
	}
	a, ok := r.attendances[id]
	if !ok {
		return repository.ErrAttendanceNotFound
	}
	if a.CheckedOutAt != nil {
		return repository.ErrAlreadyCheckedOut
	}
	a.CheckedOutAt = &checkedOutAt
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *inMemoryAttendanceRepo) List(ctx context.Context, q model.AttendanceListQuery) ([]model.Attendance, int, error) {
	var list []model.Attendance
	for _, a := range r.attendances {
		if q.UserID != "" && a.UserID != q.UserID {
			continue
		}
		if q.SessionID != "" && a.SessionID != q.SessionID {
			continue
		}
		list = append(list, *a)
	}
	return list, len(list), nil
}

type attendanceTestContext struct {
	app            *fiber.App
	userRepo       *inMemoryUserRepo
	sessionRepo    *inMemorySessionRepo
	faceRepo       *inMemoryFaceRepo
	attendanceRepo *inMemoryAttendanceRepo
	faceClient     *mockFaceClient
	helper         *AuthServiceTestHelper
}

func setupAttendanceTestApp() *attendanceTestContext {
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

	// Seed template pendaftaran aktif untuk student-uuid-1
	key := []byte(cfg.BiometricEncryptionKey)
	dummyVec := make([]float64, 512)
	for i := range dummyVec {
		dummyVec[i] = 0.05
	}
	cipherBytes, _ := crypto.EncryptEmbedding(key, dummyVec)

	tmpl := &model.FaceTemplate{
		ID:             "tmpl-student-1",
		UserID:         "student-uuid-1",
		ModelName:      "ArcFace",
		ModelVersion:   "1.0.0",
		ConsentVersion: "v1.0",
		ConsentedAt:    time.Now().UTC(),
		Status:         "ACTIVE",
		EnrolledAt:     time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	samples := []model.FaceTemplateSample{
		{FaceTemplateID: tmpl.ID, Pose: "FRONT", EmbeddingCiphertext: cipherBytes, QualityStatus: "VALID"},
		{FaceTemplateID: tmpl.ID, Pose: "LEFT", EmbeddingCiphertext: cipherBytes, QualityStatus: "VALID"},
		{FaceTemplateID: tmpl.ID, Pose: "RIGHT", EmbeddingCiphertext: cipherBytes, QualityStatus: "VALID"},
	}
	_ = faceRepo.SaveEnrollment(context.Background(), tmpl, samples, false)

	return &attendanceTestContext{
		app:            app,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		faceRepo:       faceRepo,
		attendanceRepo: attendanceRepo,
		faceClient:     faceClient,
		helper:         helper,
	}
}

func createCheckInRequest(url string, token string, sessionID string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("session_id", sessionID)

	dummyImage := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image"; filename="live.jpg"`)
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, _ = part.Write(dummyImage)

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.AddCookie(&http.Cookie{Name: "token", Value: token})
	}

	return req, nil
}

// 1. Test Check-in Success
func TestAttendanceCheckInSuccess(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// sess-1 berstatus ACTIVE
	req, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(b))
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var res model.CheckInResponse
	_ = json.Unmarshal(bodyBytes, &res)

	if res.Attendance.VerificationStatus != "VERIFIED_FACE" {
		t.Errorf("expected VERIFIED_FACE status, got %s", res.Attendance.VerificationStatus)
	}
	if res.Attendance.CheckedInAt == "" {
		t.Errorf("expected server-determined checked_in_at timestamp")
	}

	// Pastikan tidak ada embedding atau citra yang bocor
	bodyStr := string(bodyBytes)
	if strings.Contains(bodyStr, "embedding") || strings.Contains(bodyStr, "ciphertext") {
		t.Errorf("response body leaked biometric data: %s", bodyStr)
	}
}

// 2. Test Check-in Session Not Active
func TestAttendanceCheckInSessionNotActive(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// sess-2 berstatus SCHEDULED (belum aktif)
	req, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-2")
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for non-active session, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(bodyBytes, &errResp)

	if errResp.Error.Code != "SESSION_NOT_ACTIVE" {
		t.Errorf("expected SESSION_NOT_ACTIVE, got %s", errResp.Error.Code)
	}
}

// 3. Test Check-in Missing Face Template
func TestAttendanceCheckInMissingTemplate(t *testing.T) {
	tc := setupAttendanceTestApp()
	// user student-uuid-2 belum memiliki template terdaftar
	_ = tc.userRepo.UpdateStatus(context.Background(), "student-uuid-2", "ACTIVE")
	tokenNoTemplate := tc.helper.getTokenFor("student-uuid-2")

	req, _ := createCheckInRequest("/api/v1/attendance/check-in", tokenNoTemplate, "sess-1")
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request when user has no face template, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(bodyBytes, &errResp)

	if errResp.Error.Code != "FACE_NOT_ENROLLED" {
		t.Errorf("expected FACE_NOT_ENROLLED, got %s", errResp.Error.Code)
	}
}

// 4. Test Check-in Face Verification Failed
func TestAttendanceCheckInVerificationFailed(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")
	tc.faceClient.simulateVerifyFail = true

	req, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on verification fail, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(bodyBytes, &errResp)

	if errResp.Error.Code != "VERIFICATION_FAILED" {
		t.Errorf("expected VERIFICATION_FAILED, got %s", errResp.Error.Code)
	}

	// Pastikan data attendance TIDAK tersimpan
	att, err := tc.attendanceRepo.GetByUserAndSession(context.Background(), "student-uuid-1", "sess-1")
	if err == nil || att != nil {
		t.Errorf("attendance must NOT be created when verification fails")
	}
}

// 5. Test Duplicate Check-in / Race Condition / Unique Constraint
func TestAttendanceDuplicateCheckIn(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// Check-in pertama -> 200 OK
	req1, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("first check-in failed: %d", resp1.StatusCode)
	}

	// Check-in kedua pada sesi yang sama -> 409 Conflict
	req2, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict on duplicate check-in, got %d", resp2.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp2.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(bodyBytes, &errResp)

	if errResp.Error.Code != "ALREADY_CHECKED_IN" {
		t.Errorf("expected ALREADY_CHECKED_IN, got %s", errResp.Error.Code)
	}
}

// 6. Test Check-out Flow (Before check-in, success, repeated)
func TestAttendanceCheckOutLifecycle(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// a. Check-out sebelum check-in -> 404 Not Found (ATTENDANCE_NOT_FOUND)
	payload := `{"session_id":"sess-1"}`
	reqBefore, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/check-out", bytes.NewBufferString(payload))
	reqBefore.Header.Set("Content-Type", "application/json")
	reqBefore.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	respBefore, _ := tc.app.Test(reqBefore)

	if respBefore.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found on check-out before check-in, got %d", respBefore.StatusCode)
	}

	// b. Lakukan Check-in
	reqCheckIn, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	respCheckIn, _ := tc.app.Test(reqCheckIn)
	if respCheckIn.StatusCode != http.StatusOK {
		t.Fatalf("check-in failed: %d", respCheckIn.StatusCode)
	}

	// c. Check-out sukses -> 200 OK
	reqOut, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/check-out", bytes.NewBufferString(payload))
	reqOut.Header.Set("Content-Type", "application/json")
	reqOut.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	respOut, _ := tc.app.Test(reqOut)

	if respOut.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on check-out, got %d", respOut.StatusCode)
	}

	var resOut model.CheckOutResponse
	bOut, _ := io.ReadAll(respOut.Body)
	_ = json.Unmarshal(bOut, &resOut)
	if resOut.Attendance.CheckedOutAt == nil || *resOut.Attendance.CheckedOutAt == "" {
		t.Errorf("expected server-determined checked_out_at timestamp")
	}

	// d. Check-out berulang (repeated check-out) -> 409 Conflict (ALREADY_CHECKED_OUT)
	reqOutAgain, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/check-out", bytes.NewBufferString(payload))
	reqOutAgain.Header.Set("Content-Type", "application/json")
	reqOutAgain.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	respOutAgain, _ := tc.app.Test(reqOutAgain)

	if respOutAgain.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict on repeated check-out, got %d", respOutAgain.StatusCode)
	}
	bAgain, _ := io.ReadAll(respOutAgain.Body)
	var errAgain model.ErrorResponse
	_ = json.Unmarshal(bAgain, &errAgain)
	if errAgain.Error.Code != "ALREADY_CHECKED_OUT" {
		t.Errorf("expected ALREADY_CHECKED_OUT, got %s", errAgain.Error.Code)
	}
}

// 7. Test Role Access & Scoping
func TestAttendanceRoleAccess(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")
	adminToken := tc.helper.getTokenFor("admin-uuid-1")

	// a. Student mengakses endpoint list admin GET /api/v1/attendance -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/attendance", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls GET /api/v1/attendance, got %d", resp1.StatusCode)
	}

	// b. Admin mengakses GET /api/v1/attendance -> 200 OK
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/attendance", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when admin calls GET /api/v1/attendance, got %d", resp2.StatusCode)
	}

	// c. Student mengakses GET /api/v1/attendance/history miliknya -> 200 OK
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/attendance/history", nil)
	req3.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp3, _ := tc.app.Test(req3)
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student views history, got %d", resp3.StatusCode)
	}
}

// 8. Test Rollback Bila DB Gagal
func TestAttendanceRollbackOnDBError(t *testing.T) {
	tc := setupAttendanceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")
	tc.attendanceRepo.simulateDBError = true

	req, _ := createCheckInRequest("/api/v1/attendance/check-in", studentToken, "sess-1")
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 on database failure, got %d", resp.StatusCode)
	}
}
