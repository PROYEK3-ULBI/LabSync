package router

import (
	"bytes"
	"context"
	"encoding/json"
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
)

// Mock Face Client
type mockFaceClient struct {
	simulateNoFace       bool
	simulateMultipleFace bool
	simulateTimeout      bool
	simulateUnavailable  bool
	simulateVerifyFail   bool
	extractedPoses       []string
}

func (m *mockFaceClient) Extract(ctx context.Context, imgBytes []byte, filename, contentType, pose, requestID string) (*model.PythonExtractResponse, error) {
	if m.simulateTimeout {
		return nil, &model.AppError{
			StatusCode: http.StatusGatewayTimeout,
			Code:       "FACE_SERVICE_TIMEOUT",
			Message:    "Waktu tunggu pemrosesan biometrik wajah terlampaui",
		}
	}
	if m.simulateUnavailable {
		return nil, &model.AppError{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "FACE_SERVICE_UNAVAILABLE",
			Message:    "Layanan internal biometrik wajah tidak dapat dihubungi",
		}
	}
	if m.simulateNoFace {
		return nil, &model.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       "NO_FACE_DETECTED",
			Message:    "Tidak ada wajah terdeteksi pada citra",
		}
	}
	if m.simulateMultipleFace {
		return nil, &model.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       "MULTIPLE_FACES_DETECTED",
			Message:    "Terdeteksi lebih dari satu wajah pada citra",
		}
	}

	m.extractedPoses = append(m.extractedPoses, pose)

	// Vektor sintetis 512 dimensi
	embedding := make([]float64, 512)
	for i := range embedding {
		embedding[i] = 0.05
	}

	return &model.PythonExtractResponse{
		ModelName:      "ArcFace",
		ModelVersion:   "1.0.0",
		Pose:           pose,
		Embedding:      embedding,
		Dim:            512,
		FaceConfidence: 0.99,
	}, nil
}

func (m *mockFaceClient) Verify(ctx context.Context, imgBytes []byte, filename, contentType string, enrolledSamples []model.EnrolledSamplePayload, threshold *float64, strategy, requestID string) (*model.PythonVerifyResponse, error) {
	if m.simulateTimeout {
		return nil, &model.AppError{
			StatusCode: http.StatusGatewayTimeout,
			Code:       "FACE_SERVICE_TIMEOUT",
			Message:    "Waktu tunggu pemrosesan verifikasi wajah terlampaui",
		}
	}
	if m.simulateUnavailable {
		return nil, &model.AppError{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "FACE_SERVICE_UNAVAILABLE",
			Message:    "Layanan internal biometrik wajah tidak dapat dihubungi",
		}
	}
	if m.simulateNoFace {
		return nil, &model.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       "NO_FACE_DETECTED",
			Message:    "Tidak ada wajah terdeteksi pada citra",
		}
	}
	if m.simulateMultipleFace {
		return nil, &model.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       "MULTIPLE_FACES_DETECTED",
			Message:    "Terdeteksi lebih dari satu wajah pada citra",
		}
	}

	verified := !m.simulateVerifyFail
	dist := 0.25
	rule := "MIN_DISTANCE_THRESHOLD_PASSED"
	if !verified {
		dist = 0.85
		rule = "MIN_DISTANCE_THRESHOLD_EXCEEDED"
	}

	return &model.PythonVerifyResponse{
		Verified:        verified,
		Distance:        dist,
		Threshold:       0.40,
		Strategy:        strategy,
		BestMatchedPose: "FRONT",
		DistancesByPose: map[string]float64{"FRONT": dist, "LEFT": dist + 0.1, "RIGHT": dist + 0.1},
		IsCalibrated:    true,
		DecisionRule:    rule,
		ModelName:       "ArcFace",
	}, nil
}

// In-Memory Face Repository
type inMemoryFaceRepo struct {
	templates map[string]*model.FaceTemplate
	samples   map[string][]model.FaceTemplateSample
}

func newInMemoryFaceRepo() *inMemoryFaceRepo {
	return &inMemoryFaceRepo{
		templates: make(map[string]*model.FaceTemplate),
		samples:   make(map[string][]model.FaceTemplateSample),
	}
}

func (r *inMemoryFaceRepo) GetActiveTemplateByUserID(ctx context.Context, userID string) (*model.FaceTemplate, []model.FaceTemplateSample, error) {
	for _, t := range r.templates {
		if t.UserID == userID && t.Status == "ACTIVE" {
			smps := r.samples[t.ID]
			copyT := *t
			return &copyT, smps, nil
		}
	}
	return nil, nil, repository.ErrFaceTemplateNotFound
}

