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

// In-Memory Test Repository
type inMemoryUserRepo struct {
	users         map[string]*model.User
	roles         map[string]*model.Role
	faceTemplates map[string]bool
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	repo := &inMemoryUserRepo{
		users:         make(map[string]*model.User),
		roles:         make(map[string]*model.Role),
		faceTemplates: make(map[string]bool),
	}

	repo.roles["ADMIN"] = &model.Role{ID: "role-admin-id", Code: "ADMIN", Name: "Administrator"}
	repo.roles["STUDENT"] = &model.Role{ID: "role-student-id", Code: "STUDENT", Name: "Mahasiswa"}

	// Seed Admin
	adminHash, _ := service.HashPassword("AdminPass123!")
	repo.users["admin-uuid-1"] = &model.User{
		ID:           "admin-uuid-1",
		RoleID:       "role-admin-id",
		RoleCode:     "ADMIN",
		Identifier:   "ADM001",
		Email:        "admin@labsync.ac.id",
		FullName:     "Super Admin",
		PasswordHash: adminHash,
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Seed Student Active
	mhsHash, _ := service.HashPassword("StudentPass123!")
	repo.users["student-uuid-1"] = &model.User{
		ID:           "student-uuid-1",
		RoleID:       "role-student-id",
		RoleCode:     "STUDENT",
		Identifier:   "TI240001",
		Email:        "mhs1@univ.ac.id",
		FullName:     "Budi Mahasiswa",
		PasswordHash: mhsHash,
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Seed Student Inactive
	repo.users["student-uuid-2"] = &model.User{
		ID:           "student-uuid-2",
		RoleID:       "role-student-id",
		RoleCode:     "STUDENT",
		Identifier:   "TI240002",
		Email:        "mhs2_inactive@univ.ac.id",
		FullName:     "Doni Nonaktif",
		PasswordHash: mhsHash,
		Status:       "INACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return repo
}

func (m *inMemoryUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		copyUser := *u
		return &copyUser, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *inMemoryUserRepo) GetByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	for _, u := range m.users {
		if strings.EqualFold(u.Identifier, identifier) {
			copyUser := *u
			return &copyUser, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *inMemoryUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, u := range m.users {
		if strings.EqualFold(u.Email, email) {
			copyUser := *u
			return &copyUser, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *inMemoryUserRepo) List(ctx context.Context, q model.UserListQuery) ([]model.User, int, error) {
	var filtered []model.User
	for _, u := range m.users {
		if q.Role != "" && !strings.EqualFold(u.RoleCode, q.Role) {
			continue
		}
		if q.Status != "" && !strings.EqualFold(u.Status, q.Status) {
			continue
		}
		if q.Search != "" {
			s := strings.ToLower(q.Search)
			if !strings.Contains(strings.ToLower(u.Identifier), s) &&
				!strings.Contains(strings.ToLower(u.FullName), s) &&
				!strings.Contains(strings.ToLower(u.Email), s) {
				continue
			}
		}
		filtered = append(filtered, *u)
	}
	return filtered, len(filtered), nil
}

func (m *inMemoryUserRepo) Create(ctx context.Context, u *model.User) error {
	u.ID = fmt.Sprintf("user-%d", len(m.users)+1)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	m.users[u.ID] = u
	return nil
}

func (m *inMemoryUserRepo) Update(ctx context.Context, u *model.User) error {
	if _, ok := m.users[u.ID]; !ok {
		return repository.ErrUserNotFound
	}
	u.UpdatedAt = time.Now()
	m.users[u.ID] = u
	return nil
}

func (m *inMemoryUserRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	if u, ok := m.users[id]; ok {
		u.Status = status
		u.UpdatedAt = time.Now()
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *inMemoryUserRepo) GetRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	if r, ok := m.roles[strings.ToUpper(code)]; ok {
		return r, nil
	}
	return nil, repository.ErrRoleNotFound
}

func (m *inMemoryUserRepo) HasFaceTemplate(ctx context.Context, userID string) (bool, error) {
	return m.faceTemplates[userID], nil
}

type inMemoryAuditRepo struct {
	logs []model.AuditLog
}

func (m *inMemoryAuditRepo) Create(ctx context.Context, log *model.AuditLog) error {
	m.logs = append(m.logs, *log)
	return nil
}

func (m *inMemoryAuditRepo) List(ctx context.Context, q model.AuditLogListQuery) ([]model.AuditLog, int, error) {
	var list []model.AuditLog
	for _, l := range m.logs {
		if q.Action != "" && !strings.EqualFold(l.Action, q.Action) {
			continue
		}
		if q.EntityType != "" && !strings.EqualFold(l.EntityType, q.EntityType) {
			continue
		}
		if q.UserID != "" && (l.UserID == nil || *l.UserID != q.UserID) {
			continue
		}
		list = append(list, l)
	}
	return list, len(list), nil
}

func setupTestApp() (*fiber.App, *AuthServiceTestHelper) {
	cfg := &config.Config{
		AppEnv:           "development",
		Port:             "8000",
		CORSAllowOrigins: "http://localhost:8443",
		BodyLimitMB:      10,
		RateLimitMax:     100,
		JWTSecret:        "test-super-secret-key-1234567890",
		JWTExpiryHours:   24,
		CookieSecure:     false,
		CookieSameSite:   "Lax",
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
	faceService := service.NewFaceService(userRepo, faceRepo, auditRepo, faceClient, "01234567890123456789012345678901")
	attendanceService := service.NewAttendanceService(attendanceRepo, sessionRepo, faceRepo, userRepo, auditRepo, faceClient, "01234567890123456789012345678901")
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

	return app, &AuthServiceTestHelper{
		authService:    authService,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		sessionRepo:    sessionRepo,
		sessionService: sessionService,
	}
}

type AuthServiceTestHelper struct {
	authService    *service.AuthService
	userRepo       *inMemoryUserRepo
	auditRepo      *inMemoryAuditRepo
	sessionRepo    *inMemorySessionRepo
	sessionService *service.SessionService
}

func (h *AuthServiceTestHelper) getTokenFor(userID string) string {
	u, _ := h.userRepo.GetByID(context.Background(), userID)
	token, _ := h.authService.GenerateToken(u)
	return token
}

// 1. Test Login Valid
func TestLoginValid(t *testing.T) {
	app, _ := setupTestApp()

	payload := `{"identifier":"TI240001","password":"StudentPass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Cek cookie HttpOnly
	cookies := resp.Cookies()
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}

	if tokenCookie == nil {
		t.Fatalf("expected Set-Cookie with token")
	}
	if !tokenCookie.HttpOnly {
		t.Errorf("cookie must be HttpOnly")
	}
	if tokenCookie.Value == "" {
		t.Errorf("cookie value must not be empty")
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Pastikan password dan hash TIDAK PERNAH bocor
	if strings.Contains(bodyStr, "password") {
		t.Errorf("response body leaked password key: %s", bodyStr)
	}
	if strings.Contains(bodyStr, "$2a$") {
		t.Errorf("response body leaked bcrypt hash: %s", bodyStr)
	}
}

// 2. Test Login Invalid (Wrong password & Unknown identifier)
func TestLoginInvalidCredentials(t *testing.T) {
	app, _ := setupTestApp()

	// Kasus password salah
	payloadWrongPass := `{"identifier":"TI240001","password":"WrongPassword!"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(payloadWrongPass))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 on wrong password, got %d", resp1.StatusCode)
	}
	body1, _ := io.ReadAll(resp1.Body)
	var errResp1 model.ErrorResponse
	_ = json.Unmarshal(body1, &errResp1)
	if errResp1.Error.Code != "UNAUTHORIZED" {
		t.Errorf("expected UNAUTHORIZED, got %s", errResp1.Error.Code)
	}

	// Kasus identifier tidak ada
	payloadUnknownID := `{"identifier":"UNKNOWN999","password":"SomePassword!"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(payloadUnknownID))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 on unknown identifier, got %d", resp2.StatusCode)
	}
	body2, _ := io.ReadAll(resp2.Body)
	var errResp2 model.ErrorResponse
	_ = json.Unmarshal(body2, &errResp2)

	// Pesan error harus generik dan konsisten agar tidak membocorkan keberadaan identifier
	if errResp1.Error.Message != errResp2.Error.Message {
		t.Errorf("error messages must be identical for wrong pass and unknown user to prevent user enumeration: %q vs %q",
			errResp1.Error.Message, errResp2.Error.Message)
	}
}

// 3. Test Login Inactive Account
func TestLoginInactiveAccount(t *testing.T) {
	app, _ := setupTestApp()

	payload := `{"identifier":"TI240002","password":"StudentPass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 on inactive account login, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	if errResp.Error.Code != "FORBIDDEN" {
		t.Errorf("expected FORBIDDEN code, got %s", errResp.Error.Code)
	}
}

// 4. Test Logout
func TestLogout(t *testing.T) {
	app, helper := setupTestApp()

	token := helper.getTokenFor("student-uuid-1")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on logout, got %d", resp.StatusCode)
	}

	// Pastikan cookie dihapus (MaxAge <= 0)
	var tokenCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}
	if tokenCookie == nil || tokenCookie.MaxAge > 0 {
		t.Errorf("expected cookie to be invalidated on logout")
	}
}

// 5. Test Auth Without Cookie
func TestAuthWithoutCookie(t *testing.T) {
	app, _ := setupTestApp()

	endpoints := []string{
		"/api/v1/auth/me",
		"/api/v1/users",
		"/api/v1/users/admin-uuid-1",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request %s failed: %v", ep, err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401 on %s without cookie, got %d", ep, resp.StatusCode)
		}
	}
}

// 6. Test RBAC: Student vs Admin
func TestRoleStudentVsAdmin(t *testing.T) {
	app, helper := setupTestApp()

	studentToken := helper.getTokenFor("student-uuid-1")
	adminToken := helper.getTokenFor("admin-uuid-1")

	// a. Student mengakses daftar users (ADMIN only) -> 403 Forbidden
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls GET /api/v1/users, got %d", resp1.StatusCode)
	}

	// b. Admin mengakses daftar users -> 200 OK
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when admin calls GET /api/v1/users, got %d", resp2.StatusCode)
	}

	// c. Student melihat data user lain -> 403 Forbidden
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/users/admin-uuid-1", nil)
	req3.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student inspects other user, got %d", resp3.StatusCode)
	}

	// d. Student melihat data dirinya sendiri -> 200 OK
	req4 := httptest.NewRequest(http.MethodGet, "/api/v1/users/student-uuid-1", nil)
	req4.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student inspects self, got %d", resp4.StatusCode)
	}

	// e. Student mencoba mengubah role dirinya sendiri menjadi ADMIN -> 403 Forbidden
	req5 := httptest.NewRequest(http.MethodPut, "/api/v1/users/student-uuid-1", bytes.NewBufferString(`{"role":"ADMIN"}`))
	req5.Header.Set("Content-Type", "application/json")
	req5.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp5, _ := app.Test(req5)
	if resp5.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student attempts role escalation, got %d", resp5.StatusCode)
	}
}

