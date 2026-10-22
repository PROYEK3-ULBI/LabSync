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

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type inMemorySessionRepo struct {
	sessions   map[string]*model.LabSession
	attendance map[string]bool
}

func newInMemorySessionRepo() *inMemorySessionRepo {
	repo := &inMemorySessionRepo{
		sessions:   make(map[string]*model.LabSession),
		attendance: make(map[string]bool),
	}

	start1, _ := time.Parse(time.RFC3339, "2026-09-18T13:00:00Z")
	end1, _ := time.Parse(time.RFC3339, "2026-09-18T15:30:00Z")
	c1 := "IF2101"
	adminID := "admin-uuid-1"

	repo.sessions["sess-1"] = &model.LabSession{
		ID:          "sess-1",
		Code:        "PWB-01",
		Title:       "Praktikum Pemrograman Web",
		CourseCode:  &c1,
		Location:    "Lab Komputer 1",
		SessionDate: "2026-09-18",
		StartAt:     start1,
		EndAt:       end1,
		Status:      model.SessionStatusActive,
		CreatedBy:   &adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	start2, _ := time.Parse(time.RFC3339, "2026-09-19T08:00:00Z")
	end2, _ := time.Parse(time.RFC3339, "2026-09-19T10:30:00Z")
	repo.sessions["sess-2"] = &model.LabSession{
		ID:          "sess-2",
		Code:        "BDT-01",
		Title:       "Praktikum Basis Data",
		Location:    "Lab Komputer 2",
		SessionDate: "2026-09-19",
		StartAt:     start2,
		EndAt:       end2,
		Status:      model.SessionStatusScheduled,
		CreatedBy:   &adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Sesi dengan attendance terdaftar
	repo.attendance["sess-1"] = true

	return repo
}

func (m *inMemorySessionRepo) GetByID(ctx context.Context, id string) (*model.LabSession, error) {
	if s, ok := m.sessions[id]; ok {
		copySess := *s
		return &copySess, nil
	}
	return nil, repository.ErrSessionNotFound
}

func (m *inMemorySessionRepo) GetByCode(ctx context.Context, code string) (*model.LabSession, error) {
	for _, s := range m.sessions {
		if strings.EqualFold(s.Code, code) {
			copySess := *s
			return &copySess, nil
		}
	}
	return nil, repository.ErrSessionNotFound
}

func (m *inMemorySessionRepo) List(ctx context.Context, q model.SessionListQuery) ([]model.LabSession, int, error) {
	var filtered []model.LabSession
	for _, s := range m.sessions {
		if q.Status != "" && !strings.EqualFold(string(s.Status), q.Status) {
			continue
		}
		if q.SessionDate != "" && s.SessionDate != q.SessionDate {
			continue
		}
		if q.Search != "" {
			target := strings.ToLower(q.Search)
			match := strings.Contains(strings.ToLower(s.Code), target) ||
				strings.Contains(strings.ToLower(s.Title), target) ||
				strings.Contains(strings.ToLower(s.Location), target)
			if s.CourseCode != nil && strings.Contains(strings.ToLower(*s.CourseCode), target) {
				match = true
			}
			if !match {
				continue
			}
		}
		filtered = append(filtered, *s)
	}
	return filtered, len(filtered), nil
}

func (m *inMemorySessionRepo) GetActive(ctx context.Context) ([]model.LabSession, error) {
	var active []model.LabSession
	for _, s := range m.sessions {
		if s.Status == model.SessionStatusActive {
			active = append(active, *s)
		}
	}
	return active, nil
}

func (m *inMemorySessionRepo) Create(ctx context.Context, s *model.LabSession) error {
	s.ID = fmt.Sprintf("sess-%d", len(m.sessions)+1)
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	m.sessions[s.ID] = s
	return nil
}

func (m *inMemorySessionRepo) Update(ctx context.Context, s *model.LabSession) error {
	if _, ok := m.sessions[s.ID]; !ok {
		return repository.ErrSessionNotFound
	}
	s.UpdatedAt = time.Now()
	m.sessions[s.ID] = s
	return nil
}

func (m *inMemorySessionRepo) UpdateStatus(ctx context.Context, id string, status model.SessionStatus) error {
	if s, ok := m.sessions[id]; ok {
		s.Status = status
		s.UpdatedAt = time.Now()
		return nil
	}
	return repository.ErrSessionNotFound
}

func (m *inMemorySessionRepo) HasAttendance(ctx context.Context, sessionID string) (bool, error) {
	return m.attendance[sessionID], nil
}

func (m *inMemorySessionRepo) Delete(ctx context.Context, id string) error {
	if m.attendance[id] {
		return repository.ErrSessionHasAttendance
	}
	if _, ok := m.sessions[id]; !ok {
		return repository.ErrSessionNotFound
	}
	delete(m.sessions, id)
	return nil
}

// ==========================================
// TEST SUITE: LAB SESSIONS
// ==========================================

// 1. Admin Create Session Valid
func TestAdminCreateSessionValid(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	payload := `{
		"code": "JAR-01",
		"title": "Praktikum Jaringan Komputer",
		"course_code": "IF3102",
		"location": "Lab Komputer 3",
		"session_date": "2026-09-20",
		"start_at": "2026-09-20T10:00:00Z",
		"end_at": "2026-09-20T12:30:00Z",
		"status": "SCHEDULED"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201 Created, got %d: %s", resp.StatusCode, string(body))
	}

	var res model.SessionResponse
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &res)

	if res.Code != "JAR-01" {
		t.Errorf("expected code JAR-01, got %s", res.Code)
	}
	if res.Status != model.SessionStatusScheduled {
		t.Errorf("expected status SCHEDULED, got %s", res.Status)
	}
	if res.Location != "Lab Komputer 3" {
		t.Errorf("expected location Lab Komputer 3, got %s", res.Location)
	}
}

// 2. Validation: Time Range, Formats, and Duplicate Code
func TestCreateSessionValidation(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	// a. Waktu mulai setelah waktu selesai (start_at >= end_at)
	invalidTimePayload := `{
		"code": "ERR-01",
		"title": "Praktikum Salah Waktu",
		"location": "Lab 1",
		"session_date": "2026-09-20",
		"start_at": "2026-09-20T14:00:00Z",
		"end_at": "2026-09-20T13:00:00Z"
	}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(invalidTimePayload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for start_at >= end_at, got %d", resp1.StatusCode)
	}
	body1, _ := io.ReadAll(resp1.Body)
	var errResp1 model.ErrorResponse
	_ = json.Unmarshal(body1, &errResp1)
	if !strings.Contains(errResp1.Error.Message, "setelah") {
		t.Errorf("expected error message explaining time range, got: %s", errResp1.Error.Message)
	}

	// b. Format tanggal tidak valid
	invalidDatePayload := `{
		"code": "ERR-02",
		"title": "Praktikum Tanggal Salah",
		"location": "Lab 1",
		"session_date": "20-09-2026",
		"start_at": "2026-09-20T10:00:00Z",
		"end_at": "2026-09-20T12:00:00Z"
	}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(invalidDatePayload))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid date, got %d", resp2.StatusCode)
	}

	// c. Duplicate code (PWB-01 sudah ada)
	dupCodePayload := `{
		"code": "PWB-01",
		"title": "Praktikum Duplikat Kode",
		"location": "Lab 1",
		"session_date": "2026-09-20",
		"start_at": "2026-09-20T10:00:00Z",
		"end_at": "2026-09-20T12:00:00Z"
	}`
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(dupCodePayload))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp3, _ := app.Test(req3)

	if resp3.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate session code, got %d", resp3.StatusCode)
	}
}

// 3. RBAC: Student Cannot Manage Sessions
func TestStudentCannotManageSessions(t *testing.T) {
	app, helper := setupTestApp()
	studentToken := helper.getTokenFor("student-uuid-1")

	// a. Student POST /api/v1/sessions -> 403
	createPayload := `{
		"code": "STU-01",
		"title": "Sesi Mahasiswa",
		"location": "Lab 1",
		"session_date": "2026-09-20",
		"start_at": "2026-09-20T10:00:00Z",
		"end_at": "2026-09-20T12:00:00Z"
	}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewBufferString(createPayload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls POST /api/v1/sessions, got %d", resp1.StatusCode)
	}

	// b. Student PUT /api/v1/sessions/:id -> 403
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/sessions/sess-1", bytes.NewBufferString(`{"title":"Edit Sesi"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls PUT /api/v1/sessions/sess-1, got %d", resp2.StatusCode)
	}

	// c. Student PATCH /api/v1/sessions/:id/status -> 403
	req3 := httptest.NewRequest(http.MethodPatch, "/api/v1/sessions/sess-1/status", bytes.NewBufferString(`{"status":"CLOSED"}`))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls PATCH /api/v1/sessions/sess-1/status, got %d", resp3.StatusCode)
	}

	// d. Student DELETE /api/v1/sessions/:id -> 403
	req4 := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/sess-1", nil)
	req4.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp4, _ := app.Test(req4)
	if resp4.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when student calls DELETE /api/v1/sessions/sess-1, got %d", resp4.StatusCode)
	}

	// e. Student GET /api/v1/sessions -> 200 OK (Read-only diizinkan)
	req5 := httptest.NewRequest(http.MethodGet, "/api/v1/sessions", nil)
	req5.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp5, _ := app.Test(req5)
	if resp5.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student reads sessions list, got %d", resp5.StatusCode)
	}

	// f. Student GET /api/v1/sessions/active -> 200 OK
	req6 := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/active", nil)
	req6.AddCookie(&http.Cookie{Name: "token", Value: studentToken})
	resp6, _ := app.Test(req6)
	if resp6.StatusCode != http.StatusOK {
		t.Errorf("expected 200 when student reads active sessions, got %d", resp6.StatusCode)
	}
}