func (r *inMemoryFaceRepo) SaveEnrollment(ctx context.Context, template *model.FaceTemplate, samples []model.FaceTemplateSample, revokeOld bool) error {
	for _, t := range r.templates {
		if t.UserID == template.UserID && t.Status == "ACTIVE" {
			if !revokeOld {
				return repository.ErrFaceAlreadyEnrolled
			}
			t.Status = "REVOKED"
			t.UpdatedAt = time.Now()
		}
	}

	template.ID = fmt.Sprintf("tmpl-%d", len(r.templates)+1)
	template.EnrolledAt = time.Now()
	template.UpdatedAt = time.Now()
	r.templates[template.ID] = template

	for i := range samples {
		samples[i].ID = fmt.Sprintf("sample-%d-%d", len(r.templates), i+1)
		samples[i].FaceTemplateID = template.ID
		samples[i].CreatedAt = time.Now()
	}
	r.samples[template.ID] = samples

	return nil
}

func (r *inMemoryFaceRepo) RevokeTemplate(ctx context.Context, userID string) error {
	for _, t := range r.templates {
		if t.UserID == userID && t.Status == "ACTIVE" {
			t.Status = "REVOKED"
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return repository.ErrFaceTemplateNotFound
}

type faceTestContext struct {
	app        *fiber.App
	userRepo   *inMemoryUserRepo
	faceRepo   *inMemoryFaceRepo
	faceClient *mockFaceClient
	helper     *AuthServiceTestHelper
}

func setupFaceTestApp() *faceTestContext {
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

	return &faceTestContext{
		app:        app,
		userRepo:   userRepo,
		faceRepo:   faceRepo,
		faceClient: faceClient,
		helper:     helper,
	}
}

func createMultipartEnrollRequest(url string, token string, consentAgreed bool, reEnroll bool) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if consentAgreed {
		_ = writer.WriteField("consent_agreed", "true")
	} else {
		_ = writer.WriteField("consent_agreed", "false")
	}

	if reEnroll {
		_ = writer.WriteField("re_enroll", "true")
	}

	dummyImage := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46} // JPEG header magic bytes

	for _, field := range []string{"front_image", "left_image", "right_image"} {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s.jpg"`, field, field))
		h.Set("Content-Type", "image/jpeg")
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, err
		}
		_, _ = part.Write(dummyImage)
	}

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

// 1. Test Enrollment Success
func TestEnrollmentSuccess(t *testing.T) {
	tc := setupFaceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	req, err := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := tc.app.Test(req)
	if err != nil {
		t.Fatalf("unexpected app.Test error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201 Created, got %d: %s", resp.StatusCode, string(body))
	}

	var res model.FaceEnrollResponse
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &res)

	if res.SampleCount != 3 {
		t.Errorf("expected 3 samples, got %d", res.SampleCount)
	}
	if res.ModelName != "ArcFace" {
		t.Errorf("expected ArcFace, got %s", res.ModelName)
	}

	// Cek status enrollment di DB
	tmpl, samples, err := tc.faceRepo.GetActiveTemplateByUserID(context.Background(), "student-uuid-1")
	if err != nil {
		t.Fatalf("expected active template in DB: %v", err)
	}
	if tmpl.Status != "ACTIVE" {
		t.Errorf("expected template ACTIVE, got %s", tmpl.Status)
	}
	if len(samples) != 3 {
		t.Fatalf("expected 3 samples in DB, got %d", len(samples))
	}

	// Pastikan ciphertext terenkripsi dan bukan plaintext
	for _, s := range samples {
		if len(s.EmbeddingCiphertext) < 32 {
			t.Errorf("sample %s has invalid ciphertext length: %d", s.Pose, len(s.EmbeddingCiphertext))
		}
	}
}

// 2. Test No Face Detected
func TestEnrollmentNoFaceDetected(t *testing.T) {
	tc := setupFaceTestApp()
	tc.faceClient.simulateNoFace = true
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Unprocessable Entity, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(body, &errResp)

	if errResp.Error.Code != "NO_FACE_DETECTED" {
		t.Errorf("expected code NO_FACE_DETECTED, got %s", errResp.Error.Code)
	}
}

// 3. Test Multiple Faces Detected
func TestEnrollmentMultipleFacesDetected(t *testing.T) {
	tc := setupFaceTestApp()
	tc.faceClient.simulateMultipleFace = true
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Unprocessable Entity, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(body, &errResp)

	if errResp.Error.Code != "MULTIPLE_FACES_DETECTED" {
		t.Errorf("expected code MULTIPLE_FACES_DETECTED, got %s", errResp.Error.Code)
	}
}