// 7. Test Duplicate Identifier / Email
func TestDuplicateEmailOrIdentifier(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	// Duplikat Identifier
	dupIDPayload := `{"identifier":"TI240001","email":"unik1@univ.ac.id","full_name":"Mahasiswa Baru","password":"Password123!","role":"STUDENT"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(dupIDPayload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict on duplicate identifier, got %d", resp1.StatusCode)
	}

	// Duplikat Email
	dupEmailPayload := `{"identifier":"TI240099","email":"mhs1@univ.ac.id","full_name":"Mahasiswa Baru 2","password":"Password123!","role":"STUDENT"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(dupEmailPayload))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict on duplicate email, got %d", resp2.StatusCode)
	}
}

// 8. Test Password Never Leaked
func TestPasswordNeverLeaked(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	// Create user
	newUserPayload := `{"identifier":"TI240088","email":"testleak@univ.ac.id","full_name":"Test No Leak","password":"SecretPassword123!","role":"STUDENT"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(newUserPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respCreate, _ := app.Test(reqCreate)
	bodyCreate, _ := io.ReadAll(respCreate.Body)

	if strings.Contains(string(bodyCreate), "SecretPassword123!") || strings.Contains(string(bodyCreate), "password_hash") {
		t.Errorf("user create leaked password or hash: %s", string(bodyCreate))
	}

	// List users
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	reqList.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respList, _ := app.Test(reqList)
	bodyList, _ := io.ReadAll(respList.Body)

	if strings.Contains(string(bodyList), "$2a$") || strings.Contains(string(bodyList), "password_hash") {
		t.Errorf("user list leaked password hash: %s", string(bodyList))
	}

	// Auth /me
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respMe, _ := app.Test(reqMe)
	bodyMe, _ := io.ReadAll(respMe.Body)

	if strings.Contains(string(bodyMe), "$2a$") || strings.Contains(string(bodyMe), "password_hash") {
		t.Errorf("auth/me leaked password hash: %s", string(bodyMe))
	}
}