// 4. Session Status Lifecycle & Update
func TestSessionStatusLifecycle(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	// a. Ubah status sesi dari SCHEDULED menjadi ACTIVE
	patchPayload := `{"status":"ACTIVE"}`
	req1 := httptest.NewRequest(http.MethodPatch, "/api/v1/sessions/sess-2/status", bytes.NewBufferString(patchPayload))
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on status patch, got %d", resp1.StatusCode)
	}
	var res1 model.SessionResponse
	b1, _ := io.ReadAll(resp1.Body)
	_ = json.Unmarshal(b1, &res1)
	if res1.Status != model.SessionStatusActive {
		t.Errorf("expected status ACTIVE, got %s", res1.Status)
	}

	// b. Cek endpoint GET /api/v1/sessions/active kini mencakup sess-2
	reqActive := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/active", nil)
	reqActive.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	respActive, _ := app.Test(reqActive)
	var activeList []model.SessionResponse
	bActive, _ := io.ReadAll(respActive.Body)
	_ = json.Unmarshal(bActive, &activeList)

	foundSess2 := false
	for _, s := range activeList {
		if s.ID == "sess-2" {
			foundSess2 = true
			break
		}
	}
	if !foundSess2 {
		t.Errorf("expected sess-2 in active sessions list")
	}

	// c. Ubah status menjadi CLOSED
	patchClosed := `{"status":"CLOSED"}`
	req2 := httptest.NewRequest(http.MethodPatch, "/api/v1/sessions/sess-2/status", bytes.NewBufferString(patchClosed))
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := app.Test(req2)
	var res2 model.SessionResponse
	b2, _ := io.ReadAll(resp2.Body)
	_ = json.Unmarshal(b2, &res2)
	if res2.Status != model.SessionStatusClosed {
		t.Errorf("expected status CLOSED, got %s", res2.Status)
	}
}

// 5. Deletion Rule: Session Referenced by Attendance Cannot be Physically Deleted
func TestSessionDeletionWithAttendance(t *testing.T) {
	app, helper := setupTestApp()
	adminToken := helper.getTokenFor("admin-uuid-1")

	// sess-1 memiliki data attendance (m.attendance["sess-1"] = true)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/sess-1", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict when deleting session with attendance, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp model.ErrorResponse
	_ = json.Unmarshal(body, &errResp)

	if !strings.Contains(errResp.Error.Message, "CANCELLED") && !strings.Contains(errResp.Error.Message, "CLOSED") {
		t.Errorf("expected error message to guide admin to use CANCELLED or CLOSED, got: %s", errResp.Error.Message)
	}

	// sess-2 tidak memiliki data attendance -> deletion diizinkan
	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/sess-2", nil)
	req2.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK when deleting session without attendance, got %d", resp2.StatusCode)
	}
}