// 4. Test Unauthenticated / No Session
func TestEnrollmentUnauthenticated(t *testing.T) {
	tc := setupFaceTestApp()

	// Tanpa token cookie
	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", "", true, false)
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized without session, got %d", resp.StatusCode)
	}
}

// 5. Test Inactive User Cannot Enroll
func TestEnrollmentInactiveUser(t *testing.T) {
	tc := setupFaceTestApp()
	inactiveToken := tc.helper.getTokenFor("student-uuid-2") // student-uuid-2 is INACTIVE

	// Mencoba enroll dengan token akun nonaktif
	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", inactiveToken, true, false)
	resp, _ := tc.app.Test(req)

	// Auth middleware / service menolak akun INACTIVE
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403 Forbidden or 401 Unauthorized for inactive user, got %d", resp.StatusCode)
	}
}

// 6. Test Replace / Re-Enroll Rules
func TestEnrollmentReplaceAndReEnroll(t *testing.T) {
	tc := setupFaceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// a. Pendaftaran pertama kali
	req1, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("first enroll failed: %d", resp1.StatusCode)
	}

	// b. Pendaftaran kedua tanpa re_enroll=true -> Harus 409 Conflict (FACE_ALREADY_ENROLLED)
	req2, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict when enrolling again without re_enroll, got %d", resp2.StatusCode)
	}
	body2, _ := io.ReadAll(resp2.Body)
	var errResp2 model.ErrorResponse
	_ = json.Unmarshal(body2, &errResp2)
	if errResp2.Error.Code != "FACE_ALREADY_ENROLLED" {
		t.Errorf("expected FACE_ALREADY_ENROLLED, got %s", errResp2.Error.Code)
	}

	// c. Pendaftaran ulang dengan re_enroll=true -> Sukses 201 Created & template lama di-revoke
	req3, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, true)
	resp3, _ := tc.app.Test(req3)
	if resp3.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on valid re-enroll, got %d", resp3.StatusCode)
	}

	// Cek bahwa hanya ada 1 template yang berstatus ACTIVE
	activeCount := 0
	for _, tmpl := range tc.faceRepo.templates {
		if tmpl.UserID == "student-uuid-1" && tmpl.Status == "ACTIVE" {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Errorf("expected exactly 1 ACTIVE template, got %d", activeCount)
	}
}

// 7. Test No Template or Embedding Leakage
func TestNoTemplateOrEmbeddingLeakage(t *testing.T) {
	tc := setupFaceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp, _ := tc.app.Test(req)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Pastikan respons enroll tidak memuat embedding
	if strings.Contains(bodyStr, "embedding") || strings.Contains(bodyStr, "ciphertext") {
		t.Errorf("enroll response leaked embedding data: %s", bodyStr)
	}

	// Cek endpoint status
	reqStatus := httptest.NewRequest(http.MethodGet, "/api/v1/face/status", nil)
	reqStatus.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	respStatus, _ := tc.app.Test(reqStatus)
	bodyStatusBytes, _ := io.ReadAll(respStatus.Body)
	bodyStatusStr := string(bodyStatusBytes)

	if strings.Contains(bodyStatusStr, "embedding") || strings.Contains(bodyStatusStr, "ciphertext") {
		t.Errorf("status response leaked embedding data: %s", bodyStatusStr)
	}
}

// 8. Test Service Timeout & Unavailable
func TestFaceServiceTimeoutAndUnavailable(t *testing.T) {
	tc := setupFaceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// a. Service Timeout (504)
	tc.faceClient.simulateTimeout = true
	req1, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp1, _ := tc.app.Test(req1)
	if resp1.StatusCode != http.StatusGatewayTimeout {
		t.Errorf("expected 504 Gateway Timeout, got %d", resp1.StatusCode)
	}

	// b. Service Unavailable (503)
	tc.faceClient.simulateTimeout = false
	tc.faceClient.simulateUnavailable = true
	req2, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, true, false)
	resp2, _ := tc.app.Test(req2)
	if resp2.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503 Service Unavailable, got %d", resp2.StatusCode)
	}
}

// 9. Test Consent Mandatory
func TestEnrollmentConsentMandatory(t *testing.T) {
	tc := setupFaceTestApp()
	studentToken := tc.helper.getTokenFor("student-uuid-1")

	// consentAgreed = false
	req, _ := createMultipartEnrollRequest("/api/v1/face/enroll", studentToken, false, false)
	resp, _ := tc.app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request when consent is false, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	if !strings.Contains(errResp.Error.Message, "consent_agreed") {
		t.Errorf("expected consent_agreed validation error, got: %s", errResp.Error.Message)
	}
}
